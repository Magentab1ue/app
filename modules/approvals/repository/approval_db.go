package repository

import (
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
