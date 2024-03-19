package kvStore

// interface for all stores
type Store interface {
	Set(key string, value []byte, expiry int) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
