package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Logger interface {
	Log()
}

type Log struct {
	logFiles map[string]*os.File
}

func (this *Log) Log(key string, message string) {
	if logFile, ok := this.logFiles[key]; ok {
		fmt.Fprintf(logFile, "%s\n", message)
	} else {
		log.Println("Log file not found: %s", message)
	}
}

func (this *Log) Access(r *http.Request) {
	message := fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	this.Log("access", message)
}

func (this *Log) Error(r *http.Request) {
	message := fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	this.Log("error", message)
}

func (this *Log) Close(key string) {
	if logFile, ok := this.logFiles[key]; ok {
		logFile.Close()
	}
}

func (this *Log) CloseAll() {
	for _, log := range this.logFiles {
		log.Close()
	}
}

func (this *Log) AddLog(path string, key string) {
	var err error
	this.logFiles[key], err = os.Create(path)
	if err != nil {
		log.Fatal("Log file create:", err)
		return
	}
}
