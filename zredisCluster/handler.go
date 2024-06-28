package zredisCluster

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// 檢查是否存在的錯誤
func (h *Handler) IsNullKeyError(err error) bool {
	return errors.Is(redis.Nil, err)
}

// 更新時間
func (h *Handler) UpdateTTL(ctx context.Context, key string, second int64) error {
	d := time.Duration(second) * time.Second
	return h.client.Expire(ctx, key, d).Err()
}

// 執行lua語法
func (h *Handler) RunLua(ctx context.Context, luaScript string, keys []string, args ...interface{}) (interface{}, error) {
	r, err := h.client.Eval(ctx, luaScript, keys, args).Result()
	return r, err
}

// 取得
func (h *Handler) Get(ctx context.Context, key string) (string, error) {
	return h.client.Get(ctx, key).Result()
}

// 設置
func (h *Handler) Set(ctx context.Context, key string, value interface{}, second int64) error {
	if second == 0 {
		second = -1
	}

	d := time.Duration(second) * time.Second
	return h.client.Set(ctx, key, value, d).Err()
}

func (h *Handler) IncrBy(ctx context.Context, key string, value int64) error {
	return h.client.IncrBy(ctx, key, value).Err()
}

// 刪除
func (h *Handler) Del(ctx context.Context, key string) error {
	return h.client.Del(ctx, key).Err()
}

// 是否存在
func (h *Handler) Exists(ctx context.Context, key string) (bool, error) {
	ok, err := h.client.Exists(ctx, key).Result()
	if ok == 1 {
		return true, err
	}

	return false, err
}

func (h *Handler) HGet(ctx context.Context, key string, fields ...string) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	redisData, err := h.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return result, err
	}

	for i, v := range fields {
		result[v] = redisData[i]
	}

	return result, nil
}

func (h *Handler) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return h.client.HGetAll(ctx, key).Result()
}

// hset 跟 hmset 都使用map 可以更新多筆
func (h *Handler) HSet(ctx context.Context, key string, value map[string]interface{}) error {
	return h.client.HSet(ctx, key, value).Err()
}

func (h *Handler) HDel(ctx context.Context, key string, fields ...string) error {
	return h.client.HDel(ctx, key, fields...).Err()
}

func (h *Handler) HKeys(ctx context.Context, key string) ([]string, error) {
	return h.client.HKeys(ctx, key).Result()
}

func (h *Handler) HIncrBy(ctx context.Context, key string, field string, value int64) (int64, error) {
	return h.client.HIncrBy(ctx, key, field, value).Result()
}

// 增加 list 資料 從最後新增
func (h *Handler) SetListFromLast(ctx context.Context, key string, values ...interface{}) error {
	return h.client.RPush(ctx, key, values).Err()
}

// 增加 list 資料 從頭新增
func (h *Handler) SetListFromFirst(ctx context.Context, key string, values ...interface{}) error {
	return h.client.LPush(ctx, key, values).Err()
}

// 取得 list 資料
// start 負數代表從後面開始 stop = -1 代表取得最後
func (h *Handler) GetList(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return h.client.LRange(ctx, key, start, stop).Result()
}

// 新增 / 更新 score
func (h *Handler) AddAndUpdateZset(ctx context.Context, key string, score float64, member string) error {
	return h.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// 取得 member 的 score
func (h *Handler) GetZsetMember(ctx context.Context, key string, member string) (float64, error) {
	return h.client.ZScore(ctx, key, member).Result()
}

// 取得 範圍內的 member
func (h *Handler) GetZsetRange(ctx context.Context, key string, min, max string, values ...int64) ([]string, error) {

	opt := &redis.ZRangeBy{
		Min: min,
		Max: max,
	}

	switch len(values) {
	case 1:
		opt.Offset = values[0]
	case 2:
		opt.Offset = values[0]
		opt.Count = values[1]
	}

	return h.client.ZRangeByScore(ctx, key, opt).Result()
}

// 刪除 members
func (h *Handler) DelZsetMember(ctx context.Context, key string, member []string) error {
	return h.client.ZRem(ctx, key, member).Err()
}
