package common

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateXid(name string) string {
	xidNS := uuid.NewSHA1(uuid.NameSpaceURL, []byte("xid-protocol"))
	normalized := strings.ToLower(strings.TrimSpace(name))
	xid := uuid.NewSHA1(xidNS, []byte(normalized))
	return xid.String()
}
