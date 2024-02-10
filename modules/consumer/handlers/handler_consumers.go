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
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logs.Info("Kafka message channel was closed")
				return nil
			}
			logs.Info(fmt.Sprintf("Consumed message  claimed from value = %s, timestamp = %v, topic = %s, offset = %d", string(message.Value), message.Timestamp, message.Topic, message.Offset))
			err := consumer.eventHandler.Handle(message.Topic, message.Value)
			if err != nil {
				logs.Error(fmt.Sprintf("Error handling message from topic %s with partition: %d and offset: %d Error : %s", message.Topic, message.Partition, message.Offset, err.Error()))
				return nil
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
