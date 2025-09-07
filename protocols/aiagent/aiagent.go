package aiagent

import (
	"github.com/xid-protocol/common"
	"github.com/xid-protocol/xidp/protocols"
)

func NewInfo(AgentName string, desc string) protocols.Info {
	return protocols.Info{
		ID:   AgentName,
		Type: "ai-agent",
		Desc: desc,
		Tags: append([]string{}, "ai-agent"),
	}
}

func NewMetadata(operation protocols.OperationType, path string, contentType string) protocols.Metadata {
	return protocols.Metadata{
		CreatedAt:   common.GetTimestamp(),
		CardId:      common.GenerateID(),
		Operation:   operation,
		Path:        path,
		ContentType: contentType,
	}
}

func NewXID[T any](info *protocols.Info, metadata *protocols.Metadata, payload T) *protocols.XID[T] {
	return &protocols.XID[T]{
		Name:     "xid-protocol",
		Xid:      common.GenerateXid(info.ID),
		Info:     info,
		Version:  protocols.XIDVersion,
		Metadata: metadata,
		Payload:  payload,
	}
}
