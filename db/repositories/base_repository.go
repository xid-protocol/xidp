package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, id string, entity *T) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
	Exists(ctx context.Context, id string) (bool, error)
	Upsert(ctx context.Context, entity *T) error
}

// MongoRepository MongoDB通用仓储实现
type MongoRepository[T any] struct {
	collection *mongo.Collection
}

// NewMongoRepository 创建MongoDB仓储
func NewMongoRepository[T any](collection *mongo.Collection) *MongoRepository[T] {
	return &MongoRepository[T]{
		collection: collection,
	}
}

// Create 创建实体
func (r *MongoRepository[T]) Create(ctx context.Context, entity *T) error {
	_, err := r.collection.InsertOne(ctx, entity)
	return err
}

// Update 更新实体
func (r *MongoRepository[T]) Update(ctx context.Context, id string, entity *T) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": entity,
		"$currentDate": bson.M{
			"updated_at": true,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete 删除实体
func (r *MongoRepository[T]) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// FindByID 根据ID查找实体
func (r *MongoRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	filter := bson.M{"_id": id}
	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// FindAll 查找所有实体
func (r *MongoRepository[T]) FindAll(ctx context.Context) ([]*T, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []*T
	for cursor.Next(ctx) {
		var entity T
		if err := cursor.Decode(&entity); err != nil {
			continue
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}

// Exists 检查实体是否存在
func (r *MongoRepository[T]) Exists(ctx context.Context, id string) (bool, error) {
	filter := bson.M{"_id": id}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Upsert 插入或更新实体
func (r *MongoRepository[T]) Upsert(ctx context.Context, entity *T) error {
	// 需要从entity中提取ID，这里假设有GetID方法
	// 实际实现中可能需要使用反射或接口
	filter := bson.M{"_id": r.extractID(entity)}
	update := bson.M{
		"$set": entity,
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// extractID 从实体中提取ID（需要根据具体实体类型实现）
func (r *MongoRepository[T]) extractID(entity *T) interface{} {
	// 这里需要使用反射或类型断言来获取ID
	// 简化实现，实际使用时需要更复杂的逻辑
	return nil
}
