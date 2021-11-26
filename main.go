package main

import (
	"os"
	"sync"

	"github.com/eflem00/go-example-app/controllers"
	"github.com/eflem00/go-example-app/controllers/http"
	"github.com/eflem00/go-example-app/controllers/queue"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func main() {

	loadEnv()

	configLogger()

	log.Info().Msg("starting app...")

	// start a slice of blocking functions in concurrent go routines
	// functions implement IController
	waitGroup := new(sync.WaitGroup)
	contrs := []controllers.IController{
		http.HttpController{},
		queue.QueueController{},
	}

	for i, contr := range contrs {
		waitGroup.Add(i)

		go func(contr controllers.IController) {
			defer waitGroup.Done()
			contr.Start()
		}(contr)
	}

	waitGroup.Wait()
}
