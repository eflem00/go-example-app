package http

import (
	"time"

	"github.com/rs/zerolog/log"
)

type HttpController struct{}

func (controller HttpController) Start() {
	for {
		log.Info().Msg("Recieving http...")
		time.Sleep(time.Second)
	}
}
