--你的验证码在 Redis 上的 key
-- phone_code:login:152xxxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码最多重复三次，这个记录还可以验证几次
-- phone_code:login:152xxxxxxxx:cnt
local cntKey = key..":cnt"
-- 验证码 123456
local val= ARGV[1]
local lifeDuration=tonumber(ARGV[2])
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    --    key 存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl <  lifeDuration - 60 then
    redis.call("set", key, val)
    redis.call("expire", key, lifeDuration)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, lifeDuration)
    -- 完美，符合预期
    return 0
else
    -- 发送太频繁
    return -1
end