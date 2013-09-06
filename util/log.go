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

func NewLog() *Log {
	log := Log{make(map[string]*os.File)}
	return &log
}

func (l *Log) Log(key string, message string) {
	go func() {
		if logFile, ok := l.logFiles[key]; ok {
			fmt.Fprintf(logFile, "%s", message)
		} else {
			log.Println("Log file not found: %s", message)
		}
	}()
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

func (l *Log) AddLog(key string, path string) {
	var err error
	l.logFiles[key], err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Log file create:", err)
		return
	}
}
