package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ComputationTime/finesse-api/graph/model"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var database = InitDatabase()

func InitDatabase() *mongo.Database {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

	database = client.Database("finesse")

    return database
}

func CreateContent(source string, url string) error {
	contentObj := model.Content{
		Source:    source,
		URL:       url,
	}

	_, err = database.Collection("content").InsertOne(context.Background(), contentObj)
	return err
}

func CreateNContent(documents []*model.NewContent) (int, error) {


	// do we need to change the input type?
	documents := []interface{}{
        bson.D{
            {"name", "John"},
            {"age", 30},
        },
        bson.D{
            {"name", "Jane"},
            {"age", 25},
        },
        bson.D{
            {"name", "Bob"},
            {"age", 50},
        },
    }

    _, err = collection.InsertMany(context.Background(), documents)
    if err != nil {
        log.Fatal(err)
		return nil, err
    }

	return "", nil
}

func TimeUUID() int {
	var ts = time.Now().UTC().Unix() << 32
	return int(ts) + int(rand.Uint32())
}

// hopefully won't need this func and we can make a general get function
func GetContent(n int32) ([]*model.Content, error){
	
}