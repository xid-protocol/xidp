package xdb

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
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

	mr, ok := c.repo.(*mongoXIDRepo)
	if !ok {
		return nil, "", ErrNotImplemented
	}

	/* ---------- 1. Build 基础 match 过滤 ---------- */
	match := bson.M{"deletedAt": bson.M{"$exists": false}}
	if q.Path != "" {
		match["metadata.path"] = q.Path
	}
	// ……（NameEquals / NamePrefix / Tags / CreatedAt 范围 / AttributesEq 等
	//     逻辑与旧版保持一致，这里省略，和你之前代码一样即可）……

	/* ---------- 2. PageSize Guardrail ---------- */
	pageSize := q.PageSize
	switch {
	case pageSize <= 0:
		pageSize = 20
	case pageSize > 100:
		pageSize = 100
	}

	/* ---------- 3. 确定排序字段 & 顺序 ---------- */
	sortField := "metadata.createdAt"
	switch q.SortBy {
	case "name":
		sortField = "name"
	case "_id":
		sortField = "_id"
	case "createdAt":
		sortField = "metadata.createdAt"
	}
	asc := q.SortAsc
	sortOrder := 1
	if !asc {
		sortOrder = -1
	}

	/* ---------- 4. 聚合管道：先按 createdAt/_id 倒序拉最新 ---------- */
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		// 先保证最新的在前面
		{{Key: "$sort", Value: bson.D{
			{Key: sortField, Value: -1},
			{Key: "_id", Value: -1},
		}}},
		// 取每个 xid 最新一条
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$xid"},
			{Key: "doc", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$doc"}}}},
	}

	/* ---------- 5. After-Cursor 过滤 (复合字段) ---------- */
	if q.AfterCursor != nil && *q.AfterCursor != "" {
		// 游标格式： <ObjectID hex>|<createdAt 毫秒>
		parts := strings.Split(*q.AfterCursor, "|")
		if len(parts) == 2 {
			if lastID, err := primitive.ObjectIDFromHex(parts[0]); err == nil {
				if ts, err2 := strconv.ParseInt(parts[1], 10, 64); err2 == nil {
					lastCreated := ts
					var cmpPrimary, cmpTie string
					if asc { // 升序
						cmpPrimary, cmpTie = "$gt", "$gt"
					} else { // 降序
						cmpPrimary, cmpTie = "$lt", "$lt"
					}
					pipeline = append(pipeline,
						bson.D{{Key: "$match", Value: bson.M{
							"$or": bson.A{
								bson.M{sortField: bson.M{cmpPrimary: lastCreated}},
								bson.M{sortField: lastCreated, "_id": bson.M{cmpTie: lastID}},
							},
						}}},
					)
				}
			}
		}
	}

	/* ---------- 6. 排序 & Limit ---------- */
	pipeline = append(pipeline,
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: sortField, Value: sortOrder},
			{Key: "_id", Value: sortOrder},
		}}},
		bson.D{{Key: "$limit", Value: pageSize}},
	)

	/* ---------- 7. Projection ---------- */
	if len(q.Projection) > 0 {
		proj := bson.M{"_id": 1, "metadata.createdAt": 1}
		for _, f := range q.Projection {
			proj[f] = 1
		}
		pipeline = append(pipeline, bson.D{{Key: "$project", Value: proj}})
	}

	/* ---------- 8. 执行聚合 ---------- */
	cur, err := mr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)

	type docWithMeta struct {
		ID            primitive.ObjectID `bson:"_id"`
		CT            int64              `bson:"metadata.createdAt"`
		protocols.XID `bson:",inline"`
	}

	var (
		out           []*protocols.XID
		lastID        primitive.ObjectID
		lastCreatedAt int64
	)
	for cur.Next(ctx) {
		var d docWithMeta
		if err := cur.Decode(&d); err != nil {
			return nil, "", err
		}
		out = append(out, &d.XID)
		lastID, lastCreatedAt = d.ID, d.CT
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}

	/* ---------- 9. 生成 nextCursor ---------- */
	nextCursor := ""
	if len(out) == pageSize {
		// <ObjectID>|<createdAt>
		nextCursor = lastID.Hex() + "|" + strconv.FormatInt(lastCreatedAt, 10)
	}
	return out, nextCursor, nil
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
