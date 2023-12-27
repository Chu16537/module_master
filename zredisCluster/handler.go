package zredisCluster

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) Del(key string) error {
	return h.client.Del(h.ctx, key).Err()
}

// 是否為 key 不存在的錯誤
func IsNullKeyError(err error) bool {
	return errors.Is(redis.Nil, err)
}

func (h *Handler) Get(key string) (string, error) {
	return h.client.Get(h.ctx, key).Result()
}

func (h *Handler) Exists(key string) (int64, error) {
	return h.client.Exists(h.ctx, key).Result()
}

func (h *Handler) Set(key string, value interface{}, expiration time.Duration) error {
	return h.client.Set(h.ctx, key, value, expiration).Err()
}

func (h *Handler) SetEX(key string, value interface{}, expiration time.Duration) (string, error) {
	return h.client.SetEx(h.ctx, key, value, expiration).Result()
}

func (h *Handler) HDel(key string, fields ...string) error {
	return h.client.HDel(h.ctx, key, fields...).Err()
}

func (h *Handler) HGet(key, field string) (string, error) {
	return h.client.HGet(h.ctx, key, field).Result()
}

func (h *Handler) HGetAll(key string) (map[string]string, error) {
	return h.client.HGetAll(h.ctx, key).Result()
}

func (h *Handler) HKeys(key string) ([]string, error) {
	return h.client.HKeys(h.ctx, key).Result()
}

func (h *Handler) HSet(key string, values ...interface{}) error {
	return h.client.HSet(h.ctx, key, values...).Err()
}

func (h *Handler) HIncrBy(key string, field string, value int64) (int64, error) {
	return h.client.HIncrBy(h.ctx, key, field, value).Result()
}

func (h *Handler) HMSet(key string, values ...interface{}) error {
	return h.client.HMSet(h.ctx, key, values...).Err()
}

func (h *Handler) HExists(key, field string) (bool, error) {
	return h.client.HExists(h.ctx, key, field).Result()
}

func (h *Handler) LRange(key string, start, stop int64) ([]string, error) {
	return h.client.LRange(h.ctx, key, start, stop).Result()
}

func (h *Handler) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return h.client.Eval(h.ctx, script, keys, args...).Result()
}
