package mongo

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoManager struct{}

var SearchLimit int64

func (m MongoManager) GetDatabases(mongoUri string) ([]string, error) {
	mongoClient := getMongoClient(mongoUri)

	return mongoClient.ListDatabaseNames(context.TODO(), bson.D{})
}

func (m MongoManager) GetCollections(mongoUri string, databaseName string) ([]string, error) {
	mongoClient := getMongoClient(mongoUri)
	return mongoClient.Database(databaseName).ListCollectionNames(context.TODO(), bson.D{})
}

func (m MongoManager) Find(mongoUri string, databaseName string, collectionName string, query string) ([]any, error) {
	mongoClient := getMongoClient(mongoUri)
	optsFind := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(SearchLimit)

	parsedQuery := getParsedQuery(query)

	cursor, err := mongoClient.Database(databaseName).Collection(collectionName).Find(context.TODO(), parsedQuery, optsFind)
	if err != nil {
		panic("Error finding documents")
	}

	var data []any
	err = cursor.All(context.TODO(), &data)

	return data, err
}

func (m MongoManager) Insert(mongoUri string, databaseName string, collectionName string, documents []any, clearCollection bool) error {
	mongoClient := getMongoClient(mongoUri)

	var err error

	if clearCollection {
		err = mongoClient.Database(databaseName).Collection(collectionName).Drop(context.TODO())
		if err != nil {
			panic("Error dropping collection")
		}
	}

	_, err = mongoClient.Database(databaseName).Collection(collectionName).InsertMany(context.TODO(), documents)

	return err
}

func getMongoClient(mongoUri string) *mongo.Client {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic("Error connecting to mongo: " + mongoUri)
	}

	return mongoClient
}

func getParsedQuery(query string) bson.D {
	if query == "" {
		return bson.D{}
	}

	var c map[string]any
	if err := json.Unmarshal([]byte(query), &c); err != nil {
		panic("Error parsing query")
	}

	var parsedQuery bson.D
	for key, value := range c {
		parsedQuery = append(parsedQuery, bson.E{Key: key, Value: value})
	}

	return parsedQuery
}
