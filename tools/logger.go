package tools

import (
	"log"
	"os"
)

var Loggz = make(chan interface{}, 500)

type Logger interface {
	File([]interface{})
	StdOut([]interface{})
}

type LogOut struct {
	LogToFile bool
	FilePath  string
}

func (r *LogOut) File(line []interface{}) {
	Loggz <- line
}

func (r *LogOut) StdOut(line []interface{}) {
	log.Println(line)
}

var Logging = &LogOut{
	LogToFile: true,
	FilePath:  "",
}

func WriteLogs(line ...interface{}) {
	switch Logging.LogToFile {
	case false:
		Logging.StdOut(line)
	case true:
		Logging.File(line)
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
		line := <-Loggz
		logger := log.New(f, "", log.LstdFlags)
		logger.Println(line)
	}
}

// -------------------------------------------------------------------------------------------------------- //

//package wrados
//
//import (
//	"configs"
//	"fmt"
//	"log"
//	"os"
//)
//
//var Loggz = make(chan interface{}, 500)
//
//func Writelog(line ...interface{}) {
//	switch configs.Conf.LogStdout {
//	case false:
//		log.Println(line)
//	case true:
//		Loggz <- line
//	}
//}
//
//func LogToFile() {
//	if configs.Conf.LogStdout {
//		fmt.Println("WebRados is started.", "Config file", configs.Cfgfile+",", "Logfile", configs.Conf.Logfile)
//	}
//
//	f, err := os.OpenFile(configs.Conf.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Println(err)
//	}
//	defer f.Close()
//
//	for {
//		line := <-Loggz
//		logger := log.New(f, "", log.LstdFlags)
//		logger.Println(line)
//	}
//}
