package mredisCluster

const (
	// 取得node
	LuaGetNode = `
	local zset_key = KEYS[1]
	local expire_time = tonumber(ARGV[1])
	local increment = tonumber(ARGV[2])
	
	local interval = 100
	local min_score = 0
	local max_score = interval - 1
	
	
	while true do
		-- 取得範圍內的成員和分數
		local members = redis.call("ZRANGEBYSCORE", zset_key, min_score, max_score, "WITHSCORES")
		for i = 1, #members, 2 do
			local member = tonumber(members[i])
			-- 提取小數點前的數字部分(時間)
			local integer_part = tonumber(string.match(member, "^(%d+)"))

			local score = tonumber(members[i + 1])

			-- 判斷 member 是否小於 expire_time
			if integer_part < expire_time-increment then
			-- 刪除指定 score 的成員
			redis.call("ZREMRANGEBYSCORE", zset_key, score, score)  
			-- 更新分數
				local new_member = tostring(expire_time).. "." ..tostring(score)
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
				
				-- 設定新的成員，score 為最後一個 score + 1，member 設為 expire_time.score
				local new_score = last_score + 1
				local new_member = tostring(expire_time).. "." ..tostring(new_score)
				redis.call("ZADD", zset_key, new_score, new_member)
				return new_score
			else
				-- 如果沒有符合條件的成員，直接創建新的 member
				local new_score = 0
				local new_member = tostring(expire_time).. "." ..tostring(new_score)
				redis.call("ZADD", zset_key, new_score, new_member)
				return new_score
			end
		end
	end
	`

	LuaGetScore = `
	local zset_key = KEYS[1]
	local current_time = tonumber(ARGV[1])
	local expire_duration = tonumber(ARGV[2]) -- 傳入的過期秒數
	local sort_descending = ARGV[3] == "true" -- 判斷是否按大到小排序
	local num_results = tonumber(ARGV[4]) -- 需要返回的結果數量
	local expire_time = current_time - expire_duration

	-- 遍歷 ZSET 成員，刪除過期的成員
	local to_remove = {}
	local zset_members = redis.call("ZRANGE", zset_key, 0, -1)

	for _, member in ipairs(zset_members) do
	    -- 提取 member 的 UNIX 時間部分
	    local unix_timestamp = tonumber(string.match(member, "^(%d+)%."))
	    if unix_timestamp and unix_timestamp < expire_time then
	        table.insert(to_remove, member)
	    end
	end

	-- 刪除過期成員
	if #to_remove > 0 then
	    redis.call("ZREM", zset_key, unpack(to_remove))
	end

	-- 根據排序方向獲取數據
	local range_result
	if sort_descending then
	    range_result = redis.call("ZREVRANGE", zset_key, 0, num_results - 1, "WITHSCORES")
	else
	    range_result = redis.call("ZRANGE", zset_key, 0, num_results - 1, "WITHSCORES")
	end

	-- 調整回傳格式：移除 UNIX 前部分，只返回小數點後部分和分數
	local formatted_result = {}
	for i = 1, #range_result, 2 do
	    local member = range_result[i]
	    local score = range_result[i + 1]

	    -- 提取小數點後部分
	    local member_suffix = string.match(member, "%.(.+)$")
	    if member_suffix then
	        table.insert(formatted_result, member_suffix)
	        table.insert(formatted_result, score)
	    end
	end

	return formatted_result
	`
)
