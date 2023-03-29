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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Connect() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
        o.Region = "ca-central-1"
        return nil
    })
	if err != nil {
        panic(err)
    }
	return dynamodb.NewFromConfig(cfg)
}

func accessTable(tableName string) TableStruct {
	return TableStruct{Connect(), tableName}
}

func CreateTable(tableName string) (*types.TableDescription, error) {
	tableAccess := accessTable(tableName)
	var tableDesc *types.TableDescription
	table, err := tableAccess.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
		// 	{
		// 	AttributeName: aws.String("year"),
		// 	AttributeType: types.ScalarAttributeTypeN,
		// }, {
		// 	AttributeName: aws.String("title"),
		// 	AttributeType: types.ScalarAttributeTypeS,
		// }},
		// KeySchema: []types.KeySchemaElement{{
		// 	AttributeName: aws.String("year"),
		// 	KeyType:       types.KeyTypeHash,
		// }, {
		// 	AttributeName: aws.String("title"),
		// 	KeyType:       types.KeyTypeRange,
		// }
	},
		TableName: aws.String(tableAccess.TableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableAccess.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(tableAccess.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableAccess.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

func CreateContent(source string, url string) (string, error) {
	svc := Connect()
	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String("content"),
        Item: map[string]types.AttributeValue{
            "content_id":    &types.AttributeValueMemberN{Value: strconv.Itoa(TimeUUID())},
            "source":  &types.AttributeValueMemberS{Value: source},
            "url": &types.AttributeValueMemberS{Value: url},
        },
    })

	contentJson, err := json.Marshal(out.Attributes)
	contentString := string(contentJson)

	return contentString, err
}

func CreateNContent(array []*model.NewContent) (int, error) {
	basics := accessTable("content") 
	var err error
	var item map[string]types.AttributeValue
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < len(array) {
		var writeReqs []types.WriteRequest
		if end > len(array) {
			end = len(array)
		}
		for _, movie := range array[start:end] {
			item, err = attributevalue.MarshalMap(movie)
			if err != nil {
				log.Printf("Couldn't marshal for batch writing. Here's why: %v\n", err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}
		_, err = basics.DynamoDbClient.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{basics.TableName: writeReqs}})
		if err != nil {
			log.Printf("Couldn't add a batch of movies to %v. Here's why: %v\n", basics.TableName, err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	return written, err
}

func TimeUUID() int {
	var ts = time.Now().UTC().Unix() << 32
	return int(ts) + int(rand.Uint32())
}

func GetContent(n int32) ([]*model.Content, error){
	var response []*model.Content
	svc := Connect()
	input := &dynamodb.ScanInput{
		TableName: aws.String("content"),
		Limit: aws.Int32(n),
    }

    result, err := svc.Scan(context.TODO(), input)
    if err != nil {
        fmt.Println(err.Error())
    }

    attributevalue.UnmarshalListOfMaps(result.Items, &response)
	
	return response, err
}


type TableStruct struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}