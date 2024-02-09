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

func (consumer *handlerConsumeGroup) Setup(session sarama.ConsumerGroupSession) error {
	logs.Info("kafka Setup")
	return nil
}

func (consumer *handlerConsumeGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	logs.Info("kafka Cleanup")
	return nil
}

func (consumer *handlerConsumeGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logs.Info(fmt.Sprintf("Subscribed topics: %s", claim.Topic()))
	for msg := range claim.Messages() {
		logs.Info(fmt.Sprintf("Consumed message from topic %s with partition: %d and offset: %d", msg.Topic, msg.Partition, msg.Offset))
		err := consumer.eventHandler.Handle(msg.Topic, msg.Value)
		session.MarkMessage(msg, "")
		if err != nil {
			logs.Error(fmt.Sprintf("Error handling message from topic %s with partition: %d and offset: %d Error : %s", msg.Topic, msg.Partition, msg.Offset, err.Error()))
			return err
		}

	}
	return nil
}
