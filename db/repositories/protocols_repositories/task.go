package repositories

import (
	"context"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		collection: db.GetCollection("task"),
	}
}

// 检测/protocols/task是否存在
func (r *TaskRepository) CheckTaskExists(ctx context.Context, xid string, path string) (bool, error) {
	logx.Infof("xid: %s, path: %s", xid, path)
	filter := bson.M{
		"xid":           xid,
		"metadata.path": path,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// 插入新记录
func (r *TaskRepository) Insert(ctx context.Context, xid *protocols.XID) error {
	_, err := r.collection.InsertOne(ctx, xid)
	return err
}

// 更新记录
func (r *TaskRepository) UpdateTask(ctx context.Context, xid string, path string, xidData *protocols.XID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"xid":           xid,
		"metadata.path": path}, bson.M{"$set": xidData})
	return err
}

func (r *TaskRepository) FindOneByXid(ctx context.Context, xid string) (*protocols.XID, error) {
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

func (r *TaskRepository) FindByXID(ctx context.Context, xid string) ([]*protocols.XID, error) {
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

func (r *TaskRepository) FindOneByXidAndPath(ctx context.Context, id string, path string) (*protocols.XID, error) {
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

func (r *TaskRepository) FindAllByPath(ctx context.Context, path string) ([]*protocols.XID, error) {
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

// 根据任务ID查找任务
func (r *TaskRepository) FindByTaskID(ctx context.Context, taskID string) (*protocols.XID, error) {
	filter := bson.M{
		"payload.id": taskID,
	}
	var xidRecord protocols.XID
	err := r.collection.FindOne(ctx, filter).Decode(&xidRecord)
	if err != nil {
		return nil, err
	}
	return &xidRecord, nil
}

// 根据任务状态查找任务
func (r *TaskRepository) FindByStatus(ctx context.Context, status string) ([]*protocols.XID, error) {
	filter := bson.M{
		"payload.status": status,
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

// 根据任务类型查找任务
func (r *TaskRepository) FindByType(ctx context.Context, taskType string) ([]*protocols.XID, error) {
	filter := bson.M{
		"payload.type": taskType,
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
