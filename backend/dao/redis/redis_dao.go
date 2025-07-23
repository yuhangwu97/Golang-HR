package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDAO struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisDAO(client *redis.Client) *RedisDAO {
	return &RedisDAO{
		client:    client,
		keyPrefix: "app:",
	}
}

func NewRedisDAOWithPrefix(client *redis.Client, prefix string) *RedisDAO {
	return &RedisDAO{
		client:    client,
		keyPrefix: prefix + ":",
	}
}

// Basic Cache Operations
func (dao *RedisDAO) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fullKey := dao.keyPrefix + key
	
	// Serialize value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	return dao.client.Set(ctx, fullKey, jsonValue, expiration).Err()
}

func (dao *RedisDAO) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := dao.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}
	
	// Deserialize JSON to destination
	if err := json.Unmarshal([]byte(jsonValue), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	
	return nil
}

func (dao *RedisDAO) Delete(ctx context.Context, key string) error {
	fullKey := dao.keyPrefix + key
	return dao.client.Del(ctx, fullKey).Err()
}

func (dao *RedisDAO) DeletePattern(ctx context.Context, pattern string) error {
	fullPattern := dao.keyPrefix + pattern
	
	// Get all keys matching the pattern
	keys, err := dao.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys with pattern %s: %w", fullPattern, err)
	}
	
	if len(keys) == 0 {
		return nil
	}
	
	// Delete all matching keys
	return dao.client.Del(ctx, keys...).Err()
}

func (dao *RedisDAO) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := dao.keyPrefix + key
	count, err := dao.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Advanced Cache Operations
func (dao *RedisDAO) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}
	
	result, err := dao.client.SetNX(ctx, fullKey, jsonValue, expiration).Result()
	if err != nil {
		return false, err
	}
	
	return result, nil
}

func (dao *RedisDAO) Increment(ctx context.Context, key string) (int64, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.Incr(ctx, fullKey).Result()
}

func (dao *RedisDAO) Decrement(ctx context.Context, key string) (int64, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.Decr(ctx, fullKey).Result()
}

func (dao *RedisDAO) Expire(ctx context.Context, key string, expiration time.Duration) error {
	fullKey := dao.keyPrefix + key
	return dao.client.Expire(ctx, fullKey, expiration).Err()
}

func (dao *RedisDAO) TTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.TTL(ctx, fullKey).Result()
}

// Hash Operations
func (dao *RedisDAO) HSet(ctx context.Context, key, field string, value interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	return dao.client.HSet(ctx, fullKey, field, jsonValue).Err()
}

func (dao *RedisDAO) HGet(ctx context.Context, key, field string, dest interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := dao.client.HGet(ctx, fullKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("field not found: %s in key: %s", field, key)
		}
		return fmt.Errorf("failed to get hash value: %w", err)
	}
	
	if err := json.Unmarshal([]byte(jsonValue), dest); err != nil {
		return fmt.Errorf("failed to unmarshal hash value: %w", err)
	}
	
	return nil
}

func (dao *RedisDAO) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.HGetAll(ctx, fullKey).Result()
}

func (dao *RedisDAO) HDel(ctx context.Context, key string, fields ...string) error {
	fullKey := dao.keyPrefix + key
	return dao.client.HDel(ctx, fullKey, fields...).Err()
}

// List Operations
func (dao *RedisDAO) LPush(ctx context.Context, key string, values ...interface{}) error {
	fullKey := dao.keyPrefix + key
	
	// Serialize all values to JSON
	jsonValues := make([]interface{}, len(values))
	for i, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		jsonValues[i] = jsonValue
	}
	
	return dao.client.LPush(ctx, fullKey, jsonValues...).Err()
}

func (dao *RedisDAO) RPush(ctx context.Context, key string, values ...interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValues := make([]interface{}, len(values))
	for i, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		jsonValues[i] = jsonValue
	}
	
	return dao.client.RPush(ctx, fullKey, jsonValues...).Err()
}

func (dao *RedisDAO) LPop(ctx context.Context, key string, dest interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := dao.client.LPop(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("list is empty: %s", key)
		}
		return fmt.Errorf("failed to pop from list: %w", err)
	}
	
	if err := json.Unmarshal([]byte(jsonValue), dest); err != nil {
		return fmt.Errorf("failed to unmarshal popped value: %w", err)
	}
	
	return nil
}

func (dao *RedisDAO) RPop(ctx context.Context, key string, dest interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := dao.client.RPop(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("list is empty: %s", key)
		}
		return fmt.Errorf("failed to pop from list: %w", err)
	}
	
	if err := json.Unmarshal([]byte(jsonValue), dest); err != nil {
		return fmt.Errorf("failed to unmarshal popped value: %w", err)
	}
	
	return nil
}

func (dao *RedisDAO) LLen(ctx context.Context, key string) (int64, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.LLen(ctx, fullKey).Result()
}

// Set Operations
func (dao *RedisDAO) SAdd(ctx context.Context, key string, members ...interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonMembers := make([]interface{}, len(members))
	for i, member := range members {
		jsonValue, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("failed to marshal member at index %d: %w", i, err)
		}
		jsonMembers[i] = jsonValue
	}
	
	return dao.client.SAdd(ctx, fullKey, jsonMembers...).Err()
}

func (dao *RedisDAO) SRem(ctx context.Context, key string, members ...interface{}) error {
	fullKey := dao.keyPrefix + key
	
	jsonMembers := make([]interface{}, len(members))
	for i, member := range members {
		jsonValue, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("failed to marshal member at index %d: %w", i, err)
		}
		jsonMembers[i] = jsonValue
	}
	
	return dao.client.SRem(ctx, fullKey, jsonMembers...).Err()
}

func (dao *RedisDAO) SMembers(ctx context.Context, key string) ([]string, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.SMembers(ctx, fullKey).Result()
}

func (dao *RedisDAO) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	fullKey := dao.keyPrefix + key
	
	jsonValue, err := json.Marshal(member)
	if err != nil {
		return false, fmt.Errorf("failed to marshal member: %w", err)
	}
	
	return dao.client.SIsMember(ctx, fullKey, jsonValue).Result()
}

// Additional Redis-specific methods
func (dao *RedisDAO) Pipeline() redis.Pipeliner {
	return dao.client.Pipeline()
}

func (dao *RedisDAO) TxPipeline() redis.Pipeliner {
	return dao.client.TxPipeline()
}

func (dao *RedisDAO) Ping(ctx context.Context) error {
	return dao.client.Ping(ctx).Err()
}

func (dao *RedisDAO) HealthCheck(ctx context.Context) error {
	return dao.Ping(ctx)
}

func (dao *RedisDAO) Close() error {
	return dao.client.Close()
}

func (dao *RedisDAO) FlushDB(ctx context.Context) error {
	return dao.client.FlushDB(ctx).Err()
}

func (dao *RedisDAO) Info(ctx context.Context, section ...string) (string, error) {
	return dao.client.Info(ctx, section...).Result()
}

// Bulk operations for better performance
func (dao *RedisDAO) MSet(ctx context.Context, pairs map[string]interface{}) error {
	args := make([]interface{}, 0, len(pairs)*2)
	
	for key, value := range pairs {
		fullKey := dao.keyPrefix + key
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		args = append(args, fullKey, jsonValue)
	}
	
	return dao.client.MSet(ctx, args...).Err()
}

func (dao *RedisDAO) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = dao.keyPrefix + key
	}
	
	return dao.client.MGet(ctx, fullKeys...).Result()
}

// Advanced search capabilities
func (dao *RedisDAO) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	fullMatch := dao.keyPrefix + match
	keys, nextCursor, err := dao.client.Scan(ctx, cursor, fullMatch, count).Result()
	if err != nil {
		return nil, 0, err
	}
	
	// Remove prefix from returned keys
	cleanKeys := make([]string, len(keys))
	for i, key := range keys {
		cleanKeys[i] = strings.TrimPrefix(key, dao.keyPrefix)
	}
	
	return cleanKeys, nextCursor, nil
}

// Lock operations for distributed locking
func (dao *RedisDAO) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	fullKey := dao.keyPrefix + "lock:" + key
	return dao.client.SetNX(ctx, fullKey, "locked", expiration).Result()
}

func (dao *RedisDAO) Unlock(ctx context.Context, key string) error {
	fullKey := dao.keyPrefix + "lock:" + key
	return dao.client.Del(ctx, fullKey).Err()
}

// Atomic operations
func (dao *RedisDAO) Watch(ctx context.Context, keys ...string) error {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = dao.keyPrefix + key
	}
	return dao.client.Watch(ctx, func(tx *redis.Tx) error {
		return nil
	}, fullKeys...)
}

// Rate limiting support
func (dao *RedisDAO) IncrementWithExpire(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	fullKey := dao.keyPrefix + key
	
	pipe := dao.client.TxPipeline()
	incrCmd := pipe.Incr(ctx, fullKey)
	pipe.Expire(ctx, fullKey, expiration)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return incrCmd.Val(), nil
}

// Leaderboard/sorted set operations
func (dao *RedisDAO) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	fullKey := dao.keyPrefix + key
	return dao.client.ZAdd(ctx, fullKey, members...).Err()
}

func (dao *RedisDAO) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.ZRange(ctx, fullKey, start, stop).Result()
}

func (dao *RedisDAO) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.ZRangeWithScores(ctx, fullKey, start, stop).Result()
}

func (dao *RedisDAO) ZRem(ctx context.Context, key string, members ...interface{}) error {
	fullKey := dao.keyPrefix + key
	return dao.client.ZRem(ctx, fullKey, members...).Err()
}

func (dao *RedisDAO) ZScore(ctx context.Context, key, member string) (float64, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.ZScore(ctx, fullKey, member).Result()
}

func (dao *RedisDAO) ZCard(ctx context.Context, key string) (int64, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.ZCard(ctx, fullKey).Result()
}

// Geospatial operations
func (dao *RedisDAO) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error {
	fullKey := dao.keyPrefix + key
	return dao.client.GeoAdd(ctx, fullKey, geoLocation...).Err()
}

func (dao *RedisDAO) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	fullKey := dao.keyPrefix + key
	return dao.client.GeoRadius(ctx, fullKey, longitude, latitude, query).Result()
}

// HyperLogLog operations for counting unique items
func (dao *RedisDAO) PFAdd(ctx context.Context, key string, els ...interface{}) error {
	fullKey := dao.keyPrefix + key
	return dao.client.PFAdd(ctx, fullKey, els...).Err()
}

func (dao *RedisDAO) PFCount(ctx context.Context, keys ...string) (int64, error) {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = dao.keyPrefix + key
	}
	return dao.client.PFCount(ctx, fullKeys...).Result()
}

// Stream operations (Redis Streams)
func (dao *RedisDAO) XAdd(ctx context.Context, streamKey string, values map[string]interface{}) (string, error) {
	fullKey := dao.keyPrefix + streamKey
	return dao.client.XAdd(ctx, &redis.XAddArgs{
		Stream: fullKey,
		Values: values,
	}).Result()
}

func (dao *RedisDAO) XRead(ctx context.Context, streams map[string]string) ([]redis.XStream, error) {
	fullStreams := make([]string, 0, len(streams)*2)
	for stream, id := range streams {
		fullStreams = append(fullStreams, dao.keyPrefix+stream, id)
	}
	
	return dao.client.XRead(ctx, &redis.XReadArgs{
		Streams: fullStreams,
		Block:   0,
	}).Result()
}