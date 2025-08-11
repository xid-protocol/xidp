package db

import (
	"context"

	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type XidInfoRepository struct {
	collection *mongo.Collection
}

func NewXidInfoRepository() *XidInfoRepository {
	return &XidInfoRepository{
		collection: GetCollection("xid_info"), // 你的collection
	}
}

// 检测/info/sealsuite是否存在
// func (r *XidInfoRepository) CheckXidExists(ctx context.Context, xid string, path string) (bool, int, error) {
// 	logx.Infof("xid: %s, path: %s", xid, path)
// 	filter := bson.M{
// 		"xid":           xid,
// 		"metadata.path": path,
// 	}

// 	count, err := r.collection.CountDocuments(ctx, filter)
// 	if count == 0 {
// 		return false, 0, err
// 	}
// 	return true, int(count), nil
// }

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

func (r *XidInfoRepository) FindAllByXid(ctx context.Context, xid string) ([]*protocols.XID, error) {
	filter := bson.M{
		"xid": xid,
	}
	var xidRecords []*protocols.XID
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var xidRecord protocols.XID
		if err := cursor.Decode(&xidRecord); err != nil {
			return nil, err
		}
		xidRecords = append(xidRecords, &xidRecord)
	}
	return xidRecords, nil
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
