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

func getHtmlPage(url, userAgent string) ([]byte, int, error) {

}

func main() {
	var connection string

	// Ключи для командной строки
	flag.StringVar(&connection, "connection", "", "URI of Mongodb server")

	flag.Parse()

	//var collection *mongo.Collection
	var dbs_name []string
	var cols_name []string

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	dbs_name, err = client.ListDatabaseNames(
		ctx,
		bson.D{primitive.E{Key: "empty", Value: false}})
	if err != nil {
		log.Panic(err)
	}
	dbname := dbs_name[0]
	db := client.Database(dbname)

	cols_name, err = db.ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		log.Panic(err)
	}

	result := db.RunCommand(ctx, bson.M{"collStats": cols_name[0]})

	var document bson.M
	err = result.Decode(&document)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Database: %s, collection: %s", dbname, cols_name[0])
	fmt.Printf(" - Collection size: %v Bytes\n", document["size"])
	fmt.Printf(" - Average object size: %v Bytes\n", document["avgObjSize"])
	fmt.Printf(" - Storage size: %v Bytes\n", document["storageSize"])
	fmt.Printf(" - Total index size: %v Bytes\n", document["totalIndexSize"])

}
