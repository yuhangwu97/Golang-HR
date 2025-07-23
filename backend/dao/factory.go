package dao

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gin-project/dao/cache"
	"gin-project/dao/interfaces"
	"gin-project/dao/mysql"
	"gin-project/dao/pool"
	redisDAO "gin-project/dao/redis"
	"gin-project/dao/transaction"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// DataSourceType represents different data source types
type DataSourceType string

const (
	MySQL DataSourceType = "mysql"
	Redis DataSourceType = "redis"
	Cache DataSourceType = "cache"
)

// DataSourceConfig holds configuration for data sources
type DataSourceConfig struct {
	Type     DataSourceType
	Master   interface{}
	Slaves   []interface{}
	ReadOnly interface{}
	Cache    interface{}
}

// DataSourceManager implements the DataSourceManager interface
type DataSourceManager struct {
	connectionPool   *pool.ConnectionPool
	cacheManager     *cache.CacheManager
	transactionMgr   *transaction.TransactionManager
	mu               sync.RWMutex
	config           *DataSourceConfig
}

func NewDataSourceManager(db *gorm.DB, redisClient *redis.Client) *DataSourceManager {
	// Create connection pool
	slaveDBs := []*gorm.DB{db} // In production, you'd have separate slave connections
	connectionPool := pool.NewConnectionPool(db, slaveDBs, pool.DefaultConnectionPoolConfig())
	
	// Create cache manager with Redis
	redisDAOInstance := redisDAO.NewRedisDAO(redisClient)
	cacheManager := cache.NewCacheManager(redisDAOInstance, cache.DefaultCacheConfig())
	
	// Create transaction manager
	transactionMgr := transaction.NewTransactionManager(db, 30*time.Second)
	
	return &DataSourceManager{
		connectionPool: connectionPool,
		cacheManager:   cacheManager,
		transactionMgr: transactionMgr,
		config: &DataSourceConfig{
			Type:   MySQL,
			Master: db,
			Slaves: []interface{}{db},
			Cache:  redisClient,
		},
	}
}

func (dsm *DataSourceManager) GetMasterDB() interface{} {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.connectionPool.GetMasterDB()
}

func (dsm *DataSourceManager) GetSlaveDB() interface{} {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.connectionPool.GetSlaveDB()
}

func (dsm *DataSourceManager) GetReadOnlyDB() interface{} {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.connectionPool.GetReadOnlyDB()
}

func (dsm *DataSourceManager) GetCache() interfaces.CacheDAO {
	return dsm.cacheManager.GetRedisDAO()
}

func (dsm *DataSourceManager) GetTransactionManager() *transaction.TransactionManager {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.transactionMgr
}

func (dsm *DataSourceManager) GetConnectionPool() *pool.ConnectionPool {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.connectionPool
}

func (dsm *DataSourceManager) GetCacheManager() *cache.CacheManager {
	dsm.mu.RLock()
	defer dsm.mu.RUnlock()
	return dsm.cacheManager
}

func (dsm *DataSourceManager) HealthCheck(ctx context.Context) error {
	// Check connection pool
	if err := dsm.connectionPool.HealthCheck(ctx); err != nil {
		return fmt.Errorf("connection pool health check failed: %w", err)
	}
	
	// Check cache manager
	if err := dsm.cacheManager.HealthCheck(ctx); err != nil {
		return fmt.Errorf("cache manager health check failed: %w", err)
	}
	
	return nil
}

func (dsm *DataSourceManager) Close() error {
	var errs []error
	
	// Close connection pool
	if err := dsm.connectionPool.Close(); err != nil {
		errs = append(errs, fmt.Errorf("connection pool close error: %w", err))
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("errors closing data source manager: %v", errs)
	}
	
	return nil
}

// DAOFactory implements the factory pattern for creating DAO instances
type DAOFactory struct {
	dataSourceManager interfaces.DataSourceManager
	daoInstances      map[string]interface{}
	mu                sync.RWMutex
}

func NewDAOFactory(dsm interfaces.DataSourceManager) *DAOFactory {
	return &DAOFactory{
		dataSourceManager: dsm,
		daoInstances:      make(map[string]interface{}),
	}
}

func (f *DAOFactory) CreateDAO(entityType interface{}) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	typeName := fmt.Sprintf("%T", entityType)
	
	if dao, exists := f.daoInstances[typeName]; exists {
		return dao
	}
	
	// Create MySQL DAO for the entity type with enhanced architecture
	masterDB := f.dataSourceManager.GetMasterDB().(*gorm.DB)
	slaveDB := f.dataSourceManager.GetSlaveDB().(*gorm.DB)
	cache := f.dataSourceManager.GetCache()
	
	// Create enhanced MySQL DAO
	dao := mysql.NewMySQLDAO[interface{}](masterDB, slaveDB, cache, typeName)
	f.daoInstances[typeName] = dao
	
	return dao
}

func (f *DAOFactory) CreateDAOWithTransaction(entityType interface{}, tx *gorm.DB) interface{} {
	typeName := fmt.Sprintf("%T", entityType)
	cache := f.dataSourceManager.GetCache()
	
	// Create DAO with transaction
	dao := mysql.NewMySQLDAO[interface{}](tx, tx, cache, typeName)
	return dao
}

func (f *DAOFactory) CreateUserDAO() interfaces.UserDAO {
	return f.CreateDAO(&struct{}{}).(interfaces.UserDAO)
}

func (f *DAOFactory) CreateCacheDAO() interfaces.CacheDAO {
	return f.dataSourceManager.GetCache()
}

func (f *DAOFactory) CreateTransactionManager() *transaction.TransactionManager {
	if dsm, ok := f.dataSourceManager.(*DataSourceManager); ok {
		return dsm.GetTransactionManager()
	}
	return nil
}

func (f *DAOFactory) CreateConnectionPool() *pool.ConnectionPool {
	if dsm, ok := f.dataSourceManager.(*DataSourceManager); ok {
		return dsm.GetConnectionPool()
	}
	return nil
}

func (f *DAOFactory) CreateCacheManager() *cache.CacheManager {
	if dsm, ok := f.dataSourceManager.(*DataSourceManager); ok {
		return dsm.GetCacheManager()
	}
	return nil
}

// Singleton instances
var (
	dataSourceManagerInstance *DataSourceManager
	daoFactoryInstance        *DAOFactory
	once                      sync.Once
)

// GetDataSourceManager returns the singleton data source manager
func GetDataSourceManager(db *gorm.DB, redisClient *redis.Client) *DataSourceManager {
	once.Do(func() {
		dataSourceManagerInstance = NewDataSourceManager(db, redisClient)
		daoFactoryInstance = NewDAOFactory(dataSourceManagerInstance)
	})
	return dataSourceManagerInstance
}

// GetDAOFactory returns the singleton DAO factory
func GetDAOFactory(db *gorm.DB, redisClient *redis.Client) *DAOFactory {
	once.Do(func() {
		dataSourceManagerInstance = NewDataSourceManager(db, redisClient)
		daoFactoryInstance = NewDAOFactory(dataSourceManagerInstance)
	})
	return daoFactoryInstance
}

// DAORegistry manages registered DAO types
type DAORegistry struct {
	registeredTypes map[string]func() interface{}
	mu              sync.RWMutex
}

func NewDAORegistry() *DAORegistry {
	return &DAORegistry{
		registeredTypes: make(map[string]func() interface{}),
	}
}

func (r *DAORegistry) Register(typeName string, factory func() interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registeredTypes[typeName] = factory
}

func (r *DAORegistry) Create(typeName string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	factory, exists := r.registeredTypes[typeName]
	if !exists {
		return nil, fmt.Errorf("DAO type %s not registered", typeName)
	}
	
	return factory(), nil
}

func (r *DAORegistry) IsRegistered(typeName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.registeredTypes[typeName]
	return exists
}

// Global registry instance
var globalRegistry = NewDAORegistry()

func GetDAORegistry() *DAORegistry {
	return globalRegistry
}