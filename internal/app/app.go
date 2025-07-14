package app

import (
	"context"
	"keeplo/config"
	"keeplo/internal/adapter/rest/router"
	"keeplo/pkg/db/postgresql"
	"keeplo/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run() {
	config.Init()
	logger.Init()
	postgresql.Init()

	logger.Log.Debug("Initializing service",
		zap.String("service mode", config.AppConfig.Mode),
		zap.String("port", config.AppConfig.Port),
		zap.String("log_level", config.AppConfig.LogLevel),
		zap.String("db_address", config.AppConfig.DB.DSN()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listenForShutdown(cancel)
	start(ctx)
}

func start(ctx context.Context) {
	router.Run(ctx)
}

func listenForShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Log.Info("shutdown service...")
	cancelFunc()
}
