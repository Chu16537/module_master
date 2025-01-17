package mredis

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

	LuaGetGameServerTimeoutMember = `
	local zset_key = KEYS[1]
	local current_time = tonumber(ARGV[1])
	local expire_duration = tonumber(ARGV[2]) -- 傳入的過期秒數
	local expire_time = current_time - expire_duration
	
	-- 用來記錄所有刪除的數值（小數點後部分）
	local to_remove_values = {}
	
	-- 用來記錄每個成員的 unix_timestamp 和 decimal_value
	local member_info = {}
	
	-- 遍歷 ZSET 成員
	local zset_members = redis.call("ZRANGE", zset_key, 0, -1)
	
	for _, member in ipairs(zset_members) do
		-- 提取 member 的 UNIX 時間部分和小數點後部分
		local unix_timestamp, decimal_value = string.match(member, "^(%d+)%.(.+)$")
		if unix_timestamp then
			unix_timestamp = tonumber(unix_timestamp)
			-- 比較 UNIX 時間是否過期
			if unix_timestamp < expire_time then
				-- 記錄小數點後部分
				table.insert(to_remove_values, decimal_value)
				-- 記錄為字串格式的 unix_timestamp 和 decimal_value
				table.insert(member_info, tostring(member))  -- 使用 tostring() 轉換為字符串
			end
		end
	end
	
	-- 回傳過期的 unix_timestamp 和 decimal_value 作為字符串
	return member_info	
	`

	LuaTest = `
	local zset_key = KEYS[1]
	local member = ARGV[1]
	local can_update_score = tonumber(ARGV[2]) -- 转换为数字
	local new_score = tonumber(ARGV[3])       -- 转换为数字
	
	-- 获取指定 member 的 score
	local current_score = redis.call("ZSCORE", zset_key, member)
	
	if current_score then
		current_score = tonumber(current_score) -- 转换为数字
		if current_score >= can_update_score then
			return false
		end
	end
	
	-- 更新 member 的 score 为 new_score
	redis.call("ZADD", zset_key, new_score, member)
	return true
	`
)
