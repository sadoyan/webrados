package wrados

import "C"
import (
	"configs"
	"github.com/ceph/go-ceph/rados"
	"log"
	"reflect"
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
		log.Println("Error when invoke a new connection:", err)
	}
	err = conn.ReadDefaultConfigFile()
	if err != nil {
		log.Println("Error when read default config file:", err)
	}
	err = conn.Connect()
	if err != nil {
		log.Println("Error when connect:", err)
	}
	log.Println("Adding Rados connection to pool, Connected to Ceph cluster")
	Rconnect.Connection = append(Rconnect.Connection, conn)
}

func LsPools() {
	n := 0
	for {
		if len(Rconnect.Connection) < configs.Conf.Radoconns {
			for n <= configs.Conf.Radoconns {
				RadoConnect()
				n = n + 1
			}
		}
		//randindex := rand.Intn(configs.Conf.Radoconns)
		//pools, _ := Rconnect.Connection[randindex].ListPools()
		vsyo, _ := rados.NewConn()
		_ = vsyo.ReadDefaultConfigFile()
		_ = vsyo.Connect()
		pools, _ := vsyo.ListPools()
		polos := map[string]bool{}
		for p := range pools {
			o := pools[p]
			//switch Rconnect.Poolnames[o] {
			//case false:
			//	if o != "device_health_metrics" {
			//		log.Println("Enabling new pool:", o)
			//		Rconnect.Poolnames[o] = true
			//	}
			//}
			if o != "device_health_metrics" {
				polos[o] = true
			}
		}
		eq := reflect.DeepEqual(Rconnect.Poolnames, polos)
		//log.Println(eq)
		switch eq {
		case false:
			Rconnect.Poolnames = polos
			log.Println("Changes in RADOS pools are detected syncing. New pool list is:", Rconnect.Poolnames)
		}

		//log.Println(Rconnect.Poolnames)
		//log.Println(polos)

		vsyo.Shutdown()
		time.Sleep(20 * time.Second)
	}
}

//func PutData(pool string, name string, input []byte) {
//	ioctx, _ := Rconnect.Connection.OpenIOContext(pool)
//	_ = ioctx.Write(name, input, 0)
//}
//
//func GetData(pool string, name string) []byte {
//	if _, ok := Rconnect.Poolnames[pool]; ok {
//		ioctx, e := Rconnect.Connection.OpenIOContext(pool)
//		if e != nil {
//			log.Println(e)
//		}
//		xo, _ := ioctx.Stat(name)
//		bytesOut := make([]byte, xo.Size)
//		out, _ := ioctx.Read(name, bytesOut, 0)
//		log.Println(out, pool, name, xo.Size)
//		return bytesOut
//	} else {
//		log.Println("Pool " + pool + " does not exists")
//		return nil
//	}
//}
