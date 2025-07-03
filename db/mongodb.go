package db

import (
	"context"
	"time"

	"github.com/colin-404/logx"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() error {
	// 从配置文件读取MongoDB设置
	mongoURI := viper.GetString("mongodb.uri")
	dbName := viper.GetString("mongodb.database")

	// 设置默认值
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	if dbName == "" {
		dbName = "xid_protocol"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logx.Errorf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// 测试连接
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
