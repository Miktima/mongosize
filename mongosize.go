package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var connection string

	// Ключи для командной строки
	flag.StringVar(&connection, "connection", "", "URI of Mongodb server")

	flag.Parse()

	//var collection *mongo.Collection
	var dbs []string

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	dbs, err = client.ListDatabaseNames(
		ctx,
		bson.D{primitive.E{Key: "empty", Value: false}})
	if err != nil {
		log.Panic(err)
	}

	for _, db := range dbs {
		fmt.Println(db)
	}

}
