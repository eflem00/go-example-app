package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/eflem00/go-example-app/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type HttpController struct{}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func getResultsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := usecases.GetResultById(r.Context(), id)

	if err != nil {
		http.Error(w, "Error reading key", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, val)
}

func (controller HttpController) Start() error {
	log.Info().Msg("Starting http controller")

	port := os.Getenv("PORT")

	log.Info().Msg(fmt.Sprintf("listening on %v", port))

	r := chi.NewRouter()
	r.Get("/", health)
	r.Get("/health", health)
	r.Get("/getresults/{id}", getResultsById)

	// this is essentially a blocking call
	err := http.ListenAndServe(port, r)

	log.Error().Err(err).Msg("error in http controller")

	return err
}

func (controller HttpController) Exit() {
	log.Error().Msg("detected exit in http controller")
}
