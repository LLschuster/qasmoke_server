package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//InitDb initialize the mongo db connection
func InitDb() (*mongo.Database, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not load env variables, there should be a .env file in the root directory\n")
	}

	mongoUser := os.Getenv("mongo-user")
	mongoPassword := os.Getenv("mongo-password")

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%v:%v@rythim-sphere-sg9gd.mongodb.net/test?retryWrites=true&w=majority", mongoUser, mongoPassword)))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	db := client.Database("rythim_sphere")
	fmt.Println("mongo db is connected")
	return db, err
}
