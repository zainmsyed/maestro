package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"maestro/frontend"
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
	apiServer := api.New(repos)

	distFS, err := fs.Sub(frontend.Dist, "dist")
	if err != nil {
		log.Fatalf("open embedded frontend/dist: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.HandleFunc("/", serveSPA(distFS))

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Maestro database initialized at %s (configured port: %d)\n", cfg.DBPath, cfg.Port)
	log.Printf("Maestro server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func serveSPA(dist fs.FS) http.HandlerFunc {
	files := http.FileServer(http.FS(dist))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		// If the file exists in the embedded FS, serve it directly.
		if f, err := dist.Open(r.URL.Path[1:]); err == nil {
			f.Close()
			files.ServeHTTP(w, r)
			return
		}

		// Otherwise fall back to index.html for SPA routing.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		index, err := dist.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer index.Close()
		stat, _ := index.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), index.(io.ReadSeeker))
	}
}
