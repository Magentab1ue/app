package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
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

func (u approvalService) UpdateStatus(id uint, req *models.UpdateStatusReq) (*models.Approvals, error) {
	statusList := []string{models.Approve, models.Pending, models.Reject}

	hasStatus := stringInSlice(req.Status, statusList)
	if !hasStatus {
		return nil, errors.New("format status incorrect, status should have approved , pending or reject")
	}
	approvalRes, err := u.approvalRepo.UpdateStatus(id, req)
	if err != nil {
		return nil, err
	}
	event := events.ApprovalUpdatedEvent{
		ID:           approvalRes.ID,
		RequestID:    approvalRes.RequestID,
		To:           approvalRes.To,
		Approver:     approvalRes.Approver,
		Status:       approvalRes.Status,
		Project:      approvalRes.Project,
		CreationDate: approvalRes.CreationDate,
		RequestUser:  approvalRes.RequestUser,
		IsSignature:  approvalRes.IsSignature,
		Task:         approvalRes.Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) GetReceiveRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approvals, err error) {
	keyRedis := fmt.Sprintf("GetReceiveRequest:%d,optionnals:%v", id, optional)

	approvalResJson, err := u.Redis.Get(context.Background(), keyRedis).Result()
	if err == nil {
		if json.Unmarshal([]byte(approvalResJson), &approvalRes) == nil {
			logs.Info("Read data from: redis")
			return approvalRes, nil
		}

	}
	if err != nil {
		logs.Warn("Redis error", zap.Error(err))
	}
	approvalRes, err = u.approvalRepo.GetReceiveRequest(id, optional)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(approvalRes); err == nil {
		u.Redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}
	return approvalRes, nil
}

func (u approvalService) GetSendRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approvals, err error) {
	keyRedis := fmt.Sprintf("GetSendRequest:%d,optionnals:%v", id, optional)
	approvalResJson, err := u.Redis.Get(context.Background(), keyRedis).Result()
	if err == nil {
		if json.Unmarshal([]byte(approvalResJson), &approvalRes) == nil {
			logs.Info("Read data from: redis")
			return approvalRes, nil
		}

	}
	if err != nil {
		logs.Warn("Redis error", zap.Error(err))
	}
	approvalRes, err = u.approvalRepo.GetSendRequest(id, optional)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(approvalRes); err == nil {
		u.Redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}

	return approvalRes, nil
}

func (u approvalService) DeleteApproval(id uint) error {
	approval, err := u.approvalRepo.DeleteApproval(id)
	if err != nil {
		return err
	}

	event := events.ApprovalDeletedEvent{
		Task: approval.Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return err
	}
	return nil
}

func (u approvalService) GetByID(id uint) (appprove *models.Approvals, err error) {
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

func (u approvalService) SentRequest(id uint, req *models.RequestSentRequest) (*models.Approvals, error) {

	request, err := u.approvalRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	event := events.RequestCreatedEvent{
		ID:           request.ID,
		RequestID:    request.RequestID,
		To:           request.To,
		Status:       request.Status,
		Project:      request.Project,
		CreationDate: request.CreationDate,
		RequestUser:  request.RequestUser,
		Task:         request.Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}

	res, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:    request.RequestID,
		Status:       "pending",
		Project:      request.Project,
		To:           req.To,
		CreationDate: req.CreationDate,
		RequestUser:  req.RequestUser,
		Task:         request.Task,
		IsSignature:  req.IsSignature,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u approvalService) GetAll(optional map[string]interface{}) (appprove []models.Approvals, err error) {

	keyRedis := fmt.Sprintf("GetApprovals:optionnals:%v", optional)
	//redis get
	if approvalJson, err := u.Redis.Get(context.Background(), keyRedis).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	log.Debug("Read data from database")
	approvalRes, err := u.approvalRepo.GetAll(optional)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(approvalRes); err == nil {
		u.Redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}
	return approvalRes, nil
}

func (u approvalService) GetByUserID(id uint, optional map[string]interface{}) (appprove []models.Approvals, err error) {
	keyRedis := fmt.Sprintf("GetApprovalByUserId:%v:optionnals:%v", id, optional)
	//redis get
	if approvalJson, err := u.Redis.Get(context.Background(), keyRedis).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	log.Debug("Read data from database")
	approvalRes, err := u.approvalRepo.GetByUserID(id, optional)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(approvalRes); err == nil {
		u.Redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}
	return approvalRes, nil
}

func stringInSlice(str string, list []string) bool {
	for _, val := range list {
		if val == str {
			return true
		}
	}
	return false
}

func (u approvalService) CreateRequest(req *models.RequestReq) (*models.Approvals, error) {
	newRequest, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:    uuid.New(),
		To:           req.To,
		Status:       "pending",
		Project:      req.Project,
		CreationDate: req.CreationDate,
		RequestUser:  req.RequestUser,
		Task:         req.Task,
	})
	if err != nil {
		logs.Error(err)
		return nil, errors.New("can't create request")
	}

	event := events.RequestCreatedEvent{
		ID:           newRequest.ID,
		RequestID:    newRequest.RequestID,
		To:           newRequest.To,
		Status:       newRequest.Status,
		Project:      newRequest.Project,
		CreationDate: newRequest.CreationDate,
		RequestUser:  newRequest.RequestUser,
		Task:         newRequest.Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}

	log.Info("create request Successfuly")
	return newRequest, nil
}
