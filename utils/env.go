package utils

import (
	"kovaja/sun-forecast/utils/logger"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		logger.LogError("Failed to load env file", err)
	}

	return os.Getenv(key), nil
}

func IsDev() bool {
	dev, err := GetEnvVariable("DEV")
	if err != nil {
		return false
	}

	if dev == "1" {
		logger.Log("Running in dev environment")
		return true
	}

	return false
}
