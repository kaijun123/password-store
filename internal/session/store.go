package session

// interface for all stores
type SessionStore interface {
	Set(key string, value []byte, expiry int) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
