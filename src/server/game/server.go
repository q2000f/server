package game

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func Start() {
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
		header := HttpHeader{
			Opcode: opcode,
			Host: c.Request.URL.Host,
			SID: c.Request.Header.Get("sid"),
		}
		ret := HttpDo(header, data)
		c.String(200, string(ret))
	})
	r.Run(":8888")
}
