package metadata

import (
	"configs"
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

//var useMySQL bool
//var useRedis bool
//var useRados bool
//var myConns = 5

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

//func DBConnect() {
//	switch configs.Conf.DBType {
//	case "ceph":
//		log.Println("[Using Rados file xattrs for metadata]")
//		useRados = true
//	case "mysql":
//		log.Println("[Using MySQL as metadata server]")
//		useMySQL = true
//		n := 0
//		for n <= myConns {
//			//conn, err := sql.Open(datasource, username+":"+password+"@tcp("+hostname+")/"+dbname)
//			conn, err := sql.Open("mysql", configs.Conf.MySQLUser+":"+configs.Conf.MySQLPassword+"@tcp("+configs.Conf.MySQLServer+")/"+configs.Conf.MySQLDB)
//			if err != nil {
//				log.Println("Error when invoke a new connection:", err)
//			}
//			MySQLConnection = append(MySQLConnection, conn)
//			n = n + 1
//		}
//	case "redis":
//		log.Println("[Using Redis as metadata server]")
//		useRedis = true
//		redpool = &redis.Pool{
//			MaxIdle:   100,
//			MaxActive: 10000,
//			Dial: func() (redis.Conn, error) {
//				conn, err := redis.Dial("tcp", configs.Conf.RedisServer)
//				if err != nil {
//					log.Println("ERROR: fail init redis pool:", err.Error())
//				}
//				return conn, err
//			},
//		}
//	default:
//		panic("[Please set database in main section of config.ini]")
//	}
//}

func DBClient(filename string, ops string, id string) (string, error) {
	switch ops {
	case "get":
		f, err := Cache.Get(filename)
		file := string(f)
		//stattemp := "Cached:"
		if err != nil {
			file, err = cephget(filename)
			if err == nil {
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
			return "Error updating", err
		}
		return id, nil
	case "del":
		_ = Cache.Delete(filename)

		return "Done", nil
	}
	return "GGG", nil
}
