package session

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/url"
	"session/session/provider"
	"sync"
	"time"
)

//var defaultSessionProvider = MemoryProviderType
//
//// GlobalSessionManager
//var GlobalSessionManager *Manager

//func init() {
//	GlobalSessionManager, _ = NewManager(defaultSessionProvider, "sessionid", 3600)
//	go GlobalSessionManager.GC()
//}

// Manager session manager
type Manager struct {
	cookieName  string
	provider    provider.Provider
	lock        sync.Mutex
	maxLifeTime int64
}

// NewManager
func NewManager(conf provider.Config, cookieName string, maxLifeTime int64) (*Manager, error) {
	p, err := provider.NewProvider(conf)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		cookieName:  cookieName,
		provider:    p,
		lock:        sync.Mutex{},
		maxLifeTime: maxLifeTime,
	}
	return m, nil
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session provider.Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	sid, _ := url.QueryUnescape(cookie.Value)
	if err := manager.provider.SessionDestroy(sid); err != nil {
		return
	}
	expiration := time.Now()
	newCookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, &newCookie)

}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
