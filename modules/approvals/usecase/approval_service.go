package usecase

import (
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"

	"approval-service/logs"
)

type approvalService struct {
	approvalRepo models.ApprovalRepository
	produce      models.EventProducer
	Redis        *redis.Client
}

func NewApprovalService(
	approvalRepo models.ApprovalRepository,
	produce models.EventProducer,
	Redis *redis.Client,
) models.ApprovalUsecase {
	return &approvalService{approvalRepo, produce, Redis}
}

func (u approvalService) UpdateStatus(id uint, req *models.UpdateStatusReq) (*models.Approval, error) {
	approvalRes, err := u.approvalRepo.UpdateStatus(id, req)
	if err != nil {
		return nil, err
	}
	event := events.ApprovalUpdatedEvent{
		RequestId: id,
		Approver:  req.Approver,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) ReceiveRequest(id uint, optional map[string]interface{}) ([]models.Approval, error) {
	approvalRes, err := u.approvalRepo.GetReceiveRequest(id,optional)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) SendRequest(id uint, optional map[string]interface{}) ([]models.Approval, error) {
	approvalRes, err := u.approvalRepo.GetSendRequest(id, optional)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) DeleteApproval(id uint) error {
	approval, err := u.approvalRepo.DeleteApproval(id)
	if err != nil {
		return err
	}

	event := events.ApprovalDeletedEvent{
		Task: approval[0].Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return err
	}
	return nil
}

func (u approvalService) GetByID(id uint) (appprove *models.Approval, err error) {
	key := fmt.Sprintf("service:GetApprovalByID%v", id)
	//redis get
	if approvalJson, err := u.Redis.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	// Data not found in cache, fetch from the database
	log.Debug("Read data from database")
	approvalDB, err := u.approvalRepo.GetByID(id)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("couldn't get profile data")
	}

	//redis set
	if data, err := json.Marshal(approvalDB); err == nil {
		u.Redis.Set(context.Background(), key, string(data), time.Minute*1)
	}

	return approvalDB, nil
}

func (u approvalService) SentRequest(id uint, req *models.RequestSentRequest) (*models.Approval, error) {
	
	request,err := u.approvalRepo.GetByID(id)
	if err != nil {
		return nil,err
	}
	
	res,err :=u.approvalRepo.Create(&models.Approval{
		RequestID: request.RequestID,
		Status: "pending",
		Project: request.Project,
		To: req.To,
		CreationDate: req.CreationDate,
		RequestUser: req.RequestUser,
		Task: request.Task,
	})
	if err != nil{
		return nil,err
	}
	return res, nil
}


// func (u approvalService) GetAll(optional map[string]interface{}) ([]models.Approval, error) {
// 	approvalRes, err := u.approvalRepo.GetReceiveRequest(id,optional)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return approvalRes, nil
// }
