package protocols

import (
	"fmt"

	"github.com/xid-protocol/xidp/common"
)

const (
	XIDVersion = "0.1.0"
)

type XIDType string

const (
	XIDTypeAWSInstance XIDType = "aws-instanceid"
	XIDTypeEmail       XIDType = "email"
)

type Metadata struct {
	CreatedAt   int64  `json:"createdAt" bson:"createdAt"`
	Operation   string `json:"operation" bson:"operation"`
	CardId      string `json:"cardId" bson:"cardId"`
	Path        string `json:"path" bson:"path"`
	ContentType string `json:"contentType" bson:"contentType"`
}

type XID struct {
	Name     string      `json:"name" bson:"name"`
	Xid      string      `json:"xid" bson:"xid"`
	Type     XIDType     `json:"type" bson:"type"`
	Version  string      `json:"version" bson:"version"`
	Metadata Metadata    `json:"metadata" bson:"metadata"`
	Payload  interface{} `json:"payload" bson:"payload"`
}

// 转换XIDType
func ConvertXIDType(typeStr string) XIDType {
	return XIDType(typeStr)
}

// MapToMetadata converts a generic map (e.g., JSON body) to a Metadata struct.
// It validates that required string fields (operation, path, contentType) are present.
// Returns an error if any required field is missing or not a string.
func MapToMetadata(m map[string]interface{}) (Metadata, error) {
	var md Metadata
	if m == nil {
		return md, fmt.Errorf("metadata is required")
	}

	// helper closure to extract a string field
	getStr := func(key string) (string, error) {
		v, ok := m[key]
		if !ok {
			return "", fmt.Errorf("metadata.%s is required", key)
		}
		s, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("metadata.%s must be a string", key)
		}
		return s, nil
	}

	var err error
	if md.Operation, err = getStr("operation"); err != nil {
		return md, err
	}
	if md.Path, err = getStr("path"); err != nil {
		return md, err
	}
	if md.ContentType, err = getStr("contentType"); err != nil {
		return md, err
	}

	return md, nil
}

func NewXID(id string, xidType XIDType, metadata Metadata, payload interface{}) *XID {
	newXID := XID{
		Name:    "xid-protocol",
		Xid:     id,
		Type:    xidType,
		Version: XIDVersion,
		Metadata: Metadata{
			CreatedAt:   common.GetTimestamp(),
			CardId:      common.GenerateCardId(),
			Operation:   metadata.Operation,
			Path:        metadata.Path,
			ContentType: metadata.ContentType,
		},
		Payload: payload,
	}
	return &newXID
}
