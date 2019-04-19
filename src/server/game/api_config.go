package game

import (
	"server/game/handler"
)

type Callback func(data []byte) (int64, interface{})
type ApiConfig struct {
	Opcode string
	Do Callback
}

var handlers map[string]ApiConfig
var apiConfigs []ApiConfig

func init() {
	apiConfigs = []ApiConfig {
		{
			"login",
			handler.API_Login,
		},
		{
			"getPlayer",
			handler.API_GetPlayer,
		},
	}

	handlers = map[string]ApiConfig{}
	for _, apiConfig := range apiConfigs {
		handlers[apiConfig.Opcode] = apiConfig
	}
}
