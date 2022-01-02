package module

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateDBClient(e *Environments) {
	mongoURI := "mongodb://" + e.DB_USER + ":" + e.DB_PASS + "@" + e.DB_HOST + ":" + e.DB_PORT + ""
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var result bson.D
	collection := client.Database(e.DB_NAME).Collection(e.DB_COLLECTION)

	err = collection.FindOne(ctx, bson.D{{"facultyname", "工学部"}}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	client.Disconnect(ctx)
}
