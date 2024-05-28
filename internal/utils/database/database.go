package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetGormDB initializes a connection to the database and returns a handle
func GetGormDB(log *zap.Logger) (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("LOG_LEVEL") == "info" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(getConnectString()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Error(fmt.Sprintf("GetGormDB:%s\n", err))
		return nil, err
	}

	// Existence previously checked with getConnectString
	host, _ := os.LookupEnv("DBHOST")
	name, _ := os.LookupEnv("DBNAME")
	user, _ := os.LookupEnv("DBUSER")
	param, _ := os.LookupEnv("DBPARAM")
	log.Info(fmt.Sprintf("Successfully connected to 'database' %s on host %s as user '%s' (%s)", name, host, user, param))

	sqlDB, err := db.DB()
	if err != nil {
		log.Error(fmt.Sprintf("GetGormDB:%s\n", err))
		return nil, err
	}
	iMaxIdle := 5
	if os.Getenv("DBMAXIDLE") != "" {
		iMaxIdle, _ = strconv.Atoi(os.Getenv("DBMAXIDLE"))
	}
	iMaxOpen := 10
	if os.Getenv("DBMAXOPEN") != "" {
		iMaxOpen, _ = strconv.Atoi(os.Getenv("DBMAXOPEN"))
	}
	sqlDB.SetMaxIdleConns(iMaxIdle)
	sqlDB.SetMaxOpenConns(iMaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	return db, nil
}

func getConnectString() string {
	dbHost, ok := os.LookupEnv("DBHOST")
	if !ok {
		log.Panic("DBHOST environment variable required but not set")
	}
	dbPort, ok := os.LookupEnv("DBPORT")
	if !ok {
		log.Panic("DBPORT environment variable required but not set")
	}
	dbUser, ok := os.LookupEnv("DBUSER")
	if !ok {
		log.Panic("DBUSER environment variable required but not set")
	}
	dbPassword, ok := os.LookupEnv("DBPASS")
	if !ok {
		log.Panic("DBPASS environment variable required but not set")
	}
	dbName, ok := os.LookupEnv("DBNAME")
	if !ok {
		log.Panic("DBNAME environment variable required but not set")
	}
	dbParam, _ := os.LookupEnv("DBPARAM")

	var dsn string
	if dbParam != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbParam)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	}

	return dsn
}
