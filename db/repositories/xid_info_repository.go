package repositories

import (
	"context"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type XidInfoRepository struct {
	collection *mongo.Collection
}

func NewXidInfoRepository() *XidInfoRepository {
	return &XidInfoRepository{
		collection: db.GetCollection("xid_info"), // 你的collection
	}
}

// 检测/info/sealsuite是否存在
func (r *XidInfoRepository) CheckXidInfoExists(ctx context.Context, xid string, path string) (bool, error) {
	logx.Infof("xid: %s, path: %s", xid, path)
	filter := bson.M{
		"xid":           xid,
		"metadata.path": path,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// 插入新记录
func (r *XidInfoRepository) Insert(ctx context.Context, xid *protocols.XID) error {
	_, err := r.collection.InsertOne(ctx, xid)
	return err
}

// 更新记录
func (r *XidInfoRepository) UpdateXidInfo(ctx context.Context, xid string, path string, xidData *protocols.XID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"xid":           xid,
		"metadata.path": path}, bson.M{"$set": xidData})
	return err
}

func (r *XidInfoRepository) FindByName(ctx context.Context, name string, path string) (*protocols.XID, error) {
	filter := bson.M{
		"name": name,
		"path": path,
	}
	var xidRecord protocols.XID
	err := r.collection.FindOne(ctx, filter).Decode(&xidRecord)
	if err != nil {
		return nil, err
	}
	return &xidRecord, nil
}

func (r *XidInfoRepository) FindByXid(ctx context.Context, xid string) (*protocols.XID, error) {
	filter := bson.M{
		"xid": xid,
	}
	var xidRecord protocols.XID
	err := r.collection.FindOne(ctx, filter).Decode(&xidRecord)
	if err != nil {
		return nil, err
	}
	return &xidRecord, nil
}

func (r *XidInfoRepository) FindOneByXidAndPath(ctx context.Context, id string, path string) (*protocols.XID, error) {
	filter := bson.M{
		"xid":           id,
		"metadata.path": path,
	}
	var xidRecord protocols.XID
	err := r.collection.FindOne(ctx, filter).Decode(&xidRecord)
	if err != nil {
		return nil, err
	}

	return &xidRecord, nil
}
