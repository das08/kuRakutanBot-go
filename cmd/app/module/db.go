package module

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/das08/kuRakutanBot-go/models/rakutan"
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

func FindOne(e *Environments, m *MongoDB) {
	result := models.RakutanInfo{}
	collection := m.Client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)

	err := collection.FindOne(m.Ctx, bson.D{{"lecture_name", "半導体工学"}}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %#v", result)
}
