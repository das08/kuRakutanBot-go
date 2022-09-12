package module

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/goccy/go-json"
	"log"
	"strconv"
	"time"
)

type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

func CreateRedisClient() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Network:  "unix",
		Addr:     "/tmp/docker/redis.sock",
		Password: "",
		DB:       0,
	})
	rdb.Options().PoolSize = 50
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

func (r *Redis) SAddRedis(key string, values []interface{}) {
	err := r.Client.SAdd(r.Ctx, key, values...).Err()
	if err != nil {
		log.Println("[Redis] Error SADD:", err)
	} else {
		log.Printf("[Redis] SADD %s to redis", key)
	}
}

func (r *Redis) DelRedis(key string) {
	err := r.Client.Del(r.Ctx, key).Err()
	if err != nil {
		log.Println("[Redis] Error DEL:", err)
	} else {
		log.Printf("[Redis] DEL %s from redis", key)
	}
}

func (r *Redis) GetRandomOmikuji(types OmikujiType, count int64) ([]int, bool) {
	var ids []int
	key := fmt.Sprintf("set:%s", types)
	result, err := r.Client.SRandMemberN(r.Ctx, key, count).Result()
	if err != nil {
		log.Println("[Redis] Error SRANDMEMBER:", err)
		return ids, false
	}
	for _, v := range result {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Println("[Redis] Error Atoi:", err)
			return ids, false
		}
		ids = append(ids, i)
	}
	return ids, true
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
