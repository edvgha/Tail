package main

import (
	"os"
	"tail.server/app/optimizer/code/app"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Optimizer starting...")
	if err := app.App(); err != nil {
		log.Info().Msg("Optimizer exit")
		os.Exit(1)
	}
	log.Info().Msg("Optimizer exit")
	os.Exit(0)
}
