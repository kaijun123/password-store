package kvStore

import "time"

type Idempotency struct {
	Key       string    // user-generated key
	SessionId string    // session id corresponding to the user
	Status    string    // request status
	Expiry    time.Time // exact expiry time
	Request   []byte    // request body converted to bytes
	Hash      []byte    // hash of the request. Used to verify that request made matches the particular idempotency key
}

func CreateIdempotency(key string, sessionId string, status string, expiry int, request []byte, hash []byte) Idempotency {
	return Idempotency{
		Key:       key,
		SessionId: sessionId,
		Status:    status,
		Expiry:    time.Now().Add(time.Duration(expiry) * time.Second), // calculate the expiry time
		Request:   request,
		Hash:      hash,
	}
}

func (s *Idempotency) IsExpired() bool {
	return time.Now().After(s.Expiry)
}
