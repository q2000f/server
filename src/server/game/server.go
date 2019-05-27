package game

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"server/game/common"
	"strconv"
	"strings"
	"util"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "-config=config.json")
}

func Start() {
	flag.Parse()

	common.InitConfig(configFile)
	util.InitUUID(func() (u uint16, e error) {
		return common.GConfig.SeverID, nil
	})

	r := gin.Default()
	r.GET("/game/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/game/:opcode", func(c *gin.Context) {
		fmt.Println("path:", c.Request.URL.Path, "host:", c.Request.URL.Host)
		data, _ := ioutil.ReadAll(c.Request.Body)
		opcode := strings.Replace(c.Request.URL.Path, "/game/", "", -1)

		sidHeader := c.Request.Header.Get("sid")
		var sid uint64
		if sidHeader != "" {
			var err error
			sid, err = strconv.ParseUint(sidHeader, 10, 64)
			if err != nil {
				c.String(200, err.Error())
				return
			}
		}

		header := HttpHeader{
			Opcode: opcode,
			Host:   c.Request.URL.Host,
			SID:    sid,
		}
		ret := HttpDo(header, data)
		c.String(200, string(ret))
	})
	url := fmt.Sprintf(":%d", common.GConfig.Port)
	r.Run(url)
}
