package db

import (
	"time"

	"github.com/eflem00/go-example-app/common"
	"github.com/eflem00/go-example-app/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RetryConnect(settings *common.Settings, config *gorm.Config, logger *common.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(settings.PgConnectionString), config)

	if err != nil {
		logger.Err(err, "Failed to connect to db, sleeping and retying in 5 seconds...")
		time.Sleep(time.Second * 5)
		return RetryConnect(settings, config, logger)
	}

	return db
}

func NewDb(settings *common.Settings, logger *common.Logger) *gorm.DB {
	config := gorm.Config{}

	db := RetryConnect(settings, &config, logger)

	if settings.IsDev() {
		logger.Info("Automigrating DB...")

		db.AutoMigrate(&entities.Result{})
	}

	return db
}
