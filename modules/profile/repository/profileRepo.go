package repository

import (
	"approval-service/logs"
	"approval-service/modules/entities/models"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type profileRepositoryDB struct {
	db *gorm.DB
}

func NewprofileRepositoryDB(db *gorm.DB) models.ProfileRepositoryDB {
	// Check if the table exists
	if !db.Migrator().HasTable(&models.UserProfile{}) {
		err := db.AutoMigrate(models.UserProfile{})
		if err != nil {
			panic(err)
		}
	}
	return profileRepositoryDB{db: db}
}

func (r profileRepositoryDB) Create(request *models.UserProfile) error {
	if err := r.db.Create(request).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating Project with request ID %d: %v", request.ID, err), zap.Error(err))
		return fmt.Errorf("failed to create Project: %v", err)
	}

	return nil
}
func (r profileRepositoryDB) Update(req *models.UserProfile) error {

	//update to database
	if err := r.db.Save(&req).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating approval with request ID %d: %v", req.ID, err), zap.Error(err))
		return fmt.Errorf("error cant't updating approval with request ID %d", req.ID)
	}
	return nil
}

func (r profileRepositoryDB) Delete(Id uint) error {

	userProfile := new(models.UserProfile)
	if err := r.db.First(&userProfile, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't finding Profile with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't finding Profile with  ID %d: %v", Id, err)
	}

	if err := r.db.Delete(&userProfile, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete Profile with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't delete Profile with  ID %d: %v", Id, err)
	}

	return nil
}

func (r profileRepositoryDB) GetByID(id uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	result := r.db.First(&profile, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return &profile, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}