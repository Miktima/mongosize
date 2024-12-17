package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
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
	var dbpattern string
	var sizeOnDisk string

	// Ключи для командной строки
	flag.StringVar(&connection, "connection", "", "URI of Mongodb server")
	flag.StringVar(&dbpattern, "dbpattern", "", "DBname filter in REgex")
	flag.StringVar(&sizeOnDisk, "size", "", "Minimum Collection Size on Disk in bytes")
	iscolls := flag.Bool("colls", false, "Output collection size")

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

	filter := bson.M{}
	filter["empty"] = false
	if len(dbpattern) > 0 {
		filter["name"] = primitive.Regex{Pattern: dbpattern}
	}
	if len(sizeOnDisk) > 0 {
		minsize, err := strconv.ParseInt(sizeOnDisk, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		filter["sizeOnDisk"] = int32(minsize)
	}
	fmt.Printf("> Filter: %s\n", filter)
	dbs_name, err = client.ListDatabaseNames(ctx, filter)
	if err != nil {
		log.Panic(err)
	}
	for _, dbname := range dbs_name {
		db := client.Database(dbname)
		fmt.Printf("> Database: %s\n", dbname)

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
			if *iscolls {
				fmt.Printf("      > Collection: %s\n", coll)
				fmt.Printf("         - Size: %s\n", byteCount(int64(sizeByte)))
				fmt.Printf("         - Storage size: %s\n", byteCount(int64(storageSizeByte)))
			}
		}
		fmt.Printf("   - DB size: %s\n", byteCount(dbSizeBite))
		fmt.Printf("   - DB Storage size: %s\n", byteCount(dbStorageSizeByte))
	}
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Total collection size: %s\n", byteCount(totalSizeBite))
	fmt.Printf("Total storage size: %s\n", byteCount(totalStorageSizeByte))

	err = client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n\nConnection to MongoDB closed.\n")
}
