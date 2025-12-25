-- KEYS[1] = redis key
-- ARGV[1] = capacity
-- ARGV[2] = leak_rate_per_sec
-- ARGV[3] = current_time_ms

local capacity = tonumber(ARGV[1])
local leak_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- Get previous state
local data = redis.call("HMGET", KEYS[1], "water", "last_ts")
local water = tonumber(data[1]) or 0
local last_ts = tonumber(data[2]) or now

-- Leak water over time
local elapsed = (now - last_ts) / 1000
water = water - (elapsed * leak_rate)
if water < 0 then water = 0 end

-- Check capacity
if water + 1 > capacity then
    return 0
end

-- Increment
water = water + 1
redis.call("HMSET", KEYS[1], "water", water, "last_ts", now)
redis.call("EXPIRE", KEYS[1], 600)

return 1
