package mredisCluster

const (
	// 取得node
	LuaGetNode = `
	local zset_key = KEYS[1]
	local start_time = tonumber(ARGV[1])
	local increment = tonumber(ARGV[2])
	
	local interval = 100
	local min_score = 0
	local max_score = interval - 1
	
	-- 使用哈希標籤來確保操作的是同一哈希槽
	-- zset_key = zset_key .. "{node_slot}"  -- 為 zset_key 加上哈希標籤
	
	local new_member = tostring(start_time + increment)

	while true do
		-- 取得範圍內的成員和分數
		local members = redis.call("ZRANGEBYSCORE", zset_key, min_score, max_score, "WITHSCORES")
		for i = 1, #members, 2 do
			local member = tonumber(members[i])
			local score = tonumber(members[i + 1])

			-- 判斷 member 是否小於 start_time
			if member < start_time then
				-- 刪除原有的成員，然後更新分數
				redis.call("ZREMRANGEBYSCORE", zset_key, score, score)  -- 刪除指定 score 的成員
				redis.call("ZADD", zset_key, score, new_member)
				return score
			end
		end
	
		-- 若無符合條件的成員，擴大搜尋範圍
		min_score = min_score + interval
		max_score = max_score + interval
	
		-- 如果沒有更多符合條件的成員，獲取最後一個成員的 score
		if #members == 0 then
			local last_member_with_score = redis.call("ZREVRANGE", zset_key, 0, 0, "WITHSCORES")
			
			-- 檢查 ZREVRANGE 是否有返回結果
			if #last_member_with_score > 0 then
				local last_score = tonumber(last_member_with_score[2])
				
				-- 設定新的成員，score 為最後一個 score + 1，member 設為 start_time + increment
				local new_score = last_score + 1
				redis.call("ZADD", zset_key, new_score, new_member)
				return new_score
			else
				-- 如果沒有符合條件的成員，直接創建新的 member
				local new_score = 0
				redis.call("ZADD", zset_key, new_score, new_member)
				return new_score
			end
		end
	end
		`
)