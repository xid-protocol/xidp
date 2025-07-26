package mcptask

import (
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/protocols"
)

type TaskState string

const (
	TaskStateSubmitted     TaskState = "submitted"
	TaskStateInit          TaskState = "init"
	TaskStateWorking       TaskState = "working"
	TaskStateInputRequired TaskState = "input-required"
	TaskStateCompleted     TaskState = "completed"
	TaskStateCanceled      TaskState = "canceled"
	TaskStateFailed        TaskState = "failed"
	TaskStateRejected      TaskState = "rejected"
	TaskStateAuthRequired  TaskState = "auth-required"
	TaskStateUnknown       TaskState = "unknown"
)

type Task struct {
	ThreadID string         `json:"threadId" bson:"threadId"`
	Status   TaskState      `json:"status" bson:"status"`
	DataType string         `json:"dataType" bson:"dataType"`
	Data     map[string]any `json:"payload" bson:"payload"`
	Result   string         `json:"result" bson:"result"`
	Error    string         `json:"error" bson:"error"`
}

func NewTask(threadId string) *protocols.XID {
	taskid := common.GenerateId()

	info := protocols.NewInfo(taskid, "mcptask")

	var md protocols.Metadata
	//添加必填字段
	md = protocols.NewMetadata(
		protocols.OperationCreate,
		"/protocols/mcptask",
		"application/json")
	payload := &Task{
		ThreadID: threadId,
		Status:   TaskStateInit,
	}
	xid := protocols.NewXID(&info, &md, payload)
	return xid
}

func TaskEvent(taskid string, threadId string, dataType string, data map[string]any) *protocols.XID {
	task := &Task{
		ThreadID: threadId,
		Status:   TaskStateWorking,
		DataType: dataType,
		Data:     data,
	}
	info := protocols.NewInfo(taskid, "mcptask")
	metadata := protocols.NewMetadata(protocols.OperationUpdate, "/protocols/mcptask", "application/json")
	xid := protocols.NewXID(&info, &metadata, task)
	return xid
}
