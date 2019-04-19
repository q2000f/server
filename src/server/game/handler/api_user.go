package handler

import (
	"encoding/json"
	"server/game/logic"
	"server/game/proto"
)

func API_GetPlayer(data []byte) (int64, interface{}) {
	iGetUser := &proto.IGetPlayer{}
	json.Unmarshal(data, iGetUser)
	return logic.GetPlayer(iGetUser)
}
