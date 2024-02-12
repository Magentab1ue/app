package repository

import (
	"approval-service/logs"
	"approval-service/modules/entities/models"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type taskRepositoryDB struct {
	db *gorm.DB
}

func NewproTaskRepositoryDB(db *gorm.DB) models.TaskRepositoryDB {
	// Check if the table exists
	if !db.Migrator().HasTable(&models.Task{}) {
		err := db.AutoMigrate(models.Task{})
		if err != nil {
			panic(err)
		}
	}
	return taskRepositoryDB{db: db}
}

func (r taskRepositoryDB) Create(req *models.Task) error {

	if err := r.db.Create(req).Error; err != nil {
		logs.Error(fmt.Sprintf("Fail create Project with ID %d: %v", req.ID, err), zap.Error(err))
	}
	logs.Error(fmt.Sprintf("Attempting update Project with ID %d", req.ID))

	if err := r.db.Updates(req).Error; err != nil {
		logs.Error(fmt.Sprintf("Fail create Project with ID %d: %v", req.ID, err), zap.Error(err))
		return fmt.Errorf("fail to create Project with ID %d: %v", req.ID, err)
	}

	return nil
}
func (r taskRepositoryDB) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return &task, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}

func (r taskRepositoryDB) Update(req *models.Task) error {

	if err := r.db.Updates(&req).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating project with ID %d: %v", req.ID, err), zap.Error(err))
		return fmt.Errorf("error cant't updating project with ID %d", req.ID)
	}
	return nil
}

func (r taskRepositoryDB) Delete(Id uint) error {

	task := new(models.Task)
	if err := r.db.First(&task, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't finding project with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't finding project with  ID %d: %v", Id, err)
	}

	if err := r.db.Delete(&task, Id).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete project with  ID %d: %v", Id, err), zap.Error(err))
		return fmt.Errorf("error cant't delete project with  ID %d: %v", Id, err)
	}

	return nil
}
