package logic

import (
	"server/game/common"
	"server/game/proto"
	"util"
)

func GetPlayer(ctx common.Context, iGetPlayer *proto.IGetPlayer) (code int64, oGetPlayer proto.OGetPlayer) {
	player := common.GetPlayerByAid(ctx.Sess.AID)
	if player != nil {
		tmpPlayer := proto.Player{}
		util.Copy(player, &tmpPlayer)
		oGetPlayer.Players = append(oGetPlayer.Players, tmpPlayer)
	}
	return code, oGetPlayer
}

func CreatePlayer(ctx common.Context, in *proto.ICreatePlayer) (code int64, out proto.OCreatePlayer) {
	//check name
	id := util.GetNewID()
	user := common.GetUser(id)
	user.Lock()
	defer user.Unlock()
	player := common.Player{
		ID:    id,
		AID:   ctx.Sess.AID,
		Name:  in.Name,
		Level: 1,
	}
	user.Player = &player

	out.PlayerID = id
	out.Name = in.Name
	return 0, out
}
