package app

import (
	"context"
	"keeplo/config"
	"keeplo/internal/adapter/rest/router"
	"keeplo/pkg/db/postgresql"
	"keeplo/pkg/logging"

	"go.uber.org/zap"
)

func Run() {
	config.Init()
	logging.Init()
	postgresql.Init()

	logging.Log.Debug("Initializing service",
		zap.String("service mode", config.AppConfig.Mode),
		zap.String("port", config.AppConfig.Port),
		zap.String("log_level", config.AppConfig.LogLevel),
		zap.String("db_address", config.AppConfig.DB.DSN()),
	)

	ctx := context.Background()
	router.Run(ctx)
}
