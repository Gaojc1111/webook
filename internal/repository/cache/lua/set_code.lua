-- 验证码效验和更新函数
-- 参数:
-- KEYS[1]: 验证码的键名
-- ARGV[1]: 要校验或设置的验证码值
-- 返回值:
-- -2: 验证码已失效
-- 0: 验证码成功更新
-- -1: 验证码仍有效，无需更新

local key = KEYS[1] -- 获取传入的键名

local cntKey = key..":cnt" -- 构造计数器键名

local cal = ARGV[1] -- 获取传入的验证码值

local ttl = tonumber(redis.call("ttl", key)) -- 获取键的剩余时间

-- 检查验证码是否已失效
if ttl == -1 then
    return -2 -- 验证码已失效
elseif ttl == -2 or ttl < 540 then
    -- 验证码过期或剩余时间不足540秒，进行更新
    redis.call("set", key, cal)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0 -- 验证码成功更新
else
    return -1 -- 验证码仍有效，无需更新
end
