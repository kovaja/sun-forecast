package db

import (
	"database/sql"
	"fmt"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"

	_ "github.com/lib/pq"
)

var db *sql.DB

func getDbConnectionString() (string, error) {
	usr, err := utils.GetEnvVariable("DB_USR")
	pwd, err := utils.GetEnvVariable("DB_PWD")
	dbHost, err := utils.GetEnvVariable("DB_HOST")
	dbName, err := utils.GetEnvVariable("DB_NAME")

	if err != nil {
		return "", utils.CustomError("Failed to load env variables for DB", err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s", usr, pwd, dbHost, dbName), nil
}

func InitializeDatabase() error {
	conStr, err := getDbConnectionString()
	if err != nil {
		logger.LogError("Failed to initialize database", err)
		return err
	}

	database, err := sql.Open("postgres", conStr)
	if err != nil {
		logger.LogError("Failed to initialize database", err)
		return err
	}

	db = database
	logger.Log("Database initialized...")
	return nil
}

func GetDb() *sql.DB {
	return db
}
