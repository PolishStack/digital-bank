package db

import (
	"fmt"
	"time"

	"bitka/db/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(dsnUrl string) (*gorm.DB, error) {
	dsn, err := utils.SafeEncodeDSN(dsnUrl)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s", dsn)
	}
	// set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
