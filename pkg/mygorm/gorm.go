package mygorm

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")

	count := 0

	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("PostgreSQL not ready...")
			count++
		} else {
			log.Println("Connected to database!")
			return db
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
