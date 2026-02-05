package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitry/taskmanager/internal/config"
	"github.com/dmitry/taskmanager/internal/database"
	"github.com/dmitry/taskmanager/internal/handler"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/internal/router"
	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/dmitry/taskmanager/pkg/validator"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)

	// Проверка JWT secret
	if cfg.JWTSecret == "" {
		log.Fatal("Переменная окружения JWT_SECRET обязательна")
	}

	log.Info("Запуск Task Manager API", "environment", cfg.Environment)

	// Подключение к базе данных
	db, err := database.NewPostgres(cfg, log)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных", "error", err)
	}
	defer db.Close()

	// Подключение к Redis
	redis, err := database.NewRedis(cfg, log)
	if err != nil {
		log.Fatal("Не удалось подключиться к Redis", "error", err)
	}
	defer redis.Close()

	// Инициализация репозиториев
	employeeRepo := repository.NewEmployeeRepository(db.DB)
	taskRepo := repository.NewTaskRepository(db.DB)
	participantRepo := repository.NewTaskParticipantRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	timeEntryRepo := repository.NewTimeEntryRepository(db.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db.DB)

	// JWT сервис
	jwtService := service.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTAccessExpiryMin,
		cfg.JWTRefreshExpiryDays,
	)

	// Инициализация сервисов
	employeeService := service.NewEmployeeService(employeeRepo, log)
	authService := service.NewAuthService(employeeRepo, refreshTokenRepo, jwtService, log)
	taskService := service.NewTaskService(taskRepo, participantRepo, messageRepo, employeeRepo, db.DB, log)
	messageService := service.NewMessageService(messageRepo, log)
	timeEntryService := service.NewTimeEntryService(timeEntryRepo, taskRepo, log)

	_ = messageService
	_ = timeEntryService

	// Инициализация handlers
	v := validator.New()
	isProduction := cfg.Environment == "production"
	authHandler := handler.NewAuthHandler(authService, v, isProduction)
	employeeHandler := handler.NewEmployeeHandler(employeeService, v)
	taskHandler := handler.NewTaskHandler(taskService, v)

	// Настройка роутинга
	r := router.NewRouter(authHandler, employeeHandler, taskHandler, jwtService, cfg.FrontendURL, log)

	_ = redis // Redis будет использоваться для rate limiting позже

	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск горутины для очистки просроченных токенов
	go cleanupExpiredTokens(refreshTokenRepo, log)

	go func() {
		log.Info("Сервер запускается", "address", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Не удалось запустить сервер", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Сервер завершает работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Сервер принудительно остановлен", "error", err)
	}

	log.Info("Сервер остановлен")
}

// cleanupExpiredTokens выполняется ежедневно для удаления просроченных refresh-токенов
func cleanupExpiredTokens(repo repository.RefreshTokenRepository, log *logger.Logger) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := repo.DeleteExpired(ctx); err != nil {
			log.Error("Не удалось очистить просроченные токены", "error", err)
		} else {
			log.Info("Просроченные refresh-токены очищены")
		}
		cancel()
	}
}
