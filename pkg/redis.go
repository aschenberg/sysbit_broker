package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"sysbitBroker/config"

	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct{
	Rc   *redis.Client
	
}

func NewRedis(cfg *config.Config) (*Redis, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
		// DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		// ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		// WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		// PoolSize:     cfg.Redis.PoolSize,
		// PoolTimeout:  cfg.Redis.PoolTimeout,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("rdb.Ping",err)
	}

	return &Redis{Rc: rdb}, nil
}

func Set[T any](ctx context.Context, c *redis.Client, key string, value T, duration time.Duration) error {
	// ct, cancel := context.WithTimeout(ctx, 30)
	// defer cancel()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, v, duration).Err()
}

func Get[T any](ctx context.Context, c *redis.Client, key string) (T, error) {
	// ct, cancel := context.WithTimeout(ctx, 30)
	// defer cancel()
	var dest T = *new(T)
	v, err := c.Get(ctx, key).Result()
	if err != nil{
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func Del(ctx context.Context, c *redis.Client, key string) error {
	// ct, cancel := context.WithTimeout(ctx, 30)
	// defer cancel()

	return c.Del(ctx, key).Err()
}

func DellAll(ctx context.Context, c *redis.Client)error{
	// ct, cancel := context.WithTimeout(ctx, 30)
	// defer cancel()
	err:= c.FlushDB(ctx)
	if err != nil {
        return fmt.Errorf("Error flushing current database: %s", err)
    } 
	 return nil
}