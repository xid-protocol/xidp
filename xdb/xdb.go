package xdb

import (
	"context"
	"time"

	"github.com/xid-protocol/xidp/protocols"
)

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

// XIDRepo defines persistence operations for XID documents.
type XIDRepo interface {
	Exists(ctx context.Context, xid, path string) (bool, error)
	Insert(ctx context.Context, doc *protocols.XID[any]) error
	List(ctx context.Context, q Query) ([]*protocols.XID[any], string, error)
	InsertIdempotent(ctx context.Context, doc *protocols.XID[any], idempotencyKey string) error
	Upsert(ctx context.Context, xid, path string, doc any) error
	Replace(ctx context.Context, xid, path string, doc *protocols.XID[any]) error
	UpdateFields(ctx context.Context, xid, path string, fields map[string]any) error
	DeleteSoft(ctx context.Context, xid, path string, deletedAt int64) error
	DeleteHard(ctx context.Context, xid, path string) error
	FindByXid(ctx context.Context, xid, path string) (*protocols.XID[any], error)
}
