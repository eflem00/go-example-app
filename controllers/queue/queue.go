package queue

import (
	"time"

	"github.com/rs/zerolog/log"
)

type QueueController struct {
}

func (controller QueueController) Start() {
	for {
		log.Info().Msg("Polling queue...")
		time.Sleep(time.Second)
	}
}
