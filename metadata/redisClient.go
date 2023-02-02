package metadata

//
//import (
//	"log"
//	//"github.com/go-redis/redis"
//	"github.com/gomodule/redigo/redis"
//)
//
//var redpool *redis.Pool
//
//func redset(key string, val string) error {
//	conn := redpool.Get()
//	defer conn.Close()
//	_, err := conn.Do("SET", key, val)
//	if err != nil {
//		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func reddel(key string) error {
//	conn := redpool.Get()
//	defer conn.Close()
//	_, err := conn.Do("DEL", key)
//	if err != nil {
//		log.Printf("ERROR: fail delete key %s, val %s, error %s", key, err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func redget(key string) (string, error) {
//	conn := redpool.Get()
//	defer conn.Close()
//	s, err := redis.String(conn.Do("GET", key))
//	if err != nil {
//		return "", err
//	}
//	return s, nil
//}
