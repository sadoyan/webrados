package configs

import (
	"flag"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

type CfgType struct {
	HttpAddress      string
	MonAddress       string
	Monenabled       bool
	DispatchersCount int
	//ServerAuth       bool
	AuthRead      bool
	AuthWrite     bool
	ServerUser    string
	ServerPass    string
	ClientAuth    bool
	ClientUser    string
	ClientPass    string
	InternalQueue bool
	queue         chan string
	Uploadmaxpart int
	Radoconns     int
	DangeZone     bool
	Readonly      bool
	Monuser       string
	Monpass       string
	Logfile       string
	LogStdout     bool
	AllPools      bool
	PoolList      []string
}

var Conf = &CfgType{
	HttpAddress:      "127.0.0.1:8080",
	MonAddress:       "127.0.0.1:8989",
	DispatchersCount: 20,
	//ServerAuth:       false,
	AuthRead:      false,
	AuthWrite:     false,
	ServerUser:    "",
	ServerPass:    "",
	ClientAuth:    false,
	ClientUser:    "",
	ClientPass:    "",
	InternalQueue: false,
	Monenabled:    false,
	Monuser:       "",
	Monpass:       "",
	Uploadmaxpart: 0,
	Radoconns:     0,
	DangeZone:     false,
	Readonly:      false,
	LogStdout:     true,
	Logfile:       "",
	AllPools:      true,
	PoolList:      []string{},
}

var authorized = make(map[string]string, 10)

func stringTObool(key string, value string) bool {
	switch value {
	case "yes":
		return true
	case "no":
		return false
	default:
		log.Fatal("\n Value for " + key + " should be  'yes' or 'no' \n")
	}
	return false
}

var Cfgfile = "config.ini"

func SetVarsik() {

	if len(os.Args) >= 2 {
		Cfgfile = os.Args[1]
	}

	cfgFile := flag.String("config", Cfgfile, "a string")
	flag.Parse()

	cfg, err := ini.Load(*cfgFile)
	if err != nil {
		log.Fatal("Fail to read config file: %v", err)
	}

	Conf.HttpAddress = cfg.Section("main").Key("listen").String()
	Conf.DispatchersCount, _ = cfg.Section("main").Key("dispatchers").Int()
	Conf.InternalQueue, _ = cfg.Section("main").Key("internalqueue").Bool()
	qs, _ := cfg.Section("main").Key("queuesize").Int()
	Conf.queue = make(chan string, qs)

	Conf.Uploadmaxpart, err = cfg.Section("main").Key("uploadmaxpart").Int()
	if err != nil {
		log.Fatal("Please set numeric value to Uploadmaxpart")
	}

	Conf.Radoconns, err = cfg.Section("main").Key("radoconns").Int()
	if err != nil {
		log.Fatal("Please set numeric value to Radoconns")
	}

	Conf.LogStdout = stringTObool("dangerzone", strings.ToLower(cfg.Section("main").Key("logfile").String()))
	Conf.Logfile = cfg.Section("main").Key("logpath").String()
	Conf.ServerUser = cfg.Section("main").Key("serveruser").String()
	Conf.ServerPass = cfg.Section("main").Key("serverpass").String()

	Conf.ClientAuth, _ = cfg.Section("client").Key("clientauth").Bool()
	Conf.ClientUser = cfg.Section("client").Key("clientuser").String()
	Conf.ClientPass = cfg.Section("client").Key("clientpass").String()

	Conf.Monenabled, _ = cfg.Section("monitoring").Key("enabled").Bool()
	Conf.MonAddress = cfg.Section("monitoring").Key("url").String()
	Conf.Monuser = cfg.Section("monitoring").Key("user").String()
	Conf.Monpass = cfg.Section("monitoring").Key("pass").String()

	authorized["main"] = Conf.ServerUser + ":" + Conf.ServerPass
	authorized["mon"] = Conf.Monuser + ":" + Conf.Monpass

	Conf.AuthWrite = stringTObool("authwrite", strings.ToLower(cfg.Section("main").Key("authwrite").String()))
	Conf.AuthRead = stringTObool("authread", strings.ToLower(cfg.Section("main").Key("authread").String()))
	Conf.Readonly = stringTObool("readonly", strings.ToLower(cfg.Section("main").Key("readonly").String()))
	Conf.AllPools = stringTObool("allpools", strings.ToLower(cfg.Section("main").Key("allpools").String()))

	if !Conf.AllPools {
		x := cfg.Section("main").Key("poollist").String()
		z := strings.Replace(x, " ", "", -1)
		zs := strings.Split(z, ",")
		Conf.PoolList = zs
	}

	switch stringTObool("dangerzone", strings.ToLower(cfg.Section("main").Key("dangerzone").String())) {
	case true:
		if Conf.Readonly {
			log.Fatal("Running in read only mode cannot enable dangerous commands")
		} else {
			Conf.DangeZone = stringTObool("dangerzone", strings.ToLower(cfg.Section("main").Key("dangerzone").String()))
		}
	case false:
		Conf.DangeZone = stringTObool("dangerzone", strings.ToLower(cfg.Section("main").Key("dangerzone").String()))
	}

}
