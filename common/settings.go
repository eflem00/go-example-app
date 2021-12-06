package common

import (
	"os"

	"github.com/joho/godotenv"
)

type Settings struct {
	Env       string
	Port      string
	QueueName string
	RedisPort string
}

func NewSettings() *Settings {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading settings")
	}

	return &Settings{
		Env:       os.Getenv("ENV"),
		Port:      os.Getenv("PORT"),
		QueueName: os.Getenv("QUEUE_NAME"),
		RedisPort: os.Getenv("REDIS_PORT"),
	}
}

func (settings *Settings) IsDev() bool {
	return settings.Env == "dev"
}
