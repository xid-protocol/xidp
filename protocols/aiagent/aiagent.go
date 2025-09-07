package aiagent

import (
	"github.com/xid-protocol/common"
	"github.com/xid-protocol/xidp/protocols"
)

type Config struct {
	tools any
}

func NewInfo(AgentName string, systemPrompt string) protocols.Info {
	return protocols.Info{
		ID:   AgentName,
		Type: "AI_Agent_Name",
		Desc: systemPrompt,
		Tags: append([]string{}, "AI Agent"),
	}
}

func NewMetadata(operation protocols.OperationType) protocols.Metadata {
	return protocols.Metadata{
		CreatedAt:   common.GetTimestamp(),
		CardId:      common.GenerateID(),
		Operation:   operation,
		Path:        "/protocols/aiagent",
		ContentType: "application/json",
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
