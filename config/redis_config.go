// config/redis_config.go
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// コンテキスト作成（Redis操作で必要）
var (
	Ctx = context.Background()
	rdb *redis.Client
)

// Redisクライアントの初期化
func InitRedisClient() *redis.Client {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_DOMAIN"), os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
		if err := rdb.Ping(Ctx).Err(); err != nil {
			panic(fmt.Sprintf("Redis接続エラー: %s", err.Error()))
		}
	}
	return rdb
}

func RedisSet(key string, value interface{}, duration time.Duration) error {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("値のJSON変換エラー: %s", err.Error())
		}
		data = string(bytes)
	}

	if err := InitRedisClient().Set(Ctx, key, data, duration).Err(); err != nil {
		return fmt.Errorf("保存エラー: %w", err)
	}
	return nil
}

func RedisGet(key string) (string, error) {
	value, err := InitRedisClient().Get(Ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("キーが存在しません: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("取得エラー: %w", err)
	}
	return value, nil
}

func RedisDel(key string) error {
	deleted, err := InitRedisClient().Del(Ctx, key).Result()
	if err != nil {
		return fmt.Errorf("削除エラー: %w", err)
	}
	if deleted == 0 {
		return fmt.Errorf("キーが存在しないため削除されませんでした: %s", key)
	}
	return nil
}
