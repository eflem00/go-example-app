package db

import "fmt"

// TODO: fill in the persistant storage piece

type ResultRepository struct{}

func NewResultRepository() *ResultRepository {
	return &ResultRepository{}
}

func (resultRepository *ResultRepository) GetResultById(id string) (string, error) {
	return fmt.Sprintf("somee value for %v", id), nil
}

func (resultRepository *ResultRepository) WriteResult(id string, value string) error {
	return nil
}
