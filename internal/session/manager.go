package session

import (
	"encoding/json"
	"password_store/internal/util"
	"time"
)

type SessionManager struct {
	cookieName                  string
	expirationDurationInSeconds int
	randomStringLength          int
	sessionStore                SessionStore
}

func NewDefaultManager(sessionStore SessionStore) *SessionManager {
	return &SessionManager{
		cookieName:                  "default_cookie_name",
		expirationDurationInSeconds: 5 * 3600,
		randomStringLength:          20,
		sessionStore:                sessionStore,
	}
}

func (m *SessionManager) GetCookieName() string {
	return m.cookieName
}

func (m *SessionManager) GetExpiryDuration() int {
	return m.expirationDurationInSeconds
}

func (m *SessionManager) GetRandomStringLength() int {
	return m.randomStringLength
}

func (m *SessionManager) SetCookieName(name string) {
	m.cookieName = name
}

func (m *SessionManager) SetExpiryDuration(duration int) {
	m.expirationDurationInSeconds = duration
}

func (m *SessionManager) SetRandomStringLength(length int) {
	m.randomStringLength = length
}

func (m *SessionManager) createSessionId(username string) string {
	time := time.Now().String()
	randomString := util.GenerateRandomString(m.randomStringLength)

	// Format: session_manager:<time>:<username>:<random_string>
	sessionId := "session_manager:" + time + ":" + username + ":" + randomString
	// fmt.Println("sessionId: " + sessionId)
	hashedSessionId := util.Hash([]byte(sessionId))
	// fmt.Println("hashedSessionId: " + string(hashedSessionId))
	return string(hashedSessionId)
}

func (m *SessionManager) SetSession(username string) (string, error) {
	cookie := CreateSession(username, m.expirationDurationInSeconds)
	cookieBytes, err := json.Marshal(cookie)
	if err != nil {
		return "", err
	}

	// Right a mechanism to generate the sessionId
	sessionId := m.createSessionId(username)
	if err = m.sessionStore.Set(sessionId, cookieBytes, m.expirationDurationInSeconds); err != nil {
		return "", err
	}
	return sessionId, nil
}

func (m *SessionManager) GetSession(sessionId string) (Session, error) {
	cookie := CreateSession("", 0)

	cookieBytes, err := m.sessionStore.Get(sessionId)
	if err != nil {
		return cookie, err
	}

	if err = json.Unmarshal(cookieBytes, &cookie); err != nil {
		return cookie, err
	}
	return cookie, nil
}

func (m *SessionManager) DeleteSession(sessionId string) error {
	return m.sessionStore.Delete(sessionId)
}
