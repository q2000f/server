package logic

import (
	"fmt"
	"server/game/code"
	"server/game/proto"
)

func Login(login *proto.ILogin) (int64, proto.OLogin) {
	fmt.Println("Login success!")
	return code.OK, proto.OLogin{PID: login.AID, SID: "sid"}
}

func GetPlayer(pid string) bool {
	user := GetUser(pid)
	user.Lock()
	defer user.Unlock()
	if user.State != UserState_Init {
		return false
	}

	user.State = UserState_Login
	//ToDO: load user data
	user.State = UserState_InGame
	return true
}
