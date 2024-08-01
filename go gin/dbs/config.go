package config

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var KeyCollection *mongo.Collection

func ConnectMongoDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    Client = client
    UserCollection = client.Database("myapp").Collection("users")
    KeyCollection = client.Database("myapp").Collection("Keys")

    log.Println("Connected to MongoDB!")
}
