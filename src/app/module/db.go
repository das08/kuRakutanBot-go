package module

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

type Clients struct {
	Postgres *Postgres
	Redis    *Redis
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

type FindByMethod int

const (
	Name FindByMethod = iota
	ID
	Omikuji
)

func GetRakutanInfo(c Clients, method FindByMethod, value interface{}) (ExecStatus[RakutanInfos], bool) {
	var ok bool
	var status ExecStatus[RakutanInfos]

	switch method {
	case ID:
		status, ok = c.Postgres.GetRakutanInfoByID(value.(int))
	case Name:
		var subStringSearch bool
		searchWord := value.(string)
		if search := []rune(value.(string)); string(search[:1]) == "%" || string(search[:1]) == "％" {
			subStringSearch = true
			searchWord = string(search[1:])
		}
		status, ok = c.Postgres.GetRakutanInfoByLectureName(searchWord, subStringSearch)
	case Omikuji:
		status, ok = c.Postgres.GetRakutanInfoByOmikuji(value.(OmikujiType))
	}

	// Set isVerified, isFavorite and kakomonURL
	//if ok && len(result) == 1 {
	//	isVerified := IsVerified(c, env, uid)
	//	result[0].IsVerified = isVerified
	//	result[0].IsFavorite = exist(env, c.Mongo, env.DB_COLLECTION.Favorites, []KV{{Key: "uid", Value: uid}, {Key: "id", Value: result[0].ID}})
	//
	//	if isVerified && result[0].URL == "" {
	//		redisKey := fmt.Sprintf("#%d", result[0].ID)
	//		if redisStatus, cacheURL := getRedisKakomonURL(c, redisKey); redisStatus.Success {
	//			result[0].URL = cacheURL
	//		} else {
	//			kuWikiStatus := GetKakomonURL(env, result[0].LectureName)
	//			if kuWikiStatus.Success {
	//				result[0].URL = kuWikiStatus.Result
	//				setRedis(c, redisKey, kuWikiStatus.Result, time.Hour*72)
	//			} else {
	//				result[0].KUWikiErr = kuWikiStatus.Result
	//			}
	//		}
	//	}
	//}

	return status, ok
}
