package databases

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"approval-service/configs"
	"approval-service/logs"
	"approval-service/pkg/utils"
)

func NewPostgresConnection(cfg *configs.Config) (*gorm.DB, error) {

	dsn, err := utils.UrlBuilder("postgres", cfg)
	fmt.Printf("%v\n", dsn)
	if err != nil {
		logs.Error("Can't build url: ", zap.Error(err))
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logs.Error("Failed to connect to database: ", zap.Error(err))
		return nil, err
	}
	logs.Info("postgreSQL database has been connected 🐘")
	return db, nil
}

func NewPostgresConnectionX(cfg *configs.Config) (*gorm.DB, error) {

	dsn, err := utils.UrlBuilder("postgresx", cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logs.Error("Failed to connect to database: ", zap.Error(err))
	}
	logs.Info("postgreSQL database has been connected 🐘")
	return db, nil
}
