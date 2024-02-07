package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"

	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)

type eventHandler struct {
	consumer models.ConsumerUsecase
}

func NewEventHandler(consumer models.ConsumerUsecase) models.EventHandlerConsume {
	return &eventHandler{consumer}
}

func (obj *eventHandler) Handle(topic string, eventBytes []byte) error {
	log.Info("consume topic:", topic)
	switch topic {
	case events.UserProfile{}.TopicCreate():
		event := events.UserProfile{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to create profile")
		err = obj.consumer.CreateProfile(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("create profile successfully")

	case events.UserProfile{}.TopicUpdate():
		event := events.UserProfile{}

		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to update profile")
		err = obj.consumer.UpdateProfile(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("update profile successfully")
	case events.UserProfileDeleted{}.TopicDelete():
		event := events.UserProfileDeleted{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to delete profile")
		err = obj.consumer.DeleteProfile(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Delete profile successfully")

	case events.ProjectEvent{}.TopicCreate():
		event := events.ProjectEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to create project")
		err = obj.consumer.CreateProject(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Create project successfully")
	case events.ProjectEvent{}.TopicUpdate():
		event := events.ProjectEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to update project")
		err = obj.consumer.UpdateProject(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("update project successfully")
	case events.ProjectEventDeleted{}.TopicDelete():
		event := events.ProjectEventDeleted{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to Delete project")
		err = obj.consumer.DeleteProject(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Delete project successfully")
	case events.TaskEvent{}.TaskEventCreated():
		event := events.TaskEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to create task")
		err = obj.consumer.CreateTask(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Create task successfully")
	case events.TaskEvent{}.TaskEventUpdated():
		event := events.TaskEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to update Task")
		err = obj.consumer.UpdateTask(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("update Task successfully")
	case events.TaskEvent{}.TaskEventDeleted():
		event := events.TaskEvent{}
		err := json.Unmarshal(eventBytes, &event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("Attempting to delete Task")
		err = obj.consumer.DeleteTask(event)
		if err != nil {
			log.Error(err)
			return err
		}
		logs.Info("deleted Task successfully")
	}
	return nil
}

// func (obj *eventHandler) CheckMessage(msg *sarama.ConsumerMessage) error {
// 	err := obj.consumer.CheckOffsetMessage(msg.Topic, msg.Offset, msg.Partition)
// 	if err == nil {
// 		return nil
// 	}
// 	return err
// }
