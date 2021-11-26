package queue

import (
	"time"

	"github.com/rs/zerolog/log"
)

type QueueController struct {
}

func (controller QueueController) Start() {
	log.Info().Msg("Starting queue")

	// TODO: Meaningful queue implementation
	for {
		log.Info().Msg("Polling queue")
		time.Sleep(time.Second)
		panic("ahh!!!")
	}
}

func (controller QueueController) Exit() {
	log.Info().Msg("caught panic in queue controller")
}
