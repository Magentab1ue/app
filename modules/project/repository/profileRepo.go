package repository

import (
	"approval-service/logs"
	"approval-service/modules/entities/models"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type projectRepositoryDB struct {
	db *gorm.DB
}

func NewproProjectRepositoryDB(db *gorm.DB) models.ProjectRepositoryDB {
	// Check if the table exists
	if !db.Migrator().HasTable(&models.Project{}) {
		err := db.AutoMigrate(models.Project{})
		if err != nil {
			panic(err)
		}
	}
	return projectRepositoryDB{db: db}
}

func (r projectRepositoryDB) Create(request *models.Project) error {
	if err := r.db.Create(request).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating Project with ID %d: %v", request.ID, err), zap.Error(err))
		return fmt.Errorf("failed to create Project: %v", err)
	}

	return nil
}
func (r projectRepositoryDB) GetByID(id uint) (*models.Project ,error) {
	var project models.Project
	result := r.db.First(&project, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return &project, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}

func (r projectRepositoryDB) Update(req *models.Project) error {

	if err := r.db.Save(&req).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating project with ID %d: %v", req.ID, err), zap.Error(err))
		return fmt.Errorf("error cant't updating project with ID %d", req.ID)
	}
	return nil
}

func (r projectRepositoryDB) Delete(Id uint) error {

	project := new(models.Project)
	if err := r.db.First(&project, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't finding project with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't finding project with  ID %d: %v", Id, err)
	}

	if err := r.db.Delete(&project, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete project with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't delete project with  ID %d: %v", Id, err)
	}

	return nil
}
