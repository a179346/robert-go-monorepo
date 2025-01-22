package console

import (
	"log"
	"os"
)

const reset = "\033[0m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"

const flag = log.Ldate | log.Ltime | log.Lshortfile

var infoLog = log.New(os.Stdout, green+"INFO "+reset, flag)
var warningLog = log.New(os.Stdout, yellow+"WARNING "+reset, flag)
var errorLog = log.New(os.Stderr, red+"ERROR "+reset, flag)

func Info(v ...any) {
	infoLog.Println(v...)
}

func Infof(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Warn(v ...any) {
	warningLog.Println(v...)
}

func Warnf(format string, v ...any) {
	warningLog.Printf(format, v...)
}

func Error(v ...any) {
	errorLog.Println(v...)
}

func Errorf(format string, v ...any) {
	errorLog.Printf(format, v...)
}
