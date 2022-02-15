package config

import (
	"fmt"
	"for_avito_tech_with_gin/pkg/repository"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

func Init() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "error initializing config")
	}
	if err := godotenv.Load(); err != nil {
		return errors.Wrap(err, "error loading env variables")
	}

	return nil
}

func InitLogger() error {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if viper.GetString("log.output") != "" {
		currentTime := viper.GetString("log.output") + time.Now().In(time.UTC).Format(time.RFC3339)
		f, err := os.OpenFile(currentTime, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return errors.Wrap(err, "error opening file")
		}
		logrus.SetOutput(f)
	}

	return nil
}

func GetPostgresConfig() repository.Config {
	return repository.Config{
		Username: viper.GetString("db.username"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func GetAddress() string {
	return fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))
}
