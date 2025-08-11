package securityevent

type AttackGraph struct {
	TaskID      string `json:"event_id" bson:"eventID"`
	Trigger     string `json:"trigger" bson:"trigger"`
	EventName   string `json:"event_name" bson:"eventName"`
	Asset       string `json:"asset" bson:"asset"`
	Thread      string `json:"thread" bson:"thread"`
	EventTime   string `json:"event_time" bson:"eventTime"`
	EventSource string `json:"event_source" bson:"eventSource"`
}
