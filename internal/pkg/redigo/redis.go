package redigo

import (
	"log"
	"time"

	"abishar/internal/config"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	pool redis.Pool
}

type Result struct {
	Value interface{}
	Error error
}

// New Redis module
func New(redisCfg config.Redis) *Redis {
	// Set default 10 seconds timeout
	if redisCfg.Timeout == 0 {
		redisCfg.Timeout = 10
	}

	log.Print("Redis Connected")
	// Open connection to redis server
	return &Redis{
		pool: redis.Pool{
			MaxIdle:     redisCfg.MaxIdle,
			IdleTimeout: time.Duration(redisCfg.Timeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					redisCfg.Endpoint,
				)
			},
		},
	}

}

// Get redis command
func (rgo *Redis) Get(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("GET", args...)
}

// MGet redis command
func (rgo *Redis) MGet(keys ...string) *Result {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return rgo.cmd("MGET", args...)
}

// Setex redis command
func (rgo *Redis) Setex(key string, expireTime int, value interface{}) error {
	args := []interface{}{key, expireTime, value}
	return rgo.cmd("SETEX", args...).Error
}

// Del redis command
func (rgo *Redis) Del(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return rgo.cmd("DEL", args...).Error
}

// Keys redis command
func (rgo *Redis) Keys(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("KEYS", args...)
}

// Expire redis command
func (rgo *Redis) Expire(key string, seconds int) error {
	args := []interface{}{key, seconds}
	return rgo.cmd("EXPIRE", args...).Error
}

// Incr redis command
func (rgo *Redis) Incr(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return rgo.cmd("INCR", args...).Error
}

// Decr redis command
func (rgo *Redis) Decr(keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, v := range keys {
		args[i] = v
	}
	return rgo.cmd("DECR", args...).Error
}

// HIncrBy redis command
func (rgo *Redis) HIncrBy(key string, field string, value int64) error {
	args := []interface{}{key, field, value}
	return rgo.cmd("HINCRBY", args...).Error
}

// HDel redis command
func (rgo *Redis) HDel(key string, fields ...string) error {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, v := range fields {
		args[i+1] = v
	}
	return rgo.cmd("HDEL", args...).Error
}

// HSet redis command
func (rgo *Redis) HSet(key, field string, value interface{}) error {
	args := []interface{}{key, field, value}
	return rgo.cmd("HSET", args...).Error
}

// HGet redis command
func (rgo *Redis) HGet(key, field string) *Result {
	args := []interface{}{key, field}
	return rgo.cmd("HGET", args...)
}

// HKeys redis command
func (rgo *Redis) HKeys(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("HKEYS", args...)
}

// HVals redis command
func (rgo *Redis) HVals(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("HVALS", args...)
}

// HGetAll redis command
func (rgo *Redis) HGetAll(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("HGETALL", args...)
}

// ZRange redis command
func (rgo *Redis) ZRange(key string, start int, end int) *Result {
	args := []interface{}{key, start, end}
	return rgo.cmd("ZRANGE", args...)
}

// ZRangeByScore redis command
func (rgo *Redis) ZRangeByScore(key, min, max string, limit int) *Result {
	args := []interface{}{key, min, max}
	if limit > 0 {
		args = []interface{}{key, min, max, "LIMIT", 0, limit}
	}
	return rgo.cmd("ZRANGEBYSCORE", args...)
}

// TTL get expiry time
func (rgo *Redis) TTL(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("TTL", args...)
}

// Exists check whether a key is exists
func (rgo *Redis) Exists(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("EXISTS", args...)
}

// HMSet redis command
func (rgo *Redis) HMSet(key string, values map[string]interface{}) error {
	args := []interface{}{key}

	for key, val := range values {
		args = append(args, key, val)
	}

	result := rgo.cmd("HMSET", args...)
	return result.Error
}

// Rename a key
func (rgo *Redis) Rename(key, newKey string) *Result {
	return rgo.cmd("RENAME", key, newKey)
}

// HExists check whether a field is exists in a particular key
func (rgo *Redis) HExists(key string, fieldKey string) *Result {
	return rgo.cmd("HEXISTS", key, fieldKey)
}

func (rgo *Redis) GetPoolConnection() redis.Conn {
	return rgo.pool.Get()
}

func (rgo *Redis) cmd(command string, args ...interface{}) *Result {
	result := &Result{}
	conn := rgo.pool.Get()
	defer conn.Close()

	data, err := conn.Do(command, args...)
	if err != nil {
		// Retry mechanism
		retryConn := rgo.pool.Get()
		defer retryConn.Close()
		data, err = retryConn.Do(command, args...)
		if err != nil {
			result.Error = err
			return result
		}
	}
	result.Value = data

	return result
}

// Set a value without expiration
func (rgo *Redis) Set(key, value interface{}, args ...interface{}) *Result {
	args = append([]interface{}{key, value}, args...)
	return rgo.cmd("SET", args...)
}

// LPush Insert all the specified values at the Head of the list stored at key
func (rgo *Redis) LPush(key string, value interface{}) error {
	args := []interface{}{key, value}
	return rgo.cmd("LPUSH", args...).Error
}

// RPush Insert all the specified values at the tail of the list stored at key
func (rgo *Redis) RPush(key string, value interface{}) error {
	args := []interface{}{key, value}
	return rgo.cmd("RPUSH", args...).Error
}

// LPop Removes and returns the first element of the list stored at key
func (rgo *Redis) LPop(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("LPOP", args...)
}

// LLen return length of element of the list
func (rgo *Redis) LLen(key string) *Result {
	args := []interface{}{key}
	return rgo.cmd("LLEN", args...)
}
