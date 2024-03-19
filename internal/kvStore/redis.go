package kvStore

import (
	"context"
	"log"
	"password_store/internal/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func (r *Redis) CreateClient() {
	redisHost, err := util.GetEnv("REDIS_HOST")
	if err != nil {
		log.Fatalln(err.Error())
	}
	redisPort, err := util.GetEnv("REDIS_PORT")
	if err != nil {
		log.Fatalln(err.Error())
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r.ctx = context.Background()
}

// Sets an expiry when adding the key-value pair into redis
// Redis will remove the key-value pair when the expiry date reaches
func (r *Redis) Set(key string, value []byte, expiry int) error {
	expiryDuration := time.Duration(expiry) * time.Minute
	if err := r.client.Set(r.ctx, key, value, expiryDuration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	value, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return []byte(""), err
	}
	return []byte(value), nil
}

func (r *Redis) Delete(key string) error {
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// Compile-time check to ensure *Redis implements SessionStore interface
var _ Store = &Redis{}
