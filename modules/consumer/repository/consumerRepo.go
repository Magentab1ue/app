package repository

// import (
// 	"approval-service/logs"
// 	"approval-service/modules/entities/models"
// 	"fmt"

// 	"go.uber.org/zap"
// 	"gorm.io/gorm"
// )

// type consumerRepository struct {
// 	db *gorm.DB
// }

// func NewConsumerRepository(db *gorm.DB) models.ConsumerRepository {
// 	// Check if the table exists
// 	if !db.Migrator().HasTable(&models.ConsumerOffset{}) {
// 		err := db.AutoMigrate(models.ConsumerOffset{})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	return consumerRepository{db: db}
// }

// func (r consumerRepository) Create(req *models.ConsumerOffset) error {
// 	if err := r.db.Create(req).Error; err != nil {
// 		logs.Error(fmt.Sprintf("Error save offset with  ID %d: %v", req.ID, err), zap.Error(err))
// 		return fmt.Errorf("failed to create offset: %v", err)
// 	}

// 	return nil
// }

// func (r consumerRepository) Get(req *models.ConsumerOffset) (*models.ConsumerOffset, error) {
// 	optionals := map[string]interface{}{}
// 	optionals["Topic"] = req.Topic
// 	optionals["Offset"] = req.Offset
// 	optionals["Partition"] = req.Partition
// 	res := new(models.ConsumerOffset)
// 	//update to database
// 	if err := r.db.Where(optionals).Last(res); err != nil {
// 		return nil , fmt.Errorf("not found offset")
// 	}
// 	return res,nil
// }
