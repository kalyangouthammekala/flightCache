package redis

import (
	"awesomeProject1/models"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/magiconair/properties"
)

type redisClient interface {
	Query(key string) *models.CacheEntry
	AddEntry(key, value string)
}

//Addr:	  "flightcachetest.i4pew7.0001.use2.cache.amazonaws.com:6379",
//Addr:     "localhost:6379",
func getRedisClient(p *properties.Properties) *redis.Client {
	redisAddr, _ := p.Get("redis-addr-port-AWS")
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}

func Query(key string, p *properties.Properties) *models.CacheEntry {
	client := getRedisClient(p)

	val, err := client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return &models.CacheEntry{
		Key:   key,
		Value: val,
	}
}

func AddEntry(key, value string, p *properties.Properties) {
	client := getRedisClient(p)

	keys := client.Keys("*")
	fmt.Println(keys)
	/*json, err := json.Marshal([]byte(value))
	if err != nil {
		fmt.Println(err)
	}*/

	err := client.Set(key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

//func RemoveAllEntries(p *properties.Properties) {
//	client := getRedisClient(p)
//	fmt.Println("***** Removing all entries in the cache ********")
//	client.FlushAll()
//}

func RemoveEntry(key string, p *properties.Properties) {
	client := getRedisClient(p)
	fmt.Println("***** Removing entry for ", key, " in the cache ********")
	client.Del(key)
}
