package configs

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type CfgType struct {
	HttpAddress             string
	MonAddress              string
	Monenabled              bool
	DispatchersCount        int
	AuthApi                 bool
	AuthBasic               bool
	AuthJWT                 bool
	JWTSecret               []byte
	UsersFile               string
	AuthRead                bool
	AuthWrite               bool
	ServerUser              string
	ServerPass              string
	ClientAuth              bool
	ClientUser              string
	ClientPass              string
	InternalQueue           bool
	queue                   chan string
	Uploadmaxpart           int
	Radoconns               int
	DangeZone               bool
	Readonly                bool
	Monuser                 string
	Monpass                 string
	Logfile                 string
	LogStdout               bool
	AllPools                bool
	PoolList                []string
	OSDMaxObjectSize        int
	CacheShards             int
	CacheLifeWindow         int
	CacheCleanWindow        int
	CacheMaxEntriesInWindow int
	CacheMaxEntrySize       int
	CacheHardMaxCacheSize   int
	sync.RWMutex
}

var Conf = &CfgType{
	HttpAddress:             "127.0.0.1:8080",
	MonAddress:              "127.0.0.1:8989",
	DispatchersCount:        20,
	AuthApi:                 false,
	AuthBasic:               false,
	AuthJWT:                 false,
	JWTSecret:               []byte(os.Getenv("JWTSECRET")),
	AuthRead:                false,
	AuthWrite:               false,
	UsersFile:               "",
	ServerUser:              "",
	ServerPass:              "",
	ClientAuth:              false,
	ClientUser:              "",
	ClientPass:              "",
	InternalQueue:           false,
	Monenabled:              false,
	Monuser:                 "",
	Monpass:                 "",
	Uploadmaxpart:           0,
	Radoconns:               0,
	DangeZone:               false,
	Readonly:                false,
	LogStdout:               true,
	Logfile:                 "",
	AllPools:                true,
	PoolList:                []string{},
	OSDMaxObjectSize:        0,
	CacheShards:             1024,
	CacheLifeWindow:         1024,
	CacheCleanWindow:        20,
	CacheMaxEntriesInWindow: 600000,
	CacheMaxEntrySize:       5000,
	CacheHardMaxCacheSize:   1024,
	RWMutex:                 sync.RWMutex{},
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

var Cfgfile string

func SetVarsik() {

	cfgFile := flag.String("config", "config.yml", "a string")
	flag.Parse()

	data := make(map[interface{}]map[interface{}]interface{})

	yfile, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal("Cant read config file:", err)
	}

	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal("Cant parse config file:", err2)
	}

	Conf.HttpAddress = data["main"]["listen"].(string)
	Conf.DispatchersCount = data[("main")]["dispatchers"].(int)
	Conf.InternalQueue, _ = data["main"]["internalqueue"].(bool)
	qs, _ := data["main"]["queuesize"].(int)
	Conf.queue = make(chan string, qs)
	Conf.Radoconns = data["main"]["radoconns"].(int)
	Conf.LogStdout = stringTObool("logfile", strings.ToLower(data["main"]["logfile"].(string)))
	Conf.Logfile = data["main"]["logpath"].(string)
	Conf.UsersFile = data["main"]["usersfile"].(string)
	Conf.ServerUser = data["main"]["serveruser"].(string)
	Conf.ServerPass = data["main"]["serverpass"].(string)

	Conf.Monenabled = stringTObool("monenabled", strings.ToLower(data["monitoring"]["enabled"].(string)))
	Conf.MonAddress = data["monitoring"]["url"].(string)
	Conf.Monuser = data["monitoring"]["user"].(string)
	Conf.Monpass = data["monitoring"]["pass"].(string)

	authorized["main"] = Conf.ServerUser + ":" + Conf.ServerPass
	authorized["mon"] = Conf.Monuser + ":" + Conf.Monpass

	Conf.AuthWrite = stringTObool("authwrite", strings.ToLower(data["main"]["authwrite"].(string)))
	Conf.AuthRead = stringTObool("authread", strings.ToLower(data["main"]["authread"].(string)))
	Conf.Readonly = stringTObool("readonly", strings.ToLower(data["main"]["readonly"].(string)))
	Conf.AllPools = stringTObool("allpools", strings.ToLower(data["main"]["allpools"].(string)))

	authtype := data["main"]["authtype"].(string)

	switch strings.ToLower(authtype) {
	case "basic":
		log.Println("[Using HTTP basic authentication]")
		Conf.AuthBasic = true
	case "jwt":
		log.Println("[Using JWT authentication]")
		Conf.AuthJWT = true
	case "apikey":
		log.Println("[Using ApiKey authentication]")
		Conf.AuthApi = true
	}

	if !Conf.AllPools {
		for _, pool := range data["main"]["poollist"].([]interface{}) {
			Conf.PoolList = append(Conf.PoolList, pool.(string))
		}
	}

	switch stringTObool("dangerzone", strings.ToLower(data["main"]["dangerzone"].(string))) {
	case true:
		if Conf.Readonly {
			log.Fatal("Running in read only mode cannot enable dangerous commands")
		} else {
			Conf.DangeZone = stringTObool("dangerzone", strings.ToLower(data["main"]["dangerzone"].(string)))
		}
	case false:
		Conf.DangeZone = stringTObool("dangerzone", strings.ToLower(data["main"]["dangerzone"].(string)))
	}

	Conf.CacheShards = data["cache"]["shards"].(int)
	Conf.CacheLifeWindow = data["cache"]["lifewindow"].(int)
	Conf.CacheCleanWindow = data["cache"]["cleanwindow"].(int)
	Conf.CacheMaxEntriesInWindow = data["cache"]["maxrntriesinwindow"].(int)
	Conf.CacheMaxEntrySize = data["cache"]["maxentrysize"].(int)
	Conf.CacheHardMaxCacheSize = data["cache"]["maxcachemb"].(int)

}
