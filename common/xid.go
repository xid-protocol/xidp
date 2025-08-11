package common

import (
	"errors"

	"github.com/google/uuid"
	uxid "github.com/rs/xid"
)

// 传入明文生成XID
func GenerateXid(name string) string {
	xidNS := uuid.NewSHA1(uuid.NameSpaceURL, []byte("xid-protocol"))
	// normalized := strings.ToLower(strings.TrimSpace(name))
	xid := uuid.NewSHA1(xidNS, []byte(name))
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

func GenerateSHA1(text string) string {
	xidNS := uuid.NewSHA1(uuid.NameSpaceURL, []byte(text))
	xid := uuid.NewSHA1(xidNS, []byte(text))
	return xid.String()
}

// 生成唯一随机ID
func GenerateRandomId(idType string) (string, error) {
	switch idType {
	case "xid":
		uxid := uxid.New()
		return uxid.String(), nil
	case "uuid":
		xid := uuid.New()
		return xid.String(), nil
	default:
		return "", errors.New("invalid type")
	}
}
