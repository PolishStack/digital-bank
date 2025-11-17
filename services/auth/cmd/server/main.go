package main

import (
	"log"

	"bitka/auth-service/internal/config"
	"bitka/auth-service/internal/server"
)

func main() {
	cfg := config.Load() // load env-based config (see internal/config)
	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
