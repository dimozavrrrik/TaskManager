package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/dmitry/taskmanager/internal/config"
	"github.com/dmitry/taskmanager/pkg/logger"
)

type Database struct {
	DB     *sql.DB
	logger *logger.Logger
}

func NewPostgres(cfg *config.Config, log *logger.Logger) (*Database, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)

	maxIdleTime, err := time.ParseDuration(cfg.DBMaxIdleTime)
	if err != nil {
		maxIdleTime = 15 * time.Minute
	}
	db.SetConnMaxIdleTime(maxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	log.Info("Подключение к базе данных установлено")

	return &Database{
		DB:     db,
		logger: log,
	}, nil
}

func (d *Database) Close() error {
	d.logger.Info("Закрытие подключения к базе данных")
	return d.DB.Close()
}

func (d *Database) HealthCheck() error {
	return d.DB.Ping()
}
