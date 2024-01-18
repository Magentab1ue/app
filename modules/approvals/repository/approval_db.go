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

	if optional == nil {
		// Case when optional map is nil, retrieve all data
		err := r.db.Find(&approval).Error
		if err != nil {
			return nil, errors.New("Failed to retrieve data")
		}
		return approval, nil
	}

	// Check if the map is not nil and has at least one key
	if len(optional) > 0 {

		if optional["status"] != nil && optional["to"] == nil {
			err := r.db.Where("status = ?", optional["status"]).First(&approval).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("approval not found")
				}
				return nil, errors.New("Error getting active approval by ID")
			}
			return approval, nil
		}

		// Example: if "to" key is present and "status" key is absent
		if optional["to"] != nil && optional["status"] == nil {
			err := r.db.Where("to = ?", optional["to"]).First(&approval).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("approval not found")
				}
				return nil, errors.New("Error getting active approval by ID")
			}
			return approval, nil
		}

		err := r.db.Where("to = ? AND status = ? ", optional["to"], optional["status"]).First(&approval).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("approval not found")
			}
			return nil, errors.New("Error getting active approval by ID")
		}
		return approval, nil
	}

	// Default case when optional map is not nil but doesn't match any specific condition
	return nil, errors.New("Invalid optional parameters")
}

func (r approvalRepositoryDB) GetByUserID(id uint, optional map[string]interface{}) ([]models.Approvals, error) {
	approval := []models.Approvals{}

	if optional == nil {
		// Case when optional map is nil, retrieve all data
		err := r.db.Where("request_user = ?", id).First(&approval).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("approval not found")
			}
			return nil, errors.New("Error getting active approval by ID")
		}
		return approval, nil
	}

	// Check if the map is not nil and has at least one key
	if len(optional) > 0 {
		if optional["status"] != nil && optional["to"] == nil {
			err := r.db.Where("request_user = ? AND status = ?", id, optional["status"]).First(&approval).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("approval not found")
				}
				return nil, errors.New("Error getting active approval by ID")
			}
			return approval, nil
		}

		// Example: if "to" key is present and "status" key is absent
		if optional["to"] != nil && optional["status"] == nil {
			err := r.db.Where("request_user = ? AND to = ?", id, optional["to"]).First(&approval).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("approval not found")
				}
				return nil, errors.New("Error getting active approval by ID")
			}
			return approval, nil
		}

		err := r.db.Where("request_user = ? AND to = ? AND status = ? ", id, optional["to"], optional["status"]).First(&approval).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("approval not found")
			}
			return nil, errors.New("Error getting active approval by ID")
		}
		return approval, nil
	}

	// Default case when optional map is not nil but doesn't match any specific condition
	return nil, errors.New("Invalid optional parameters")
}
