package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eflem00/go-example-app/controllers"
	"github.com/eflem00/go-example-app/controllers/http"
	"github.com/eflem00/go-example-app/controllers/queue"
	"github.com/eflem00/go-example-app/gateways/cache"
	"github.com/eflem00/go-example-app/gateways/db"
	"github.com/eflem00/go-example-app/usecases"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
}

func configLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env := os.Getenv("ENV"); env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func awaitSigterm() {
	log.Info().Msg("awaiting sigterm")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan

	log.Info().Msgf("caught sigterm %v", sig)
}

func startController(controller controllers.IController) {
	// start is intended to be a blocking call
	// if Exit() is called, we have caught a panic
	// if start returns, one of our controllers is no longer active and thus we should force a panic
	defer controller.Exit()
	err := controller.Start()
	panic(err)
}

func main() {

	loadEnv()

	configLogger()

	log.Info().Msg("starting app")

	container := dig.New()

	container.Provide(cache.NewCache)
	container.Provide(db.NewResultRepository)
	container.Provide(usecases.NewResultUseCase)
	container.Provide(http.NewHttpController)
	container.Provide(queue.NewQueueController)

	// start controllers in concurrent go routine

	err := container.Invoke(func(controller *http.HttpController) {
		go startController(controller)
	})

	if err != nil {
		panic(err)
	}

	err = container.Invoke(func(controller *queue.QueueController) {
		go startController(controller)
	})

	if err != nil {
		panic(err)
	}

	// blocking call in main routine to await sigterm
	awaitSigterm()

	// TODO: Shutdown gracefully below
}
