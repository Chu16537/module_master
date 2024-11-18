package hmrediscluster

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mrediscluster"
	"github.com/redis/go-redis/v9"
)

const (
	gameServerRank = "game_server_rank"
	gameServerIP   = "game_server_ip"

	GameServerRankExpireDuration = 300
)

// 取得 game server 排名
func GetGameServerRank(h *mrediscluster.Handler, ctx context.Context, unix int64, isMin bool, count int) ([]mrediscluster.GetScoreInfo, *errorcode.Error) {
	gsis, err := h.GetScore(ctx, gameServerRank, unix, GameServerRankExpireDuration, isMin, count)
	if err != nil {
		return nil, errorcode.New(errorcode.Redis_Error, err)
	}

	return gsis, nil
}

// 取得 game server ip
func GetGameServerIP(h *mrediscluster.Handler, ctx context.Context, nodeIDs []string) (map[string]string, *errorcode.Error) {
	data, err := h.GetClient().HMGet(ctx, gameServerIP, nodeIDs...).Result()
	if err != nil {
		return nil, errorcode.New(errorcode.Redis_Error, err)
	}

	result := map[string]string{}
	for i, v := range data {
		if v != nil {
			result[nodeIDs[i]] = v.(string)
		}
	}

	return result, nil
}

// 刪除 不使用的 game server ip
func DelGameServerIP(h *mrediscluster.Handler, ctx context.Context, unix int64) *errorcode.Error {
	result, err := h.RunLua(ctx, mrediscluster.LuaGetGameServerTimeoutMember, []string{gameServerRank}, unix, GameServerRankExpireDuration)
	if err != nil {
		return errorcode.New(errorcode.Redis_Error, err)
	}

	// 用來存儲小數點後的部分
	members := make([]string, len(result.([]interface{})))
	nodeIDs := make([]string, len(result.([]interface{})))

	for i := 0; i < len(result.([]interface{})); i += 1 {
		member := result.([]interface{})[i].(string)
		parts := strings.Split(member, ".")
		if len(parts) > 1 {
			// parts[1] 就是小數點後的部分
			members[i] = member
			nodeIDs[i] = parts[1]
		}
	}

	pipe := h.Pipe()

	// 刪除score
	pipe.ZRem(ctx, gameServerRank, members)
	// 刪除 hash ip表
	pipe.HDel(ctx, gameServerIP, nodeIDs...)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return errorcode.New(errorcode.Redis_Error, err)
	}

	return nil
}

// 更新 game server score
func UpdateGameServerRank(h *mrediscluster.Handler, ctx context.Context, oldMember string, nodeId int64, unix int64, ip string) *errorcode.Error {
	pipe := h.Pipe()

	// 刪除舊資料
	pipe.ZRem(ctx, gameServerRank, oldMember)

	// 更新
	addData := redis.Z{Score: float64(nodeId), Member: fmt.Sprintf("%v.%v", unix, nodeId)}
	pipe.ZAdd(ctx, gameServerRank, addData)

	// 記錄ip
	setData := map[string]interface{}{
		strconv.Itoa(int(nodeId)): ip,
	}
	fmt.Println("setData", setData)
	pipe.HSet(ctx, gameServerIP, setData)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errorcode.New(errorcode.Redis_Error, err)
	}

	return nil
}
