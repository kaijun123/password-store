package session

import (
	"password_store/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

const key string = "foo"
const value string = "bar"

func TestRedis(t *testing.T) {
	util.SetEnv("REDIS_HOST", "localhost")
	util.SetEnv("REDIS_PORT", "6379")

	redis := &Redis{}
	redis.CreateClient()
	err := redis.Set(key, []byte(value), 10)
	assert.Nil(t, err)

	value, err := redis.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, []byte(value))

	err = redis.Delete(key)
	assert.Nil(t, err)

	// Expect for there to be an error since the key-value pair should be deleted
	_, err = redis.Get(key)
	assert.NotNil(t, err)
}
