package module

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

func CreateRedisClient() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	return &Redis{Client: rdb, Ctx: ctx}
}

//func setRedis(c Clients, key string, value interface{}, cacheTime time.Duration) {
//	resultJson, _ := json.Marshal(value)
//	err := c.Redis.Client.Set(c.Redis.Ctx, key, resultJson, cacheTime).Err()
//	if err != nil {
//		log.Println("[Redis] Error:", err)
//	} else {
//		log.Printf("[Redis] Saved %s to redis", key)
//	}
//}
//
//func getRedisRakutanInfo(c Clients, key string) (ExecStatus, RakutanInfos) {
//	data, err := c.Redis.Client.Get(c.Redis.Ctx, key).Result()
//	if err != nil {
//		return ExecStatus{Success: false}, nil
//	}
//
//	rakutanInfo := new(RakutanInfos)
//	err = json.Unmarshal([]byte(data), rakutanInfo)
//	if err != nil {
//		return ExecStatus{Success: false}, nil
//	}
//	log.Printf("[Redis] Fetched RakutanInfo from redis")
//	return ExecStatus{Success: true}, *rakutanInfo
//}
//
//func getRedisKakomonURL(c Clients, key string) (ExecStatus, string) {
//	data, err := c.Redis.Client.Get(c.Redis.Ctx, key).Result()
//	if err != nil {
//		return ExecStatus{Success: false}, ""
//	}
//
//	kakomonURL := new(string)
//	err = json.Unmarshal([]byte(data), kakomonURL)
//	if err != nil {
//		return ExecStatus{Success: false}, ""
//	}
//	log.Printf("[Redis] Fetched %s from redis", key)
//	return ExecStatus{Success: true}, *kakomonURL
//}
