package module

import (
	"context"
	"fmt"
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

type KV struct {
	Key   string
	Value interface{}
}

type QueryStatus struct {
	Success bool
	Message string
}

func CreateDBClient(e *Environments) *MongoDB {
	mongoURI := "mongodb://" + e.DB_USER + ":" + e.DB_PASS + "@" + e.DB_HOST + ":" + e.DB_PORT + "/?authSource=" + e.DB_NAME
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

func findOne(e *Environments, m *MongoDB, col Collection, kvs []KV) *mongo.SingleResult {
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	entry := generateBsonD(kvs)
	singleResult := collection.FindOne(m.Ctx, entry) //.Decode(&result)
	return singleResult
}

func insertOne(e *Environments, m *MongoDB, col Collection, kvs []KV) QueryStatus {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	filter := generateBsonD(kvs)
	_, err := collection.InsertOne(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		queryStatus = QueryStatus{false, "[i]DB接続でエラーが起きました。"}
		fmt.Println(err)
	}
	return queryStatus
}

func deleteOne(e *Environments, m *MongoDB, col Collection, kvs []KV) QueryStatus {
	var queryStatus QueryStatus
	collection := m.Client.Database(e.DB_NAME).Collection(col)
	filter := generateBsonD(kvs)
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

func FindByLectureID(e *Environments, m *MongoDB, lectureID int) (QueryStatus, []rakutan.RakutanInfo) {
	var queryStatus QueryStatus
	result := rakutan.RakutanInfo{}
	singleResult := findOne(e, m, e.DB_COLLECTION.Rakutan, []KV{{Key: "id", Value: lectureID}})

	err := singleResult.Decode(&result)
	if err != nil {
		queryStatus = QueryStatus{false, "[f]DB接続でエラーが起きました。"}
		fmt.Println(err)
	} else {
		queryStatus.Success = true
	}
	return queryStatus, []rakutan.RakutanInfo{result}
}

// TODO: Add error message to query status
func FindByLectureName(e *Environments, m *MongoDB, lectureName string) (QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Rakutan)
	var queryStatus QueryStatus

	filterCursor, err := collection.Find(m.Ctx, bson.D{{"lecture_name", primitive.Regex{Pattern: "^" + lectureName, Options: "i"}}})
	queryStatus.Success = true

	if err != nil {
		fmt.Println(err)
		queryStatus.Success = false
	}
	if err = filterCursor.All(m.Ctx, &result); err != nil {
		queryStatus.Success = false
	}
	return queryStatus, result
}

func FindByOmikuji(e *Environments, m *MongoDB, omikujiType string) (QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Rakutan)
	var queryStatus QueryStatus

	filter := bson.D{{"omikuji", omikujiType}}
	filterCursor, err := collection.Find(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		queryStatus.Success = false
	}
	if err = filterCursor.All(m.Ctx, &result); err != nil {
		queryStatus.Success = false
	}

	randomIdx := randomIndex(len(result))
	return queryStatus, []rakutan.RakutanInfo{result[randomIdx]}
}

func FindByFav(e *Environments, m *MongoDB, uid string) (QueryStatus, []rakutan.Favorite) {
	var result []rakutan.Favorite
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION.Favorites)
	var queryStatus QueryStatus

	filter := bson.D{{"uid", uid}}
	filterCursor, err := collection.Find(m.Ctx, filter)
	queryStatus.Success = true

	if err != nil {
		return QueryStatus{false, "お気に入りの取得に失敗しました。"}, result
	}

	err = filterCursor.All(m.Ctx, &result)
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

func GetRakutanInfo(env *Environments, method FindByMethod, value interface{}) (QueryStatus, []rakutan.RakutanInfo) {
	mongoDB := CreateDBClient(env)
	defer mongoDB.Cancel()
	defer func() {
		fmt.Println("connection closed")
		if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
			panic(err)
		}
	}()
	var queryStatus QueryStatus
	var result []rakutan.RakutanInfo

	switch method {
	case ID:
		queryStatus, result = FindByLectureID(env, mongoDB, value.(int))
	case Name:
		queryStatus, result = FindByLectureName(env, mongoDB, value.(string))
	case Omikuji:
		queryStatus, result = FindByOmikuji(env, mongoDB, value.(string))
	}

	return queryStatus, result
}

func GetFavorites(env *Environments, uid string) (QueryStatus, []rakutan.Favorite) {
	mongoDB := CreateDBClient(env)
	defer mongoDB.Cancel()
	defer func() {
		fmt.Println("connection closed")
		if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
			panic(err)
		}
	}()
	var queryStatus QueryStatus
	var result []rakutan.Favorite
	queryStatus, result = FindByFav(env, mongoDB, uid)

	return queryStatus, result
}

func InsertFavorite(env *Environments, pbe PostbackEntry) QueryStatus {
	mongoDB := CreateDBClient(env)
	defer mongoDB.Cancel()
	defer func() {
		fmt.Println("connection closed")
		if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
			panic(err)
		}
	}()
	kvs := []KV{{"uid", pbe.Uid}, {"id", pbe.Param.ID}}
	singleResult := findOne(env, mongoDB, env.DB_COLLECTION.Favorites, kvs)

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
		deleteStatus := deleteOne(env, mongoDB, env.DB_COLLECTION.Favorites, kvs)
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
		countStatus, favCount := count(env, mongoDB, env.DB_COLLECTION.Favorites, []KV{{"uid", pbe.Uid}})
		fmt.Println("favCnt", favCount)
		switch {
		case countStatus.Success && favCount < 50:
			kvs = append(kvs, KV{"lecture_name", pbe.Param.LectureName})
			insertStatus := insertOne(env, mongoDB, env.DB_COLLECTION.Favorites, kvs)
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

func DeleteFavorite(env *Environments, pbe PostbackEntry) QueryStatus {
	mongoDB := CreateDBClient(env)
	defer mongoDB.Cancel()
	defer func() {
		fmt.Println("connection closed")
		if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
			panic(err)
		}
	}()
	kvs := []KV{{"uid", pbe.Uid}, {"id", pbe.Param.ID}}
	deleteStatus := deleteOne(env, mongoDB, env.DB_COLLECTION.Favorites, kvs)
	if deleteStatus.Success {
		return QueryStatus{Success: true, Message: fmt.Sprintf("「%s」をお気に入りから削除しました！", pbe.Param.LectureName)}
	} else {
		return QueryStatus{Success: false, Message: "お気に入りの削除に失敗しました。"}
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
