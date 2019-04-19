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
