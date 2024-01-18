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
	optional["request_user"] = userId

	if _, ok := optional["to"]; !ok {
		if err := r.db.Where("request_user = ?", userId).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	} else {
		to := optional["to"]
		delete(optional, "to")
		if err := r.db.Where("\"to\" IN  \"?\"", to).Where(optional).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}

	if optional != nil {
		if err := r.db.Where(optional).Where(" \"to\" IN (?)", []uint{userId}).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	} else {

		if err := r.db.Where("to IN ?", []uint{userId}).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	}

	return approval, nil
}

func (r approvalRepositoryDB) DeleteApproval(requestId uint) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	if err := r.db.Find(&approval, requestId).Error; err != nil {
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
	mock1 := models.Approvals{
		RequestID:    uuid.New(),
		To:           pq.Int64Array([]int64{1, 2, 3, 4, 5}),
		Status:       "pending",
		Project:      json.RawMessage{},
		CreationDate: time.Now(),
		RequestUser:  uint(5),
		Task:         json.RawMessage{},
	}
	mock2 := models.Approvals{
		RequestID:    uuid.New(),
		To:           pq.Int64Array([]int64{1, 2, 3, 4, 8}),
		Status:       "pending",
		Project:      json.RawMessage{},
		CreationDate: time.Now(),
		RequestUser:  uint(9),
		Task:         json.RawMessage{},
	}

	list = append(list, mock1)
	list = append(list, mock2)
	return list
}

// func (r approvalRepositoryDB) GetAll(optional map[string]interface{}) ([]models.Approval, error) {

// 	approval := []models.Approval{}

// 	if optional != nil {

// 	err := r.db.Find(&approval).Error
// 	if err != nil {
// 		return nil, errors.New("Get dog not found")
// 	}
// 	return approval, nil

// 		// if err := r.db.Where(optional).Where("to IN (?)", []uint{userId}).Find(&approval).Error; err != nil {
// 		// 	logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
// 		// 	return nil, err
// 		// }
// 	} else {

// 		if err := r.db.Where("to IN (?)", []uint{userId}).Find(&approval).Error; err != nil {
// 			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
// 			return nil, err
// 		}
// 	}

// 	return approval, nil
// }
