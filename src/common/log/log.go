package log

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func fileAndLineFromStack(frame int) string {
	_, caller, line, _ := runtime.Caller(frame)
	callerFile := filepath.Base(caller)
	fileAndLine := callerFile + ":" + strconv.Itoa(line)
	return fileAndLine
}

func args(args []interface{}) []interface{} {
	stackLength := 3

	if len(args) > 1 {
		lastEl := args[len(args)-1]
		switch lastEl.(type) {
		case int:
			stackLength += lastEl.(int)
			args = args[:len(args)-1]
		}
	}

	fileAndLine := fileAndLineFromStack(stackLength)

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
	s := args(v)
	log.Println(append([]interface{}{"INFO:"}, s...))
}

func Warn(v ...interface{}) {
	s := args(v)
	log.Println(append([]interface{}{"WARN:"}, s))
}

func Error(v ...interface{}) {
	var stackBuf bytes.Buffer

	for i := 1; i <= 10; i++ {
		fileAndLine := fileAndLineFromStack(i)
		if fileAndLine == ".:0" {
			break
		}

		stackBuf.WriteString("\r\n\t" + fileAndLine)
	}

	errorArgs := args(v)
	finalErrorArgs := make([]interface{}, 1)
	finalErrorArgs[0] = "Error: "

	s := fmt.Sprintf("%v, %v", append(finalErrorArgs, errorArgs...), "[Stack: "+stackBuf.String()+"]")
	log.Println(s)
}
