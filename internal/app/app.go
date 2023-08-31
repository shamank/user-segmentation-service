package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/shamank/user-segmentation-service/internal/config"
	"github.com/shamank/user-segmentation-service/internal/delivery/http"
	"github.com/shamank/user-segmentation-service/internal/repository"
	"github.com/shamank/user-segmentation-service/internal/server"
	"github.com/shamank/user-segmentation-service/internal/service"
	"github.com/shamank/user-segmentation-service/pkg/csver"
	"github.com/shamank/user-segmentation-service/pkg/logger/sl"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Dynamic Segment Service
// @version 1.0
// @description REST API for Dynamic Segment Service

// @host localhost:8000
// @BasePath /api/v1/

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func Run(configDir string) {

	//init config
	cfg := config.InitConfig(configDir)

	//init logger
	logger := sl.SetupLogger(cfg.Env)

	//init db
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresConfig.User, cfg.PostgresConfig.Password, cfg.PostgresConfig.Host, cfg.PostgresConfig.Port, cfg.PostgresConfig.DBName, cfg.PostgresConfig.SSL))
	if err != nil {
		logger.Error("error occurred in connection to database!", sl.Err(err))
	}

	//init repository
	repo := repository.NewRepositories(db)

	//init services
	services := service.NewService(service.Dependencies{
		Logger:     logger,
		Repo:       repo,
		CSVManager: csver.NewCSVManager(cfg.UserLogsDir),
	})

	//init handlers
	basePath := "http://" + cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port + "/"
	handler := http.NewHandler(services, logger, basePath)

	//init http server
	serv := server.NewServer(cfg.HTTPServer, handler.InitRoute())

	//start server
	go func() {
		logger.Info(fmt.Sprintf("HTTP-server is start up on %s:%s!", cfg.HTTPServer.Host, cfg.HTTPServer.Port))
		if err := serv.Start(); err != nil {
			logger.Error("error occurred when starting the HTTP-server", sl.Err(err))
			return
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := serv.Stop(ctx); err != nil {
		logger.Error("error occurred when stopping the HTTP server", sl.Err(err))
		return
	}

	logger.Info("the server has shut down")
}
