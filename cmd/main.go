package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "morse-app: ", log.LstdFlags|log.Lshortfile)

	srv := server.New(logger)

	if err := srv.Start(); err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}
