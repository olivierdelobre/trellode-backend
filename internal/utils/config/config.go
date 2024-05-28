package config

import (
	"bufio"
	"fmt"
	"log"
	"trellode-go/internal/utils/database"

	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type Config struct {
	Log *zap.Logger
	Db  *gorm.DB
}

// Init
func GetConfig() Config {
	// Get a new logger
	level := zap.InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = zap.DebugLevel
	}
	if os.Getenv("LOG_LEVEL") == "error" {
		level = zap.ErrorLevel
	}
	if os.Getenv("LOG_LEVEL") == "warn" {
		level = zap.WarnLevel
	}
	if os.Getenv("LOG_LEVEL") == "fatal" {
		level = zap.FatalLevel
	}

	// Get a new logger
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timegenerated"
	encoderCfg.LevelKey = "log.level"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			//			"pid": os.Getpid(),
		},
	}

	logger := zap.Must(config.Build())

	// Load environment variables from .env file
	err := godotenv.Load("/home/trellode/conf/.env")
	if err != nil {
		logger.Info(fmt.Sprintf("Unable to load /home/trellode/conf/.env file: %s", err))
	}

	readFile, err := os.Open("/home/trellode/conf/.env")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		// Should check if contains PASS and not display
		logger.Info(fileScanner.Text())
	}
	readFile.Close()

	_, ok := os.LookupEnv("DBHOST")
	if !ok {
		log.Panic("DBHOST environment variable required but not set")
	}
	_, ok = os.LookupEnv("DBPORT")
	if !ok {
		log.Panic("DBPORT environment variable required but not set")
	}
	_, ok = os.LookupEnv("DBPASS")
	if !ok {
		log.Panic("DBPASS environment variable required but not set")
	}
	_, ok = os.LookupEnv("DBUSER")
	if !ok {
		log.Panic("DBUSER environment variable required but not set")
	}
	_, ok = os.LookupEnv("DBNAME")
	if !ok {
		log.Panic("DBNAME environment variable required but not set")
	}

	//connect to mysql database
	//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := database.GetGormDB(logger)
	if err != nil {
		panic(err)
	}

	return Config{logger, db}
}

func GetTestConfig() Config {
	// Get a new logger
	log := zap.Must(zap.NewProduction())

	return Config{log, nil}
}
