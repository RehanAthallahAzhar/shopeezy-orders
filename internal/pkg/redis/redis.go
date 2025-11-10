package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf("Connecting to Redis at %s", redisAddr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	return &RedisClient{Client: rdb}, nil
}

func (rc *RedisClient) Close() {
	if rc.Client != nil {
		log.Println("Menutup koneksi Redis...")
		err := rc.Client.Close()
		if err != nil {
			log.Printf("Gagal menutup koneksi Redis: %v", err)
		}
	}
}
func (rc *RedisClient) Ping(ctx context.Context) error {
	if err := rc.Client.Ping(ctx).Err(); err != nil {
		log.Printf("Gagal ping ke Redis: %v", err)
		return err
	}
	log.Println("Koneksi Redis berhasil di-ping.")
	return nil
}

// Anda dapat menambahkan metode utilitas Redis umum di sini jika diperlukan,
// seperti Get, Set, Delete, dll.
// func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
// 	return rc.Client.Get(ctx, key).Result()
// }
// func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
// 	return rc.Client.Set(ctx, key, value, expiration).Err()
// }
