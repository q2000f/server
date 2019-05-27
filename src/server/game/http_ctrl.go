package game

import (
	"encoding/json"
	"log"
	"server/game/code"
	"server/game/common"
)

type HttpHeader struct {
	Opcode string
	Host   string
	SID    uint64
}

type HttpResponse struct {
	ErrorCode int64
	Body      interface{}
}

func HttpDo(header HttpHeader, data []byte) (result string) {
	defer func() {
		//if r := recover(); r != nil {
		//	response := HttpResponse{
		//		ErrorCode: code.LogicPanic,
		//	}
		//	bytes, _ := json.Marshal(response)
		//	result = string(bytes)
		//	log.Println(r)
		//}
	}()
	log.Print("opcode:", header.Opcode, " sid:", header.SID, " host:", header.Host)
	apiConfig, ok := handlers[header.Opcode]
	if !ok {
		response := HttpResponse{
			ErrorCode: code.InvalidOpcode,
		}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	var sess *common.Session
	if !apiConfig.SkipCheckSession {
		sess = common.GSessionMap.GetSession(header.SID)
		if sess == nil {
			log.Println("error sess == nil")
			return
		}
	}

	ctx := common.Context{Sess: sess}

	errCode, ret := apiConfig.Do(ctx, data)
	response := HttpResponse{
		ErrorCode: errCode,
	}

	if errCode != code.OK {
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	response.Body = ret
	bytes, _ := json.Marshal(response)
	return string(bytes)
}
