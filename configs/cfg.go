package configs

import (
	"flag"
	"github.com/ceph/go-ceph/rados"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type CfgType struct {
	HttpAddress             string
	MonAddress              string
	Monenabled              bool
	DispatchersCount        int
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
}

var Conf = &CfgType{
	HttpAddress:             "127.0.0.1:8080",
	MonAddress:              "127.0.0.1:8989",
	DispatchersCount:        20,
	AuthRead:                false,
	AuthWrite:               false,
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

var Cfgfile = "config.yml"

func SetVarsik() {

	if len(os.Args) >= 2 {
		Cfgfile = os.Args[1]
	}

	cfgFile := flag.String("config", Cfgfile, "a string")
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
	vsyo, _ := rados.NewConn()
	_ = vsyo.ReadDefaultConfigFile()
	_ = vsyo.Connect()
	osdMaxObjectSize, _ := vsyo.GetConfigOption("osd max object size")
	s, _ := strconv.Atoi(osdMaxObjectSize)
	Conf.Uploadmaxpart = s

	Conf.Radoconns = data["main"]["radoconns"].(int)

	Conf.LogStdout = stringTObool("logfile", strings.ToLower(data["main"]["logfile"].(string)))
	Conf.Logfile = data["main"]["logpath"].(string)
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
