package metadata

import (
	"configs"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"log"
)

var useMySQL bool
var useRedis bool
var myConns = 5

func DBConnect() {
	switch configs.Conf.DBType {
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
	if useRedis {
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
	if useMySQL {
		switch ops {
		case "get":
			_ = id
			file, err := myget(filename)
			if err == nil {
				return file, err
			} else {
				return "MySQL Error", err
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
