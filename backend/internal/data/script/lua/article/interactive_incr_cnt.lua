-- hashmap_key -> article:{articleId}:{fieldName}
local key = KEYS[1]
-- obj_key -> "read_cnt" / "collect_cnt" / "like_cnt"
local cntKey = ARGV[1]
-- +1 or -1
local delta = tonumber(ARGV[2])
local exists = redis.call("EXISTS", key)
if exists == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    return 1
else
    return 0
end