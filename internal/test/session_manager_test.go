package test

import (
	"password_store/internal/kvStore"
	"password_store/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookieName string = "test_cookie_name"
const cookieExpiryDuration int = 10
const cookieRandomStringLength int = 10
const username string = "test_username"

func TestManager(t *testing.T) {
	util.SetEnv("REDIS_HOST", "localhost")
	util.SetEnv("REDIS_PORT", "6379")

	redis := &kvStore.Redis{}
	redis.CreateClient()

	sessionManager := kvStore.NewSessionManager(redis)
	sessionManager.SetCookieName(cookieName)
	sessionManager.SetExpiryDuration(cookieExpiryDuration)
	sessionManager.SetRandomStringLength(cookieRandomStringLength)

	sessionId, err := sessionManager.SetSession(username)
	assert.Nil(t, err)
	session, err := sessionManager.GetSession(sessionId)
	assert.Nil(t, err)
	assert.Equal(t, session.Username, username)

	err = sessionManager.DeleteSession(sessionId)
	assert.Nil(t, err)

	_, err = sessionManager.GetSession(sessionId)
	assert.NotNil(t, err)
}
