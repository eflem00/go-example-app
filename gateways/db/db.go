package db

import (
	"github.com/eflem00/go-example-app/common"
	"github.com/eflem00/go-example-app/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb(settings *common.Settings, logger *common.Logger) *gorm.DB {
	config := gorm.Config{}

	db, err := gorm.Open(postgres.Open(settings.PgConnectionString), &config)

	if err != nil {
		panic(err)
	}

	if settings.IsDev() {
		logger.Info("Automigrating DB...")

		db.AutoMigrate(&entities.Result{})
	}

	return db
}
