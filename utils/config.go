package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
		return err
	}
	return nil
}
