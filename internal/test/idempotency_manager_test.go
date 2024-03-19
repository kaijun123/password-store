package test

import (
	"encoding/json"
	"password_store/internal/kvStore"
	"password_store/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testIdempotencyPrefix                 string = "test-idempotency-key://"
	testIdempotencyExpirationDuration     int    = 10
	testIdempotencyUserGeneratedKeyLength int    = 10
)

type mockRequest struct {
	From   string
	To     string
	Amount float32
}

func TestIdempotency(t *testing.T) {
	util.SetEnv("REDIS_HOST", "localhost")
	util.SetEnv("REDIS_PORT", "6379")

	redis := &kvStore.Redis{}
	redis.CreateClient()

	idempotencyManager := kvStore.NewIdempotencyManager(redis)
	idempotencyManager.SetExpiryDuration(testIdempotencyExpirationDuration)
	idempotencyManager.SetUserGeneratedKeyLength(testIdempotencyUserGeneratedKeyLength)

	testUserGeneratedKey := util.GenerateRandomString(testIdempotencyUserGeneratedKeyLength)
	testSessionId := util.GenerateRandomString(10)

	// Create Request
	req := mockRequest{
		From:   "alice",
		To:     "bob",
		Amount: 1000,
	}

	reqBytes, err := json.Marshal(req)
	assert.Nil(t, err)
	calculatedHash := util.Hash(reqBytes)
	expectedIdempotency := kvStore.CreateIdempotency(testUserGeneratedKey, testSessionId, kvStore.Pending, testIdempotencyExpirationDuration, reqBytes, calculatedHash)

	// Successful Set
	err = idempotencyManager.SetIdempotency(testUserGeneratedKey, testSessionId, reqBytes)
	assert.Nil(t, err)

	// Successful Get
	actualIdempotency, err := idempotencyManager.GetIdempotency(testUserGeneratedKey, testSessionId, calculatedHash)
	assert.Nil(t, err)
	assert.Equal(t, expectedIdempotency.Key, actualIdempotency.Key)
	assert.Equal(t, expectedIdempotency.SessionId, actualIdempotency.SessionId)
	assert.Equal(t, expectedIdempotency.Status, actualIdempotency.Status)
	assert.Equal(t, expectedIdempotency.Request, actualIdempotency.Request)
	assert.Equal(t, expectedIdempotency.Hash, actualIdempotency.Hash)

	// Failed Get (Wrong Hash)
	_, err = idempotencyManager.GetIdempotency(testUserGeneratedKey, testSessionId, []byte("test"))
	assert.NotNil(t, err)

	// Failed Get (Wrong sessionId)
	_, err = idempotencyManager.GetIdempotency(testUserGeneratedKey, "test", calculatedHash)
	assert.NotNil(t, err)

	// Delete
	err = idempotencyManager.DeleteIdempotency(testUserGeneratedKey)
	assert.Nil(t, err)

	// Failed Get (kv does not exist)
	_, err = idempotencyManager.GetIdempotency(testUserGeneratedKey, testSessionId, calculatedHash)
	assert.NotNil(t, err)
}
