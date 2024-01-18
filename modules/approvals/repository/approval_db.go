package repository

import (
	"errors"
	"fmt"

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
		if err := r.db.Where(optional).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	} else {
		if err := r.db.Where("to IN (?)", optional["to"]).Where(optional).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}

	if optional != nil {
		if err := r.db.Where(optional).Where("to IN (?)", []uint{userId}).Find(&approval).Error; err != nil {
			logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d: %v", userId, err), zap.Error(err))
			return nil, err
		}
	} else {

		if err := r.db.Where("to IN (?)", []uint{userId}).Find(&approval).Error; err != nil {
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
		// Case when there are keys in the optional map
		// You can add more cases based on your specific conditions

		// Example: if "status" key is present and "to" key is absent
		if optional["status"] != nil && optional["to"] == nil {
			// Add logic to handle this case
			// For now, just return the approval slice
			return approval, nil
		}

		// Example: if "to" key is present and "status" key is absent
		if optional["to"] != nil && optional["status"] == nil {
			// Add logic to handle this case
			// For now, just return the approval slice
			return approval, nil
		}

		// Add more cases based on your specific conditions
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
		// Case when there are keys in the optional map
		// You can add more cases based on your specific conditions

		// Example: if "status" key is present and "to" key is absent
		if optional["status"] != nil && optional["to"] == nil {
			// Add logic to handle this case
			// For now, just return the approval slice
			return approval, nil
		}

		// Example: if "to" key is present and "status" key is absent
		if optional["to"] != nil && optional["status"] == nil {
			// Add logic to handle this case
			// For now, just return the approval slice
			return approval, nil
		}

		// Add more cases based on your specific conditions
	}

	// Default case when optional map is not nil but doesn't match any specific condition
	return nil, errors.New("Invalid optional parameters")
}
