package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"approval-service/logs"
	"approval-service/modules/entities/models"
)

type approvalRepositoryDB struct {
	db *gorm.DB
}

func NewapprovalRepositoryDB(db *gorm.DB) approvalRepositoryDB {
	err := db.AutoMigrate(models.Approvals{})
	if err != nil {
		panic(err)
	}
	mock := mockdata()
	db.Save(mock)
	return approvalRepositoryDB{db: db}
}

func (r approvalRepositoryDB) UpdateStatus(requestId uint, req *models.UpdateStatusReq) (*models.Approvals, error) {

	approval := new(models.Approvals)

	if err := r.db.First(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error finding approval for update with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, err
	}
	//update data
	approval.Status = req.Status
	approval.Approver = req.Approver

	//update to database
	if err := r.db.Save(&approval).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, err
	}
	return approval, nil
}

func (r approvalRepositoryDB) GetSendRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where("request_user = ?", userId)
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where(optional).Where("? = ANY(\"to\")", to)
	}
	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where(optional).Where("project @> ?", projectId).Find(&approval)
	}
	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval Send request user with ID %d with optional %s : %v", userId, optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approval Send request user with ID %d with optional %s ", userId, optionalStr)
	}
	return approval, nil
}

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where(optional).Where("? = ANY(\"to\")", userId)

	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where(optional).Where("project @> ?", projectId).Find(&approval)
	}
	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d with optinoal %s %v", userId, optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approval receive request user with ID %d with optinoal %s", userId, optionalStr)
	}
	return approval, nil
}

func (r approvalRepositoryDB) DeleteApproval(requestId uint) (*models.Approvals, error) {

	approval := new(models.Approvals)
	if err := r.db.First(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't finding approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, err
	}

	if err := r.db.Delete(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, err
	}

	return approval, nil
}

func (r approvalRepositoryDB) Create(request *models.Approvals) (*models.Approvals, error) {
	if err := r.db.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return request, nil
}

func (r approvalRepositoryDB) GetByID(id uint) (*models.Approvals, error) {
	var approval models.Approvals
	result := r.db.First(&approval, id)

	if result.Error == nil && result.RowsAffected > 0 {
		return &approval, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}

func mockdata() (list []models.Approvals) {
	optional := map[string]interface{}{
		"id": 1,
	}
	json, _ := json.Marshal(optional)
	mock1 := models.Approvals{
		RequestID:    uuid.New(),
		To:           pq.Int64Array([]int64{1, 2, 3, 4, 5}),
		Status:       "pending",
		Project:      json,
		CreationDate: time.Now(),
		RequestUser:  uint(5),
		Task:         json,
	}

	mock2 := models.Approvals{
		RequestID:    uuid.New(),
		To:           pq.Int64Array([]int64{1, 2, 3, 4, 8}),
		Status:       "pending",
		Project:      json,
		CreationDate: time.Now(),
		RequestUser:  uint(9),
		Task:         json,
	}

	mock3 := models.Approvals{
		RequestID:    uuid.New(),
		To:           pq.Int64Array([]int64{1, 2, 3, 4, 9}),
		Status:       "pending",
		Project:      json,
		CreationDate: time.Now(),
		RequestUser:  uint(20),
		Task:         json,
	}

	list = append(list, mock1)
	list = append(list, mock2)
	list = append(list, mock3)
	return list
}

func (r approvalRepositoryDB) GetAll(optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db

	if _, ok := optional["status"]; ok {
		condition = condition.Where("status = ?", optional["status"])
		
	}
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where(optional).Where("? = ANY(\"to\")", to)
	}
	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where(optional).Where("project @> ?", projectId).Find(&approval)
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals with optional %s : %v", optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals optional %s ",  optionalStr)
	}
	return approval, nil

}

func (r approvalRepositoryDB) GetByUserID(id uint, optional map[string]interface{}) ([]models.Approvals, error) {
	
	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where("request_user = ?", id)

	if _, ok := optional["status"]; ok {
		condition = condition.Where("status = ?", optional["status"])
		
	}
	if _, ok := optional["to"]; ok {
		condition = condition.Where("? = ANY(\"to\")", optional["to"])
		
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by userid %d with optional %s : %v",id, optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals by userid %d optional %s ",id,  optionalStr)
	}

	return approval, nil
}
