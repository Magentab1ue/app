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

func (r approvalRepositoryDB) GetReceiveRequest(userId int) ([]models.Approval, error) {

	approval := []models.Approval{}

	if err := r.db.Where("request_user = ?", userId).First(&approval).Error; err != nil {
		logs.Error(fmt.Sprintf("Error finding approval request user with ID %d: %v", userId, err), zap.Error(err))
		return nil, err
	}

	return approval, nil
}

func (r approvalRepositoryDB) Create(request *models.Approval) (*models.Approval, error) {
	if err := r.db.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return request, nil
}

func (r approvalRepositoryDB) GetByID(id uint) (*models.Approval, error) {
	var approval models.Approval
	result := r.db.First(&approval, id)

	if result.Error == nil && result.RowsAffected > 0 {
		return &approval, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}
