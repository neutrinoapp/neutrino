package log

import (
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func args() []interface{} {
	_, caller, line, _ := runtime.Caller(2)
	callerFile := filepath.Base(caller)
	fileAndLine := callerFile + ":" + strconv.Itoa(line)

	a := make([]interface{}, 0)
	a = append(a, fileAndLine)

	return a
}

func Info(v ...interface{}) {
	log.Println(append(args(), v)...)
}

func Error(v ...interface{}) {
	log.Println(append(args(), v)...)
}
