package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Failed to load env file, %v\n", err)
		return "", err
	}

	return os.Getenv(key), nil
}
