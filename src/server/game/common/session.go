package common

import (
	"sync"
	"time"
)

//ToDO: replace by session
type Session struct {
	ID       uint64
	AID      string
	PID      uint64
	CreateAt int64
	UpdateAt int64
}

type SessionMap struct {
	sync.Map
}

var GSessionMap SessionMap

func (sessions *SessionMap) GetSession(id uint64) *Session {
	sess, ok := sessions.Load(id)
	if ok {
		return sess.(*Session)
	}
	return nil
}

func (sessions *SessionMap) CreateSession(sess *Session) {
	sessions.Store(sess.ID, sess)
}

func (sessions *SessionMap) RemoveSession(id uint64) {
	sessions.Delete(id)
}

func (sessions *SessionMap) Update() {
	var expires []uint64
	now := time.Now().Unix()
	sessions.Range(func(key, value interface{}) bool {
		sess := value.(*Session)
		if now-sess.UpdateAt > 60 {
			expires = append(expires, key.(uint64))
		}
		return true
	})
	for _, id := range expires {
		sessions.RemoveSession(id)
	}
}
