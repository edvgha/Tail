package main

import (
	"os"
	"tail.server/app/optimizer/code/app"
)

func main() {
	if err := app.App(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
