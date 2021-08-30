package database

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	// Import package for side effects
	_ "github.com/lib/pq"
)

// DB global variable
var DB *sql.DB

// Config represents the database configuration
type Config struct {
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     int    `env:"DB_PORT" envDefault:"5432"`
	DbName     string `env:"DB_NAME"`
}

// Connect initializes a connnection to the database thanks to a config object
func Connect(cfg Config) {
	var err error
	dbURL := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	fmt.Println(cfg.DbPassword)
	DB, err = sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to database!")
}
