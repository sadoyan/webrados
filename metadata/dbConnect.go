package metadata

import (
	"configs"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"log"
)

var useMySQL bool
var useRedis bool
var useRados bool
var myConns = 5
var Cache LRUCache

func DBConnect() {
	Cache = CacheConstructor(configs.Conf.Cache)
	go Cache.Evict()
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
			cache := &Cache
			file, err := cache.Get(filename)
			//stattemp := "Cached:"
			if err != nil {
				file, err = cephget(filename)
				if err == nil {
					cache.Put(filename, file)
					//stattemp = "Fresh:"
				} else {
					return "", err
				}
			}
			//fmt.Println(stattemp, cache.Size(), len(Cache.Items), file)
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
