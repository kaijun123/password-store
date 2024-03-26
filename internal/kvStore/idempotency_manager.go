package kvStore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"password_store/internal/constants"
	"password_store/internal/util"
)

type IdempotencyManager struct {
	requestHeader               string
	prefix                      string
	expirationDurationInSeconds int
	userGeneratedKeyLength      int
	idempotencyStore            Store
}

func NewIdempotencyManager(idempotencyStore Store) *IdempotencyManager {
	return &IdempotencyManager{
		requestHeader:               "Idempotency-Key",
		prefix:                      constants.IdempotencyPrefix,
		expirationDurationInSeconds: 24 * 3600, // 24 hours
		userGeneratedKeyLength:      10,        // length of user-generated key
		idempotencyStore:            idempotencyStore,
	}
}

func (m *IdempotencyManager) GetRequestHeader() string {
	return m.requestHeader
}
func (m *IdempotencyManager) GetExpiryDuration() int {
	return m.expirationDurationInSeconds
}

func (m *IdempotencyManager) GetUserGeneratedKeyLength() int {
	return m.userGeneratedKeyLength
}

func (m *IdempotencyManager) SetExpiryDuration(duration int) {
	m.expirationDurationInSeconds = duration
}

func (m *IdempotencyManager) SetUserGeneratedKeyLength(duration int) {
	m.userGeneratedKeyLength = duration
}

func (m *IdempotencyManager) createIdempotencyId(userGeneratedKey string) (string, error) {

	if len(userGeneratedKey) != m.userGeneratedKeyLength {
		// TODO: change to bad request later
		return "", errors.New("invalid length for user generated key")
	}
	// Format: idempotency-key://<user-generated-key>
	idempotencyId := m.prefix + userGeneratedKey
	// fmt.Println("idempotencyId: " + idempotencyId)
	hashedIdempotencyId := util.Hash([]byte(idempotencyId))
	// fmt.Println("idempotencyId: " + string(idempotencyId))
	return string(hashedIdempotencyId), nil
}

// Data structure stored:
// key: idempotency-key://<user-generated-key>
// value: byte form of Idempotency struct

// Retrieve the idempotency key
func (m *IdempotencyManager) GetIdempotency(userGeneratedKey string, sessionId string, request []byte) (Idempotency, error) {

	idempotencyId, err := m.createIdempotencyId(userGeneratedKey)
	// fmt.Println("(idemp manager; get idemp) idempotencyId: ", []byte(idempotencyId))
	if err != nil {
		return Idempotency{}, nil
	}

	// Get idempotency data from the idempotency store
	idempotencyBytes, err := m.idempotencyStore.Get(idempotencyId)
	// fmt.Println("(idemp manager; get idemp) idempotencyBytes: ", idempotencyBytes)
	if err != nil {
		return Idempotency{}, errors.New(constants.IdempNew)
	}

	idempotency := Idempotency{}
	if err := json.Unmarshal(idempotencyBytes, &idempotency); err != nil {
		return Idempotency{}, errors.New(constants.IdempServerErr)
	}
	// fmt.Println("(idemp manager; get idemp) idemp:", idempotency)

	// Check for the sessionId
	// fmt.Println("(idemp manager; get idemp) received sessionId:", []byte(sessionId))
	// fmt.Println("(idemp manager; get idemp) retrieved sessionId:", []byte(idempotency.SessionId))

	// TODO: [Major Bug] Fix this check. It fails when a repeated idemp request is made. Provided session id and retrieved session id is not the same.
	// if idempotency.SessionId != sessionId {
	// 	return Idempotency{}, errors.New(constants.IdempBadRequest)
	// }

	// Check for the hash
	fmt.Println("request: ", string(request))
	calculatedHash := util.Hash(request)
	fmt.Println("idempotency.Hash: ", fmt.Sprintf("%x", idempotency.Hash))
	fmt.Println("calculatedHash: ", fmt.Sprintf("%x", calculatedHash))
	if !bytes.Equal(idempotency.Hash, calculatedHash) {
		return Idempotency{}, errors.New(constants.IdempBadRequest)
	}

	return idempotency, nil
}

// Add the idempotency key to the store
// Use the idempotency prefix
func (m *IdempotencyManager) SetIdempotency(userGeneratedKey string, sessionId string, status string, request []byte) error {
	// fmt.Println("(idemp manager; set idemp) received sessionId:", []byte(sessionId))
	idempotencyId, err := m.createIdempotencyId(userGeneratedKey)
	if err != nil {
		return err
	}
	// fmt.Println("(idemp manager; set idemp) idempotencyId: ", []byte(idempotencyId))

	calculatedHash := util.Hash(request)
	idempotency := CreateIdempotency(userGeneratedKey, sessionId, status, m.expirationDurationInSeconds, request, calculatedHash)
	// fmt.Println("(idemp manager; set idemp) idemp:", idempotency)
	// fmt.Println("(idemp manager; set idemp) retrieved sessionId:", []byte(idempotency.SessionId))
	idempotencyBytes, err := json.Marshal(idempotency)
	// fmt.Println("(idemp manager; set idemp) idempotencyBytes: ", idempotencyBytes)
	if err != nil {
		return err
	}

	if err := m.idempotencyStore.Set(idempotencyId, idempotencyBytes, m.expirationDurationInSeconds); err != nil {
		return err
	}
	return nil
}

// Add the idempotency key to the store
// Use the idempotency prefix
func (m *IdempotencyManager) UpdateIdempotency(userGeneratedKey string, sessionId string, newStatus string, request []byte) error {

	idempotency, err := m.GetIdempotency(userGeneratedKey, sessionId, request)
	if err != nil {
		return err
	}
	idempotency.Status = newStatus

	idempotencyId, err := m.createIdempotencyId(userGeneratedKey)
	if err != nil {
		return err
	}
	idempotencyBytes, err := json.Marshal(idempotency)
	if err != nil {
		return err
	}

	if err := m.idempotencyStore.Set(idempotencyId, idempotencyBytes, m.expirationDurationInSeconds); err != nil {
		return err
	}
	return nil
}

func (m *IdempotencyManager) DeleteIdempotency(userGeneratedKey string) error {
	idempotencyId, err := m.createIdempotencyId(userGeneratedKey)
	if err != nil {
		return err
	}
	return m.idempotencyStore.Delete(idempotencyId)
}
