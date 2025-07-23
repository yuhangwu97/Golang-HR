package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ConnectionPoolConfig holds connection pool configuration
type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	HealthCheckInterval time.Duration
	ConnectTimeout     time.Duration
	QueryTimeout       time.Duration
}

// DefaultConnectionPoolConfig returns default connection pool configuration
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxOpenConns:        100,
		MaxIdleConns:        10,
		ConnMaxLifetime:     time.Hour,
		ConnMaxIdleTime:     10 * time.Minute,
		HealthCheckInterval: 30 * time.Second,
		ConnectTimeout:      5 * time.Second,
		QueryTimeout:        30 * time.Second,
	}
}

// ConnectionPool manages database connections
type ConnectionPool struct {
	masterDB   *gorm.DB
	slaveDBs   []*gorm.DB
	config     *ConnectionPoolConfig
	mu         sync.RWMutex
	closed     bool
	healthTicker *time.Ticker
	stats      *PoolStats
}

// PoolStats holds connection pool statistics
type PoolStats struct {
	OpenConnections     int
	InUseConnections    int
	IdleConnections     int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
	MaxOpenConnections int
	mu                 sync.RWMutex
}

// GetStats returns a copy of the current pool statistics
func (ps *PoolStats) GetStats() PoolStats {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return *ps
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(masterDB *gorm.DB, slaveDBs []*gorm.DB, config *ConnectionPoolConfig) *ConnectionPool {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}
	
	pool := &ConnectionPool{
		masterDB: masterDB,
		slaveDBs: slaveDBs,
		config:   config,
		stats:    &PoolStats{MaxOpenConnections: config.MaxOpenConns},
	}
	
	// Configure connection pools
	pool.configurePools()
	
	// Start health check
	pool.startHealthCheck()
	
	return pool
}

// configurePools configures the underlying database connection pools
func (cp *ConnectionPool) configurePools() {
	// Configure master DB
	if sqlDB, err := cp.masterDB.DB(); err == nil {
		sqlDB.SetMaxOpenConns(cp.config.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cp.config.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cp.config.ConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(cp.config.ConnMaxIdleTime)
	}
	
	// Configure slave DBs
	for _, slaveDB := range cp.slaveDBs {
		if sqlDB, err := slaveDB.DB(); err == nil {
			sqlDB.SetMaxOpenConns(cp.config.MaxOpenConns)
			sqlDB.SetMaxIdleConns(cp.config.MaxIdleConns)
			sqlDB.SetConnMaxLifetime(cp.config.ConnMaxLifetime)
			sqlDB.SetConnMaxIdleTime(cp.config.ConnMaxIdleTime)
		}
	}
}

// GetMasterDB returns the master database connection
func (cp *ConnectionPool) GetMasterDB() *gorm.DB {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	
	if cp.closed {
		return nil
	}
	
	return cp.masterDB
}

// GetSlaveDB returns a slave database connection using round-robin
func (cp *ConnectionPool) GetSlaveDB() *gorm.DB {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	
	if cp.closed {
		return nil
	}
	
	if len(cp.slaveDBs) == 0 {
		return cp.masterDB
	}
	
	// Simple round-robin selection
	// In production, you might want more sophisticated load balancing
	index := time.Now().UnixNano() % int64(len(cp.slaveDBs))
	return cp.slaveDBs[index]
}

// GetReadOnlyDB returns a read-only database connection
func (cp *ConnectionPool) GetReadOnlyDB() *gorm.DB {
	return cp.GetSlaveDB()
}

// WithContext returns a database connection with context
func (cp *ConnectionPool) WithContext(ctx context.Context, readOnly bool) *gorm.DB {
	var db *gorm.DB
	
	if readOnly {
		db = cp.GetSlaveDB()
	} else {
		db = cp.GetMasterDB()
	}
	
	if db == nil {
		return nil
	}
	
	return db.WithContext(ctx)
}

// ExecuteWithTimeout executes a function with timeout
func (cp *ConnectionPool) ExecuteWithTimeout(ctx context.Context, timeout time.Duration, readOnly bool, fn func(*gorm.DB) error) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	db := cp.WithContext(timeoutCtx, readOnly)
	if db == nil {
		return fmt.Errorf("database connection is unavailable")
	}
	
	return fn(db)
}

// HealthCheck checks the health of all database connections
func (cp *ConnectionPool) HealthCheck(ctx context.Context) error {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	
	if cp.closed {
		return fmt.Errorf("connection pool is closed")
	}
	
	// Check master database
	if err := cp.checkConnection(ctx, cp.masterDB, "master"); err != nil {
		return err
	}
	
	// Check slave databases
	for i, slaveDB := range cp.slaveDBs {
		if err := cp.checkConnection(ctx, slaveDB, fmt.Sprintf("slave-%d", i)); err != nil {
			return err
		}
	}
	
	return nil
}

// checkConnection checks a single database connection
func (cp *ConnectionPool) checkConnection(ctx context.Context, db *gorm.DB, name string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("%s database error: %w", name, err)
	}
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("%s database ping failed: %w", name, err)
	}
	
	return nil
}

// startHealthCheck starts the health check routine
func (cp *ConnectionPool) startHealthCheck() {
	cp.healthTicker = time.NewTicker(cp.config.HealthCheckInterval)
	
	go func() {
		for range cp.healthTicker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := cp.HealthCheck(ctx); err != nil {
				// Log error (in a real system, you'd use a proper logger)
				fmt.Printf("Health check failed: %v\n", err)
			}
			cancel()
		}
	}()
}

// GetStats returns connection pool statistics
func (cp *ConnectionPool) GetStats() *PoolStats {
	cp.updateStats()
	return cp.stats
}

// updateStats updates the pool statistics
func (cp *ConnectionPool) updateStats() {
	cp.stats.mu.Lock()
	defer cp.stats.mu.Unlock()
	
	// Update master DB stats
	if sqlDB, err := cp.masterDB.DB(); err == nil {
		stats := sqlDB.Stats()
		cp.stats.OpenConnections = stats.OpenConnections
		cp.stats.InUseConnections = stats.InUse
		cp.stats.IdleConnections = stats.Idle
		cp.stats.WaitCount = stats.WaitCount
		cp.stats.WaitDuration = stats.WaitDuration
		cp.stats.MaxIdleClosed = stats.MaxIdleClosed
		cp.stats.MaxLifetimeClosed = stats.MaxLifetimeClosed
	}
}

// Close closes the connection pool
func (cp *ConnectionPool) Close() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	
	if cp.closed {
		return nil
	}
	
	cp.closed = true
	
	// Stop health check
	if cp.healthTicker != nil {
		cp.healthTicker.Stop()
	}
	
	var errors []error
	
	// Close master database
	if sqlDB, err := cp.masterDB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close master DB: %w", err))
		}
	}
	
	// Close slave databases
	for i, slaveDB := range cp.slaveDBs {
		if sqlDB, err := slaveDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close slave DB %d: %w", i, err))
			}
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("errors closing connection pool: %v", errors)
	}
	
	return nil
}

// LoadBalancer manages load balancing across slave databases
type LoadBalancer struct {
	strategy LoadBalancingStrategy
	slaves   []*gorm.DB
	weights  []int
	mu       sync.RWMutex
	current  int
}

// LoadBalancingStrategy defines load balancing strategies
type LoadBalancingStrategy int

const (
	RoundRobin LoadBalancingStrategy = iota
	WeightedRoundRobin
	LeastConnections
	Random
)

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(strategy LoadBalancingStrategy, slaves []*gorm.DB, weights []int) *LoadBalancer {
	if weights == nil {
		weights = make([]int, len(slaves))
		for i := range weights {
			weights[i] = 1 // Equal weights
		}
	}
	
	return &LoadBalancer{
		strategy: strategy,
		slaves:   slaves,
		weights:  weights,
	}
}

// GetConnection returns a database connection based on the load balancing strategy
func (lb *LoadBalancer) GetConnection() *gorm.DB {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	if len(lb.slaves) == 0 {
		return nil
	}
	
	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin()
	case WeightedRoundRobin:
		return lb.weightedRoundRobin()
	case LeastConnections:
		return lb.leastConnections()
	case Random:
		return lb.random()
	default:
		return lb.roundRobin()
	}
}

// roundRobin implements round-robin load balancing
func (lb *LoadBalancer) roundRobin() *gorm.DB {
	db := lb.slaves[lb.current]
	lb.current = (lb.current + 1) % len(lb.slaves)
	return db
}

// weightedRoundRobin implements weighted round-robin load balancing
func (lb *LoadBalancer) weightedRoundRobin() *gorm.DB {
	// Simplified weighted round-robin implementation
	totalWeight := 0
	for _, weight := range lb.weights {
		totalWeight += weight
	}
	
	target := lb.current % totalWeight
	currentWeight := 0
	
	for i, weight := range lb.weights {
		currentWeight += weight
		if target < currentWeight {
			lb.current++
			return lb.slaves[i]
		}
	}
	
	lb.current++
	return lb.slaves[0]
}

// leastConnections implements least connections load balancing
func (lb *LoadBalancer) leastConnections() *gorm.DB {
	minConnections := int(^uint(0) >> 1) // Max int
	var selectedDB *gorm.DB
	
	for _, db := range lb.slaves {
		if sqlDB, err := db.DB(); err == nil {
			connections := sqlDB.Stats().InUse
			if connections < minConnections {
				minConnections = connections
				selectedDB = db
			}
		}
	}
	
	if selectedDB == nil {
		return lb.slaves[0]
	}
	
	return selectedDB
}

// random implements random load balancing
func (lb *LoadBalancer) random() *gorm.DB {
	index := time.Now().UnixNano() % int64(len(lb.slaves))
	return lb.slaves[index]
}

// ConnectionMonitor monitors database connections
type ConnectionMonitor struct {
	pools    []*ConnectionPool
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan bool
	mu       sync.RWMutex
}

// NewConnectionMonitor creates a new connection monitor
func NewConnectionMonitor(interval time.Duration) *ConnectionMonitor {
	return &ConnectionMonitor{
		pools:    make([]*ConnectionPool, 0),
		interval: interval,
		stopCh:   make(chan bool),
	}
}

// AddPool adds a connection pool to monitor
func (cm *ConnectionMonitor) AddPool(pool *ConnectionPool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.pools = append(cm.pools, pool)
}

// Start starts the connection monitor
func (cm *ConnectionMonitor) Start() {
	cm.ticker = time.NewTicker(cm.interval)
	
	go func() {
		for {
			select {
			case <-cm.ticker.C:
				cm.monitor()
			case <-cm.stopCh:
				return
			}
		}
	}()
}

// Stop stops the connection monitor
func (cm *ConnectionMonitor) Stop() {
	if cm.ticker != nil {
		cm.ticker.Stop()
	}
	cm.stopCh <- true
}

// monitor performs the monitoring check
func (cm *ConnectionMonitor) monitor() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	for i, pool := range cm.pools {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		if err := pool.HealthCheck(ctx); err != nil {
			fmt.Printf("Pool %d health check failed: %v\n", i, err)
		}
		
		stats := pool.GetStats()
		fmt.Printf("Pool %d stats: Open=%d, InUse=%d, Idle=%d\n", 
			i, stats.OpenConnections, stats.InUseConnections, stats.IdleConnections)
		
		cancel()
	}
}