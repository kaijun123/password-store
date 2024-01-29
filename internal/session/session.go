package session

import "time"

type Session struct {
	Username string
	Expiry   time.Time // exact expiry time
}

func CreateSession(username string, expiry int) Session {
	cookie := Session{
		Username: username,
		Expiry:   time.Now().Add(time.Duration(expiry) * time.Second), // calculate the expiry time
	}
	return cookie
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expiry)
}
