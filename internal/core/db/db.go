package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB = gorm.DB

func Connect(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return conn, nil
}

func Migrate(conn *gorm.DB, models ...any) error {
	if err := conn.AutoMigrate(models...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

func Close(conn *gorm.DB) error {
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
