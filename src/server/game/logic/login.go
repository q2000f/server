package logic

import (
	"fmt"
	"server/game/code"
	"server/game/common"
	"server/game/proto"
	"time"
	"util"
)

func Login(ctx common.Context, login *proto.ILogin) (errorCode int64, oLogin proto.OLogin) {
	fmt.Println("Login success!")
	sid := util.GetNewID()

	sess := &common.Session{
		ID:       sid,
		AID:      login.AID,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}

	common.GSessionMap.CreateSession(sess)

	return code.OK, proto.OLogin{PID: login.AID, SID: sid}
}
