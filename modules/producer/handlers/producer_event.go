package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2/log"

	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)

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
		log.Error(err)
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	p, o, err := obj.producer.SendMessage(&msg)
	if err != nil {
		log.Error(err)
		return err
	}
	logs.Info(fmt.Sprintf("sent to topic: %v, partition: %v, offset %v, value = %s\n", topic, p, o, string(value)))
	return nil
}
