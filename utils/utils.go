package utils

import (
	"encoding/json"
	"errors"
	"kovaja/sun-forecast/utils/logger"
	"os"
)

func WriteJson(filePath string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

func ReadJson(filePath string, target any) error {
	bytes, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, target)

	if err != nil {
		return err
	}

	return nil
}

func ReturnStringResultOrError(str string, err error) (string, error) {
	if err != nil {
		return "", err
	}

	return str, nil
}

func CustomError(str string, err error) error {
	logger.LogError(str, err)
	return errors.New(str)
}
