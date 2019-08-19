package logHelper

import (
	"MQClient/config"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type logInfo struct {
	logObj *log.Logger
	file *os.File
}

func(this *logInfo) prepareLog() error{
	logPath := "./log"
	_, err := os.Stat(logPath)
	if err != nil {
		os.Mkdir(logPath, os.ModePerm)
	}
	fileName := strings.ReplaceAll(time.Now().Format(config.TIME_FORMAT_SHORT), "-", "") + ".log"
	file, err := os.OpenFile(logPath + "/log-" + fileName, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		return err
	}
	this.file = file
	multiWrite := io.MultiWriter(os.Stdout, this.file)
	this.logObj = log.New(multiWrite, "Info:", log.Ldate | log.Ltime | log.Lshortfile)
	return nil
}
func(this *logInfo) closeFile(){
	this.file.Close()
}
func LogDefault(prefix, msg string){
	l := &logInfo{}
	err := l.prepareLog()
	if err != nil {
		return
	}
	l.logObj.Printf("[%s]%s\n", prefix, msg)
	l.closeFile()
}

func FailOnError(prefix string, err error, msg string){
	if err != nil {
		l := &logInfo{}
		err := l.prepareLog()
		if err != nil {
			return
		}
		l.logObj.Printf("[%s] %s : %s \n", prefix, msg, err)
		l.closeFile()
	}
}

func PanicOnError(err error){
	if err != nil {
		panic(err.Error())
	}
}
