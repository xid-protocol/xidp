package repositories

import (
	"context"
	"time"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JumpServerUserRepository JumpServer用户仓储
type JumpServerUserRepository struct {
	collection *mongo.Collection
}

// NewJumpServerUserRepository 创建JumpServer用户仓储
func NewJumpServerUserRepository() *JumpServerUserRepository {
	collection := db.GetCollection("jumpserver_users")

	// 创建索引
	r := &JumpServerUserRepository{collection: collection}
	r.createIndexes()

	return r
}

// createIndexes 创建索引
func (r *JumpServerUserRepository) createIndexes() {
	ctx := context.Background()

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "source", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logx.Errorf("Failed to create indexes for jumpserver_users: %v", err)
	}
}

// Create 创建新用户
func (r *JumpServerUserRepository) Create(ctx context.Context, user *models.JumpServerUser) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update 更新用户
func (r *JumpServerUserRepository) Update(ctx context.Context, user *models.JumpServerUser) error {
	filter := bson.M{"_id": user.ID}

	user.UpdatedAt = time.Now()
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Upsert 插入或更新用户
func (r *JumpServerUserRepository) Upsert(ctx context.Context, user *models.JumpServerUser) error {
	filter := bson.M{"_id": user.ID}

	// now := time.Now()
	// user.UpdatedAt = now

	// update := bson.M{
	// 	"$set": user,
	// 	"$setOnInsert": bson.M{
	// 		"created_at": now,
	// 	},
	// }

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, opts)
	return err
}

// FindByID 根据ID查找用户
func (r *JumpServerUserRepository) FindByID(ctx context.Context, id string) (*models.JumpServerUser, error) {
	filter := bson.M{"_id": id}

	var user models.JumpServerUser
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *JumpServerUserRepository) FindByUsername(ctx context.Context, username string) (*models.JumpServerUser, error) {
	filter := bson.M{"username": username}

	var user models.JumpServerUser
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindAll 查找所有用户
func (r *JumpServerUserRepository) FindAll(ctx context.Context) ([]*models.JumpServerUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.JumpServerUser
	for cursor.Next(ctx) {
		var user models.JumpServerUser
		if err := cursor.Decode(&user); err != nil {
			logx.Errorf("Failed to decode user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// FindActiveUsers 查找活跃用户
func (r *JumpServerUserRepository) FindActiveUsers(ctx context.Context) ([]*models.JumpServerUser, error) {
	filter := bson.M{"is_active": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.JumpServerUser
	for cursor.Next(ctx) {
		var user models.JumpServerUser
		if err := cursor.Decode(&user); err != nil {
			logx.Errorf("Failed to decode user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// Exists 检查用户是否存在
func (r *JumpServerUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	filter := bson.M{"_id": id}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete 删除用户
func (r *JumpServerUserRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// Count 统计用户数量
func (r *JumpServerUserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountBySource 根据来源统计用户数量
func (r *JumpServerUserRepository) CountBySource(ctx context.Context, source string) (int64, error) {
	filter := bson.M{"source": source}
	return r.collection.CountDocuments(ctx, filter)
}
