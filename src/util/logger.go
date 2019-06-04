package util

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

//example
//func main() {
//	log := NewLogger("./data", "item")
//	myLog.Println("hello world")
//}
const (
	LOG_TO_SCREEN = 1 << iota
	LOG_TO_FILE
)

const (
	log_debug = 1 << iota
	log_info
	log_warn
	log_error
	log_fatal
)

var ErrorLog Logger

type Logger struct {
	logger   *log.Logger // 日志组件
	curDay   int
	fileName string
	logPath  string
	flag     int
	isClose  bool
	level    int
	tag      int
	prefix   string
}

func NewLogger(logPath string, fileName string) *Logger {
	c := Logger{}
	c.logPath = logPath
	c.fileName = fileName
	c.flag = log.Ltime
	log.SetFlags(log.Ltime)
	c.tag = LOG_TO_SCREEN | LOG_TO_FILE

	//in docker can be ignore
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		fmt.Println("os.MkdirAll ", err)
	}

	c.curDay = -1
	c.OnChangDay()
	return &c
}

func (c *Logger) SetToTag(logToTage int) {
	c.tag = logToTage
}

func (c *Logger) SetLevel(logLevel int) {
	c.level = logLevel
}

func (c *Logger) CloseLog() {
	c.isClose = true
}

func (c *Logger) OpenLog() {
	c.isClose = false
}

func (c *Logger) OnChangDay() {
	now := time.Now()
	curDay := now.Day()
	if c.curDay == curDay {
		return
	}
	fileName := c.logPath + "/" + c.fileName + "_" + now.Format("20060102.log")
	_, err := os.Stat(fileName)
	logfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("open log file fail!", err.Error())
		return
	}
	c.logger = log.New(logfile, "", c.flag)
	c.curDay = curDay
}

func GetShortFileName(longFileName string) string {
	short := longFileName
	for i := len(longFileName) - 1; i > 0; i-- {
		if longFileName[i] == '/' {
			short = longFileName[i+1:]
			break
		}
	}
	return short
}

func GetShortFuncName(longFuncName string) string {
	short := longFuncName
	for i := len(longFuncName) - 1; i > 0; i-- {
		if longFuncName[i] == '.' {
			short = longFuncName[i+1:]
			break
		}
	}
	return short
}

func GetStackInfo(depth int) (fileName string, line int, funcName string) {
	pc, fileLongName, line, ok := runtime.Caller(depth)
	if !ok {
		return "??", 0, ""
	}
	funcHandle := runtime.FuncForPC(pc)
	return GetShortFileName(fileLongName), line, GetShortFuncName(funcHandle.Name())
}

func (c *Logger) Println(level int, v ...interface{}) {
	if c.level > level {
		return
	}
	if c.isClose {
		return
	}

	fileName, line, funcName := GetStackInfo(3)
	msg := fmt.Sprint(fileName, ":", line, ":"+funcName+" ") + fmt.Sprint(v...)

	//screen
	if c.tag&LOG_TO_SCREEN != 0 {
		log.Println(msg)
	}

	//save file
	if c.tag&LOG_TO_FILE != 0 {
		c.OnChangDay()
		c.logger.Output(2, msg)
	}
}

func (c *Logger) Debug(v ...interface{}) {
	c.Println(log_debug, v...)

}

func (c *Logger) Info(v ...interface{}) {
	c.Println(log_info, v...)
}

func (c *Logger) Warn(v ...interface{}) {
	c.Println(log_warn, v...)
}

func (c *Logger) Error(v ...interface{}) {
	c.Println(log_error, v...)
}

func (c *Logger) Fatal(v ...interface{}) {
	c.Println(log_fatal, v...)
}
