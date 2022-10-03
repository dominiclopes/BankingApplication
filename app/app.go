package app

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"example.com/banking/config"
)

var (
	db     *sqlx.DB
	logger *zap.SugaredLogger
)

func Init() {
	InitLogger()

	if err := initDB(); err != nil {
		panic(err)
	}
}

func InitLogger() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

func initDB() (err error) {

	dbConfig := config.Database()

	db, err = sqlx.Open(dbConfig.Driver(), dbConfig.ConnectionURL())
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTimeMins()))
	db.SetMaxIdleConns(dbConfig.MaxPoolSize())
	db.SetMaxOpenConns(dbConfig.MaxOpenCons())

	return
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
