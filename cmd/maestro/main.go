package main

import (
	"fmt"
	"log"
	"os"

	"maestro/internal/config"
	"maestro/internal/db"
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

	fmt.Printf("Maestro database initialized at %s (configured port: %d)\n", cfg.DBPath, cfg.Port)
}
