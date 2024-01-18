package models

import "approval-service/modules/entities/events"

type ApprovalUsecase interface {
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approval, error)
	GetSendRequest(uint, map[string]interface{}) ([]Approval, error)
	DeleteApproval(id uint) error
}

type ApprovalRepository interface {
	Create(*Approval) (*Approval, error)
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
	GetReceiveRequest(userId uint, optional map[string]interface{}) ([]Approval, error)
	GetSendRequest(userId uint, optional map[string]interface{}) ([]Approval, error)
	DeleteApproval(requestId uint) ([]Approval, error)
}

type ProducerProfile interface {
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
