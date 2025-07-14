package repositories

import (
	"context"

	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type XIDRepository struct {
	collection *mongo.Collection
}

func NewXIDRepository() *XIDRepository {

	return &XIDRepository{
		collection: db.GetCollection("xids"), // 你的collection
	}
}

// 检查email是否存在
func (r *XIDRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	filter := bson.M{
		"info.email": email, // 根据你的数据结构调整
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// 插入新记录
func (r *XIDRepository) Insert(ctx context.Context, xid *protocols.XID) error {
	_, err := r.collection.InsertOne(ctx, xid)
	return err
}

func (r *XIDRepository) FindByName(ctx context.Context, name string) (*protocols.XID, error) {
	filter := bson.M{
		"name": name,
	}
	var xidInfo protocols.XID
	err := r.collection.FindOne(ctx, filter).Decode(&xidInfo)
	if err != nil {
		return nil, err
	}
	return &xidInfo, nil
}
