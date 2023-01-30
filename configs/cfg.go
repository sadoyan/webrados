package configs

import (
	"flag"
	"github.com/ceph/go-ceph/rados"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strconv"
	"strings"
)

type CfgType struct {
	HttpAddress      string
	MonAddress       string
	Monenabled       bool
	DispatchersCount int
	//ServerAuth       bool
	AuthRead         bool
	AuthWrite        bool
	ServerUser       string
	ServerPass       string
	ClientAuth       bool
	ClientUser       string
	ClientPass       string
	InternalQueue    bool
	queue            chan string
	Uploadmaxpart    int
	Radoconns        int
	DangeZone        bool
	Readonly         bool
	Monuser          string
	Monpass          string
	Logfile          string
	LogStdout        bool
	AllPools         bool
	PoolList         []string
	OSDMaxObjectSize int
	RedisServer      string
	RedisUser        string
	RedisPass        string
	RedisDB          int
	MySQLServer      string
	MySQLDB          string
	MySQLUser        string
	MySQLPassword    string
	DBType           string
	//CacheItems         int
	//CacheTTL           time.Duration
	CacheShards             int
	CacheLifeWindow         int
	CacheCleanWindow        int
	CacheMaxEntriesInWindow int
	CacheMaxEntrySize       int
	CacheHardMaxCacheSize   int
}

var Conf = &CfgType{
	HttpAddress:      "127.0.0.1:8080",
	MonAddress:       "127.0.0.1:8989",
	DispatchersCount: 20,
	//ServerAuth:       false,
	AuthRead:         false,
	AuthWrite:        false,
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
	Radoconns:        0,
	DangeZone:        false,
	Readonly:         false,
	LogStdout:        true,
	Logfile:          "",
	AllPools:         true,
	PoolList:         []string{},
	OSDMaxObjectSize: 0,
	RedisServer:      "127.0.0.1:6379",
	RedisUser:        "",
	RedisPass:        "",
	RedisDB:          0,
	MySQLServer:      "127.0.0.1:3108",
	MySQLDB:          "",
	MySQLUser:        "",
	MySQLPassword:    "",
	DBType:           "ceph",
	//CacheItems:       0,
	//CacheTTL:         math.MaxInt,
	CacheShards:             1024,
	CacheLifeWindow:         1024,
	CacheCleanWindow:        20,
	CacheMaxEntriesInWindow: 600000,
	CacheMaxEntrySize:       5000,
	CacheHardMaxCacheSize:   1024,
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

func stringTOint(key string) int {
	a, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal("\n Value for " + key + " should be  numeric\n")
	}
	return a
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
	Conf.DispatchersCount = stringTOint(cfg.Section("main").Key("dispatchers").String())
	Conf.InternalQueue, _ = cfg.Section("main").Key("internalqueue").Bool()
	qs, _ := cfg.Section("main").Key("queuesize").Int()
	Conf.queue = make(chan string, qs)
	vsyo, _ := rados.NewConn()
	_ = vsyo.ReadDefaultConfigFile()
	_ = vsyo.Connect()
	osdMaxObjectSize, _ := vsyo.GetConfigOption("osd max object size")
	s, _ := strconv.Atoi(osdMaxObjectSize)
	Conf.Uploadmaxpart = s

	Conf.Radoconns = stringTOint(cfg.Section("main").Key("radoconns").String())

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

	Conf.CacheShards = stringTOint(cfg.Section("cache").Key("shards").String())
	Conf.CacheLifeWindow = stringTOint(cfg.Section("cache").Key("lifewindow").String())
	Conf.CacheCleanWindow = stringTOint(cfg.Section("cache").Key("cleanwindow").String())
	Conf.CacheMaxEntriesInWindow = stringTOint(cfg.Section("cache").Key("maxrntriesinwindow").String())
	Conf.CacheMaxEntrySize = stringTOint(cfg.Section("cache").Key("maxentrysize").String())
	Conf.CacheHardMaxCacheSize = stringTOint(cfg.Section("cache").Key("maxcachemb").String())

	Conf.DBType = cfg.Section("database").Key("type").String()
	switch Conf.DBType {
	case "redis":
		Conf.RedisServer = cfg.Section("database").Key("server").String()
		redisUser := cfg.Section("database").Key("username").String()
		redisPass := cfg.Section("database").Key("password").String()
		if strings.ToLower(redisUser) != "none" {
			Conf.RedisUser = redisUser
		}
		if strings.ToLower(redisPass) != "none" {
			Conf.RedisPass = redisPass
		}

		redisdb, rederr := cfg.Section("database").Key("database").Int()
		if rederr != nil {
			log.Fatal("Redis database name should be numeric", rederr)
		} else {
			Conf.RedisDB = redisdb
		}
	case "mysql":
		Conf.MySQLDB = cfg.Section("database").Key("database").String()
		Conf.MySQLServer = cfg.Section("database").Key("server").String()
		Conf.MySQLUser = cfg.Section("database").Key("username").String()
		Conf.MySQLPassword = cfg.Section("database").Key("password").String()
	case "ceph":
		Conf.DBType = "ceph"
	}

}
