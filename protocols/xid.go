package protocols

import (
	"strings"

	"github.com/google/uuid"
	"github.com/xid-protocol/xidp/common"
)

const (
	XIDVersion = "0.1.3"
)

const (
	NoEncryption bool = false
	Encryption   bool = true
)

type Info struct {
	ID   string `json:"plainID" bson:"plainID"`
	Type string `json:"type" bson:"type"`
	//是否加密ID
	Encryption bool `json:"encryption" bson:"encryption"`
}

type Metadata struct {
	//加密算法
	EncryptionAlgorithm string `json:"encryptionAlgorithm" bson:"encryptionAlgorithm"`
	//加密key
	EncryptionKey string `json:"encryptionKey" bson:"encryptionKey"`
	//是否加密payload
	Encryption  bool   `json:"encryptionPayload" bson:"encryptionPayload"`
	CreatedAt   int64  `json:"createdAt" bson:"createdAt"`
	Operation   string `json:"operation" bson:"operation"`
	CardId      string `json:"cardId" bson:"cardId"`
	Path        string `json:"path" bson:"path"`
	ContentType string `json:"contentType" bson:"contentType"`
}

type XID struct {
	Name     string      `json:"name" bson:"name"`
	Xid      string      `json:"xid" bson:"xid"`
	Info     Info        `json:"info" bson:"info"`
	Version  string      `json:"version" bson:"version"`
	Metadata Metadata    `json:"metadata" bson:"metadata"`
	Payload  interface{} `json:"payload" bson:"payload"`
}

func NewInfo(id string, xidType string, encryption bool) Info {
	//id can't be empty
	//xidType can't be empty
	return Info{
		ID:         id,
		Type:       xidType,
		Encryption: encryption,
	}
}

func NewMetadata(operation string, path string, contentType string) Metadata {
	return Metadata{
		CreatedAt:   common.GetTimestamp(),
		CardId:      common.GenerateCardId(),
		Operation:   operation,
		Path:        path,
		ContentType: contentType,
	}
}

func NewXID(info Info, metadata Metadata, payload interface{}) *XID {

	//如果加密key不为空，
	if metadata.EncryptionKey != "" && metadata.EncryptionAlgorithm != "" {
		//加密payload
		if metadata.Encryption {
			//加密payload
		}
		if info.Encryption {
			//加密ID
		}
	}

	newXID := XID{
		Name:     "xid-protocol",
		Xid:      GenerateXid(info.ID),
		Info:     info,
		Version:  XIDVersion,
		Metadata: metadata,
		Payload:  payload,
	}
	return &newXID
}

// 传入明文，生成xid
func GenerateXid(id string) string {
	xidNS := uuid.NewSHA1(uuid.NameSpaceURL, []byte("xid-protocol"))
	normalized := strings.ToLower(strings.TrimSpace(id))
	xid := uuid.NewSHA1(xidNS, []byte(normalized))
	return xid.String()
}
