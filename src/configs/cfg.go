package configs

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

type CfgType struct {
	HttpAddress      string
	MonAddress       string
	Monenabled       bool
	DispatchersCount int
	ServerAuth       bool
	ServerUser       string
	ServerPass       string
	ClientAuth       bool
	ClientUser       string
	ClientPass       string
	InternalQueue    bool
	queue            chan string
	Uploadmaxpart    int
	DangeZone        bool
	Monurl           string
	Monuser          string
	Monpass          string
}

var Conf = &CfgType{
	HttpAddress:      "127.0.0.1:8080",
	MonAddress:       "127.0.0.1:8989",
	DispatchersCount: 20,
	ServerAuth:       false,
	ServerUser:       "",
	ServerPass:       "",
	ClientAuth:       false,
	ClientUser:       "",
	ClientPass:       "",
	InternalQueue:    false,
	Monenabled:       false,
	Monuser:          "",
	Monpass:          "",
	Uploadmaxpart:    0,
	DangeZone:        false,
}

var authorized = make(map[string]string, 10)

func SetVarsik() {

	cfgFile := flag.String("config", "config.ini", "a string")
	flag.Parse()
	fmt.Println("Using :", *cfgFile, "as config file")

	cfg, err := ini.Load(*cfgFile)
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	Conf.HttpAddress = cfg.Section("main").Key("listen").String()
	Conf.DispatchersCount, _ = cfg.Section("main").Key("dispatchers").Int()
	Conf.InternalQueue, _ = cfg.Section("main").Key("internalqueue").Bool()
	qs, _ := cfg.Section("main").Key("queuesize").Int()
	Conf.queue = make(chan string, qs)

	Conf.Uploadmaxpart, err = cfg.Section("main").Key("uploadmaxpart").Int()
	if err != nil {
		panic("Please set numeric value to Uploadmaxpart")
	}

	switch strings.ToLower(cfg.Section("main").Key("dangerzone").String()) {
	case "yes":
		Conf.DangeZone = true
	case "no":
		Conf.DangeZone = false
	default:
		fmt.Println("\n DangerZone should be 'yes' or 'no' \n")
		os.Exit(1)
	}

	Conf.ServerAuth, _ = cfg.Section("main").Key("serverauth").Bool()
	Conf.ServerUser = cfg.Section("main").Key("serveruser").String()
	Conf.ServerPass = cfg.Section("main").Key("serverpass").String()

	Conf.ClientAuth, _ = cfg.Section("client").Key("clientauth").Bool()
	Conf.ClientUser = cfg.Section("client").Key("clientuser").String()
	Conf.ClientPass = cfg.Section("client").Key("clientpass").String()

	Conf.Monenabled, _ = cfg.Section("monitoring").Key("enabled").Bool()
	Conf.Monurl = cfg.Section("monitoring").Key("url").String()
	Conf.Monuser = cfg.Section("monitoring").Key("user").String()
	Conf.Monpass = cfg.Section("monitoring").Key("pass").String()

	authorized["main"] = Conf.ServerUser + ":" + Conf.ServerPass
	authorized["mon"] = Conf.Monuser + ":" + Conf.Monpass
}
