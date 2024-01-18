package models

import "approval-service/modules/entities/events"

type ApprovalUsecase interface {
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
	ReceiveRequest(id int, optional map[string]interface{}) ([]Approval, error)
	GetByID(uint) (*Approval, error)
	SentRequest(uint, *RequestSentRequest) (*Approval, error)
	
}

type ApprovalRepository interface {
	Create(*Approval) (*Approval, error)
	GetByID(uint) (*Approval, error)
	GetReceiveRequest(int) ([]Approval, error)
	UpdateStatus(uint, *UpdateStatusReq) (*Approval, error)
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
