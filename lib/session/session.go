package session

import (
	"errors"

	"github.com/oklog/ulid/v2"
)

type session struct {
	ID     ulid.ULID
	userID ulid.ULID
}

type Session struct {
	sessions map[string]session
}

func New() *Session {
	return &Session{
		sessions: map[string]session{},
	}
}

func (s *Session) Create(userId ulid.ULID) string {
	ses := session{
		ID:     ulid.Make(),
		userID: userId,
	}
	s.sessions[ses.ID.String()] = ses

	return ses.ID.String()
}

func (s *Session) Get(sessionId string) (ulid.ULID, error) {
	session, ok := s.sessions[sessionId]
	if !ok {
		return ulid.ULID{}, errors.New("session not found")
	}

	return session.userID, nil
}

func (s *Session) Delete(sessionId string) {
	delete(s.sessions, sessionId)
}
