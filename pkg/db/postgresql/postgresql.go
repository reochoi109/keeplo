package postgresql

import (
	"keeplo/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() {
	conf := config.AppConfig.DB
	logLevel := logger.Silent // 숨김

	if config.AppConfig.Debug {
		logLevel = logger.Info
	}

	dsn := conf.DSN()
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("[GORM] Failed to connect : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[GORM] could not get sql.db : %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("[Gorm] DB not initialized")
	}
	return db
}
