package game

import (
	"server/game/common"
	"server/game/handler"
)

type Callback func(ctx common.Context, data []byte) (int64, interface{})
type ApiConfig struct {
	Opcode           string
	Do               Callback
	SkipCheckSession bool
}

var handlers map[string]ApiConfig
var apiConfigs []ApiConfig

func init() {
	apiConfigs = []ApiConfig{
		{
			Opcode:           "login",
			Do:               handler.API_Login,
			SkipCheckSession: true,
		},
		{
			Opcode: "getPlayer",
			Do:     handler.API_GetPlayer,
		},
		{
			Opcode: "createPlayer",
			Do:     handler.API_CreatePlayer,
		},
	}

	handlers = map[string]ApiConfig{}
	for _, apiConfig := range apiConfigs {
		handlers[apiConfig.Opcode] = apiConfig
	}
}
