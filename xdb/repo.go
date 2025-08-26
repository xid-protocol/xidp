package xdb

import (
	"context"

	"github.com/xid-protocol/xidp/protocols"
)

// XIDRepo defines persistence operations for XID documents.
type XIDRepo interface {
	Exists(ctx context.Context, xid, path string) (bool, error)
	Insert(ctx context.Context, doc *protocols.XID) error
	// InsertIdempotent inserts with `(metadata.path, _idk)` uniqueness; if already exists, it should not create a duplicate
	// and return nil.
	InsertIdempotent(ctx context.Context, doc *protocols.XID, idempotencyKey string) error
	Upsert(ctx context.Context, doc *protocols.XID) error
	Replace(ctx context.Context, xid, path string, doc *protocols.XID) error
	UpdateFields(ctx context.Context, xid, path string, fields map[string]any) error
	DeleteSoft(ctx context.Context, xid, path string, deletedAt int64) error
	DeleteHard(ctx context.Context, xid, path string) error
	FindByXid(ctx context.Context, xid, path string) (*protocols.XID, error)
}
