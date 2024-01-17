package events

import (
	"encoding/json"
	"time"
)

type Event interface {
	String() string
}

var SubscribedTopics = []string{
	RequestCreatedEvent{}.String(),
}

type RequestCreatedEvent struct {
	To           string          `json:"to"`
	Approver     json.RawMessage `json:"approver"`
	Status       string          `json:"status"`
	Project      json.RawMessage `json:"project"`
	CreationDate time.Time       `json:"creation_date"`
	RequestUser  uint            `json:"request_user"`
	Task         json.RawMessage `json:"task"`
}

func (RequestCreatedEvent) String() string {
	return "RequestCreated"
}
