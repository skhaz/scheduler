package database

import (
	"github.com/skhaz/scheduler/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return
	}

	if err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return
	}

	if err = db.AutoMigrate(&model.Trigger{}); err != nil {
		return
	}

	return
}
