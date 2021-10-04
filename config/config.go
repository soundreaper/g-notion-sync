package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Load Environment Variables from .env
func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Did not load variables from .env file. This is normal for CI/CD or production.")
	}
}
