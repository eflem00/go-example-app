package http

import (
	"time"

	"github.com/rs/zerolog/log"
)

type HttpController struct{}

func (controller HttpController) Start() {
	log.Info().Msg("starting http")

	// TODO: Meaningful http implementation
	for {
		log.Info().Msg("recieving http")
		time.Sleep(time.Second)
		panic("ahh!!!")
	}
}

func (controller HttpController) Exit() {
	log.Info().Msg("caught panic in http controller")
}
