package handlers

import (
	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
)

// import (
// 	"encoding/json"
// 	"fmt"
//
//

// 	"github.com/IBM/sarama"
// 	"github.com/gofiber/fiber/v2/log"
// )

type eventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) models.EventProducer {
	return &eventProducer{producer}
}

func (obj *eventProducer) Produce(event events.Event) error {
	topic := event.String()
	value, err := json.Marshal(event)
	if err != nil {
		logs.Error(fmt.Sprintln("can't convert data to json"), zap.Error(err))
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	p, o, err := obj.producer.SendMessage(&msg)
	if err != nil {
		logs.Error(fmt.Sprintf("can't produce event to %s topic", topic), zap.Error(err))
		return err
	}
	log.Info(fmt.Sprintf("sent to topic: %v, partition: %v, offset %v", topic, p, o))
	return nil
}
