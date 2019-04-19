package common

import (
	"sync"
)

type Player struct {
	ID string
	Name string
	Level string
}

type User struct {
	sync.RWMutex
	ID string
	State UserState
	Player Player
}

type UserState int
const (
	UserState_Init UserState = iota
	UserState_Login
	UserState_InGame
	UserState_OffLine
)

var Users sync.Map

func GetUser(id string) *User {
	user, ok := Users.Load(id)
	if ok {
		return user.(*User)
	}
	newUser := &User{ID:id, State: UserState_Init}
	user, _ = Users.LoadOrStore(id, newUser)
	return user.(*User)
}
