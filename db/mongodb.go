package db

import (
	"context"
	"time"

	"github.com/colin-404/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
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

	MongoClient = client
	MongoDB = client.Database(dbName)

	return nil
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return MongoClient.Disconnect(ctx)
	}
	return nil
}

// GetCollection 获取指定的集合
func GetCollection(collectionName string) *mongo.Collection {
	return MongoDB.Collection(collectionName)
}
