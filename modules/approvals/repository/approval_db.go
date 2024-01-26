package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	return approvalRepositoryDB{db: db}
}

func (r approvalRepositoryDB) Update(req *models.Approvals) (*models.Approvals, error) {

	//update to database
	if err := r.db.Save(&req).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating approval with request ID %d: %v", req.ID, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't updating approval with request ID %d", req.ID)
	}
	return req, nil
}

func (r approvalRepositoryDB) GetSendRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where("request_user = ?", userId)
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where("? = ANY(\"to\")", to)
	}
	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where("project ->> 'id' = ? ", projectId)
	}
	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval Send request user with ID %d with optional %s : %v", userId, optionalStr, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't finding approval Send request user with ID %d with optional %s", userId, optionalStr)
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approval Send request user with ID %d with optional %s ", userId, optionalStr)
	}
	return approval, nil
}

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where("? = ANY(\"to\")", userId)

	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where("project ->> 'id' = ? ", projectId)
	}
	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d with optinoal %s %v", userId, optionalStr, err), zap.Error(err))
		return nil, fmt.Errorf("error finding approval receive request user with ID %d with optinoal %s", userId, optionalStr)
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
		return nil, fmt.Errorf("error cant't finding approval with request ID %d: %v", requestId, err)
	}

	if err := r.db.Delete(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't delete approval with request ID %d: %v", requestId, err)
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
		condition = condition.Where("? = ANY(\"to\")", to)
	}
	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where("project ->> 'id' = ? ", projectId)
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals with optional %s : %v", optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals optional %s ", optionalStr)
	}
	return approval, nil

}

func (r approvalRepositoryDB) GetByUserID(id uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db

	if _, ok := optional["status"]; ok {
		condition = condition.Where("status = ?", optional["status"])

	}
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where("? = ANY(\"to\")", to)
	}
	if _, ok := optional["project"]; ok {
		projectId := optional["project"]
		delete(optional, "project")
		condition = condition.Where("project ->> 'id' = ? ", projectId)
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by userid %d with optional %s : %v", id, optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals by userid %d optional %s ", id, optionalStr)
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetByRequestID(id uuid.UUID) ([]models.Approvals, error) {
	approval := []models.Approvals{}

	err := r.db.Where("request_id = ?", id).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by request_id : %v  : %v", id, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals by request_id : %v ", id)
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetByRequestIDLast(id uuid.UUID) (*models.Approvals, error) {
	approval := new(models.Approvals)

	err := r.db.Where("request_id = ?", id).Last(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by request_id : %v  : %v to validation", id, err), zap.Error(err))
		return nil, err
	}

	return approval, nil
}

//test
