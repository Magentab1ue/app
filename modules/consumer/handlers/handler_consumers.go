package handlers

import (
	"approval-service/modules/entities/models"

	"github.com/IBM/sarama"
)

type handlerConsumeGroup struct {
	eventHandler models.EventHandlerConsume
}

// create new consumerHandler, consumerHandler implement sarama.ConsumerGroupHandler
func NewHandlerConsumeGroup(eventHandler models.EventHandlerConsume) sarama.ConsumerGroupHandler {
	return &handlerConsumeGroup{eventHandler}
}

func (consumer *handlerConsumeGroup) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *handlerConsumeGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *handlerConsumeGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		go func() {
			consumer.eventHandler.Handle(msg.Topic, msg.Value)
		}()
		session.MarkMessage(msg, "")
	}

	return nil
}
