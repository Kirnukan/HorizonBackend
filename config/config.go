package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port      string
	PgHost    string
	PgPort    string
	PgUser    string
	PgPass    string
	PgDBName  string
	PgSSLMode string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		PgHost:    os.Getenv("PG_HOST"),
		PgPort:    os.Getenv("PG_PORT"),
		PgUser:    os.Getenv("PG_USER"),
		PgPass:    os.Getenv("PG_PASS"),
		PgDBName:  os.Getenv("PG_DBNAME"),
		PgSSLMode: os.Getenv("PG_SSLMODE"),
	}, nil
}
