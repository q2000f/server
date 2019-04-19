package logic

import (
	"server/game/logic/common"
	"server/game/proto"
)

func GetPlayer(iGetPlayer *proto.IGetPlayer) (code int64, oGetPlayer proto.OGetPlayer){
	user := common.GetUser(iGetPlayer.PID)
	user.Lock()
	defer user.Unlock()

	player := proto.Player{
		ID: user.Player.ID,
		Name: user.Player.Name,
	}

	oGetPlayer.Players = append(oGetPlayer.Players, player)
	return code, oGetPlayer
}
