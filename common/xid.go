package common

import (
	"strings"

	"github.com/google/uuid"
	uxid "github.com/rs/xid"
)

// uuidv5，表示该条数据的身份ID
func GenerateId(name string) string {
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
