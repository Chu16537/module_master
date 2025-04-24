package mredis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
)

// 檢查是否存在的錯誤
func (h *Handler) IsNullKeyError(err error) bool {
	return errors.Is(redis.Nil, err)
}

// 更新時間
func (h *Handler) Expire(ctx context.Context, key string, second int64) error {
	d := time.Duration(second) * time.Second
	return h.client.Expire(ctx, key, d).Err()
}

// 執行lua語法
func (h *Handler) RunLua(ctx context.Context, luaScript string, keys []string, args ...interface{}) (interface{}, error) {
	r, err := h.client.Eval(ctx, luaScript, keys, args).Result()
	return r, err
}

// 是否存在
func (h *Handler) Exists(ctx context.Context, key string) (bool, error) {
	ok, err := h.client.Exists(ctx, key).Result()
	if ok == 1 {
		return true, err
	}

	return false, err
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

// 刪除
func (h *Handler) Del(ctx context.Context, key string) error {
	return h.client.Del(ctx, key).Err()
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

func (h *Handler) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return h.client.IncrBy(ctx, key, value).Result()
}

func (h *Handler) HIncrBy(ctx context.Context, key string, field string, value int64) (int64, error) {
	return h.client.HIncrBy(ctx, key, field, value).Result()
}

// 增加 list 資料
func (h *Handler) AddList(ctx context.Context, isFirst bool, key string, values ...interface{}) error {
	if isFirst {
		return h.client.LPush(ctx, key, values).Err()
	}
	// 從最後新增
	return h.client.RPush(ctx, key, values).Err()
}

// 取得 list 資料
// start 負數代表從後面開始 stop = -1 代表取得最後
func (h *Handler) GetList(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return h.client.LRange(ctx, key, start, stop).Result()
}

// pipe
func (h *Handler) Pipe() redis.Pipeliner {
	return h.client.TxPipeline()
}

// 新增 / 更新 指定的score(不能重複才可以使用)
func (h *Handler) AddAndUpdateZset(ctx context.Context, key string, score float64, member string) error {
	// 刪除指定 score 的 member
	_, err := h.client.ZRemRangeByScore(ctx, key, fmt.Sprintf("%f", score), fmt.Sprintf("%f", score)).Result()
	if err != nil {
		return err
	}
	err = h.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
	if err != nil {
		return err
	}
	return nil
}

// 取得 member 的 score
func (h *Handler) GetZsetMember(ctx context.Context, key string, member string) (float64, error) {
	return h.client.ZScore(ctx, key, member).Result()
}

// 取得 範圍內的 member
func (h *Handler) GetZsetRange(ctx context.Context, key string, o *ZsetRangeOpt) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min:    o.Min,
		Max:    o.Max,
		Offset: o.Offset,
	}

	if o.Count != 0 {
		opt.Count = o.Count
	}

	return h.client.ZRangeByScore(ctx, key, opt).Result()
}

// 刪除 members
func (h *Handler) DelZsetMember(ctx context.Context, key string, member []string) error {
	return h.client.ZRem(ctx, key, member).Err()
}

// 上鎖
func (h *Handler) Lock(ctx context.Context, key string, nodeID int64, second int) error {
	success, err := h.client.SetNX(ctx, key, nodeID, time.Duration(second*int(time.Second))).Result()
	if !success && err == nil {
		return errors.New("lock fail")
	}

	if err != nil {
		return err
	}

	return nil
}

// 釋放鎖
func (h *Handler) Unlock(ctx context.Context, key string, nodeID int64) error {
	value, err := h.client.Get(ctx, key).Int64()
	if err != nil {
		return err
	}
	if value == nodeID {
		_, err := h.client.Del(ctx, key).Result()
		return err
	}
	return nil // 不是持有者，不執行刪除
}

// 給每個服務更新上鎖時間
// key: server_lock:${name}  score 是更新時間  member 是nodeID
func (h *Handler) UpdateServerLockTime(ctx context.Context, name string, nodeID int64) error {
	key := fmt.Sprintf("lock:%v", name)

	err := h.AddAndUpdateZset(ctx, key, float64(time.Now().Unix()), strconv.Itoa(int(nodeID)))
	if err != nil {
		return err
	}
	return nil
}
