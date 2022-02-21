package metadata

import (
	"log"
	//"github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
)

var redpool *redis.Pool

//func RedinitPool() {
//	redpool = &redis.Pool{
//		MaxIdle:   100,
//		MaxActive: 10000,
//		Dial: func() (redis.Conn, error) {
//			conn, err := redis.Dial("tcp", configs.Conf.RedisServer)
//			if err != nil {
//				log.Println("ERROR: fail init redis pool:", err.Error())
//			}
//			return conn, err
//		},
//	}
//}

func redset(key string, val string) error {
	// get conn and put back when exit from method
	conn := redpool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func reddel(key string) error {
	// get conn and put back when exit from method
	conn := redpool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Printf("ERROR: fail delete key %s, val %s, error %s", key, err.Error())
		return err
	}

	return nil
}

func redget(key string) (string, error) {
	conn := redpool.Get()
	defer conn.Close()
	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return s, nil
}

//func RedClient(filename string, ops string, id string) (string, error) {
//	switch ops {
//	case "get":
//		file, err := redget(filename)
//		if err == nil {
//			return file, err
//		} else {
//			return "Redis Error", err
//		}
//	case "set":
//		err := redset(filename, id)
//		if err != nil {
//			return "Error updating Redis", err
//		}
//		return id, nil
//	case "del":
//		_ = reddel(filename)
//		return "Done", nil
//	}
//	return "GGG", nil
//}
