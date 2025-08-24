package main

import (
	"log"
	"os"
	"release-calendar/backend/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
