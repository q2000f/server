package main

import "sync"

type User struct {
	initLock sync.Once
}

func (u* User) Init() {
	u.initLock.Do(func() {
		u.Load()
	})
}

func (u* User) Load() {

}

var UserMap sync.Map

func main() {
	playerid := "id1"
	user, ok := UserMap.LoadOrStore(playerid, &User{})
	if !ok {
		v, isOk := user.(*User)
		if !isOk {
			panic("error type!")
		}
	}
}
