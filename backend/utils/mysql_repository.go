package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type MySQLRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewMySQLRepository(db *gorm.DB, redisClient *redis.Client) *MySQLRepository {
	return &MySQLRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *MySQLRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *MySQLRepository) Create(entity interface{}) error {
	return r.db.Create(entity).Error
}

func (r *MySQLRepository) FindByID(id interface{}, result interface{}) error {
	return r.db.First(result, id).Error
}

func (r *MySQLRepository) FindByIDWithCache(id interface{}, result interface{}, cacheKey string, expiration time.Duration) error {
	// Try to get from Redis cache first
	ctx := context.Background()
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - unmarshal and return
		return json.Unmarshal([]byte(cached), result)
	}

	// Cache miss - get from database
	if err := r.db.First(result, id).Error; err != nil {
		return err
	}

	// Cache the result
	resultJSON, err := json.Marshal(result)
	if err == nil {
		r.redis.Set(ctx, cacheKey, resultJSON, expiration)
	}

	return nil
}

func (r *MySQLRepository) Find(result interface{}, conditions ...interface{}) error {
	return r.db.Find(result, conditions...).Error
}

func (r *MySQLRepository) FindWithPreload(result interface{}, preload string, conditions ...interface{}) error {
	return r.db.Preload(preload).Find(result, conditions...).Error
}

func (r *MySQLRepository) FindOne(result interface{}, conditions ...interface{}) error {
	return r.db.Where(conditions[0], conditions[1:]...).First(result).Error
}

func (r *MySQLRepository) Update(entity interface{}, conditions ...interface{}) error {
	return r.db.Model(entity).Where(conditions[0], conditions[1:]...).Updates(entity).Error
}

func (r *MySQLRepository) UpdateByID(id interface{}, updates interface{}) error {
	return r.db.Model(updates).Where("id = ?", id).Updates(updates).Error
}

func (r *MySQLRepository) Delete(entity interface{}, conditions ...interface{}) error {
	return r.db.Where(conditions[0], conditions[1:]...).Delete(entity).Error
}

func (r *MySQLRepository) DeleteByID(entity interface{}, id interface{}) error {
	return r.db.Delete(entity, id).Error
}

func (r *MySQLRepository) Count(entity interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	err := r.db.Model(entity).Where(conditions[0], conditions[1:]...).Count(&count).Error
	return count, err
}

func (r *MySQLRepository) Exists(entity interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	err := r.db.Model(entity).Where(conditions[0], conditions[1:]...).Limit(1).Count(&count).Error
	return count > 0, err
}

func (r *MySQLRepository) FindWithPagination(result interface{}, page, pageSize int, conditions ...interface{}) error {
	offset := (page - 1) * pageSize
	query := r.db.Offset(offset).Limit(pageSize)
	
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	return query.Find(result).Error
}

func (r *MySQLRepository) Raw(sql string, values ...interface{}) *gorm.DB {
	return r.db.Raw(sql, values...)
}

func (r *MySQLRepository) Exec(sql string, values ...interface{}) error {
	return r.db.Exec(sql, values...).Error
}

// Cache methods
func (r *MySQLRepository) SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return r.redis.Set(ctx, key, valueJSON, expiration).Err()
}

func (r *MySQLRepository) GetCache(key string, result interface{}) error {
	ctx := context.Background()
	
	cached, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(cached), result)
}

func (r *MySQLRepository) DeleteCache(key string) error {
	ctx := context.Background()
	return r.redis.Del(ctx, key).Err()
}

func (r *MySQLRepository) InvalidateCache(pattern string) error {
	ctx := context.Background()
	
	keys, err := r.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return r.redis.Del(ctx, keys...).Err()
	}
	
	return nil
}

// Redis operations
func (r *MySQLRepository) RedisSet(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *MySQLRepository) RedisGet(key string) (string, error) {
	ctx := context.Background()
	return r.redis.Get(ctx, key).Result()
}

func (r *MySQLRepository) RedisIncr(key string) (int64, error) {
	ctx := context.Background()
	return r.redis.Incr(ctx, key).Result()
}

func (r *MySQLRepository) RedisExpire(key string, expiration time.Duration) error {
	ctx := context.Background()
	return r.redis.Expire(ctx, key, expiration).Err()
}

func (r *MySQLRepository) RedisHSet(key, field string, value interface{}) error {
	ctx := context.Background()
	return r.redis.HSet(ctx, key, field, value).Err()
}

func (r *MySQLRepository) RedisHGet(key, field string) (string, error) {
	ctx := context.Background()
	return r.redis.HGet(ctx, key, field).Result()
}

func (r *MySQLRepository) RedisHGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	return r.redis.HGetAll(ctx, key).Result()
}

func (r *MySQLRepository) RedisExists(key string) (bool, error) {
	ctx := context.Background()
	count, err := r.redis.Exists(ctx, key).Result()
	return count > 0, err
}

func (r *MySQLRepository) RedisDelete(key string) error {
	ctx := context.Background()
	return r.redis.Del(ctx, key).Err()
}

// Helper methods
func (r *MySQLRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *MySQLRepository) GetRedis() *redis.Client {
	return r.redis
}

func (r *MySQLRepository) GenerateCacheKey(prefix string, id interface{}) string {
	return fmt.Sprintf("%s:%v", prefix, id)
}