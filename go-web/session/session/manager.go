package session

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/learning-go/go-web/session/session/provider"
)

var defaultSessionProvider = provider.MemoryProviderType
var globalSessionManager *Manager

func init() {
	globalSessionManager, _ = NewManager(defaultSessionProvider, "sessionid", 3600)
	go globalSessionManager.GC()
}

// Manager session manager
type Manager struct {
	cookieName  string
	provider    Provider
	lock        sync.Mutex
	maxLifeTime int64
}

// NewManager
func NewManager(providerName, cookieName string, maxLifeTime int64) (*Manager, error) {
	p, ok := Provides[providerName]
	if !ok {
		return nil, fmt.Errorf("not support provider type for %s", providerName)
	}

	return &Manager{
		cookieName:  cookieName,
		provider:    p,
		lock:        sync.Mutex{},
		maxLifeTime: maxLifeTime,
	}, nil
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) sessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
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

func (manager *Manager) sessionDestroy(w http.ResponseWriter, r *http.Request) {
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
