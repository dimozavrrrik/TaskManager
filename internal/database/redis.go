package database

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitry/taskmanager/internal/config"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	logger *logger.Logger
}

func NewRedis(cfg *config.Config, log *logger.Logger) (*RedisClient, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Redis: %w", err)
	}

	log.Info("Подключение к Redis установлено")

	return &RedisClient{
		Client: client,
		logger: log,
	}, nil
}

func (r *RedisClient) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// Set сохраняет значение в Redis с TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get получает значение из Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete удаляет ключ из Redis
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists проверяет существование ключа
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Increment увеличивает счетчик
func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// Expire устанавливает TTL для ключа
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}
