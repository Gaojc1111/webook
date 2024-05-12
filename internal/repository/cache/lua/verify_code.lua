local key = KEYS[1]
local cntKey = key..":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    return -1
elseif code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    return -1
end