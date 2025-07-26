package mcptask

import (
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/protocols"
)

type State string

const (
	StateInit      State = "init"
	StateCreated   State = "created"
	StateWorking   State = "working"
	StateCompleted State = "completed"
	StateCanceled  State = "canceled"
	StateFailed    State = "failed"
	StateRejected  State = "rejected"
	StateUnknown   State = "unknown"
)

type Step struct {
	TaskID    string         `json:"taskId" bson:"taskId"`
	ThreadID  string         `json:"threadId" bson:"threadId"`
	StepID    string         `json:"stepId" bson:"stepId"`
	StepName  string         `json:"stepName,omitempty" bson:"stepName,omitempty"`
	MCPServer string         `json:"mcpServer" bson:"mcpServer"`
	ToolName  string         `json:"toolName" bson:"toolName"`
	Targets   []string       `json:"targets" bson:"targets"` //target url
	Params    map[string]any `json:"params" bson:"params"`   //tool params
	Status    State          `json:"status" bson:"status"`
	Result    string         `json:"result" bson:"result"`
	Error     string         `json:"error" bson:"error"`
}

type Task struct {
	ThreadID    string   `json:"threadId,omitempty" bson:"threadId,omitempty"` //user chat thread id
	TaskName    string   `json:"taskName,omitempty" bson:"taskName,omitempty"`
	UserInput   string   `json:"userInput,omitempty" bson:"userInput,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Targets     []string `json:"targets,omitempty" bson:"targets,omitempty"` //target url
	History     []string `json:"history" bson:"history"`                     //user input history
	Status      State    `json:"status" bson:"status"`
	Steps       []Step   `json:"steps" bson:"steps"`
	Result      string   `json:"result" bson:"result"`
	Error       string   `json:"error" bson:"error"`
}

func InitTask() *protocols.XID {
	taskid := common.GenerateId()

	//info
	info := protocols.NewInfo(taskid, "mcptask")

	//metadata
	md := protocols.NewMetadata(
		protocols.OperationInit,
		"/protocols/mcptask",
		"application/json")

	//payload
	payload := &Task{
		Status: StateInit,
	}
	xid := protocols.NewXID(&info, &md, payload)
	return xid
}

// func CreateTaskEvent(taskid string, threadId string, step Step) *protocols.XID {
// 	task := &Task{
// 		ThreadID: threadId,
// 		Steps:    []Step{step},
// 	}
// 	info := protocols.NewInfo(taskid, "mcptask")
// 	metadata := protocols.NewMetadata(protocols.OperationUpdate, "/protocols/mcptask", "application/json")
// 	xid := protocols.NewXID(&info, &metadata, task)
// 	return xid
// }

// func CreateStep(taskid string, threadId string, step Step) *protocols.XID {
// 	task := &Task{
// 		ThreadID: threadId,
// 		Steps:    []Step{step},
// 	}
// 	info := protocols.NewInfo(taskid, "mcptask")
// 	metadata := protocols.NewMetadata(protocols.OperationUpdate, "/protocols/mcptask", "application/json")
// 	xid := protocols.NewXID(&info, &metadata, task)
// 	return xid
// }
