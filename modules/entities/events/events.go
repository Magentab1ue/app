package events

import (
	"time"

	"github.com/google/uuid"
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
	ID           uint           `json:"id"`
	RequestID    uuid.UUID      `json:"request_id"`
	To           pq.Int64Array  `json:"to" gorm:"type:integer[]"`
	Approver     uint           `json:"approver"`
	Status       string         `json:"status"`
	Project      datatypes.JSON `json:"project" gorm:"type:jsonb"` // Assuming your database supports JSONB
	CreationDate time.Time      `json:"creation_date"`
	RequestUser  uint           `json:"request_user"`
	IsSignature  bool           `json:"is_signature"`
	Task         datatypes.JSON `json:"task" gorm:"type:jsonb"`
}

func (RequestCreatedEvent) String() string {
	return "tcchub-approval-approvalCreated"
}

type ApprovalUpdatedEvent struct {
	ID           uint           `json:"id"`
	RequestID    uuid.UUID      `json:"request_id"`
	To           pq.Int64Array  `json:"to" gorm:"type:integer[]"`
	Approver     uint           `json:"approver"`
	Status       string         `json:"status"`
	Project      datatypes.JSON `json:"project" gorm:"type:jsonb"` // Assuming your database supports JSONB
	CreationDate time.Time      `json:"creation_date"`
	RequestUser  uint           `json:"request_user"`
	IsSignature  bool           `json:"is_signature"`
	Task         datatypes.JSON `json:"task" gorm:"type:jsonb"`
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
