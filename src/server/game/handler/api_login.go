package handler

import (
	"encoding/json"
	"server/game/logic"
	"server/game/proto"
)

func API_Login(data []byte) (int64, interface{}) {
	iLogin := &proto.ILogin{}
	json.Unmarshal(data, iLogin)
	return logic.Login(iLogin)
}
