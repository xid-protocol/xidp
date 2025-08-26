package task

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusInit      TaskStatus = "init"      // 初始化
	TaskStatusPending   TaskStatus = "pending"   // 待执行
	TaskStatusRunning   TaskStatus = "running"   // 执行中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 执行失败
	TaskStatusCancelled TaskStatus = "cancelled" // 已取消
	TaskStatusTimeout   TaskStatus = "timeout"   // 超时
)

// TaskPriority 任务优先级枚举
// type TaskPriority int

// const (
// 	TaskPriorityLow    TaskPriority = 3 // 低优先级
// 	TaskPriorityNormal TaskPriority = 2 // 普通优先级
// 	TaskPriorityHigh   TaskPriority = 1 // 高优先级
// 	TaskPriorityUrgent TaskPriority = 0 // 紧急优先级
// )

// // TaskSchedule 任务调度信息
// type TaskSchedule struct {
// 	StartTime      *time.Time `json:"startTime,omitempty" bson:"startTime,omitempty"`           // 计划开始时间
// 	EndTime        *time.Time `json:"endTime,omitempty" bson:"endTime,omitempty"`               // 计划结束时间
// 	Timeout        int64      `json:"timeout,omitempty" bson:"timeout,omitempty"`               // 超时时间（秒）
// 	RetryCount     int        `json:"retryCount,omitempty" bson:"retryCount,omitempty"`         // 重试次数
// 	RetryDelay     int64      `json:"retryDelay,omitempty" bson:"retryDelay,omitempty"`         // 重试延迟（秒）
// 	CronExpression string     `json:"cronExpression,omitempty" bson:"cronExpression,omitempty"` // Cron表达式（定时任务）
// }

// // TaskExecution 任务执行信息
// type TaskExecution struct {
// 	ExecutionID string      `json:"executionId,omitempty" bson:"executionId,omitempty"` // 执行ID
// 	Success     int8        `json:"success" bson:"success"`                             // 是否成功
// 	StartedAt   int64       `json:"startedAt,omitempty" bson:"startedAt,omitempty"`     // 实际开始时间
// 	CompletedAt int64       `json:"completedAt,omitempty" bson:"completedAt,omitempty"` // 实际完成时间
// 	Duration    int64       `json:"duration,omitempty" bson:"duration,omitempty"`       // 执行时长（秒）
// 	RetryCount  int         `json:"retryCount,omitempty" bson:"retryCount,omitempty"`   // 已重试次数
// 	WorkerID    string      `json:"workerId,omitempty" bson:"workerId,omitempty"`       // 执行者ID
// 	WorkerName  string      `json:"workerName,omitempty" bson:"workerName,omitempty"`   // 执行者名称
// 	Logs        []TaskLog   `json:"logs,omitempty" bson:"logs,omitempty"`               // 执行日志
// 	Result      *TaskResult `json:"result,omitempty" bson:"result,omitempty"`           // 执行结果
// 	Error       *TaskError  `json:"error,omitempty" bson:"error,omitempty"`             // 错误信息
// }

// // TaskLog 任务日志
// type TaskLog struct {
// 	ExecutionID string
// 	Type        string            `json:"type" bson:"type"`                     // 日志类型
// 	Timestamp   int64             `json:"timestamp" bson:"timestamp"`           // 日志时间
// 	Level       string            `json:"level" bson:"level"`                   // 日志级别
// 	Message     map[string]string `json:"message" bson:"message"`               // 日志消息
// 	Data        any               `json:"data,omitempty" bson:"data,omitempty"` // 附加数据
// }

// // TaskResult 任务结果
// type TaskResult struct {
// 	Success int8           `json:"success" bson:"success"`                     // 是否成功
// 	Data    any            `json:"data,omitempty" bson:"data,omitempty"`       // 结果数据
// 	Summary string         `json:"summary,omitempty" bson:"summary,omitempty"` // 结果摘要
// 	Metrics map[string]any `json:"metrics,omitempty" bson:"metrics,omitempty"` // 执行指标
// }

// // TaskError 任务错误
// type TaskError struct {
// 	Code    string `json:"code" bson:"code"`                           // 错误代码
// 	Message string `json:"message" bson:"message"`                     // 错误消息
// 	Details any    `json:"details,omitempty" bson:"details,omitempty"` // 错误详情
// }

// // TaskDependency 任务依赖
// type TaskDependency struct {
// 	TaskID    string `json:"taskId" bson:"taskId"`       // 依赖的任务ID
// 	Condition string `json:"condition" bson:"condition"` // 依赖条件（success/failed/completed）
// 	Required  bool   `json:"required" bson:"required"`   // 是否必需
// }

// // Task 任务主体结构
// type Task struct {
// 	Name         string           `json:"name" bson:"name"`                                     // 任务名称
// 	Description  string           `json:"description" bson:"description"`                       // 任务描述
// 	Type         string           `json:"type" bson:"type"`                                     // 任务类型
// 	Status       TaskStatus       `json:"status" bson:"status"`                                 // 任务状态
// 	Priority     TaskPriority     `json:"priority" bson:"priority"`                             // 任务优先级
// 	Schedule     *TaskSchedule    `json:"schedule,omitempty" bson:"schedule,omitempty"`         // 调度信息
// 	Executions   []*TaskExecution `json:"executions,omitempty" bson:"executions,omitempty"`     // 执行信息
// 	Result       *TaskResult      `json:"result,omitempty" bson:"result,omitempty"`             // 任务结果
// 	Dependencies []TaskDependency `json:"dependencies,omitempty" bson:"dependencies,omitempty"` // 任务依赖
// 	Parameters   map[string]any   `json:"parameters,omitempty" bson:"parameters,omitempty"`     // 任务参数
// 	Tags         []string         `json:"tags,omitempty" bson:"tags,omitempty"`                 // 任务标签
// 	CreatedBy    string           `json:"createdBy" bson:"createdBy"`                           // 创建者
// 	UpdatedBy    string           `json:"updatedBy" bson:"updatedBy"`                           // 更新者
// 	CreatedAt    int64            `json:"createdAt" bson:"createdAt"`                           // 创建时间
// 	UpdatedAt    int64            `json:"updatedAt" bson:"updatedAt"`                           // 更新时间
// }

type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"   // 待执行
	StepStatusRunning   StepStatus = "running"   // 执行中
	StepStatusCompleted StepStatus = "completed" // 已完成
	StepStatusFailed    StepStatus = "failed"    // 执行失败
	StepStatusCancelled StepStatus = "cancelled" // 已取消
	StepStatusTimeout   StepStatus = "timeout"   // 超时
)

type TaskStep struct {
	StepID     string         `json:"stepId" bson:"stepId"`
	StepName   string         `json:"stepName,omitempty" bson:"stepName,omitempty"`
	WorkerID   string         `json:"workerId" bson:"workerId"`
	WorkerName string         `json:"workerName" bson:"workerName"`
	Params     map[string]any `json:"params" bson:"params"` //step params
	Status     StepStatus     `json:"status" bson:"status"`
	Result     map[string]any `json:"result" bson:"result"`
	Error      string         `json:"error" bson:"error"`
}

type StepEvent struct {
	StepID     string         `json:"stepId" bson:"stepId"`
	StepName   string         `json:"stepName,omitempty" bson:"stepName,omitempty"`
	WorkerID   string         `json:"worker" bson:"worker"`
	WorkerName string         `json:"workerName" bson:"workerName"`
	DataType   string         `json:"dataType" bson:"dataType"`
	Data       map[string]any `json:"data" bson:"data"`
}

type Task struct {
	TaskID      string     `json:"taskId" bson:"taskId"`
	Name        string     `json:"name,omitempty" bson:"name,omitempty"`
	TaskType    string     `json:"taskType,omitempty" bson:"taskType,omitempty"`
	UserInput   string     `json:"userInput,omitempty" bson:"userInput,omitempty"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
	Targets     []string   `json:"targets,omitempty" bson:"targets,omitempty"` //target url
	History     []string   `json:"history" bson:"history"`                     //user input history
	Status      TaskStatus `json:"status" bson:"status"`
	Steps       []TaskStep `json:"steps" bson:"steps"`
	Result      string     `json:"result" bson:"result"`
	Error       string     `json:"error" bson:"error"`
	CreatedAt   int64      `json:"createdAt" bson:"createdAt"`
	UpdatedAt   int64      `json:"updatedAt" bson:"updatedAt"`
}

// NewTask
// func CreateTask(task *Task) (*protocols.XID, error) {
// 	task.CreatedAt = common.GetTimestamp()
// 	task.UpdatedAt = common.GetTimestamp()
// 	xid := protocols.GenerateXid(task.Name)
// 	//check if XID already exists
// 	xidRepository := xdb.NewXidInfoRepository()
// 	XID, err := xidRepository.FindOneByXidAndPath(context.Background(), xid, "/protocols/task")
// 	//if not nil, return error
// 	if XID != nil {
// 		return nil, errors.New("xid already exists")
// 	}
// 	if err != nil {
// 		logx.Errorf("CreateTask Error: %v", err)
// 		return nil, err
// 	}

// 	info := protocols.NewInfo(task.Name, "taskName")
// 	info.Extra = map[string]any{
// 		"description": task.Description,
// 	}
// 	metadata := protocols.NewMetadata(protocols.OperationInit, "/protocols/task", "application/json")

// 	//create task
// 	payload := &Task{
// 		TaskType: task.TaskType,
// 		Status:   TaskStatusInit,
// 	}

// 	XID = protocols.NewXID(&info, &metadata, payload)

// 	return XID, nil
// }
