package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type config struct {
	appName       string
	appPort       int
	migrationPath string
	db            databaseConfig
}

var appConfig config

func Load() {
	viper.SetDefault("APP_NAME", "banking_application")
	viper.SetDefault("APP_PORT", 8000)

	viper.AddConfigPath("./")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()
	viper.AutomaticEnv()

	appConfig = config{
		appName:       readEnvString("APP_NAME"),
		appPort:       readEnvInt("APP_PORT"),
		migrationPath: readEnvString("MIGRATION_PATH"),
		db:            newDatabaseConfig(),
	}

}

func AppName() string {
	return appConfig.appName
}

func AppPort() int {
	return appConfig.appPort
}

func MigrationPath() string {
	return appConfig.migrationPath
}

func checkIfSet(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("key %v is not set", key))
	}
}

func readEnvInt(key string) int {
	checkIfSet(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Errorf("key %v is not a valid integer", key))
	}
	return v
}

func readEnvString(key string) string {
	checkIfSet(key)
	return viper.GetString(key)
}
