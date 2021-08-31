package redis

import (
	"awesomeProject1/models"
	"fmt"
	"github.com/go-redis/redis"
)

type redisClient interface {
	Query(key string) *models.CacheEntry
	AddEntry(key, value string)
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func Query(key string) *models.CacheEntry {
	client := getRedisClient()

	val, err := client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return &models.CacheEntry{
		Key:   key,
		Value: val,
	}
}

func AddEntry(key, value string) {
	client := getRedisClient()
	/*json, err := json.Marshal([]byte(value))
	if err != nil {
		fmt.Println(err)
	}*/

	err := client.Set(key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveAllEntries() {
	client := getRedisClient()
	fmt.Println("***** Removing all entries in the cache ********")
	client.FlushAll()
}

func RemoveEntry(key string) {
	client := getRedisClient()
	fmt.Println("***** Removing entry for ", key, " in the cache ********")
	client.Del(key)
}
