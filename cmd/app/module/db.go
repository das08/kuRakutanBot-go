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

func CreateClient() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://********:********@localhost:27017"))
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
	collection := client.Database("rakutanDB").Collection("rakutan2021")

	err = collection.FindOne(ctx, bson.D{{"facultyname", "工学部"}}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	client.Disconnect(ctx)
}
