package models

import (
	"time"

	"approval-service/modules/entities/events"
)

type ApprovalUsecase interface {
	UpdateStatus(uint, *UpdateStatusReq) (*Approvals, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approvals, error)
	GetSendRequest(uint, map[string]interface{}) ([]Approvals, error)
	DeleteApproval(uint) error
	GetByID(uint) (*Approvals, error)
	SentRequest(uint, *RequestSentRequest) (*Approvals, error)
	GetAll(map[string]interface{}) ([]Approvals, error)
	GetByUserID(uint,map[string]interface{}) ([]Approvals, error)
}

type ApprovalRepository interface {
	Create(*Approvals) (*Approvals, error)
	GetByID(uint) (*Approvals, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approvals, error)
	UpdateStatus(uint, *UpdateStatusReq) (*Approvals, error)
	GetSendRequest(userId uint, optional map[string]interface{}) ([]Approvals, error)
	DeleteApproval(requestId uint) ([]Approvals, error)
	GetAll(map[string]interface{}) ([]Approvals, error)
	GetByUserID(uint, map[string]interface{}) ([]Approvals, error)
}

type ProducerApproval interface {
	RequestCreated(*ProduceReq, time.Time) error
	ApprovalUpdated(*ProduceReq, time.Time) error
	ApprovalDeleted(uint) error
}

// consumer
type EventHandlerConsume interface {
	Handle(toppic string, eventByte []byte)
}

type ConsumerUsecase interface {
	RequestCreated(event events.RequestCreatedEvent) error
}

type EventProducer interface {
	Produce(events.Event) error
}
