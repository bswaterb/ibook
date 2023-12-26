local key = KEYS[1]
local key_exists = redis.call("exists", key)
if key_exists == 0 then
    return -3
end
-- 用户输入的验证码
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
-- 验证码转为数字
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 1. 错误次数过多，验证码已被失效 2. 验证码已被使用
    return -1
elseif expectedCode == code then
    -- 输入正确，删除验证码
    redis.call("del", key)
    redis.call("del", cntKey)
    return 0
else
    -- 本次验证码输入错误，可重试余量 - 1
    redis.call("decr", cntKey)
    return -2
end