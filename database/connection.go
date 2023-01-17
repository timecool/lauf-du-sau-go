package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var DB *mongo.Database
var Ctx = context.TODO()

func Connect() {
	credential := options.Credential{
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}
	clientOptions := options.Client().ApplyURI("mongodb://" + os.Getenv("MONGODB_URL")).SetAuth(credential)

	Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//Set Database
	DB = client.Database("run")
}

func InitUserCollection() *mongo.Collection {
	Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return DB.Collection("users")
}

func InitRunCollection() *mongo.Collection {
	Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return DB.Collection("runs")
}
