package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

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

