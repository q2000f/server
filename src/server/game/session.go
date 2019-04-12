package game

import "sync"

//ToDO: replace by session
type Session struct {
	ID string
	PID string
	CreateAt int64
	UpdateAt int64
}

type SessionMap struct {
	sync.Map
}

var Sessions SessionMap

func (sessions *SessionMap) GetSession(id string) *Session {
	sess, ok := sessions.Load(id)
	if ok {
		return sess.(*Session)
	}
	return nil
}

func (sessions *SessionMap) CreateSession(sess *Session) {
	sessions.Store(sess.ID, sess)
}

func (sessions *SessionMap) RemoveSession(id string) {
	sessions.Delete(id)
}
