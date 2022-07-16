package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Println("Failed to Get ENV File")
		log.Println(err)
	}

	getKey := os.Getenv(key)

	return getKey
}
