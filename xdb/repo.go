package xdb

import (
	"context"

	"github.com/xid-protocol/xidp/protocols"
)

// XIDRepo defines persistence operations for XID documents.
type XIDRepo interface {
	Exists(ctx context.Context, xid, path string) (bool, error)
	Insert(ctx context.Context, doc *protocols.XID) error
	List(ctx context.Context, q Query) ([]*protocols.XID, string, error)
	InsertIdempotent(ctx context.Context, doc *protocols.XID, idempotencyKey string) error
	Upsert(ctx context.Context, doc *protocols.XID) error
	Replace(ctx context.Context, xid, path string, doc *protocols.XID) error
	UpdateFields(ctx context.Context, xid, path string, fields map[string]any) error
	DeleteSoft(ctx context.Context, xid, path string, deletedAt int64) error
	DeleteHard(ctx context.Context, xid, path string) error
	FindByXid(ctx context.Context, xid, path string) (*protocols.XID, error)
}
