package config

import (
	"flag"
	"fmt"
)

const (
	DefaultPort   = 8080
	DefaultDBPath = "maestro.db"
)

type Config struct {
	Port   int
	DBPath string
}

func Defaults() Config {
	return Config{
		Port:   DefaultPort,
		DBPath: DefaultDBPath,
	}
}

func ParseFlags(args []string) (Config, error) {
	cfg := Defaults()
	fs := flag.NewFlagSet("maestro", flag.ContinueOnError)
	fs.IntVar(&cfg.Port, "port", cfg.Port, "port for the local HTTP server")
	fs.StringVar(&cfg.DBPath, "db", cfg.DBPath, "path to the SQLite database file")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("port must be between 1 and 65535")
	}
	if cfg.DBPath == "" {
		return Config{}, fmt.Errorf("db path must not be empty")
	}
	return cfg, nil
}
