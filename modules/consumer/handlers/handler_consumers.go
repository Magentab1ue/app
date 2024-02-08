package handlers

import (
	"approval-service/logs"
	"approval-service/modules/entities/models"
	"fmt"

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
		logs.Info(fmt.Sprintf("Consumed message from topic %s with partition: %d and offset: %d", msg.Topic, msg.Partition, msg.Offset))
		err := consumer.eventHandler.Handle(msg.Topic, msg.Value)
		if err != nil {
			logs.Error(fmt.Sprintf("Error handling message from topic %s with partition: %d and offset: %d Error : %s", msg.Topic, msg.Partition, msg.Offset, err.Error()))
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}