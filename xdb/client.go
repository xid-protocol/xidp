package xdb

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrConflict         = errors.New("conflict")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPermissionDenied = errors.New("permission denied")
	ErrDeadlineExceeded = errors.New("deadline exceeded")
	ErrInternal         = errors.New("internal")
	ErrNotImplemented   = errors.New("not implemented")
)

// XIDRepo is defined in repo.go

type ClientOptions struct {
	// DefaultTimeout applies when ctx has no deadline.
	DefaultTimeout time.Duration
	// EnableIdempotency controls Create idempotency behavior if supported by repo/DB layer.
	EnableIdempotency bool
}

type Client struct {
	repo XIDRepo
	opts ClientOptions
}

func NewClientWithRepo(repo XIDRepo, opts *ClientOptions) *Client {
	co := ClientOptions{}
	if opts != nil {
		co = *opts
	}
	return &Client{repo: repo, opts: co}
}

func NewClientWithMongo(collection *mongo.Collection, opts *ClientOptions) *Client {
	repo := NewMongoXIDRepo(collection)
	return NewClientWithRepo(repo, opts)
}

// withTimeout returns a context with default timeout if caller did not provide one.
func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline || c.opts.DefaultTimeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.opts.DefaultTimeout)
}

// Create inserts a new XID document. When EnableIdempotency=true and an idempotencyKey is provided,
// the implementation should guarantee at-most-once creation under the same key within a path.
func (c *Client) Create(ctx context.Context, path string, doc *protocols.XID, idempotencyKey ...string) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	if doc == nil || doc.Metadata == nil {
		return ErrInvalidArgument
	}
	doc.Metadata.Path = path
	if c.opts.EnableIdempotency && len(idempotencyKey) > 0 && idempotencyKey[0] != "" {
		return c.repo.InsertIdempotent(ctx, doc, idempotencyKey[0])
	}
	return c.repo.Insert(ctx, doc)
}

// Upsert creates or updates the document identified by (xid, path).
func (c *Client) Upsert(ctx context.Context, doc *protocols.XID) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	if doc == nil || doc.Metadata == nil {
		return ErrInvalidArgument
	}
	return c.repo.Upsert(ctx, doc)
}

func (c *Client) GetByXid(ctx context.Context, path, xid string) (*protocols.XID, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	out, err := c.repo.FindByXid(ctx, xid, path)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return out, nil
}

// UpdateReplace replaces the full document identified by (xid, path).
func (c *Client) UpdateReplace(ctx context.Context, path, xid string, doc *protocols.XID) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	if doc == nil {
		return ErrInvalidArgument
	}
	return c.repo.Replace(ctx, xid, path, doc)
}

// UpdateFields applies partial field updates to the document identified by (xid, path).
func (c *Client) UpdateFields(ctx context.Context, path, xid string, fields map[string]any) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	if len(fields) == 0 {
		return nil
	}
	return c.repo.UpdateFields(ctx, xid, path, fields)
}

// DeleteSoft marks the document as deleted without physically removing it.
func (c *Client) DeleteSoft(ctx context.Context, path, xid string) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.repo.DeleteSoft(ctx, xid, path, time.Now().UnixMilli())
}

// DeleteHard permanently removes the document.
func (c *Client) DeleteHard(ctx context.Context, path, xid string) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.repo.DeleteHard(ctx, xid, path)
}

func (c *Client) Exists(ctx context.Context, path, xid string) (bool, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.repo.Exists(ctx, xid, path)
}

type Query struct {
	Path         string
	NameEquals   *string
	NamePrefix   *string
	TagsAll      []string
	CreatedAtGTE *time.Time
	CreatedAtLT  *time.Time
	AttributesEq map[string]any
	SortBy       string // "createdAt","name","_id"
	SortAsc      bool
	PageSize     int
	AfterCursor  *string
	Projection   []string
}

func (c *Client) List(ctx context.Context, q Query) ([]*protocols.XID, string, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	// Only implemented for Mongo-backed repos for now
	mr, ok := c.repo.(*mongoXIDRepo)
	if !ok {
		return nil, "", ErrNotImplemented
	}

	// Build match filter
	match := bson.M{"deletedAt": bson.M{"$exists": false}}
	if q.Path != "" {
		match["metadata.path"] = q.Path
	}
	if q.NameEquals != nil {
		match["name"] = *q.NameEquals
	}
	if q.NamePrefix != nil {
		prefix := "^" + regexp.QuoteMeta(*q.NamePrefix)
		match["name"] = bson.M{"$regex": prefix}
	}
	if len(q.TagsAll) > 0 {
		match["info.tags"] = bson.M{"$all": q.TagsAll}
	}
	if q.CreatedAtGTE != nil || q.CreatedAtLT != nil {
		rangeCond := bson.M{}
		if q.CreatedAtGTE != nil {
			rangeCond["$gte"] = q.CreatedAtGTE.UnixMilli()
		}
		if q.CreatedAtLT != nil {
			rangeCond["$lt"] = q.CreatedAtLT.UnixMilli()
		}
		match["metadata.createdAt"] = rangeCond
	}
	if len(q.AttributesEq) > 0 {
		for k, v := range q.AttributesEq {
			match["metadata.extra."+k] = v
		}
	}

	// Page size guardrails
	pageSize := q.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Final sort selection
	sortField := "metadata.createdAt"
	switch q.SortBy {
	case "name":
		sortField = "name"
	case "_id":
		sortField = "_id"
	case "createdAt":
		sortField = "metadata.createdAt"
	}
	sortOrder := 1
	if !q.SortAsc {
		sortOrder = -1
	}

	// Build aggregation pipeline to get latest doc per xid
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		// Ensure newest first so $first in group is the latest per xid
		bson.D{{Key: "$sort", Value: bson.D{{Key: "metadata.createdAt", Value: -1}, {Key: "_id", Value: -1}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$xid"},
			{Key: "doc", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$doc"}}}},
	}

	// After-cursor on _id (applied after grouping)
	if q.AfterCursor != nil && *q.AfterCursor != "" {
		if oid, err := primitive.ObjectIDFromHex(*q.AfterCursor); err == nil {
			cmp := "$gt"
			if !q.SortAsc {
				cmp = "$lt"
			}
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"_id": bson.M{cmp: oid}}}})
		}
	}

	// Sort and limit on the latest-per-xid set
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: sortOrder}, {Key: "_id", Value: sortOrder}}}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: pageSize}})

	// Projection if provided
	if len(q.Projection) > 0 {
		proj := bson.M{"_id": 1}
		for _, f := range q.Projection {
			proj[f] = 1
		}
		pipeline = append(pipeline, bson.D{{Key: "$project", Value: proj}})
	}

	cur, err := mr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)

	type aggXIDDoc struct {
		ID            primitive.ObjectID `bson:"_id"`
		protocols.XID `bson:",inline"`
	}

	var (
		out      []*protocols.XID
		lastID   primitive.ObjectID
		haveLast bool
	)
	for cur.Next(ctx) {
		var rec aggXIDDoc
		if err := cur.Decode(&rec); err != nil {
			return nil, "", err
		}
		out = append(out, &rec.XID)
		lastID = rec.ID
		haveLast = true
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}

	next := ""
	if len(out) == pageSize && haveLast {
		next = lastID.Hex()
	}
	return out, next, nil
}

func (c *Client) Count(ctx context.Context, q Query) (int64, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	mr, ok := c.repo.(*mongoXIDRepo)
	if !ok {
		return 0, ErrNotImplemented
	}

	match := bson.M{"deletedAt": bson.M{"$exists": false}}
	if q.Path != "" {
		match["metadata.path"] = q.Path
	}
	if q.NameEquals != nil {
		match["name"] = *q.NameEquals
	}
	if q.NamePrefix != nil {
		prefix := "^" + regexp.QuoteMeta(*q.NamePrefix)
		match["name"] = bson.M{"$regex": prefix}
	}
	if len(q.TagsAll) > 0 {
		match["info.tags"] = bson.M{"$all": q.TagsAll}
	}
	if q.CreatedAtGTE != nil || q.CreatedAtLT != nil {
		rangeCond := bson.M{}
		if q.CreatedAtGTE != nil {
			rangeCond["$gte"] = q.CreatedAtGTE.UnixMilli()
		}
		if q.CreatedAtLT != nil {
			rangeCond["$lt"] = q.CreatedAtLT.UnixMilli()
		}
		match["metadata.createdAt"] = rangeCond
	}
	if len(q.AttributesEq) > 0 {
		for k, v := range q.AttributesEq {
			match["metadata.extra."+k] = v
		}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$xid"}}}},
		bson.D{{Key: "$count", Value: "total"}},
	}

	cur, err := mr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	var total int64
	if cur.Next(ctx) {
		var row struct {
			Total int64 `bson:"total"`
		}
		if err := cur.Decode(&row); err != nil {
			return 0, err
		}
		total = row.Total
	}
	if err := cur.Err(); err != nil {
		return 0, err
	}
	return total, nil
}

// ListLatest returns the latest document per xid with pagination, using the same
// semantics as List. It is an alias for clarity of intent.
// func (c *Client) ListLatest(ctx context.Context, q Query) ([]*protocols.XID, string, error) {
// 	return c.List(ctx, q)
// }
