package handler

import (
	"encoding/json"
	"server/game/common"
	"server/game/logic"
	"server/game/proto"
)

func API_GetPlayer(ctx common.Context, data []byte) (int64, interface{}) {
	iGetUser := &proto.IGetPlayer{}
	json.Unmarshal(data, iGetUser)
	return logic.GetPlayer(ctx, iGetUser)
}

func API_CreatePlayer(ctx common.Context, data []byte) (int64, interface{}) {
	in := &proto.ICreatePlayer{}
	json.Unmarshal(data, in)
	return logic.CreatePlayer(ctx, in)
}
