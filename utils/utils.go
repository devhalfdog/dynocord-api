package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Environment(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}
