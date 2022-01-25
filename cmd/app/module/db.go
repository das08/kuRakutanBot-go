package module

import (
	"context"
	"log"
	"time"

	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
	status "github.com/das08/kuRakutanBot-go/models/status"
	"go.mongodb.org/mongo-driver/bson"
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

func findOne(e *Environments, m *MongoDB, fieldName string, value interface{}) (status.QueryStatus, rakutan.RakutanInfo) {
	result := rakutan.RakutanInfo{}
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)
	var queryStatus status.QueryStatus

	err := collection.FindOne(m.Ctx, bson.D{{fieldName, value}}).Decode(&result)
	if err != nil {
		queryStatus.Success = false
	} else {
		queryStatus.Success = true
	}
	return queryStatus, result
}

func FindByLectureID(e *Environments, m *MongoDB, lectureID int) (status.QueryStatus, rakutan.RakutanInfo) {
	return findOne(e, m, "id", lectureID)
}
