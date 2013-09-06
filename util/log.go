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

func (l *Log) Log(key string, message string) {
	if logFile, ok := l.logFiles[key]; ok {
		fmt.Fprintf(logFile, "%s\n", message)
	} else {
		log.Println("Log file not found: %s", message)
	}
}

func (l *Log) Access(r *http.Request) {
	message := fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	l.Log("access", message)
}

func (l *Log) Error(r *http.Request) {
	message := fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	l.Log("error", message)
}

func (l *Log) Close(key string) {
	if logFile, ok := l.logFiles[key]; ok {
		logFile.Close()
	}
}

func (l *Log) CloseAll() {
	for _, log := range l.logFiles {
		log.Close()
	}
}

func (l *Log) AddLog(path string, key string) {
	var err error
	l.logFiles[key], err = os.Create(path)
	if err != nil {
		log.Fatal("Log file create:", err)
		return
	}
}
