package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	if req.ToRole != "Approver" && req.ToRole != "HR" {
		return nil, errors.New("to_role field should be Approver or HR")
	}
	request, err := u.approvalRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	project := new(models.Project)
	err = json.Unmarshal(request.Project, project)
	if err != nil {
		return nil, errors.New("this approval have project formatted incorrectly")
	}

	projectCheck := reflect.TypeOf(project)
	if projectCheck.Kind() == reflect.Ptr {
		projectCheck = projectCheck.Elem()
	}
	_, ok := projectCheck.FieldByName("Approvers")
	if !ok {
		return nil, errors.New("this project dont have a approvers")
	}

	var to pq.Int64Array
	for _, approver := range project.Approvers {
		for _, role := range approver.Role {
			if role == req.ToRole {
				to = append(to, int64(approver.ID))
				break
			}
		}
	}
	res, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:       request.RequestID,
		Status:          "pending",
		Project:         request.Project,
		To:              to,
		CreationDate:    time.Now(),
		RequestUser:     req.RequestUser,
		Task:            request.Task,
		IsSignature:     req.IsSignature,
		Name:            req.Name,
		Detail:          req.Detail,
		NameRequestUser: req.NameRequestUser,
		ToRole:          req.ToRole,
	})
	if err != nil {
		return nil, err
	}
	event := events.RequestCreatedEvent{
		ID:              res.ID,
		RequestID:       res.RequestID,
		To:              res.To,
		Status:          res.Status,
		Project:         res.Project,
		CreationDate:    res.CreationDate,
		RequestUser:     res.RequestUser,
		Task:            res.Task,
		Name:            res.Name,
		Detail:          res.Detail,
		NameRequestUser: res.NameRequestUser,
		ToRole:          res.ToRole,
	}
	err = u.produce.Produce(event)
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

func (u approvalService) CreateRequest(req *models.CreateReq) (*models.Approvals, error) {

	project := new(models.Project)
	err := json.Unmarshal(req.Project, project)
	if err != nil {
		return nil, errors.New("project is formatted incorrectly")
	}
	projectCheck := reflect.TypeOf(project)
	if projectCheck.Kind() == reflect.Ptr {
		projectCheck = projectCheck.Elem()
	}
	_, okApprover := projectCheck.FieldByName("Approvers")
	_, okTeamlead := projectCheck.FieldByName("TeamLeads")
	if !okTeamlead || !okApprover {
		return nil, errors.New("project should have a teamleads or approvers")
	}
	var to pq.Int64Array
	for _, teamLead := range project.TeamLeads {
		to = append(to, int64(teamLead.ID))
	}

	newRequest, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:       uuid.New(),
		To:              to,
		Status:          "pending",
		Project:         req.Project,
		CreationDate:    time.Now(),
		RequestUser:     req.RequestUser,
		Task:            req.Task,
		Name:            req.Name,
		Detail:          req.Detail,
		NameRequestUser: req.NameRequestUser,
		ToRole:          "teamlead ",
	})
	if err != nil {
		logs.Error(err)
		return nil, errors.New("can't create request")
	}

	event := events.RequestCreatedEvent{
		ID:              newRequest.ID,
		RequestID:       newRequest.RequestID,
		To:              newRequest.To,
		Status:          newRequest.Status,
		Project:         newRequest.Project,
		CreationDate:    newRequest.CreationDate,
		RequestUser:     newRequest.RequestUser,
		Task:            newRequest.Task,
		Name:            newRequest.Name,
		Detail:          newRequest.Detail,
		NameRequestUser: newRequest.NameRequestUser,
		ToRole:          newRequest.ToRole,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}

	log.Info("create request Successfuly")
	return newRequest, nil
}

func (u approvalService) GetByRequestID(id uuid.UUID) (appprove []models.Approvals, err error) {
	key := fmt.Sprintf("service:GetApprovalByRequestID%v", id)
	//redis get
	if approvalJson, err := u.Redis.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	// Data not found in cache, fetch from the database
	log.Debug("Read data from database")
	approvalDB, err := u.approvalRepo.GetByRequestID(id)
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
