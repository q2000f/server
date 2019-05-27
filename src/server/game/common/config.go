package common

import "util"

type Config struct {
	SeverID uint16
	Port    int
}

var GConfig Config

func InitConfig(filePath string) {
	err := util.LoadJsonFile(filePath, &GConfig)
	if err != nil {
		panic(err)
	}
}
