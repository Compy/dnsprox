// This file contains all logic for managing log output and warning levels
package main

import (
	"fmt"
	"log"
	"os"
)

const (
	errorLevel = 0
	warnLevel  = 1
	infoLevel  = 2
	debugLevel = 3
)

var (
	warnEnabled  bool
	infoEnabled  bool
	debugEnabled bool
	errorEnabled bool
)

//Output of logging functions
var Output *os.File

//Init the log output and logging level
func InitLogging(path string, level uint) error {
	switch path {
	case "stdout":
		Output = os.Stdout
	case "stderr":
		Output = os.Stderr
	default:
		var err error
		Output, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("error opening log file: %v", err)
		}
	}
	log.SetOutput(Output)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	warnEnabled = level >= warnLevel
	infoEnabled = level >= infoLevel
	debugEnabled = level >= debugLevel
	errorEnabled = level >= errorLevel
	return nil
}

func println(title string, msg ...interface{}) {
	arg := append([]interface{}{title}, msg...)
	log.Println(arg...)
}

func printFormat(format string, msg ...interface{}) {
	log.Printf(format, msg...)
}

//Debug level > 1 logging with "[DEBUG]" prefix
func Debug(msg ...interface{}) {
	if debugEnabled {
		println("[DEBUG]", msg...)
	}
}

//Debugf level > 1 format logging with "[DEBUG]" prefix
func Debugf(format string, msg ...interface{}) {
	if debugEnabled {
		printFormat("[DEBUG] "+format, msg...)
	}
}

//Info level > 0 logging with "[INFO]" prefix
func Info(msg ...interface{}) {
	if infoEnabled {
		println("[INFO] ", msg...)
	}
}

//Infof level > 0 format logging with "[INFO]" prefix
func Infof(format string, msg ...interface{}) {
	if infoEnabled {
		printFormat("[INFO]  "+format, msg...)
	}
}

//Warning any level logging with "[WARN]" prefix
func Warning(msg ...interface{}) {
	if warnEnabled {
		println("[WARN] ", msg...)
	}
}

//Warningf any level format logging with "[WARN]" prefix
func Warningf(format string, msg ...interface{}) {
	if warnEnabled {
		printFormat("[WARN]  "+format, msg...)
	}
}

//Error any level logging with "[ERROR]" prefix
func Error(msg ...interface{}) {
	if errorEnabled {
		println("[ERROR] ", msg...)
	}
}

//Errorf any level format logging with "[ERROR]" prefix
func Errorf(format string, msg ...interface{}) {
	if errorEnabled {
		printFormat("[ERROR] "+format, msg...)
	}
}
