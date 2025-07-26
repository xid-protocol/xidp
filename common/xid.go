package common

import (
	"strings"

	"github.com/google/uuid"
	uxid "github.com/rs/xid"
)

// 传入明文生成XID
func GenerateXid(name string) string {
	xidNS := uuid.NewSHA1(uuid.NameSpaceURL, []byte("xid-protocol"))
	normalized := strings.ToLower(strings.TrimSpace(name))
	xid := uuid.NewSHA1(xidNS, []byte(normalized))
	return xid.String()
}

// 唯一随机ID，表示该条数据的ID
func GenerateCardId() string {
	uxid := uxid.New()
	return uxid.String()
}

// 唯一随机ID，表示该条数据的ID
func GenerateId() string {
	uxid := uxid.New()
	return uxid.String()
}
