package xdb

import (
	"context"
	"time"

	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoXIDRepo struct {
	collection *mongo.Collection
}

func NewMongoXIDRepo(c *mongo.Collection) XIDRepo {
	return &mongoXIDRepo{collection: c}
}

func (r *mongoXIDRepo) List(ctx context.Context, q Query) ([]*protocols.XID, string, error) {
	return nil, "", nil
}

// EnsureXIDIndexes creates recommended indexes on the collection.
// func EnsureXIDIndexes(ctx context.Context, c *mongo.Collection) error {
// 	models := []mongo.IndexModel{
// 		{Keys: bson.D{{Key: "xid", Value: 1}, {Key: "metadata.path", Value: 1}}, Options: options.Index().SetUnique(true)},
// 		{Keys: bson.D{{Key: "idempotencyKey", Value: 1}, {Key: "metadata.path", Value: 1}}, Options: options.Index().SetUnique(true)},
// 		{Keys: bson.D{{Key: "metadata.createdAt", Value: 1}, {Key: "metadata.path", Value: 1}}},
// 	}
// 	_, err := c.Indexes().CreateMany(ctx, models)
// 	return err
// }

func (r *mongoXIDRepo) Exists(ctx context.Context, xid, path string) (bool, error) {
	filter := bson.M{"xid": xid, "metadata.path": path, "deletedAt": bson.M{"$exists": false}}
	err := r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return err == nil, err
}

func (r *mongoXIDRepo) Insert(ctx context.Context, doc *protocols.XID) error {
	if doc.Metadata != nil && doc.Metadata.CreatedAt == 0 {
		doc.Metadata.CreatedAt = time.Now().UnixMilli()
	}
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *mongoXIDRepo) InsertIdempotent(ctx context.Context, doc *protocols.XID, idempotencyKey string) error {
	if doc.Metadata != nil && doc.Metadata.CreatedAt == 0 {
		doc.Metadata.CreatedAt = time.Now().UnixMilli()
	}
	filter := bson.M{"metadata.path": doc.Metadata.Path, "idempotencyKey": idempotencyKey}
	// upsert with full doc content only on first insert
	setOnInsert := bson.M{}
	// marshal doc into map to merge with extra field
	raw, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	if err := bson.Unmarshal(raw, &setOnInsert); err != nil {
		return err
	}
	setOnInsert["idempotencyKey"] = idempotencyKey
	update := bson.M{"$setOnInsert": setOnInsert}
	_, err = r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *mongoXIDRepo) Upsert(ctx context.Context, doc *protocols.XID) error {
	filter := bson.M{"xid": doc.Xid, "metadata.path": doc.Metadata.Path}
	update := bson.M{"$set": doc}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *mongoXIDRepo) Replace(ctx context.Context, xid, path string, doc *protocols.XID) error {
	filter := bson.M{"xid": xid, "metadata.path": path}
	_, err := r.collection.ReplaceOne(ctx, filter, doc)
	return err
}

func (r *mongoXIDRepo) UpdateFields(ctx context.Context, xid, path string, fields map[string]any) error {
	filter := bson.M{"xid": xid, "metadata.path": path}
	update := bson.M{"$set": fields}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// modify

func (r *mongoXIDRepo) DeleteSoft(ctx context.Context, xid, path string, deletedAt int64) error {
	filter := bson.M{"xid": xid, "metadata.path": path}
	update := bson.M{"$set": bson.M{"deletedAt": deletedAt}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoXIDRepo) DeleteHard(ctx context.Context, xid, path string) error {
	filter := bson.M{"xid": xid, "metadata.path": path}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *mongoXIDRepo) FindByXid(ctx context.Context, xid, path string) (*protocols.XID, error) {
	filter := bson.M{"xid": xid, "metadata.path": path, "deletedAt": bson.M{"$exists": false}}
	var out protocols.XID
	if err := r.collection.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
