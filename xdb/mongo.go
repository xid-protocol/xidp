package xdb

import (
	"context"
	"time"

	"github.com/colin-404/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
)

func InitMongoDB(dbName string, mongoURI string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logx.Errorf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logx.Errorf("Failed to ping MongoDB: %v", err)
		return err
	}

	logx.Infof("Successfully connected to MongoDB")

	mongoClient = client
	mongoDB = client.Database(dbName)

	return nil
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mongoClient.Disconnect(ctx)
	}
	return nil
}

// GetCollection 获取指定的集合
func GetCollection(collectionName string) *mongo.Collection {
	return mongoDB.Collection(collectionName)
}

// write to mongodb collection
func WriteToMongoDB(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	return mongoDB.Collection(collectionName).InsertOne(context.Background(), data)
}
