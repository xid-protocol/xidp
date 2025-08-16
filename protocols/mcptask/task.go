package mcptask

import (
	"context"
	"errors"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
)

// TaskStatus 任务状态枚举
type Status string

const (
	StatusInit      Status = "init"
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
	StatusTimeout   Status = "timeout"
)

type StepEvent struct {
	ThreadID   string         `json:"threadID" bson:"threadID"`
	StepID     string         `json:"stepID" bson:"stepID"`
	StepName   string         `json:"stepName,omitempty" bson:"stepName,omitempty"`
	WorkerID   string         `json:"worker" bson:"worker"`
	WorkerName string         `json:"workerName" bson:"workerName"`
	DataType   string         `json:"dataType" bson:"dataType"`
	Data       map[string]any `json:"data" bson:"data"`
}

type Step struct {
	ThreadID   string         `json:"threadID" bson:"threadID"`
	StepID     string         `json:"stepID" bson:"stepID"`
	StepName   string         `json:"stepName,omitempty" bson:"stepName,omitempty"`
	WorkerID   string         `json:"workerID" bson:"workerID"`
	WorkerName string         `json:"workerName" bson:"workerName"`
	Params     map[string]any `json:"params" bson:"params"` //step params
	Status     Status         `json:"status" bson:"status"`
	Result     map[string]any `json:"result" bson:"result"`
	Error      string         `json:"error" bson:"error"`
}

// One thread per page, each thread contains multiple steps
type Thread struct {
	ThreadID   string `json:"threadID" bson:"threadID"`
	ThreadName string `json:"threadName" bson:"threadName"`
	Status     Status `json:"status" bson:"status"`
	Steps      []Step `json:"steps" bson:"steps"`
}

type Task struct {
	TaskID      string   `json:"taskID" bson:"taskID"`
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	TaskType    string   `json:"taskType,omitempty" bson:"taskType,omitempty"`
	UserInput   string   `json:"userInput,omitempty" bson:"userInput,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Targets     []string `json:"targets,omitempty" bson:"targets,omitempty"` //target url
	History     []string `json:"history" bson:"history"`                     //user input history
	Status      Status   `json:"status" bson:"status"`
	Threads     []Thread `json:"threads" bson:"threads"`
	Result      string   `json:"result" bson:"result"`
	Error       string   `json:"error" bson:"error"`
	CreatedAt   int64    `json:"createdAt" bson:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt" bson:"updatedAt"`
}

// NewTask
func CreateTask(task *Task) (*protocols.XID, error) {
	task.CreatedAt = common.GetTimestamp()
	task.UpdatedAt = common.GetTimestamp()
	xid := protocols.GenerateXid(task.Name)
	//check if XID already exists
	xidRepository := db.NewXidInfoRepository()
	XID, err := xidRepository.FindOneByXidAndPath(context.Background(), xid, "/protocols/task")
	//if not nil, return error
	if XID != nil {
		return nil, errors.New("xid already exists")
	}
	if err != nil {
		logx.Errorf("CreateTask Error: %v", err)
		return nil, err
	}

	info := protocols.NewInfo(task.Name, "taskName")
	info.Extra = map[string]any{
		"description": task.Description,
	}
	metadata := protocols.NewMetadata(protocols.OperationInit, "/protocols/task", "application/json")

	//create task
	payload := &Task{
		TaskType: task.TaskType,
		Status:   StatusInit,
	}

	XID = protocols.NewXID(&info, &metadata, payload)

	return XID, nil
}
