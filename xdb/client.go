package xdb

import (
	"context"
	"errors"
	"time"

	"github.com/xid-protocol/xidp/protocols"
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
	return nil, "", ErrNotImplemented
}

func (c *Client) Count(ctx context.Context, q Query) (int64, error) {
	return 0, ErrNotImplemented
}
