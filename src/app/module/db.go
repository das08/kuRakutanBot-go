package module

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
	status "github.com/das08/kuRakutanBot-go/models/status"
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

	return &MongoDB{Client: client, Ctx: ctx, Cancel: cancel}
}

func findOne(e *Environments, m *MongoDB, fieldName string, value interface{}) (status.QueryStatus, []rakutan.RakutanInfo) {
	result := rakutan.RakutanInfo{}
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)
	var queryStatus status.QueryStatus

	err := collection.FindOne(m.Ctx, bson.D{{fieldName, value}}).Decode(&result)
	if err != nil {
		queryStatus.Success = false
	} else {
		queryStatus.Success = true
	}
	return queryStatus, []rakutan.RakutanInfo{result}
}

func FindByLectureID(e *Environments, m *MongoDB, lectureID int) (status.QueryStatus, []rakutan.RakutanInfo) {
	return findOne(e, m, "id", lectureID)
}

// TODO: Add error message to query status
func FindByLectureName(e *Environments, m *MongoDB, lectureName string) (status.QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)
	var queryStatus status.QueryStatus

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

func FindByOmikuji(e *Environments, m *MongoDB, omikujiType string) (status.QueryStatus, []rakutan.RakutanInfo) {
	var result []rakutan.RakutanInfo
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)
	var queryStatus status.QueryStatus

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

type FindByMethod int

const (
	Name FindByMethod = iota
	ID
	Omikuji
)

func GetRakutanInfo(env *Environments, method FindByMethod, value interface{}) (status.QueryStatus, []rakutan.RakutanInfo) {
	mongoDB := CreateDBClient(env)
	defer mongoDB.Cancel()
	defer mongoDB.Client.Disconnect(mongoDB.Ctx)
	var queryStatus status.QueryStatus
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

func randomIndex(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
