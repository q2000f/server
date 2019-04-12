package logic

import "sync"

type User struct {
	sync.Locker
	ID string
	State UserState
}

type UserState int
const (
	UserState_Init UserState = iota
	UserState_Login
	UserState_InGame
	UserState_OffLine
)

var users sync.Map

func GetUser(id string) *User {
	user, ok := users.Load(id)
	if ok {
		return user.(*User)
	}
	newUser := &User{ID:id, State: UserState_Init}
	user, _ = users.LoadOrStore(id, newUser)
	return user.(*User)
}
