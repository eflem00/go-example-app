package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpController struct{}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func process(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second) // simulate some expensive work here
	fmt.Fprint(w, "OK")
}

func (controller HttpController) Start() error {
	log.Info().Msg("Starting http controller")

	port := os.Getenv("PORT")

	log.Info().Msg(fmt.Sprintf("listening on %v", port))

	http.HandleFunc("/", health)
	http.HandleFunc("/health", health)
	http.HandleFunc("/process", process)

	err := http.ListenAndServe(port, nil)

	log.Error().Err(err).Msg("error in http controller")

	return err
}

func (controller HttpController) Exit() {
	log.Error().Msg("detected exit in http controller")
}
