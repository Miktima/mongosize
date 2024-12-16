package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func byteCount(bytesize int64) string {
	const unit = 1024

	if bytesize < unit {
		return fmt.Sprintf("%d B", bytesize)
	}
	div, exp := int64(unit), 0
	for n := bytesize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB",
		float64(bytesize)/float64(div), "KMGTPE"[exp])

}

func main() {
	var connection string
	var mode string

	// Ключи для командной строки
	flag.StringVar(&connection, "connection", "", "URI of Mongodb server")
	flag.StringVar(&mode, "mode", "d", "Output mode: 'd' - DB size, 'c' - collection size, 'i' - collection indexes")

	flag.Parse()

	var dbs_name []string
	var cols_name []string
	var totalSizeBite int64
	var totalStorageSizeByte int64
	var dbSizeBite int64
	var dbStorageSizeByte int64

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
	for _, dbname := range dbs_name {
		db := client.Database(dbname)
		fmt.Printf("Database: %s\n", dbname)

		cols_name, err = db.ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			log.Panic(err)
		}

		dbSizeBite = 0
		dbStorageSizeByte = 0

		for _, coll := range cols_name {
			result := db.RunCommand(ctx, bson.M{"collStats": coll})

			var document bson.M
			err = result.Decode(&document)

			if err != nil {
				panic(err)
			}

			sizeByte, _ := document["size"].(int32)
			totalSizeBite += int64(sizeByte)
			dbSizeBite += int64(sizeByte)
			storageSizeByte, _ := document["storageSize"].(int32)
			totalStorageSizeByte += int64(storageSizeByte)
			dbStorageSizeByte += int64(storageSizeByte)
			if strings.Contains("c", mode) {
				fmt.Printf(" > Collection: %s\n", coll)
				//fmt.Printf("   - Collection size: %s\n", byteCount(int64(sizeByte)))
				//fmt.Printf("   - Storage size: %s\n", byteCount(int64(storageSizeByte)))
			}
		}
		if strings.Contains("d", mode) {
			fmt.Printf("   - DB size: %s\n", byteCount(dbSizeBite))
			fmt.Printf("   - DB Storage size: %s\n", byteCount(dbStorageSizeByte))
		}
	}
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Total collection size: %s\n", byteCount(totalSizeBite))
	fmt.Printf("Total storage size: %s\n", byteCount(totalStorageSizeByte))

}
