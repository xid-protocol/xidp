package protocols

import (
	"encoding/json"
	"strings"

	"github.com/colin-404/logx"
	"github.com/google/uuid"
	"github.com/xid-protocol/xidp/common"
)

const (
	XIDVersion = "0.1.5"
)

type OperationType string

const (
	OperationInit   OperationType = "init"
	OperationModify OperationType = "modify"
	OperationDelete OperationType = "delete"
	OperationCreate OperationType = "create"
	OperationUpdate OperationType = "update"
)

type Info struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type" bson:"type"`
}

type Encryption struct {
	Algorithm         string `json:"algorithm" bson:"algorithm"`
	SecretKey         string `json:"secretKey" bson:"secretKey"`
	EncryptionPayload bool   `json:"encryptionPayload" bson:"encryptionPayload"`
	EncryptionID      bool   `json:"encryptionID,omitempty" bson:"encryptionID,omitempty"`
}

type Metadata struct {
	CreatedAt   int64         `json:"createdAt" bson:"createdAt"`
	Encryption  *Encryption   `json:"encryption,omitempty" bson:"encryption,omitempty"`
	Operation   OperationType `json:"operation" bson:"operation"`
	CardId      string        `json:"cardId" bson:"cardId"`
	Path        string        `json:"path" bson:"path"`
	ContentType string        `json:"contentType" bson:"contentType"`
	//自定义key-value
	Extra map[string]any `json:"extra,omitempty" bson:"extra,omitempty"`
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

func NewMetadata(operation OperationType, path string, contentType string) Metadata {
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

func (m *Metadata) UnmarshalJSON(b []byte) error {
	// 1) 先解到同结构别名，拿到已知字段
	type Alias Metadata
	var tmp Alias
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	// 2) 再解到 map 捕获所有字段
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	// 3) 把已知字段删掉，其余的存进 Extra
	for _, k := range []string{
		"createdAt", "encryption", "operation",
		"cardId", "path", "contentType",
	} {
		delete(raw, k)
	}
	*m = Metadata(tmp)
	m.Extra = raw
	return nil
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	// 1) 先把已知字段转成 map
	out := map[string]any{
		"createdAt":   m.CreatedAt,
		"operation":   m.Operation,
		"cardId":      m.CardId,
		"path":        m.Path,
		"contentType": m.ContentType,
	}
	if m.Encryption != nil {
		out["encryption"] = m.Encryption
	}
	// 2) 再把 Extra 展开
	for k, v := range m.Extra {
		out[k] = v
	}
	return json.Marshal(out)
}
