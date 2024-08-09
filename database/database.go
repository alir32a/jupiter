package database

import (
	"fmt"
	"github.com/alir32a/jupiter/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDatabaseConnection(cfg *config.DBConfig) (*gorm.DB, error) {
	var err error

	if db != nil {
		return db, nil
	}

	db, err = gorm.Open(postgres.Open(getDSN(cfg)))
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %s", err.Error())
	}

	return db, nil
}

func getDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)
}
