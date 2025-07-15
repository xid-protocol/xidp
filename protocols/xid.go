package protocols

import (
	"strings"

	"github.com/colin-404/logx"
	"github.com/google/uuid"
	"github.com/xid-protocol/xidp/common"
)

const (
	XIDVersion = "0.1.4"
)

type Info struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type" bson:"type"`
}

type Encryption struct {
	Algorithm         string `json:"algorithm" bson:"algorithm"`
	SecretKey         string `json:"secretKey" bson:"secretKey"`
	EncryptionPayload bool   `json:"encryptionPayload" bson:"encryptionPayload"`
	EncryptionID      bool   `json:"encryptionID" bson:"encryptionID"`
}

type Metadata struct {
	CreatedAt   int64       `json:"createdAt" bson:"createdAt"`
	Encryption  *Encryption `json:"encryption" bson:"encryption"`
	Operation   string      `json:"operation" bson:"operation"`
	CardId      string      `json:"cardId" bson:"cardId"`
	Path        string      `json:"path" bson:"path"`
	ContentType string      `json:"contentType" bson:"contentType"`
}

type XID struct {
	Name     string      `json:"name" bson:"name"`
	Xid      string      `json:"xid" bson:"xid"`
	Info     *Info       `json:"info" bson:"info"`
	Version  string      `json:"version" bson:"version"`
	Metadata *Metadata   `json:"metadata" bson:"metadata"`
	Payload  interface{} `json:"payload" bson:"payload"`
}

func NewInfo(id string, xidType string) Info {
	//id can't be empty
	//xidType can't be empty
	return Info{
		ID:   id,
		Type: xidType,
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

func NewXID(info *Info, metadata *Metadata, payload interface{}) *XID {

	//如果加密key不为空，
	if metadata.Encryption != nil {
		logx.Infof("encryption: %v", metadata.Encryption)
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
