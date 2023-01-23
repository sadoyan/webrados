package wrados

import "C"
import (
	"configs"
	"github.com/ceph/go-ceph/rados"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Radcon struct {
	Connection []*rados.Conn
	Poolnames  map[string]bool
}

var Rconnect = &Radcon{
	Connection: nil,
	Poolnames:  map[string]bool{},
}

func RadoConnect() {
	conn, err := rados.NewConn()
	if err != nil {
		Writelog("Error when invoke a new connection:", err)
	}
	err = conn.ReadDefaultConfigFile()
	if err != nil {
		Writelog("Error when read default config file:", err)
	}
	err = conn.Connect()
	if err != nil {
		Writelog("Error when connect: ", err)
	}
	Rconnect.Connection = append(Rconnect.Connection, conn)
}

//var OSDMaxObjectSize int

func LsPools() {
	n := 0
	for {
		if len(Rconnect.Connection) < configs.Conf.Radoconns {
			for n < configs.Conf.Radoconns {
				RadoConnect()
				n = n + 1
			}
			Writelog("Created", strconv.Itoa(n), "connections to Ceph cluster")
		}
		vsyo, _ := rados.NewConn()
		_ = vsyo.ReadDefaultConfigFile()
		_ = vsyo.Connect()
		pools, _ := vsyo.ListPools()

		osdMaxObjectSize, _ := vsyo.GetConfigOption("osd max object size")
		s, _ := strconv.Atoi(osdMaxObjectSize)
		configs.Conf.OSDMaxObjectSize = s
		configs.Conf.Uploadmaxpart = s

		polos := map[string]bool{}
		switch configs.Conf.AllPools {
		case true:
			for p := range pools {
				o := pools[p]
				if o != "device_health_metrics" {
					polos[o] = true
				}
			}
		case false:
			for p := range configs.Conf.PoolList {
				o := configs.Conf.PoolList[p]
				polos[o] = true
			}
		}
		eq := reflect.DeepEqual(Rconnect.Poolnames, polos)
		switch eq {
		case false:
			Rconnect.Poolnames = polos
			lst := []string{}
			for t := range Rconnect.Poolnames {
				lst = append(lst, t)
			}
			Writelog("Syncing RADOS pools. New pool list is:", strings.Join(lst, ", "))
		}
		vsyo.Shutdown()
		time.Sleep(20 * time.Second)
	}
}
