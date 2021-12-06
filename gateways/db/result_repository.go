package db

import (
	"github.com/eflem00/go-example-app/entities"
	"gorm.io/gorm"
)

// TODO: fill in the persistant storage piece

type ResultRepository struct {
	db *gorm.DB
}

func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{
		db,
	}
}

func (repo *ResultRepository) GetResultById(id string) (entities.Result, error) {
	var result entities.Result
	repo.db.First(&result, id)

	return result, nil
}

func (repo *ResultRepository) WriteResult(id string, value string) error {
	repo.db.Create(&entities.Result{
		Id:    id,
		Value: value,
	})
	return nil
}
