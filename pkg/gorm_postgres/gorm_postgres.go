package GormPostgres

import (
	"github.com/sirupsen/logrus"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"

	gormv2logrus "github.com/thomas-tacquet/gormv2-logrus"
	"gorm.io/driver/postgres"
)

func NewClient(e *logrus.Entry) *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")

	gormLogger := gormv2logrus.NewGormlog(gormv2logrus.WithLogrusEntry(e))
	gormLogger.LogMode(logger.Error)

	count := 0

	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                 gormLogger,
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err != nil {
			e.Logger.Info("PostgreSQL not ready...")
			count++
		} else {
			e.Logger.Info("Connected to database!")
			return db
		}

		if count > 10 {
			e.Logger.Info(err)
			return nil
		}

		e.Logger.Info("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
