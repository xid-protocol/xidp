package securityevent

type AttackGraph struct {
	TaskID      string `json:"event_id" bson:"eventID"`
	ThreadID    string `json:"thread_id" bson:"threadID"`
	Trigger     string `json:"trigger" bson:"trigger"`
	EventName   string `json:"event_name" bson:"eventName"`
	Asset       string `json:"asset" bson:"asset"`
	Thread      string `json:"thread" bson:"thread"`
	EventTime   string `json:"event_time" bson:"eventTime"`
	EventSource string `json:"event_source" bson:"eventSource"`
}

type AttackStatus string

const (
	AttackStatusAttacking AttackStatus = "attacking"
	AttackStatusFailed    AttackStatus = "failed"
	AttackStatusCancelled AttackStatus = "cancelled"
	AttackStatusTimeout   AttackStatus = "timeout"
	AttackStatusSuccess   AttackStatus = "success"
)

type AttackEvent struct {
	AttackID   string       `json:"attack_id" bson:"attackID"`
	AttackName string       `json:"attack_name" bson:"attackName"`
	Targets    string       `json:"targets" bson:"targets"`
	Status     AttackStatus `json:"status" bson:"status"`
}
