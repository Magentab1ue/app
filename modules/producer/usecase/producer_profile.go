package usecase

// import (
// 	"time"
// 	"profile-service/modules/entities/events"
// 	"profile-service/modules/entities/models"
//
//
// )

// type producerUser struct {
// 	eventProducer models.EventProducer
// }

// func NewProducerServiceUsers(eventProducer models.EventProducer) models.ProducerProfile {
// 	return &producerUser{eventProducer}
// }

// // UserCreated implements ProducerUser.

// func (obj *producerUser) ProfileCreated(user *models.ProduceReq, timeStamp time.Time) error {

// 	return obj.eventProducer.Produce(events.ProfileCreatedEvent{})
// }

// // UserUpdated implements ProducerUser.
// func (obj *producerUser) ProfileUpdated(user *models.ProduceReq, timeStamp time.Time) error {
// 	return obj.eventProducer.Produce(events.ProfileUpdatedEvent{})
// }

// // UserDeleted implements ProducerUser.
// func (obj *producerUser) ProfileDeleted(user uint) error {
// 	return obj.eventProducer.Produce(events.ProfileDeletedEvent{})
// }
