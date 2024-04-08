package GormPostgres

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"

	gormv2logrus "github.com/thomas-tacquet/gormv2-logrus"
	"gorm.io/driver/postgres"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func NewClient(e *logrus.Entry, postgresDSN string) *gorm.DB {
	once.Do(func() {
		gormLogger := gormv2logrus.NewGormlog(gormv2logrus.WithLogrusEntry(e))
		gormLogger.LogMode(logger.Info)
		gormLogger.SlowThreshold = 100 * time.Millisecond
		gormLogger.SkipErrRecordNotFound = true
		count := 0
		for {
			db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
				Logger:                 gormLogger,
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
			})
			if err != nil {
				e.Logger.Info("PostgreSQL not ready...")
				count++
			} else {
				e.Logger.Info("Connected to database!")
				instance = db
				return
			}
			if count > 10 {
				e.Logger.Info(err)
				return
			}
			e.Logger.Info("Backing off for two seconds...")
			time.Sleep(2 * time.Second)
		}
	})
	return instance
}
