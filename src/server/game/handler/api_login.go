package handler

import (
	"encoding/json"
	"server/game/common"
	"server/game/logic"
	"server/game/proto"
)

func API_Login(ctx common.Context, data []byte) (int64, interface{}) {
	iLogin := &proto.ILogin{}
	json.Unmarshal(data, iLogin)
	return logic.Login(ctx, iLogin)
}
