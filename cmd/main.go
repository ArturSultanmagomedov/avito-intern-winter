package main

import (
	"context"
	"for_avito_tech_with_gin/config"
	"for_avito_tech_with_gin/pkg"
	"for_avito_tech_with_gin/pkg/handler"
	"for_avito_tech_with_gin/pkg/repository"
	"for_avito_tech_with_gin/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/ginS"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @tittle Balance Service
// @version 1.0
// @description REST API сервис - тестовое задание для стажировки в AvitoTech

// @host localhost:8080
// @BasePath /api/v1/

func main() {
	if err := run(); err != nil {
		logrus.Fatalf("%v", err)
	}
}

func run() error {
	if err := config.Init(); err != nil {
		return err
	}
	if err := config.InitLogger(); err != nil {
		return err
	}

	// Update currencies quotes every 6 hours
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("recovered: %v", err)
			}
		}()
		var calculator pkg.CurrencyCalculator = &pkg.DefaultCurrencyCalculator{}
		calculator.UpdateRates()
		for range time.Tick(time.Hour * 6) {
			calculator.UpdateRates()
		}
	}()

	// Initialize database
	postgres, err := repository.NewPostgresDB(config.GetPostgresConfig())
	if err != nil {
		return errors.Wrap(err, "failed to initialize db")
	}
	defer postgres.Close()

	repositories := repository.NewRepository(postgres)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	ginS.Use(gin.Logger())
	ginS.Use(gin.Recovery())

	srv := new(pkg.Server)

	go func() {
		if err := srv.Run(config.GetAddress(), handlers.InitRouters()); err != nil {
			logrus.Fatal(errors.Wrap(err, "filed to init server"))
		}
		//defer func() {
		//	if err := recover(); err != nil {}
		//}()
		// обработка паник в middleware компоненте
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Error(errors.Wrap(err, "filed to shutdown server"))
		}
	}()

	return nil
}
