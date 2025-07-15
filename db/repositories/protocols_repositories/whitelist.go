package repositories

import (
	"context"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WhitelistRepository struct {
	collection *mongo.Collection
}

func NewWhitelistRepository() *WhitelistRepository {
	return &WhitelistRepository{
		collection: db.GetCollection("whitelist"),
	}
}

// 检测/protocols/whitelist是否存在
func (r *WhitelistRepository) CheckWhitelistExists(ctx context.Context, xid string, path string) (bool, error) {
	logx.Infof("xid: %s, path: %s", xid, path)
	filter := bson.M{
		"xid":           xid,
		"metadata.path": path,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// 插入新记录
func (r *WhitelistRepository) Insert(ctx context.Context, xid *protocols.XID) error {
	_, err := r.collection.InsertOne(ctx, xid)
	return err
}

// 更新记录
func (r *WhitelistRepository) UpdateWhitelist(ctx context.Context, xid string, path string, xidData *protocols.XID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"xid":           xid,
		"metadata.path": path}, bson.M{"$set": xidData})
	return err
}

func (r *WhitelistRepository) FindOneByXid(ctx context.Context, xid string) (*protocols.XID, error) {
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

func (r *WhitelistRepository) FindOneByXidAndPath(ctx context.Context, id string, path string) (*protocols.XID, error) {
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

func (r *WhitelistRepository) FindAllByPath(ctx context.Context, path string) ([]*protocols.XID, error) {
	filter := bson.M{
		"metadata.path": path,
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
