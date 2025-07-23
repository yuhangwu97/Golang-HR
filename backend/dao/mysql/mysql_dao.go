package mysql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gin-project/dao/interfaces"

	"gorm.io/gorm"
)

// MySQLDAO implements BaseDAO interface for MySQL operations
type MySQLDAO[T any] struct {
	masterDB  *gorm.DB
	slaveDB   *gorm.DB
	cache     interfaces.CacheDAO
	tableName string
	modelType reflect.Type
}

func NewMySQLDAO[T any](masterDB, slaveDB *gorm.DB, cache interfaces.CacheDAO, tableName string) *MySQLDAO[T] {
	var zero T
	modelType := reflect.TypeOf(zero)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	return &MySQLDAO[T]{
		masterDB:  masterDB,
		slaveDB:   slaveDB,
		cache:     cache,
		tableName: tableName,
		modelType: modelType,
	}
}

func (dao *MySQLDAO[T]) Create(ctx context.Context, entity *T) error {
	result := dao.masterDB.WithContext(ctx).Create(entity)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate related cache
	dao.invalidateCache(ctx, "list")
	
	return nil
}

func (dao *MySQLDAO[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	// Try cache first
	cacheKey := dao.generateCacheKey("id", id)
	var entity T
	
	if dao.cache != nil {
		err := dao.cache.Get(ctx, cacheKey, &entity)
		if err == nil {
			return &entity, nil
		}
	}
	
	// Cache miss, query database
	result := dao.slaveDB.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Cache the result
	if dao.cache != nil {
		dao.cache.Set(ctx, cacheKey, entity, 10*time.Minute)
	}
	
	return &entity, nil
}

func (dao *MySQLDAO[T]) Find(ctx context.Context, queryBuilder interfaces.QueryBuilder) ([]*T, error) {
	var entities []*T
	
	// Build query
	sql, args := queryBuilder.ToSQL()
	
	// Try cache for simple queries
	if dao.canCache(queryBuilder) {
		cacheKey := dao.generateQueryCacheKey(sql, args)
		if dao.cache != nil {
			err := dao.cache.Get(ctx, cacheKey, &entities)
			if err == nil {
				return entities, nil
			}
		}
	}
	
	// Execute query
	result := dao.slaveDB.WithContext(ctx).Raw(sql, args...).Scan(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Cache results for simple queries
	if dao.canCache(queryBuilder) && dao.cache != nil {
		cacheKey := dao.generateQueryCacheKey(sql, args)
		dao.cache.Set(ctx, cacheKey, entities, 5*time.Minute)
	}
	
	return entities, nil
}

func (dao *MySQLDAO[T]) Update(ctx context.Context, entity *T, queryBuilder interfaces.QueryBuilder) error {
	sql, args := queryBuilder.ToSQL()
	
	result := dao.masterDB.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate cache
	dao.invalidateAllCache(ctx)
	
	return nil
}

func (dao *MySQLDAO[T]) Delete(ctx context.Context, queryBuilder interfaces.QueryBuilder) error {
	sql, args := queryBuilder.ToSQL()
	
	result := dao.masterDB.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate cache
	dao.invalidateAllCache(ctx)
	
	return nil
}

func (dao *MySQLDAO[T]) CreateMany(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	
	result := dao.masterDB.WithContext(ctx).CreateInBatches(entities, 100)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate cache
	dao.invalidateAllCache(ctx)
	
	return nil
}

func (dao *MySQLDAO[T]) UpdateMany(ctx context.Context, updates map[string]interface{}, queryBuilder interfaces.QueryBuilder) error {
	sql, args := queryBuilder.ToSQL()
	
	result := dao.masterDB.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate cache
	dao.invalidateAllCache(ctx)
	
	return nil
}

func (dao *MySQLDAO[T]) DeleteMany(ctx context.Context, queryBuilder interfaces.QueryBuilder) error {
	sql, args := queryBuilder.ToSQL()
	
	result := dao.masterDB.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		return result.Error
	}
	
	// Invalidate cache
	dao.invalidateAllCache(ctx)
	
	return nil
}

func (dao *MySQLDAO[T]) Count(ctx context.Context, queryBuilder interfaces.QueryBuilder) (int64, error) {
	sql, args := queryBuilder.ToSQL()
	
	var count int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as count_query", sql)
	
	result := dao.slaveDB.WithContext(ctx).Raw(countSQL, args...).Scan(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	
	return count, nil
}

func (dao *MySQLDAO[T]) Exists(ctx context.Context, queryBuilder interfaces.QueryBuilder) (bool, error) {
	count, err := dao.Count(ctx, queryBuilder)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *MySQLDAO[T]) FindWithPagination(ctx context.Context, queryBuilder interfaces.QueryBuilder, page, pageSize int) ([]*T, int64, error) {
	// Get total count
	total, err := dao.Count(ctx, queryBuilder)
	if err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	offset := (page - 1) * pageSize
	paginatedQuery := queryBuilder.Clone().Limit(pageSize).Offset(offset)
	
	entities, err := dao.Find(ctx, paginatedQuery)
	if err != nil {
		return nil, 0, err
	}
	
	return entities, total, nil
}

// Helper methods
func (dao *MySQLDAO[T]) generateCacheKey(prefix string, id interface{}) string {
	return fmt.Sprintf("%s:%s:%v", dao.tableName, prefix, id)
}

func (dao *MySQLDAO[T]) generateQueryCacheKey(sql string, args []interface{}) string {
	return fmt.Sprintf("%s:query:%s:%v", dao.tableName, sql, args)
}

func (dao *MySQLDAO[T]) canCache(queryBuilder interfaces.QueryBuilder) bool {
	// Simple heuristic: cache if query is not too complex
	sql, _ := queryBuilder.ToSQL()
	return len(sql) < 1000 // Arbitrary limit
}

func (dao *MySQLDAO[T]) invalidateCache(ctx context.Context, pattern string) {
	if dao.cache != nil {
		key := fmt.Sprintf("%s:%s:*", dao.tableName, pattern)
		dao.cache.DeletePattern(ctx, key)
	}
}

func (dao *MySQLDAO[T]) invalidateAllCache(ctx context.Context) {
	if dao.cache != nil {
		pattern := fmt.Sprintf("%s:*", dao.tableName)
		dao.cache.DeletePattern(ctx, pattern)
	}
}

// WithTransaction executes operations within a transaction
func (dao *MySQLDAO[T]) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return dao.masterDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// Raw query support
func (dao *MySQLDAO[T]) Raw(ctx context.Context, sql string, args ...interface{}) ([]*T, error) {
	var entities []*T
	result := dao.slaveDB.WithContext(ctx).Raw(sql, args...).Scan(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

// Exec raw command
func (dao *MySQLDAO[T]) Exec(ctx context.Context, sql string, args ...interface{}) error {
	result := dao.masterDB.WithContext(ctx).Exec(sql, args...)
	return result.Error
}

// GetDB returns the underlying database connection
func (dao *MySQLDAO[T]) GetMasterDB() *gorm.DB {
	return dao.masterDB
}

func (dao *MySQLDAO[T]) GetSlaveDB() *gorm.DB {
	return dao.slaveDB
}

// Utility function to create a new MySQL DAO without generics for dynamic typing
func NewMySQLDAODynamic(masterDB, slaveDB *gorm.DB, cache interfaces.CacheDAO, tableName string) *MySQLDAO[interface{}] {
	return &MySQLDAO[interface{}]{
		masterDB:  masterDB,
		slaveDB:   slaveDB,
		cache:     cache,
		tableName: tableName,
		modelType: reflect.TypeOf((*interface{})(nil)).Elem(),
	}
}