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
	//validation
	approvalCheck, err := u.approvalRepo.GetByID(id)
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval for update with request ID %d: %v", id, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't finding approval with request ID %d", id)
	}
	if approvalCheck.Status == req.Status {
		return nil, errors.New("this approval status already exists")
	}
	// has approver in send to
	checkPermission := intInSlice(int64(req.Approver), approvalCheck.To)
	if !checkPermission {
		return nil, fmt.Errorf("this user id %d dont have permission to update status this approval", int(req.Approver))
	}
	logs.Info("Attempting to Update approval")
	approvalCheck.Approver = req.Approver
	approvalCheck.IsSignature = req.IsSignature
	approvalCheck.Status = req.Status
	approvalRes, err := u.approvalRepo.Update(approvalCheck)
	if err != nil {
		return nil, err
	}
	event := events.ApprovalUpdatedEvent{
		ID:           approvalRes.ID,
		RequestID:    approvalRes.RequestID,
		To:           approvalRes.To,
		Status:       approvalRes.Status,
		ProjectID:    approvalRes.ProjectID,
		CreationDate: approvalRes.CreationDate,
		SenderID:     approvalRes.SenderID,
		Task:         approvalRes.Task,
		Name:         approvalRes.Name,
		Detail:       approvalRes.Detail,
		ToRole:       approvalRes.ToRole,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) GetReceiveRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approvals, err error) {
	keyRedis := fmt.Sprintf("GetReceiveRequest:%d,optionnals:%v", id, optional)
	logs.Info("Attempting to Get data from redis")
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
	logs.Info("Attempting to Get data from database")
	approvalRes, err = u.approvalRepo.GetReceiveRequest(id, optional)
	if err != nil {
		return nil, err
	}
	logs.Info(fmt.Sprintf("Attempting to Set data to redis with delay time %d second", (time.Minute*1)/100000000))
	if data, err := json.Marshal(approvalRes); err == nil {
		u.Redis.Set(context.Background(), keyRedis, string(data), time.Minute*1)
	} else {
		logs.Warn("Can't set data to redis", zap.Error(err))
	}
	return approvalRes, nil
}

func (u approvalService) GetSendRequest(id uint, optional map[string]interface{}) (approvalRes []models.Approvals, err error) {
	keyRedis := fmt.Sprintf("GetSendRequest:%d,optionnals:%v", id, optional)
	logs.Info("Attempting to Get data from redis")
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
	logs.Info("Attempting to Get data from database")
	approvalRes, err = u.approvalRepo.GetSendRequest(id, optional)
	if err != nil {
		return nil, err
	}
	logs.Info(fmt.Sprintf("Attempting to Set data to redis with delay time %d", (time.Minute*1)/100000000))
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
	logs.Info("Attempting to Delete data from database")
	event := events.ApprovalDeletedEvent{
		Task: approval.Task,
	}
	logs.Info("Attempting to produce event to kafka")
	err = u.produce.Produce(event)
	if err != nil {
		return err
	}
	return nil
}

func (u approvalService) GetByID(id uint) (appprove *models.Approvals, err error) {
	key := fmt.Sprintf("service:GetApprovalByID%v", id)
	//redis get
	logs.Info("Attempting to Get data from redis")

	if approvalJson, err := u.Redis.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	// Data not found in cache, fetch from the database
	logs.Info("Attempting to Get data from database")
	approvalDB, err := u.approvalRepo.GetByID(id)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("couldn't get profile data")
	}
	logs.Info(fmt.Sprintf("Attempting to Set data to redis with delay time %d", (time.Minute*1)/100000000))

	//redis set
	if data, err := json.Marshal(approvalDB); err == nil {
		u.Redis.Set(context.Background(), key, string(data), time.Minute*1)
	}

	return approvalDB, nil
}

func (u approvalService) SentRequest(id uint, req *models.RequestSentRequest) (*models.Approvals, error) {
	logs.Info("validation request")

	if req.ToRole != "Approver" && req.ToRole != "HR" {
		return nil, errors.New("to_role field should be Approver or HR")
	}
	// //validation
	request, err := u.approvalRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if request.Status != models.Approve {
		return nil, errors.New("this approvals has not yet been approved")
	}
	requestLast, err := u.approvalRepo.GetByRequestIDLast(request.RequestID)
	if err != nil {
		return nil, err
	}
	if requestLast.ToRole == "Approver" && requestLast.Status == models.Approve {
		return nil, errors.New("this approval has been approved by approver")
	}
	if requestLast.Status == models.Pending {
		return nil, errors.New("the approval has already been sent")
	}
	checkPermission := intInSlice(int64(req.SenderID), request.To)
	if !checkPermission {
		return nil, fmt.Errorf("this user id %d dont have permission to send this approval", int(req.SenderID))
	}

	project, err := u.approvalRepo.GetProjectById(request.ProjectID)
	if err != nil {
		return nil, err
	}
	projectJson := new(models.ProjectJson)
	err = json.Unmarshal(project.Project, projectJson)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("can't create request")
	}

	var to pq.Int64Array
	for _, approver := range projectJson.Approvers {
		for _, role := range approver.Role {
			if role == req.ToRole {
				to = append(to, int64(approver.ID))
				break
			}
		}
	}

	logs.Info(fmt.Sprintf("Attempting to Create request approval requestId %v", request.RequestID))
	res, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:    request.RequestID,
		Status:       "pending",
		ProjectID:    request.ProjectID,
		To:           to,
		CreationDate: time.Now(),
		SenderID:     req.SenderID,
		Task:         request.Task,
		IsSignature:  req.IsSignature,
		Name:         req.Name,
		Detail:       req.Detail,
		ToRole:       req.ToRole,
	})
	if err != nil {
		return nil, err
	}
	event := events.RequestCreatedEvent{
		ID:           res.ID,
		RequestID:    res.RequestID,
		To:           res.To,
		Status:       res.Status,
		ProjectID:    res.ProjectID,
		CreationDate: res.CreationDate,
		SenderID:     res.SenderID,
		Task:         res.Task,
		Name:         res.Name,
		Detail:       res.Detail,
		ToRole:       res.ToRole,
	}
	logs.Info("Attempting to produce event to kafka")
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

	// project := new(models.Project)
	// err := json.Unmarshal(req.ProjectId, project)
	// if err != nil {
	// 	return nil, errors.New("project is formatted incorrectly")
	// }
	// projectCheck := reflect.TypeOf(project)
	// if projectCheck.Kind() == reflect.Ptr {
	// 	projectCheck = projectCheck.Elem()
	// }
	// _, okApprover := projectCheck.FieldByName("Approvers")
	// _, okTeamlead := projectCheck.FieldByName("TeamLeads")
	// _, okMembers := projectCheck.FieldByName("Members")
	// if !okTeamlead || !okApprover || !okMembers {
	// 	return nil, errors.New("project should have a teamlaeds and approvers and members")
	// }
	project, err := u.approvalRepo.GetProjectById(req.ProjectId)
	if err != nil {
		return nil, err
	}
	var to pq.Int64Array
	projectJson := new(models.ProjectJson)
	err = json.Unmarshal(project.Project, projectJson)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("can't create request")
	}
	for _, teamLead := range projectJson.TeamLeads {
		to = append(to, int64(teamLead.ID))
	}
	newRequest, err := u.approvalRepo.Create(&models.Approvals{
		RequestID:    uuid.New(),
		To:           to,
		Status:       models.Pending,
		ProjectID:    req.ProjectId,
		CreationDate: time.Now(),
		SenderID:     req.SenderID,
		Task:         req.Task,
		Name:         req.Name,
		Detail:       req.Detail,
		ToRole:       "teamlead",
	})
	if err != nil {
		logs.Error(err)
		return nil, errors.New("can't create request")
	}
	logs.Info("Attempting to produce event to kafka")

	event := events.RequestCreatedEvent{
		ID:           newRequest.ID,
		RequestID:    newRequest.RequestID,
		To:           newRequest.To,
		Status:       newRequest.Status,
		ProjectID:    newRequest.ProjectID,
		CreationDate: newRequest.CreationDate,
		SenderID:     newRequest.SenderID,
		Task:         newRequest.Task,
		Name:         newRequest.Name,
		Detail:       newRequest.Detail,
		ToRole:       newRequest.ToRole,
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
	logs.Info("Attempting to Get data from redis")
	if approvalJson, err := u.Redis.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(approvalJson), &appprove) == nil {
			log.Debug("Read data from: redis")
			return appprove, nil
		}
	}

	// Data not found in cache, fetch from the database
	logs.Info("Attempting to Get data from database")
	approvalDB, err := u.approvalRepo.GetByRequestID(id)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("couldn't get approval data")
	}

	//redis set
	logs.Info("Attempting to set data to redis")
	if data, err := json.Marshal(approvalDB); err == nil {
		u.Redis.Set(context.Background(), key, string(data), time.Minute*1)
	}

	return approvalDB, nil
}

func intInSlice(variable int64, list pq.Int64Array) bool {
	for _, val := range list {
		if val == variable {
			return true
		}
	}
	return false
}
