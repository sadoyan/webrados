package tools

import (
	"log"
	"os"
)

var logchan = make(chan interface{}, 500)

type logg struct {
	LogToFile bool
	FilePath  string
}

var Logging = &logg{
	LogToFile: true,
	FilePath:  "",
}

func WriteLogs(line ...interface{}) {
	switch Logging.LogToFile {
	case false:
		log.Println(line)
	case true:
		logchan <- line
	}
}

func LogToFile() {
	if !Logging.LogToFile {
		log.Println("[WebRados is started.", "Config file Logging to stdout]")
		return
	}

	log.Println("[WebRados is started.", "Config file, Logging to", Logging.FilePath+"]")
	f, err := os.OpenFile(Logging.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for {
		line := <-logchan
		logger := log.New(f, "", log.LstdFlags)
		logger.Println(line)
	}
}
