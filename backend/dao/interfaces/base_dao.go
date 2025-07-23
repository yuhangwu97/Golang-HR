package interfaces

import (
	"context"
	"time"
)

// BaseDAO defines the basic database operations interface
type BaseDAO[T any] interface {
	// CRUD Operations
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	Find(ctx context.Context, query QueryBuilder) ([]*T, error)
	Update(ctx context.Context, entity *T, query QueryBuilder) error
	Delete(ctx context.Context, query QueryBuilder) error
	
	// Batch Operations
	CreateMany(ctx context.Context, entities []*T) error
	UpdateMany(ctx context.Context, updates map[string]interface{}, query QueryBuilder) error
	DeleteMany(ctx context.Context, query QueryBuilder) error
	
	// Query Operations
	Count(ctx context.Context, query QueryBuilder) (int64, error)
	Exists(ctx context.Context, query QueryBuilder) (bool, error)
	
	// Pagination
	FindWithPagination(ctx context.Context, query QueryBuilder, page, pageSize int) ([]*T, int64, error)
}

// CacheDAO defines caching operations interface
type CacheDAO interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
	
	// Advanced Cache Operations
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Hash Operations
	HSet(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string, dest interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	
	// List Operations
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string, dest interface{}) error
	RPop(ctx context.Context, key string, dest interface{}) error
	LLen(ctx context.Context, key string) (int64, error)
	
	// Set Operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
}

// TransactionDAO defines transaction operations interface
type TransactionDAO interface {
	Begin(ctx context.Context) (Transaction, error)
	WithTransaction(ctx context.Context, fn func(tx Transaction) error) error
}

// Transaction represents a database transaction
type Transaction interface {
	Commit() error
	Rollback() error
	GetContext() context.Context
}

// QueryBuilder defines query building interface
type QueryBuilder interface {
	Where(field string, operator string, value interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNotIn(field string, values []interface{}) QueryBuilder
	WhereBetween(field string, start, end interface{}) QueryBuilder
	WhereNull(field string) QueryBuilder
	WhereNotNull(field string) QueryBuilder
	OrWhere(field string, operator string, value interface{}) QueryBuilder
	
	// Joins
	Join(table, condition string) QueryBuilder
	LeftJoin(table, condition string) QueryBuilder
	RightJoin(table, condition string) QueryBuilder
	InnerJoin(table, condition string) QueryBuilder
	
	// Ordering and Grouping
	OrderBy(field string, direction string) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Having(condition string, args ...interface{}) QueryBuilder
	
	// Limits and Offsets
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	
	// Aggregation
	Select(fields ...string) QueryBuilder
	Distinct() QueryBuilder
	
	// Raw SQL
	Raw(sql string, args ...interface{}) QueryBuilder
	
	// Build methods
	ToSQL() (string, []interface{})
	Reset() QueryBuilder
	Clone() QueryBuilder
}

// DataSourceManager manages multiple data sources
type DataSourceManager interface {
	GetMasterDB() interface{}
	GetSlaveDB() interface{}
	GetCache() CacheDAO
	GetReadOnlyDB() interface{}
	
	// Health check
	HealthCheck(ctx context.Context) error
	
	// Connection management
	Close() error
}

// DAOFactory creates DAO instances
type DAOFactory interface {
	CreateDAO(entityType interface{}) interface{}
	CreateUserDAO() UserDAO
	// Add more specific DAO creators as needed
}

// UserDAO specific interface for user operations
type UserDAO interface {
	BaseDAO[any] // Will be replaced with actual User type
	FindByEmail(ctx context.Context, email string) (interface{}, error)
	FindByRole(ctx context.Context, role string) ([]interface{}, error)
	UpdatePassword(ctx context.Context, userID interface{}, hashedPassword string) error
}