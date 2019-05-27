package common

import (
	"sync"
)

type Player struct {
	ID    uint64
	AID   string
	Name  string
	Level uint8
}

type User struct {
	sync.RWMutex
	ID     uint64
	State  UserState
	Player *Player
}

type UserState int

const (
	UserState_Init UserState = iota
	UserState_Login
	UserState_InGame
	UserState_OffLine
)

var Users sync.Map

func GetUser(id uint64) *User {
	user, ok := Users.Load(id)
	if ok {
		return user.(*User)
	}
	newUser := &User{ID: id, State: UserState_Init}
	user, _ = Users.LoadOrStore(id, newUser)
	return user.(*User)
}

func GetPlayerByAid(aid string) *Player {
	var ret *Player
	Users.Range(func(key, value interface{}) bool {
		user := value.(*User)
		if user.Player != nil && user.Player.AID == aid {
			ret = user.Player
			return false
		}

		return true
	})
	return ret
}
