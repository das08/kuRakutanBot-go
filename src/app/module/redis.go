package module

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
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

func (r *Redis) SetRedis(key string, value interface{}, cacheTime time.Duration) {
	resultJson, _ := json.Marshal(value)
	err := r.Client.Set(r.Ctx, key, resultJson, cacheTime).Err()
	if err != nil {
		log.Println("[Redis] Error:", err)
	} else {
		log.Printf("[Redis] Saved %s to redis", key)
	}
}

func (r *Redis) GetRakutanInfoByID(id int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	var rakutanInfo RakutanInfo
	redisKey := fmt.Sprintf("rinfo:%d", id)
	result, err := r.Client.Get(r.Ctx, redisKey).Result()
	if err != nil {
		log.Println("[Redis] Error:", err)
		status.Err = ErrorMessageRedisGetFailed
		return status, false
	}
	log.Printf("[Redis] Getting %s from redis", result)
	err = json.Unmarshal([]byte(result), &rakutanInfo)
	if err != nil {
		log.Println("[Redis] Error:", err)
		status.Err = ErrorMessageRedisGetFailed
		return status, false
	}
	log.Printf("[Redis] Fetched %s from redis", redisKey)
	status.Result = RakutanInfos{rakutanInfo}
	return status, true
}

func (r *Redis) GetKakomonURL(key string) (ExecStatus[KUWikiKakomon], bool) {
	var status ExecStatus[KUWikiKakomon]
	var kakomonURL KUWikiKakomon
	data, err := r.Client.Get(r.Ctx, key).Result()
	if err != nil {
		log.Println("[Redis] Error:", err)
		status.Err = ErrorMessageRedisGetFailed
		return status, false
	}

	err = json.Unmarshal([]byte(data), &kakomonURL)
	if err != nil {
		status.Err = ErrorMessageRedisGetFailed
		return status, false
	}
	log.Printf("[Redis] Fetched %s from redis", key)
	status.Result = kakomonURL
	return status, true
}
