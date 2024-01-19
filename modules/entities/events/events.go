package events

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Event interface {
	String() string
}

var SubscribedTopics = []string{
	RequestCreatedEvent{}.String(),
}

type RequestCreatedEvent struct {
	To           pq.Int64Array  `json:"to"`
	Project      datatypes.JSON `json:"project"`
	CreationDate time.Time      `json:"creation_date"`
	RequestUser  uint           `json:"request_user"`
	Task         datatypes.JSON `json:"task"`
}

func (RequestCreatedEvent) String() string {
	return "tcchub-approval-approvalCreated"
}

type ApprovalUpdatedEvent struct {
	RequestId uint           `json:"requestId"`
	Approver  uint           `json:"approver"`
	Status    string         `json:"status"`
	Task      datatypes.JSON `json:"task"`
}

func (ApprovalUpdatedEvent) String() string {
	return "tcchub-approval-approvalUpdated"
}

type ApprovalDeletedEvent struct {
	Task datatypes.JSON `json:"task"`
}

func (ApprovalDeletedEvent) String() string {
	return "tcchub-approval-approvalDeleted"
}
