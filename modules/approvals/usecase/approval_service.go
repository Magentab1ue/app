package usecase

import (
	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type approvalService struct {
	approvalRepo models.ApprovalRepository
	produce      models.EventProducer
	redis        *redis.Client
}

func NewApprovalService(
	approvalRepo models.ApprovalRepository,
	produce models.EventProducer,
	redis *redis.Client,
) models.ApprovalUsecase {
	return &approvalService{approvalRepo, produce, redis}
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

func (u approvalService) GetReceiveRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approval, err error) {
	keyRedis := fmt.Sprintf("GetReceiveRequest:%d,optionnals:%v", id, optional)
	approvalResJson, err := u.redis.Get(context.Background(), keyRedis).Result()

	if json.Unmarshal([]byte(approvalResJson), &approvalRes); err == nil {
		logs.Debug("Read data from: redis")
		return approvalRes, nil
	}
	approvalRes, err = u.approvalRepo.GetReceiveRequest(id, optional)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(approvalRes); err == nil {
		u.redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}
	return approvalRes, nil
}

func (u approvalService) GetSendRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approval, err error) {
	keyRedis := fmt.Sprintf("GetSendRequest:%d,optionnals:%v", id, optional)
	approvalResJson, err := u.redis.Get(context.Background(), keyRedis).Result()
	if json.Unmarshal([]byte(approvalResJson), &approvalRes); err == nil {
		logs.Debug("Read data from: redis")
		return approvalRes, nil
	}

	approvalRes, err = u.approvalRepo.GetSendRequest(id, optional)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(approvalRes); err == nil {
		u.redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
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
		Task: approval[0].Task,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return err
	}
	return nil
}
