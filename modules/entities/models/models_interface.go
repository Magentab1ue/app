package models

import "approval-service/modules/entities/events"

type ApprovalUsecase interface {
}

type ApprovalRepository interface {
	Create(*Approval) (*Approval, error)
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
