package provider

import "fmt"

var RedisProviderType = "redis"
var MemoryProviderType = "memory"
var ErrNotSupportProviderType = fmt.Errorf("not support this provider type")

// Provides global providers map
var Provides = make(map[string]Provider)

// Provider
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

type Config struct {
	ProviderType string
	// option provider endpoint. such as redis cluster: 127.0.0.1:6379,127.0.0.1:2380
	Endpoints string
	// option provider auth user
	User string
	// option provider auth password
	Password string
}

func NewProvider(config Config) (Provider, error) {
	if config.ProviderType == MemoryProviderType {
		return NewMemoryProvider(), nil
	} else if config.ProviderType == RedisProviderType {
		return NewRedisProvider(config.Endpoints, config.Password), nil
	} else {
		return nil, ErrNotSupportProviderType
	}
}
