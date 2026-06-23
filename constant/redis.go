package constant

const (
	REDIS_KEY_ACCESS_TOKEN_BLACKLIST = "access_token_blacklist"
	REDIS_KEY_EVENT_REMAINING        = "event_remaining"
)

// REDIS_SCRIPT_RESERVE_TICKETS: giữ vé atomic.
// KEYS[1] = remaining counter key; ARGV[1] = quantity; ARGV[2] = số vé còn lại (khởi tạo nếu key chưa có).
// Trả -1 nếu không đủ vé, ngược lại trả số còn lại sau khi trừ.
const REDIS_SCRIPT_RESERVE_TICKETS = `
if redis.call('EXISTS', KEYS[1]) == 0 then
	redis.call('SET', KEYS[1], ARGV[2])
end
local remaining = tonumber(redis.call('GET', KEYS[1]))
local qty = tonumber(ARGV[1])
if remaining < qty then
	return -1
end
return redis.call('DECRBY', KEYS[1], qty)
`

// REDIS_SCRIPT_RELEASE_TICKETS: trả vé về counter.
// KEYS[1] = remaining counter key; ARGV[1] = quantity. Chỉ cộng nếu key đang tồn tại, ngược lại trả -1.
const REDIS_SCRIPT_RELEASE_TICKETS = `
if redis.call('EXISTS', KEYS[1]) == 1 then
	return redis.call('INCRBY', KEYS[1], ARGV[1])
end
return -1
`
