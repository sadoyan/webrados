package wrados

import (
	"configs"
	"fmt"
	"log"
	"os"
)

var Loggz = make(chan interface{}, 500)

func Writelog(line ...interface{}) {
	switch configs.Conf.LogStdout {
	case false:
		log.Println(line)
	case true:
		Loggz <- line
	}
}

func LogToFile() {
	if configs.Conf.LogStdout {
		fmt.Println("WebRados is started. Writing logs to", configs.Conf.Logfile)
	}

	f, err := os.OpenFile(configs.Conf.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for {
		line := <-Loggz
		logger := log.New(f, "", log.LstdFlags)
		logger.Println(line)
	}
}
