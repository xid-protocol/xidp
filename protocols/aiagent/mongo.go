package aiagent

import (
	"github.com/colin-404/logx"
	"github.com/xid-protocol/common"
	"github.com/xid-protocol/xidp/protocols"
	"go.mongodb.org/mongo-driver/mongo"
)

// Input unique Agent name
func InitWithMongo(collection *mongo.Collection, AgentName string) *protocols.XID[Config] {
	//check AgentName is valid
	xid := common.GenerateXid(AgentName)
	//check agentID is exist
	exist := common.CheckXidExistsWithMongo(collection, xid, "/protocols/aiagent")
	if exist {
		logx.Errorf("Agent %s already exists", AgentName)
		return nil
	}
	agentConfig := Config{
		tools: nil,
	}

	info := NewInfo(AgentName, "")
	metadata := NewMetadata(protocols.OperationInit)
	XID := NewXID(&info, &metadata, agentConfig)

	return XID
}

func UpdateAIAgentWithMongo(collection *mongo.Collection, AgentName string, config Config) *protocols.XID[Config] {
	//check AgentName is valid
	xid := common.GenerateXid(AgentName)
	//check agentID is exist
	exist := common.CheckXidExistsWithMongo(collection, xid, "/protocols/aiagent")
	if !exist {
		logx.Errorf("Agent %s not found", AgentName)
		return nil
	}
	info := NewInfo(AgentName, "")
	metadata := NewMetadata(protocols.OperationUpdate)
	XID := NewXID(&info, &metadata, config)

	return XID
}
