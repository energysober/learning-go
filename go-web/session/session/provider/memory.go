package provider

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

//var MemoryProviderType = "memory"
//
var pder = &MemoryProvider{}

type MemoryProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (p *MemoryProvider) SessionInit(sid string) (Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSess := &MemorySessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := p.list.PushBack(newSess)
	p.sessions[sid] = element
	return newSess, nil
}

func (p *MemoryProvider) SessionRead(sid string) (Session, error) {
	if element, ok := p.sessions[sid]; ok {
		return element.Value.(*MemorySessionStore), nil
	} else {
		return p.SessionInit(sid)
	}
}

func (p *MemoryProvider) SessionDestroy(sid string) error {
	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
		return nil
	}
	return nil
}

func (p *MemoryProvider) SessionGC(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		element := p.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemorySessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			p.list.Remove(element)
			delete(p.sessions, element.Value.(*MemorySessionStore).sid)
		} else {
			break
		}
	}
}

func (p *MemoryProvider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if element, ok := p.sessions[sid]; ok {
		p.sessions[sid].Value.(*MemorySessionStore).timeAccessed = time.Now()
		p.list.MoveToFront(element)
	}
	return nil
}

type MemorySessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (s *MemorySessionStore) Set(key, value interface{}) error {
	s.value[key] = value
	err := pder.SessionUpdate(s.sid)
	return err
}

func (s *MemorySessionStore) Get(key interface{}) (interface{}, error) {
	if value, ok := s.value[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("can not found")
}

func (s *MemorySessionStore) Delete(key interface{}) error {
	delete(s.value, key)
	return nil
}

func (s *MemorySessionStore) SessionId() string {
	return s.sid
}

// NewMemoryProvider
func NewMemoryProvider() Provider {
	if pder.list == nil {
		pder.list = list.New()
	}
	if pder.sessions == nil {
		pder.sessions = make(map[string]*list.Element, 0)
	}
	return pder
}
