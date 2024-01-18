package models

import (
	"time"

	"approval-service/modules/entities/events"
)

type ApprovalUsecase interface {
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approval, error)
	GetSendRequest(uint, map[string]interface{}) ([]Approval, error)
	DeleteApproval(uint) error
	GetByID(uint) (*Approval, error)
	SentRequest(uint, *RequestSentRequest) (*Approval, error)
}

type ApprovalRepository interface {
	Create(*Approval) (*Approval, error)
	GetByID(uint) (*Approval, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approval, error)
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
	GetSendRequest(userId uint, optional map[string]interface{}) ([]Approval, error)
	DeleteApproval(requestId uint) ([]Approval, error)
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
