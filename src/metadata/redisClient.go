package metadata

import (
	"configs"
	"context"
	"github.com/go-redis/redis"
)

func RedClient(filename string, ops string, id string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.Conf.RedisServer,
		Username: configs.Conf.RedisUser,
		Password: configs.Conf.RedisPass,
		DB:       configs.Conf.RedisDB,
	})

	var ctx = context.Background()

	switch ops {
	case "get":
		file, err := rdb.Get(ctx, filename).Result()
		if err == nil {
			return file, err
		} else {
			return "Redis Error", err
		}
	case "set":
		err := rdb.Set(ctx, filename, id, 0).Err()
		if err != nil {
			return "Error updating Redis", err
		}
		return id, nil
	case "del":
		_ = rdb.Del(ctx, filename)
		return "Done", nil
	}

	if ops == "get" {
		file, err := rdb.Get(ctx, filename).Result()
		if err == nil {
			return file, err
		} else {
			return "Redis Error", err
		}
	} else if ops == "set" {
		err := rdb.Set(ctx, filename, id, 0).Err()
		if err != nil {
			return "Error updating Redis", err
		}
		return id, nil
	}

	return "GGG", nil
}
