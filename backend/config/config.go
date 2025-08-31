package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system env vars")
	}
}

func GetAPIKey() string {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		fmt.Println("GEMINI_API_KEY is missing")
	}
	return key
}
