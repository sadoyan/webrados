package metadata

import (
	"configs"
	"log"
	"os"
	//"github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
)

var redpool *redis.Pool

func RedinitPool() {
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
}

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
		//log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}
	return s, nil
}
func redsadd(key string, val string) error {
	// get conn and put back when exit from method
	conn := redpool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", key, val)
	if err != nil {
		log.Printf("ERROR: fail add val %s to set %s, error %s", val, key, err.Error())
		return err
	}
	return nil
}
func redsmembers(key string) ([]string, error) {
	// get conn and put back when exit from method
	conn := redpool.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Printf("ERROR: fail get set %s , error %s", key, err.Error())
		return nil, err
	}
	return s, nil
}
func redping(conn redis.Conn) {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		log.Printf("ERROR: fail ping redis conn: %s", err.Error())
		os.Exit(1)
	}
}

func RedClient(filename string, ops string, id string) (string, error) {
	switch ops {
	case "get":
		file, err := redget(filename)
		if err == nil {
			return file, err
		} else {
			return "Redis Error", err
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

// --------------------------------------------------------------------- //
//func RedClient(filename string, ops string, id string) (string, error) {
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     configs.Conf.RedisServer,
//		Username: configs.Conf.RedisUser,
//		Password: configs.Conf.RedisPass,
//		DB:       configs.Conf.RedisDB,
//	})
//
//	var ctx = context.Background()
//
//	switch ops {
//	case "get":
//		file, err := rdb.Get(ctx, filename).Result()
//		if err == nil {
//			return file, err
//		} else {
//			return "Redis Error", err
//		}
//	case "set":
//		err := rdb.Set(ctx, filename, id, 0).Err()
//		if err != nil {
//			return "Error updating Redis", err
//		}
//		return id, nil
//	case "del":
//		_ = rdb.Del(ctx, filename)
//		return "Done", nil
//	}
//
//	if ops == "get" {
//		file, err := rdb.Get(ctx, filename).Result()
//		if err == nil {
//			return file, err
//		} else {
//			return "Redis Error", err
//		}
//	} else if ops == "set" {
//		err := rdb.Set(ctx, filename, id, 0).Err()
//		if err != nil {
//			return "Error updating Redis", err
//		}
//		return id, nil
//	}
//
//	return "GGG", nil
//}
