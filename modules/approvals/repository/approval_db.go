package repository

import (
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
	err := db.AutoMigrate(models.Approval{})
	if err != nil {
		panic(err)
	}
	return approvalRepositoryDB{db: db}
}

func (r approvalRepositoryDB) UpdateStatus(requestId int, req *models.UpdateStatusReq) (*models.Approval, error) {

	approval := new(models.Approval)

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

func (r approvalRepositoryDB) GetSendRequest(userId uint, optional map[string]interface{}) ([]models.Approval, error) {

	approval := []models.Approval{}
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

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approval, error) {

	approval := []models.Approval{}

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

func (r approvalRepositoryDB) DeleteApproval(requestId uint) ([]models.Approval, error) {

	approval := []models.Approval{}
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
