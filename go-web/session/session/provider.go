package session

// Provides global providers map
var Provides = make(map[string]Provider)

// Provider interface
type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

// Session interface
type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	Delete(key interface{}) error
	SessionId() string
}

// Registry registry provider
func Registry(name string, provider Provider) {
	if provider == nil {
		panic("session: Registry provider is nil")
	}
	if _, dup := Provides[name]; dup {
		panic("session: Registry called twice for provider " + name)
	}
	Provides[name] = provider
}
