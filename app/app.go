package app

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/dominiclopes/BankingApplication/config"
)

var (
	db     *sqlx.DB
	logger *zap.SugaredLogger
)

func Init() {
	InitLogger()

	initDB()
}

func InitLogger() {
	log.Info("Initializing the logger")

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while initializing the logger, err: %v\n", err.Error())
	}

	logger = zapLogger.Sugar()
	log.Info("Initialized the logger")
}

func initDB() {
	log.Info("Initializing the DB connection")
	dbConfig := config.Database()

	var err error
	db, err = sqlx.Open(dbConfig.Driver(), dbConfig.ConnectionURL())
	if err != nil {
		log.Fatalf("error occurred while connecting to the database, err: %v\n", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error occurred while verifying connection to the database, err: %v\n", err.Error())
	}

	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTimeMins()))
	db.SetMaxIdleConns(dbConfig.MaxPoolSize())
	db.SetMaxOpenConns(dbConfig.MaxOpenCons())
	log.Info("Initialized the DB connection")
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func GetDB() *sqlx.DB {
	return db
}

func Close() {
	logger.Sync()
	db.Close()
}
