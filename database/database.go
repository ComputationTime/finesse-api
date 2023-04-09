package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ComputationTime/finesse-api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database = InitDatabase()

func InitDatabase() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://172.31.0.2:27017/finesse")
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	var db = client.Database("finesse")

	return db
}

func CreateContent(source string, url string) error {

	contentObj := model.Content{
		Source: source,
		URL:    url,
	}

	_, err := database.Collection("content").InsertOne(context.Background(), contentObj)

	return err
}

func convertFromContent(slice []*model.NewContent) []interface{} {
	var result []interface{}
	for _, v := range slice {
		result = append(result, v)
	}
	return result
}

func convertToContent(slice []interface{}) []*model.NewContent {
	var result []*model.NewContent
	for _, v := range slice {
		result = append(result, v.(*model.NewContent))
	}
	return result
}

func CreateNContent(documents []*model.NewContent) error {
	db_documents := convertFromContent(documents)

	_, err := database.Collection("content").InsertMany(context.Background(), db_documents)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func TimeUUID() int {
	var ts = time.Now().UTC().Unix() << 32
	return int(ts) + int(rand.Uint32())
}

func GetContent(n int) ([]*model.Content, error) {
	var result []*model.Content
	opts := options.Find().SetSort(bson.D{{Key: "$natural", Value: -1}}).SetLimit(int64(n))
	cursor, err := database.Collection("content").Find(context.Background(), bson.D{{}}, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var content model.Content
		err := cursor.Decode(&content)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, &content)
	}

	fmt.Println(result)

	return result, nil
}
