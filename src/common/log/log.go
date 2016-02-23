package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func args(args []interface{}) []interface{} {
	stackLength := 2

	if len(args) > 1 {
		lastEl := args[len(args)-1]
		switch lastEl.(type) {
		case int:
			stackLength += lastEl.(int)
			args = args[:len(args)-1]
		}
	}

	_, caller, line, _ := runtime.Caller(stackLength)
	callerFile := filepath.Base(caller)
	fileAndLine := callerFile + ":" + strconv.Itoa(line)

	a := make([]interface{}, 0)

	stringifiedArgs := make([]string, 0)

	for i := range args {
		arg := args[i]
		stringifiedArgs = append(stringifiedArgs, fmt.Sprintf("%+v", arg))
	}

	a = append(a, stringifiedArgs, fileAndLine)

	return a
}

func Info(v ...interface{}) {
	log.Println(args(v)...)
}

func Error(v ...interface{}) {
	log.Println("Stack:")
	log.Println(args(v)...)
}
