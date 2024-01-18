package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"

	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)

type eventHandler struct {
	consumer models.ConsumerUsecase
}

func NewEventHandler(consumer models.ConsumerUsecase) models.EventHandlerConsume {
	return &eventHandler{consumer}
}

func (obj *eventHandler) Handle(topic string, eventBytes []byte) {
	log.Info("consume topic:", topic)
	switch topic {
	case events.RequestCreatedEvent{}.String():
		event := events.RequestCreatedEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return
		}
		err = obj.consumer.RequestCreated(event)
		if err != nil {
			log.Error(err)
			return
		}

	}
}
