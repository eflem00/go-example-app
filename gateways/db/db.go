package db

import "fmt"

// TODO: fill in the persistant storage piece

func GetResultById(id string) (string, error) {
	return fmt.Sprintf("somee value for %v", id), nil
}

func WriteResult(id string, value string) error {
	return nil
}
