package main

import (
	"fmt"
	"log"
	"os"

	"maestro/internal/api"
	"maestro/internal/config"
	"maestro/internal/db"
	"maestro/internal/repository"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("parse flags: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	repos := repository.New(database)
	server := api.New(repos)

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Maestro database initialized at %s (configured port: %d)\n", cfg.DBPath, cfg.Port)
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
