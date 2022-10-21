package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/dominiclopes/BankingApplication/config"
)

var ErrFindingDriver = errors.New("no migrate driver instance found")

func CreateFile(fileName string) (err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return
	}

	err = f.Close()
	return
}

func CreateMigrationFile(fileName string) (err error) {
	if len(fileName) == 0 {
		err = errors.New("filename is not provided")
		return
	}

	timeStamp := time.Now().Unix()
	upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.sql", config.MigrationPath(), timeStamp, fileName)
	downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.sql", config.MigrationPath(), timeStamp, fileName)

	if err = CreateFile(upMigrationFilePath); err != nil {
		return
	}
	fmt.Println("created file:", upMigrationFilePath)

	if err = CreateFile(downMigrationFilePath); err != nil {
		err = os.Remove(upMigrationFilePath)
		return
	}

	fmt.Println("created file:", downMigrationFilePath)
	return
}

func RunMigrations() (err error) {
	dbConfig := config.Database()

	db, err := sql.Open(dbConfig.Driver(), dbConfig.ConnectionURL())
	if err != nil {
		return err
	}

	driver, err := getDBDriverInstance(db, dbConfig.Driver())
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance(GetMigrationPath(), config.Database().ConnectionURL(), driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange || err == nil {
		return nil
	}

	return
}

func RollbackMigration(s string) (err error) {
	steps, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	m, err := migrate.New(GetMigrationPath(), config.Database().ConnectionURL())
	if err != nil {
		return err
	}

	err = m.Steps(-1 * steps)
	if err == migrate.ErrNoChange || err == nil {
		return nil
	}

	return err
}

func getDBDriverInstance(db *sql.DB, driverName string) (driver database.Driver, err error) {
	switch driverName {
	case "postgres":
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		return
	default:
		return nil, ErrFindingDriver
	}
}

func GetMigrationPath() string {
	return fmt.Sprintf("file://%s", config.MigrationPath())
}
