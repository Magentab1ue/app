package usecase

import (
	"time"

	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)

type producerUser struct {
	eventProducer models.EventProducer
}

func NewProducerServiceApprovals(eventProducer models.EventProducer) models.ProducerApproval {
	return &producerUser{eventProducer}
}

// UserCreated implements ProducerUser.

func (obj *producerUser) RequestCreated(user *models.ProduceReq, timeStamp time.Time) error {

	return obj.eventProducer.Produce(events.RequestCreatedEvent{})
}

// UserUpdated implements ProducerUser.
func (obj *producerUser) ApprovalUpdated(user *models.ProduceReq, timeStamp time.Time) error {
	return obj.eventProducer.Produce(events.ApprovalUpdatedEvent{})
}

// UserDeleted implements ProducerUser.
func (obj *producerUser) ApprovalDeleted(user uint) error {
	return obj.eventProducer.Produce(events.ApprovalDeletedEvent{})
}
