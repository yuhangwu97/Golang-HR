package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gin-project/dao/interfaces"
)

// CacheStrategy defines different caching strategies
type CacheStrategy string

const (
	// WriteThrough writes to cache and database synchronously
	WriteThrough CacheStrategy = "write_through"
	
	// WriteBack writes to cache immediately, database asynchronously
	WriteBack CacheStrategy = "write_back"
	
	// WriteAround writes only to database, cache on read miss
	WriteAround CacheStrategy = "write_around"
	
	// CacheAside application manages cache and database separately
	CacheAside CacheStrategy = "cache_aside"
)

// EvictionPolicy defines cache eviction policies
type EvictionPolicy string

const (
	LRU  EvictionPolicy = "lru"  // Least Recently Used
	LFU  EvictionPolicy = "lfu"  // Least Frequently Used
	FIFO EvictionPolicy = "fifo" // First In First Out
	TTL  EvictionPolicy = "ttl"  // Time To Live
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	Strategy       CacheStrategy
	EvictionPolicy EvictionPolicy
	DefaultTTL     time.Duration
	MaxSize        int64
	WriteBackDelay time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	EnableMetrics  bool
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Strategy:       CacheAside,
		EvictionPolicy: LRU,
		DefaultTTL:     10 * time.Minute,
		MaxSize:        1000,
		WriteBackDelay: 5 * time.Second,
		ReadTimeout:    100 * time.Millisecond,
		WriteTimeout:   200 * time.Millisecond,
		EnableMetrics:  true,
	}
}

// CacheMetrics holds cache performance metrics
type CacheMetrics struct {
	Hits        int64
	Misses      int64
	Writes      int64
	Deletes     int64
	Evictions   int64
	Errors      int64
	mu          sync.RWMutex
}

func (m *CacheMetrics) IncrementHits() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits++
}

func (m *CacheMetrics) IncrementMisses() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Misses++
}

func (m *CacheMetrics) IncrementWrites() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Writes++
}

func (m *CacheMetrics) IncrementDeletes() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Deletes++
}

func (m *CacheMetrics) IncrementEvictions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Evictions++
}

func (m *CacheMetrics) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
}

func (m *CacheMetrics) GetStats() (hits, misses, writes, deletes, evictions, errors int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Hits, m.Misses, m.Writes, m.Deletes, m.Evictions, m.Errors
}

func (m *CacheMetrics) GetHitRatio() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Hits) / float64(total)
}

// CacheManager manages cache operations with different strategies
type CacheManager struct {
	cache      interfaces.CacheDAO
	config     *CacheConfig
	metrics    *CacheMetrics
	writeQueue chan *writeBackItem
	mu         sync.RWMutex
}

type writeBackItem struct {
	key   string
	value interface{}
	ttl   time.Duration
}

func NewCacheManager(cache interfaces.CacheDAO, config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}
	
	manager := &CacheManager{
		cache:      cache,
		config:     config,
		metrics:    &CacheMetrics{},
		writeQueue: make(chan *writeBackItem, 1000),
	}
	
	// Start write-back worker if strategy is WriteBack
	if config.Strategy == WriteBack {
		go manager.writeBackWorker()
	}
	
	return manager
}

func (cm *CacheManager) writeBackWorker() {
	ticker := time.NewTicker(cm.config.WriteBackDelay)
	defer ticker.Stop()
	
	batch := make([]*writeBackItem, 0, 100)
	
	for {
		select {
		case item := <-cm.writeQueue:
			batch = append(batch, item)
			if len(batch) >= 100 {
				cm.flushBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				cm.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

func (cm *CacheManager) flushBatch(batch []*writeBackItem) {
	ctx, cancel := context.WithTimeout(context.Background(), cm.config.WriteTimeout)
	defer cancel()
	
	for _, item := range batch {
		if err := cm.cache.Set(ctx, item.key, item.value, item.ttl); err != nil {
			cm.metrics.IncrementErrors()
		} else {
			cm.metrics.IncrementWrites()
		}
	}
}

// Get retrieves value from cache with strategy-aware logic
func (cm *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.ReadTimeout)
	defer cancel()
	
	err := cm.cache.Get(ctxWithTimeout, key, dest)
	if err == nil {
		cm.metrics.IncrementHits()
		return nil
	}
	
	cm.metrics.IncrementMisses()
	return err
}

// Set stores value in cache with strategy-aware logic
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = cm.config.DefaultTTL
	}
	
	switch cm.config.Strategy {
	case WriteThrough, CacheAside:
		ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.WriteTimeout)
		defer cancel()
		
		err := cm.cache.Set(ctxWithTimeout, key, value, ttl)
		if err != nil {
			cm.metrics.IncrementErrors()
			return err
		}
		cm.metrics.IncrementWrites()
		return nil
		
	case WriteBack:
		select {
		case cm.writeQueue <- &writeBackItem{key: key, value: value, ttl: ttl}:
			return nil
		default:
			// Queue is full, write immediately
			ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.WriteTimeout)
			defer cancel()
			
			err := cm.cache.Set(ctxWithTimeout, key, value, ttl)
			if err != nil {
				cm.metrics.IncrementErrors()
				return err
			}
			cm.metrics.IncrementWrites()
			return nil
		}
		
	case WriteAround:
		// Don't write to cache, only to database
		return nil
		
	default:
		return fmt.Errorf("unsupported cache strategy: %s", cm.config.Strategy)
	}
}

// Delete removes value from cache
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.WriteTimeout)
	defer cancel()
	
	err := cm.cache.Delete(ctxWithTimeout, key)
	if err != nil {
		cm.metrics.IncrementErrors()
		return err
	}
	
	cm.metrics.IncrementDeletes()
	return nil
}

// DeletePattern removes keys matching pattern
func (cm *CacheManager) DeletePattern(ctx context.Context, pattern string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.WriteTimeout)
	defer cancel()
	
	err := cm.cache.DeletePattern(ctxWithTimeout, pattern)
	if err != nil {
		cm.metrics.IncrementErrors()
		return err
	}
	
	cm.metrics.IncrementDeletes()
	return nil
}

// Exists checks if key exists in cache
func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cm.config.ReadTimeout)
	defer cancel()
	
	exists, err := cm.cache.Exists(ctxWithTimeout, key)
	if err != nil {
		cm.metrics.IncrementErrors()
		return false, err
	}
	
	if exists {
		cm.metrics.IncrementHits()
	} else {
		cm.metrics.IncrementMisses()
	}
	
	return exists, nil
}

// GetMetrics returns cache performance metrics
func (cm *CacheManager) GetMetrics() *CacheMetrics {
	return cm.metrics
}

// GetConfig returns cache configuration
func (cm *CacheManager) GetConfig() *CacheConfig {
	return cm.config
}

// UpdateConfig updates cache configuration
func (cm *CacheManager) UpdateConfig(config *CacheConfig) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = config
}

// Warm warms up cache with provided data
func (cm *CacheManager) Warm(ctx context.Context, data map[string]interface{}) error {
	for key, value := range data {
		if err := cm.Set(ctx, key, value, cm.config.DefaultTTL); err != nil {
			return fmt.Errorf("failed to warm cache for key %s: %w", key, err)
		}
	}
	return nil
}

// Clear clears all cache data (use with caution)
func (cm *CacheManager) Clear(ctx context.Context) error {
	return cm.DeletePattern(ctx, "*")
}

// HealthCheck checks cache health
func (cm *CacheManager) HealthCheck(ctx context.Context) error {
	testKey := "__health_check__"
	testValue := "ok"
	
	// Test write
	if err := cm.Set(ctx, testKey, testValue, time.Minute); err != nil {
		return fmt.Errorf("cache write health check failed: %w", err)
	}
	
	// Test read
	var result string
	if err := cm.Get(ctx, testKey, &result); err != nil {
		return fmt.Errorf("cache read health check failed: %w", err)
	}
	
	if result != testValue {
		return fmt.Errorf("cache health check: expected %s, got %s", testValue, result)
	}
	
	// Clean up
	if err := cm.Delete(ctx, testKey); err != nil {
		return fmt.Errorf("cache delete health check failed: %w", err)
	}
	
	return nil
}

// Multi-level Cache Manager for L1/L2 caching
type MultiLevelCacheManager struct {
	l1Cache *CacheManager // Fast local cache (e.g., in-memory)
	l2Cache *CacheManager // Slower distributed cache (e.g., Redis)
	config  *CacheConfig
}

func NewMultiLevelCacheManager(l1Cache, l2Cache interfaces.CacheDAO, config *CacheConfig) *MultiLevelCacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}
	
	return &MultiLevelCacheManager{
		l1Cache: NewCacheManager(l1Cache, config),
		l2Cache: NewCacheManager(l2Cache, config),
		config:  config,
	}
}

func (ml *MultiLevelCacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	// Try L1 cache first
	err := ml.l1Cache.Get(ctx, key, dest)
	if err == nil {
		return nil
	}
	
	// Try L2 cache
	err = ml.l2Cache.Get(ctx, key, dest)
	if err == nil {
		// Populate L1 cache
		ml.l1Cache.Set(ctx, key, dest, ml.config.DefaultTTL)
		return nil
	}
	
	return err
}

func (ml *MultiLevelCacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Write to both caches
	l1Err := ml.l1Cache.Set(ctx, key, value, ttl)
	l2Err := ml.l2Cache.Set(ctx, key, value, ttl)
	
	if l1Err != nil && l2Err != nil {
		return fmt.Errorf("both L1 and L2 cache writes failed: L1=%v, L2=%v", l1Err, l2Err)
	}
	
	return nil
}

func (ml *MultiLevelCacheManager) Delete(ctx context.Context, key string) error {
	l1Err := ml.l1Cache.Delete(ctx, key)
	l2Err := ml.l2Cache.Delete(ctx, key)
	
	if l1Err != nil && l2Err != nil {
		return fmt.Errorf("both L1 and L2 cache deletes failed: L1=%v, L2=%v", l1Err, l2Err)
	}
	
	return nil
}

func (ml *MultiLevelCacheManager) GetMetrics() (l1Metrics, l2Metrics *CacheMetrics) {
	return ml.l1Cache.GetMetrics(), ml.l2Cache.GetMetrics()
}

func (cm *CacheManager) GetRedisDAO() interfaces.CacheDAO {
	return cm.cache
}