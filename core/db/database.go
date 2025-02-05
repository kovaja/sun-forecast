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
	dbQuery, err := utils.GetEnvVariable("DB_QUERY")

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

	if dbQuery == "" {
		return "", errors.New("Failed to load db query variable")
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", usr, pwd, dbHost, dbName, dbQuery), nil
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
