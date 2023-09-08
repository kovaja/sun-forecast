package db

import (
	"database/sql"
	"errors"
	"fmt"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"

	_ "github.com/lib/pq"
)

func getDbConnectionString() (string, error) {
	usr, err := utils.GetEnvVariable("DB_USR")
	pwd, err := utils.GetEnvVariable("DB_PWD")
	dbHost, err := utils.GetEnvVariable("DB_HOST")
	dbName, err := utils.GetEnvVariable("DB_NAME")

	if err != nil {
		return "", utils.CustomError("Failed to load env variables for DB", err)
	}

	if usr == "" {
		return "", errors.New("Failed to load db usr variable")
	}

	if pwd == "" {
		return "", errors.New("Failed to load db pwd variable")
	}

	if dbHost == "" {
		return "", errors.New("Failed to load db host variable")
	}

	if dbName == "" {
		return "", errors.New("Failed to load db name variable")
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s", usr, pwd, dbHost, dbName), nil
}

func InitializeDatabase() (*sql.DB, error) {
	conStr, err := getDbConnectionString()
	if err != nil {
		logger.LogError("Failed to initialize database", err)
		return nil, err
	}

	database, err := sql.Open("postgres", conStr)
	if err != nil {
		logger.LogError("Failed to initialize database", err)
		return nil, err
	}

	logger.Log("Database initialized...")
	return database, nil
}
