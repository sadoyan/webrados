package wrados

import "C"
import (
	"configs"
	"github.com/ceph/go-ceph/rados"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type radcon struct {
	Connection []*rados.Conn
	Poolnames  map[string]bool
}

var Rconnect = &radcon{
	Connection: nil,
	Poolnames:  map[string]bool{},
}

func (r *radcon) connect() {
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

//func radoConnect() {
//	conn, err := rados.NewConn()
//	if err != nil {
//		Writelog("Error when invoke a new connection:", err)
//	}
//	err = conn.ReadDefaultConfigFile()
//	if err != nil {
//		Writelog("Error when read default config file:", err)
//	}
//	err = conn.Connect()
//	if err != nil {
//		Writelog("Error when connect: ", err)
//	}
//	Rconnect.Connection = append(Rconnect.Connection, conn)
//}

//var OSDMaxObjectSize int

func LsPools() {
	n := 0
	for {
		if len(Rconnect.Connection) < configs.Conf.Radoconns {
			for n < configs.Conf.Radoconns {
				Rconnect.connect()
				n = n + 1
			}
			Writelog("Created", strconv.Itoa(n), "connections to Ceph cluster")
		}

		randindex := rand.Intn(len(Rconnect.Connection))
		pools, _ := Rconnect.Connection[randindex].ListPools()
		osdMaxObjectSize, _ := Rconnect.Connection[randindex].GetConfigOption("osd max object size")
		s, _ := strconv.Atoi(osdMaxObjectSize)
		if configs.Conf.OSDMaxObjectSize != s {
			configs.Conf.Lock()
			configs.Conf.OSDMaxObjectSize = s
			configs.Conf.Uploadmaxpart = s
			configs.Conf.Unlock()
			Writelog("Setting max upload part to", s, "bytes")
		}
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
			configs.Conf.Lock()
			Rconnect.Poolnames = polos
			lst := []string{}
			for t := range Rconnect.Poolnames {
				lst = append(lst, t)
			}
			configs.Conf.Unlock()
			Writelog("Syncing RADOS pools. New pool list is:", strings.Join(lst, ", "))
		}
		time.Sleep(20 * time.Second)
	}
}
