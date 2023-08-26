package config

import (
	"log"

	"github.com/joho/godotenv"
)

var Config map[string]string

func init() {
	var err error
	Config, err = godotenv.Read()
	if err != nil {
		log.Fatal("Failed to read .env")
	}
	if _, ok := Config["PORT"]; !ok {
		log.Fatal("PORT variable not specified")
	}
	if _, ok := Config["POSTGRES_HOST"]; !ok {
		log.Fatal("POSTGRES_HOST variable not specified")
	}
	if _, ok := Config["POSTGRES_PORT"]; !ok {
		log.Fatal("POSTGRES_PORT variable not specified")
	}
	if _, ok := Config["POSTGRES_USER"]; !ok {
		log.Fatal("POSTGRES_USER variable not specified")
	}
	if _, ok := Config["POSTGRES_PASSWORD"]; !ok {
		log.Fatal("POSTGRES_PASSWORD variable not specified")
	}
	if _, ok := Config["POSTGRES_DB"]; !ok {
		log.Fatal("POSTGRES_DBNAME variable not specified")
	}
	if _, ok := Config["POSTGRES_SSLMODE"]; !ok {
		log.Fatal("POSTGRES_SSLMODE variable not specified")
	}
	if _, ok := Config["AUTH_SALT"]; !ok {
		log.Fatal("AUTH_SALT variable not specified")
	}
	if _, ok := Config["AUTH_PRIVATE_KEY"]; !ok {
		log.Fatal("AUTH_PRIVATE_KEY variable not specified")
	}
}
