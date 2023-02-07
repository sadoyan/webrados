package metadata

import (
	"configs"
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

var config = bigcache.Config{
	Shards:             configs.Conf.CacheShards,
	LifeWindow:         time.Duration(configs.Conf.CacheLifeWindow) * time.Minute,
	CleanWindow:        time.Duration(configs.Conf.CacheCleanWindow) * time.Minute,
	MaxEntriesInWindow: configs.Conf.CacheMaxEntriesInWindow,
	MaxEntrySize:       configs.Conf.CacheMaxEntrySize,
	Verbose:            true,
	HardMaxCacheSize:   configs.Conf.CacheHardMaxCacheSize,
	OnRemove:           nil,
	OnRemoveWithReason: nil,
}
var Cache, _ = bigcache.New(context.Background(), config)

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
