package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"

	"approval-service/logs"
	"approval-service/modules/entities/models"
)

type approvalService struct {
	approvalRepo models.ApprovalRepository
	Redis        *redis.Client
}

func NewApprovalService(
	approvalRepo models.ApprovalRepository,
	Redis *redis.Client,
) models.ApprovalUsecase {
	return &approvalService{
		approvalRepo,
		Redis,
	}
}

func (u approvalService) GetByID(id uint) (appprove *models.Approval,err error) {
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
	return nil, nil
}
