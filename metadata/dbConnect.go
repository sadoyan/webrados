package metadata

import (
	"configs"
	"context"
	"database/sql"
	"github.com/allegro/bigcache/v3"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"time"
)

var useMySQL bool
var useRedis bool
var useRados bool
var myConns = 5

//var Cache LRUCache

var config = bigcache.Config{
	Shards:             configs.Conf.CacheShards,
	LifeWindow:         time.Duration(configs.Conf.CacheLifeWindow) * time.Minute,
	CleanWindow:        time.Duration(configs.Conf.CacheCleanWindow) * time.Minute,
	MaxEntriesInWindow: configs.Conf.CacheMaxEntriesInWindow,
	MaxEntrySize:       configs.Conf.CacheMaxEntrySize,
	Verbose:            true,
	HardMaxCacheSize:   configs.Conf.CacheHardMaxCacheSize,
	OnRemove:           nil, // callback on remove
	OnRemoveWithReason: nil, // callback on remove with reason
}
var Cache, _ = bigcache.New(context.Background(), config)

func DBConnect() {
	//go func() {
	//	for {
	//		log.Println("[Cache statistics. Len:", Cache.Len(), "Cap:", Cache.Capacity(), "Hits:", Cache.Stats().Hits, "Misses:", strconv.Itoa(int(Cache.Stats().Misses))+"]")
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	switch configs.Conf.DBType {
	case "ceph":
		log.Println("[Using Rados file xattrs for metadata]")
		useRados = true
	case "mysql":
		log.Println("[Using MySQL as metadata server]")
		useMySQL = true
		n := 0
		for n <= myConns {
			//conn, err := sql.Open(datasource, username+":"+password+"@tcp("+hostname+")/"+dbname)
			conn, err := sql.Open("mysql", configs.Conf.MySQLUser+":"+configs.Conf.MySQLPassword+"@tcp("+configs.Conf.MySQLServer+")/"+configs.Conf.MySQLDB)
			if err != nil {
				log.Println("Error when invoke a new connection:", err)
			}
			MySQLConnection = append(MySQLConnection, conn)
			n = n + 1
		}
	case "redis":
		log.Println("[Using Redis as metadata server]")
		useRedis = true
		redpool = &redis.Pool{
			MaxIdle:   100,
			MaxActive: 10000,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", configs.Conf.RedisServer)
				if err != nil {
					log.Println("ERROR: fail init redis pool:", err.Error())
				}
				return conn, err
			},
		}
	default:
		panic("[Please set database in main section of config.ini]")
	}
}

func DBClient(filename string, ops string, id string) (string, error) {
	if useRados {
		switch ops {
		case "get":
			f, err := Cache.Get(filename)
			file := string(f)
			//stattemp := "Cached:"
			if err != nil {
				file, err = cephget(filename)
				if err == nil {

					for i := 1; i <= 100000; i++ {
						_ = Cache.Set(filename+strconv.Itoa(i), []byte(file))
					}

					_ = Cache.Set(filename, []byte(file))
					//stattemp = "Fresh:"
				} else {
					return "", err
				}
			}
			//fmt.Println(stattemp, Cache.Len(), Cache.Capacity(), Cache.Stats(), filename)
			return file, err
		case "set":
			_, err := cephset(filename, id)
			if err != nil {
				return "Error updating Redis", err
			}
			return id, nil
		case "del":
			return "Done", nil
		}
		return "GGG", nil
	}
	if useRedis {
		switch ops {
		case "get":
			file, err := redget(filename)
			if err == nil {
				return file, err
				//return flir, rer
			} else {
				return "", err
			}
		case "set":
			err := redset(filename, id)
			if err != nil {
				return "Error updating Redis", err
			}
			return id, nil
		case "del":
			_ = reddel(filename)
			return "Done", nil
		}
		return "GGG", nil
	}
	if useMySQL {
		switch ops {
		case "get":
			_ = id
			file, err := myget(filename)
			if err == nil {
				return file, err
			} else {
				return "", err
			}
		case "set":
			err := myset(filename, id)
			if err != nil {
				return "Error updating MySQL", err
			}
			return id, nil
		case "del":
			_ = mydel(filename)
			return "Done", nil
		}
		return "GGG", nil
	}
	return "GGG", nil
}
