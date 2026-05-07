package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"maestro/internal/api"
	"maestro/internal/config"
	"maestro/internal/db"
	"maestro/internal/repository"
)

const frontendDist = "frontend/dist"

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
	apiServer := api.New(repos)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.HandleFunc("/", serveFrontend(frontendDist))

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Maestro database initialized at %s (configured port: %d)\n", cfg.DBPath, cfg.Port)
	log.Printf("Maestro server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func serveFrontend(staticDir string) http.HandlerFunc {
	files := http.FileServer(http.Dir(staticDir))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		path := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if path == "." {
			path = "index.html"
		}
		if path == ".." || strings.HasPrefix(path, "../") {
			http.NotFound(w, r)
			return
		}

		candidate := filepath.Join(staticDir, path)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}
}
