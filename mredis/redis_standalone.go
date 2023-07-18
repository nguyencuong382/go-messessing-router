package mredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCmd struct {
	client *redis.Client
}

// NewRedisStandalone ...
func NewRedisStandalone(config *RedisConfig) (IRedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", *config.Host, *config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &redisCmd{
		client: client,
	}, nil
}

func (_this *redisCmd) Ping(ctx context.Context) error {
	if err := _this.client.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}

func (_this *redisCmd) Set(ctx context.Context, key string, value interface{}, expireTime int64) error {
	_, err := _this.client.Set(ctx, key, value, time.Duration(expireTime)*time.Second).Result()
	return err
}

func (_this *redisCmd) Get(ctx context.Context, key string) (interface{}, error) {
	return _this.client.Get(ctx, key).Result()
}

func (_this *redisCmd) Del(ctx context.Context, key string) (int64, error) {
	return _this.client.Del(ctx, key).Result()
}

func (_this *redisCmd) Expire(ctx context.Context, key string, expire int64) (bool, error) {
	value, err := _this.client.Expire(ctx, key, time.Duration(expire)*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return value, nil
}

func (_this *redisCmd) Exist(ctx context.Context, key string) (bool, error) {
	value, err := _this.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (_this *redisCmd) Incr(ctx context.Context, key string) (int64, error) {
	result, err := _this.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (_this *redisCmd) TTL(ctx context.Context, key string) (int64, error) {
	result, err := _this.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int64(result.Seconds()), nil
}

func (_this *redisCmd) Publish(ctx context.Context, channel string, value interface{}) error {
	return _this.client.Publish(ctx, channel, value).Err()
}

func (_this *redisCmd) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return _this.client.Subscribe(ctx, channels...)
}