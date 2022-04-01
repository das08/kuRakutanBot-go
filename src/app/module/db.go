package module

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"

	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

type Clients struct {
	Mongo *MongoDB
	Redis *Redis
}

type KV struct {
	Key   string
	Value interface{}
}

type QueryStatus struct {
	Success bool
	Message string
}

func CreateDBClient(e *Environments) *MongoDB {
	mongoURI := "mongodb://" + e.DB_USER + ":" + e.DB_PASS + "@" + e.DB_HOST + ":" + e.DB_PORT + "/?authSource=admin"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection created")

	return &MongoDB{Client: client, Ctx: ctx, Cancel: cancel}
}

func CreateRedisClient() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	return &Redis{Client: rdb, Ctx: ctx}
}

func setRedis(c Clients, key string, value interface{}, cacheTime time.Duration) {
	resultJson, _ := json.Marshal(value)
	err := c.Redis.Client.Set(c.Redis.Ctx, key, resultJson, cacheTime).Err()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Saved %s to redis", key)
}

func getRedisRakutanInfo(c Clients, key string) (QueryStatus, []rakutan.RakutanInfo) {
	data, err := c.Redis.Client.Get(c.Redis.Ctx, key).Result()
	if err != nil {
		return QueryStatus{Success: false}, nil
	}

	rakutanInfo := new([]rakutan.RakutanInfo)
	err = json.Unmarshal([]byte(data), rakutanInfo)
	if err != nil {
		return QueryStatus{Success: false}, nil
	}
	return QueryStatus{Success: true}, *rakutanInfo
}

func findOne(e *Environments, m *MongoDB, col Collection, filter bson.D) *mongo.SingleResult {
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	singleResult := collection.FindOne(m.Ctx, filter) //.Decode(&result)
	return singleResult
}

func insertOne(e *Environments, m *MongoDB, col Collection, filter bson.D) QueryStatus {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	_, err := collection.InsertOne(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		queryStatus = QueryStatus{false, "[i]DB接続でエラーが起きました。"}
		fmt.Println(err)
	}
	return queryStatus
}

func findOneAndUpdate(e *Environments, m *MongoDB, col Collection, filter bson.D, update bson.D) QueryStatus {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	_, err := collection.FindOneAndUpdate(m.Ctx, filter, update).DecodeBytes()
	queryStatus.Success = true

	if err != nil {
		queryStatus = QueryStatus{false, "[fu]DB接続でエラーが起きました。"}
	}
	return queryStatus
}

func deleteOne(e *Environments, m *MongoDB, col Collection, filter bson.D) QueryStatus {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	_, err := collection.DeleteOne(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		queryStatus = QueryStatus{false, "[d]DB接続でエラーが起きました。"}
		fmt.Println(err)
	}
	return queryStatus
}

func count(e *Environments, m *MongoDB, col Collection, kvs []KV) (QueryStatus, int) {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	filter := generateBsonD(kvs)
	cnt, err := collection.CountDocuments(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		queryStatus = QueryStatus{false, "[c]DB接続でエラーが起きました。"}
		fmt.Println(err)
	}
	return queryStatus, int(cnt)
}

func exist(e *Environments, m *MongoDB, col Collection, kvs []KV) bool {
	queryStatus, cnt := count(e, m, col, kvs)
	if queryStatus.Success && cnt > 0 {
		return true
	}
	return false
}

func FindByLectureID(c Clients, e *Environments, lectureID int) (QueryStatus, []rakutan.RakutanInfo) {
	var queryStatus QueryStatus
	result := rakutan.RakutanInfo{}
	singleResult := findOne(e, c.Mongo, e.DB_COLLECTION.Rakutan, generateBsonD([]KV{{Key: "id", Value: lectureID}}))

	err := singleResult.Decode(&result)
	if err != nil {
		queryStatus = QueryStatus{false, "[f]DB接続でエラーが起きました。"}
		fmt.Println(err)
	} else {
		queryStatus.Success = true
	}
	return queryStatus, []rakutan.RakutanInfo{result}
}

func FindByUID(c Clients, e *Environments, uid string) (QueryStatus, []rakutan.UserData) {
	var result []rakutan.UserData
	collection := c.Mongo.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.User)
	var queryStatus QueryStatus

	filter := bson.D{{"uid", uid}}
	filterCursor, err := collection.Find(c.Mongo.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		return QueryStatus{false, "[f]DB接続でエラーが起きました。"}, result
	}

	if err := filterCursor.All(c.Mongo.Ctx, &result); err != nil {
		queryStatus = QueryStatus{false, "[f]DB接続でエラーが起きました。"}
	}
	return queryStatus, result
}

// TODO: Add error message to query status
func FindByLectureName(c Clients, e *Environments, lectureName string) (QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := c.Mongo.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Rakutan)
	var queryStatus QueryStatus

	filterCursor, err := collection.Find(c.Mongo.Ctx, bson.D{{"lecture_name", primitive.Regex{Pattern: "^" + lectureName, Options: "i"}}})
	queryStatus.Success = true

	if err != nil {
		return QueryStatus{false, ""}, nil
	}
	if err = filterCursor.All(c.Mongo.Ctx, &result); err != nil {
		return QueryStatus{false, ""}, nil
	}
	return queryStatus, result
}

func FindByOmikuji(c Clients, e *Environments, omikujiType string) (QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := c.Mongo.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Rakutan)
	var queryStatus QueryStatus
	queryStatus.Success = true

	redisStatus, d := getRedisRakutanInfo(c, omikujiType)
	if redisStatus.Success {
		return queryStatus, []rakutan.RakutanInfo{d[randomIndex(len(d))]}
	}

	filter := bson.D{{"omikuji", omikujiType}}
	filterCursor, err := collection.Find(c.Mongo.Ctx, filter)

	if err != nil {
		return QueryStatus{false, ""}, nil
	}
	if err = filterCursor.All(c.Mongo.Ctx, &result); err != nil || result == nil {
		return QueryStatus{false, ""}, nil
	}
	setRedis(c, omikujiType, result, time.Minute*1)

	randomIdx := randomIndex(len(result))
	return queryStatus, []rakutan.RakutanInfo{result[randomIdx]}
}

func FindByFav(c Clients, e *Environments, uid string) (QueryStatus, []rakutan.Favorite) {
	var result []rakutan.Favorite
	collection := c.Mongo.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Favorites)
	var queryStatus QueryStatus

	filter := bson.D{{"uid", uid}}
	filterCursor, err := collection.Find(c.Mongo.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		return QueryStatus{false, "お気に入りの取得に失敗しました。"}, result
	}

	err = filterCursor.All(c.Mongo.Ctx, &result)
	switch {
	case err != nil:
		queryStatus = QueryStatus{false, "お気に入りの取得に失敗しました。"}
	case len(result) == 0:
		queryStatus = QueryStatus{false, "お気に入り登録している講義はありません。講義名の左上にある★マークを押すとお気に入り登録できます！"}
	}

	return queryStatus, result
}

type FindByMethod int

const (
	Name FindByMethod = iota
	ID
	Omikuji
)

func GetRakutanInfo(c Clients, env *Environments, uid string, method FindByMethod, value interface{}) (QueryStatus, []rakutan.RakutanInfo) {
	var queryStatus QueryStatus
	var result []rakutan.RakutanInfo

	switch method {
	case ID:
		queryStatus, result = FindByLectureID(c, env, value.(int))
	case Name:
		queryStatus, result = FindByLectureName(c, env, value.(string))
	case Omikuji:
		queryStatus, result = FindByOmikuji(c, env, value.(string))
	}

	// Set isVerified, isFavorite and kakomonURL
	if queryStatus.Success && len(result) == 1 {
		isVerified := exist(env, c.Mongo, env.DB_COLLECTION.User, []KV{{Key: "uid", Value: uid}, {Key: "verified", Value: true}})
		result[0].IsVerified = isVerified
		result[0].IsFavorite = exist(env, c.Mongo, env.DB_COLLECTION.Favorites, []KV{{Key: "uid", Value: uid}})

		if isVerified {
			if kakomonURL := GetKakomonURL(env, result[0].LectureName); kakomonURL != nil {
				result[0].URL = *kakomonURL
			}
		}
	}

	return queryStatus, result
}

func GetFavorites(c Clients, env *Environments, uid string) (QueryStatus, []rakutan.Favorite) {
	var queryStatus QueryStatus
	var result []rakutan.Favorite
	queryStatus, result = FindByFav(c, env, uid)

	return queryStatus, result
}

func InsertFavorite(m *MongoDB, env *Environments, pbe PostbackEntry) QueryStatus {
	kvs := []KV{{"uid", pbe.Uid}, {"id", pbe.Param.ID}}
	singleResult := findOne(env, m, env.DB_COLLECTION.Favorites, generateBsonD(kvs))

	// すでにお気に入り登録されているか調べる
	var findStatus QueryStatus
	err := singleResult.Decode(&rakutan.Favorite{})
	switch {
	case err != nil && err == mongo.ErrNoDocuments:
		// documentがなければお気に入り登録できる
		findStatus.Success = true
	case err != nil:
		findStatus = QueryStatus{false, "[f]DB接続でエラーが起きました。"}
	default:
		deleteStatus := deleteOne(env, m, env.DB_COLLECTION.Favorites, generateBsonD(kvs))
		findStatus = QueryStatus{false, deleteStatus.Message}
		if deleteStatus.Success {
			findStatus.Message = fmt.Sprintf("「%s」をお気に入りから削除しました！", pbe.Param.LectureName)
		}
	}

	switch {
	// お気に入り登録されていた場合
	case !findStatus.Success:
		return QueryStatus{false, findStatus.Message}
	// まだお気に入り登録されていなかった場合
	default:
		countStatus, favCount := count(env, m, env.DB_COLLECTION.Favorites, []KV{{"uid", pbe.Uid}})
		fmt.Println("favCnt", favCount)
		switch {
		case countStatus.Success && favCount < 50:
			kvs = append(kvs, KV{"lecture_name", pbe.Param.LectureName})
			insertStatus := insertOne(env, m, env.DB_COLLECTION.Favorites, generateBsonD(kvs))
			if insertStatus.Success {
				insertStatus.Message = fmt.Sprintf("「%s」をお気に入り登録しました！", pbe.Param.LectureName)
			}
			return insertStatus
		case countStatus.Success && favCount >= 50:
			return QueryStatus{Success: false, Message: "お気に入り数が上限(50件)に達しています。"}
		default:
			return countStatus
		}
	}
}

func DeleteFavorite(m *MongoDB, env *Environments, pbe PostbackEntry) QueryStatus {
	kvs := []KV{{"uid", pbe.Uid}, {"id", pbe.Param.ID}}
	deleteStatus := deleteOne(env, m, env.DB_COLLECTION.Favorites, generateBsonD(kvs))
	if deleteStatus.Success {
		return QueryStatus{Success: true, Message: fmt.Sprintf("「%s」をお気に入りから削除しました！", pbe.Param.LectureName)}
	} else {
		return QueryStatus{Success: false, Message: "お気に入りの削除に失敗しました。"}
	}
}

func registerUser(env *Environments, m *MongoDB, uid string) {
	bsonD := bson.D{
		{"uid", uid},
		{"count", bson.D{{"message", 1}, {"rakutan", 0}, {"onitan", 0}}},
		{"register_time", int(time.Now().Unix())},
		{"verified", false},
	}
	insertStatus := insertOne(env, m, env.DB_COLLECTION.User, bsonD)
	if insertStatus.Success {
		fmt.Println("Register Success.")
	} else {
		fmt.Println("Register Failed.")
	}
}

func countUp(env *Environments, m *MongoDB, uid string, key string) {
	filter := generateBsonD([]KV{{"uid", uid}})
	countUpStatus := findOneAndUpdate(env, m, env.DB_COLLECTION.User, filter, bson.D{{"$inc", bson.D{{fmt.Sprintf("count.%s", key), 1}}}})
	if countUpStatus.Success {
		fmt.Println("Countup Success.")
	} else {
		fmt.Println("Countup Failed.")
	}
}

func CountMessage(c Clients, env *Environments, uid string) {
	queryStatus, result := FindByUID(c, env, uid)
	if queryStatus.Success {
		if len(result) == 0 {
			registerUser(env, c.Mongo, uid)
		} else {
			countUp(env, c.Mongo, uid, "message")
		}
	}
}

func generateBsonD(kvs []KV) bson.D {
	entry := bson.D{}
	for _, kv := range kvs {
		entry = append(entry, bson.E{Key: kv.Key, Value: kv.Value})
	}
	return entry
}

func randomIndex(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
