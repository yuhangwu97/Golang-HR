package transaction

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gin-project/dao/interfaces"

	"gorm.io/gorm"
)

// TransactionManager manages database transactions
type TransactionManager struct {
	db      *gorm.DB
	timeout time.Duration
	mu      sync.RWMutex
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB, timeout time.Duration) *TransactionManager {
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}
	
	return &TransactionManager{
		db:      db,
		timeout: timeout,
	}
}

// Transaction represents a database transaction
type Transaction struct {
	tx      *gorm.DB
	ctx     context.Context
	cancel  context.CancelFunc
	manager *TransactionManager
}

// GetContext returns the transaction context
func (t *Transaction) GetContext() context.Context {
	return t.ctx
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	defer t.cancel()
	return t.tx.Commit().Error
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	defer t.cancel()
	return t.tx.Rollback().Error
}

// GetDB returns the transaction database instance
func (t *Transaction) GetDB() *gorm.DB {
	return t.tx
}

// Begin starts a new transaction
func (tm *TransactionManager) Begin(ctx context.Context) (interfaces.Transaction, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	// Create context with timeout
	txCtx, cancel := context.WithTimeout(ctx, tm.timeout)
	
	// Begin transaction
	tx := tm.db.WithContext(txCtx).Begin()
	if tx.Error != nil {
		cancel()
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	
	return &Transaction{
		tx:      tx,
		ctx:     txCtx,
		cancel:  cancel,
		manager: tm,
	}, nil
}

// WithTransaction executes a function within a transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx interfaces.Transaction) error) error {
	tx, err := tm.Begin(ctx)
	if err != nil {
		return err
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}
	
	return tx.Commit()
}

// TransactionContext provides transaction-aware context
type TransactionContext struct {
	ctx context.Context
	tx  *gorm.DB
}

// NewTransactionContext creates a new transaction context
func NewTransactionContext(ctx context.Context, tx *gorm.DB) *TransactionContext {
	return &TransactionContext{
		ctx: ctx,
		tx:  tx,
	}
}

// Context returns the underlying context
func (tc *TransactionContext) Context() context.Context {
	return tc.ctx
}

// DB returns the transaction database
func (tc *TransactionContext) DB() *gorm.DB {
	return tc.tx
}

// DistributedTransactionManager manages distributed transactions
type DistributedTransactionManager struct {
	localTxManager *TransactionManager
	participants   []TransactionParticipant
	mu             sync.RWMutex
}

// TransactionParticipant represents a participant in a distributed transaction
type TransactionParticipant interface {
	Prepare(ctx context.Context, txID string) error
	Commit(ctx context.Context, txID string) error
	Rollback(ctx context.Context, txID string) error
}

// NewDistributedTransactionManager creates a new distributed transaction manager
func NewDistributedTransactionManager(localTxManager *TransactionManager) *DistributedTransactionManager {
	return &DistributedTransactionManager{
		localTxManager: localTxManager,
		participants:   make([]TransactionParticipant, 0),
	}
}

// AddParticipant adds a participant to the distributed transaction
func (dtm *DistributedTransactionManager) AddParticipant(participant TransactionParticipant) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()
	dtm.participants = append(dtm.participants, participant)
}

// ExecuteDistributedTransaction executes a distributed transaction using 2PC
func (dtm *DistributedTransactionManager) ExecuteDistributedTransaction(ctx context.Context, txID string, localFn func(tx interfaces.Transaction) error) error {
	dtm.mu.RLock()
	participants := make([]TransactionParticipant, len(dtm.participants))
	copy(participants, dtm.participants)
	dtm.mu.RUnlock()
	
	// Phase 1: Prepare
	for _, participant := range participants {
		if err := participant.Prepare(ctx, txID); err != nil {
			// Rollback all participants
			for _, p := range participants {
				p.Rollback(ctx, txID)
			}
			return fmt.Errorf("prepare phase failed: %w", err)
		}
	}
	
	// Execute local transaction
	err := dtm.localTxManager.WithTransaction(ctx, localFn)
	if err != nil {
		// Rollback all participants
		for _, participant := range participants {
			participant.Rollback(ctx, txID)
		}
		return fmt.Errorf("local transaction failed: %w", err)
	}
	
	// Phase 2: Commit
	for _, participant := range participants {
		if err := participant.Commit(ctx, txID); err != nil {
			// This is problematic - some participants might have committed
			// In a real system, you'd need compensation logic here
			return fmt.Errorf("commit phase failed: %w", err)
		}
	}
	
	return nil
}

// TransactionStateMachine manages transaction state
type TransactionStateMachine struct {
	state     TransactionState
	history   []TransactionState
	mu        sync.RWMutex
	callbacks map[TransactionState]func()
}

// TransactionState represents the state of a transaction
type TransactionState int

const (
	StateInitial TransactionState = iota
	StateActive
	StatePrepared
	StateCommitted
	StateRolledBack
	StateAborted
)

func (s TransactionState) String() string {
	switch s {
	case StateInitial:
		return "Initial"
	case StateActive:
		return "Active"
	case StatePrepared:
		return "Prepared"
	case StateCommitted:
		return "Committed"
	case StateRolledBack:
		return "RolledBack"
	case StateAborted:
		return "Aborted"
	default:
		return "Unknown"
	}
}

// NewTransactionStateMachine creates a new transaction state machine
func NewTransactionStateMachine() *TransactionStateMachine {
	return &TransactionStateMachine{
		state:     StateInitial,
		history:   make([]TransactionState, 0),
		callbacks: make(map[TransactionState]func()),
	}
}

// GetState returns the current state
func (tsm *TransactionStateMachine) GetState() TransactionState {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	return tsm.state
}

// SetState sets the transaction state
func (tsm *TransactionStateMachine) SetState(newState TransactionState) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	
	if !tsm.isValidTransition(tsm.state, newState) {
		return fmt.Errorf("invalid state transition from %s to %s", tsm.state, newState)
	}
	
	tsm.history = append(tsm.history, tsm.state)
	tsm.state = newState
	
	// Execute callback if registered
	if callback, exists := tsm.callbacks[newState]; exists {
		go callback()
	}
	
	return nil
}

// RegisterCallback registers a callback for a state
func (tsm *TransactionStateMachine) RegisterCallback(state TransactionState, callback func()) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	tsm.callbacks[state] = callback
}

// GetHistory returns the state history
func (tsm *TransactionStateMachine) GetHistory() []TransactionState {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	
	history := make([]TransactionState, len(tsm.history))
	copy(history, tsm.history)
	return history
}

// isValidTransition checks if a state transition is valid
func (tsm *TransactionStateMachine) isValidTransition(from, to TransactionState) bool {
	validTransitions := map[TransactionState][]TransactionState{
		StateInitial:    {StateActive},
		StateActive:     {StatePrepared, StateCommitted, StateRolledBack, StateAborted},
		StatePrepared:   {StateCommitted, StateRolledBack},
		StateCommitted:  {},
		StateRolledBack: {},
		StateAborted:    {},
	}
	
	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}
	
	for _, allowedState := range allowed {
		if allowedState == to {
			return true
		}
	}
	
	return false
}

// TransactionPool manages a pool of transactions
type TransactionPool struct {
	manager     *TransactionManager
	pool        chan interfaces.Transaction
	maxSize     int
	timeout     time.Duration
	mu          sync.RWMutex
	activeCount int
}

// NewTransactionPool creates a new transaction pool
func NewTransactionPool(manager *TransactionManager, maxSize int, timeout time.Duration) *TransactionPool {
	return &TransactionPool{
		manager: manager,
		pool:    make(chan interfaces.Transaction, maxSize),
		maxSize: maxSize,
		timeout: timeout,
	}
}

// Get gets a transaction from the pool
func (tp *TransactionPool) Get(ctx context.Context) (interfaces.Transaction, error) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	select {
	case tx := <-tp.pool:
		tp.activeCount++
		return tx, nil
	default:
		// Pool is empty, create new transaction
		tx, err := tp.manager.Begin(ctx)
		if err != nil {
			return nil, err
		}
		tp.activeCount++
		return tx, nil
	}
}

// Put returns a transaction to the pool
func (tp *TransactionPool) Put(tx interfaces.Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	tp.activeCount--
	
	select {
	case tp.pool <- tx:
		// Transaction returned to pool
	default:
		// Pool is full, commit and discard
		tx.Commit()
	}
}

// GetStats returns pool statistics
func (tp *TransactionPool) GetStats() (poolSize, activeCount int) {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return len(tp.pool), tp.activeCount
}

// Close closes the transaction pool
func (tp *TransactionPool) Close() error {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	close(tp.pool)
	
	// Commit all pooled transactions
	for tx := range tp.pool {
		tx.Commit()
	}
	
	return nil
}

// TransactionInterceptor intercepts and logs transaction operations
type TransactionInterceptor struct {
	next   interfaces.TransactionDAO
	logger func(string, ...interface{})
}

// NewTransactionInterceptor creates a new transaction interceptor
func NewTransactionInterceptor(next interfaces.TransactionDAO, logger func(string, ...interface{})) *TransactionInterceptor {
	return &TransactionInterceptor{
		next:   next,
		logger: logger,
	}
}

// Begin begins a transaction with logging
func (ti *TransactionInterceptor) Begin(ctx context.Context) (interfaces.Transaction, error) {
	ti.logger("Beginning transaction")
	
	tx, err := ti.next.Begin(ctx)
	if err != nil {
		ti.logger("Failed to begin transaction: %v", err)
		return nil, err
	}
	
	ti.logger("Transaction begun successfully")
	return &interceptedTransaction{
		Transaction: tx,
		logger:      ti.logger,
	}, nil
}

// WithTransaction executes a function within a transaction with logging
func (ti *TransactionInterceptor) WithTransaction(ctx context.Context, fn func(tx interfaces.Transaction) error) error {
	ti.logger("Executing transaction")
	
	start := time.Now()
	err := ti.next.WithTransaction(ctx, fn)
	duration := time.Since(start)
	
	if err != nil {
		ti.logger("Transaction failed after %v: %v", duration, err)
		return err
	}
	
	ti.logger("Transaction completed successfully in %v", duration)
	return nil
}

// interceptedTransaction wraps a transaction with logging
type interceptedTransaction struct {
	interfaces.Transaction
	logger func(string, ...interface{})
}

// Commit commits the transaction with logging
func (it *interceptedTransaction) Commit() error {
	it.logger("Committing transaction")
	
	err := it.Transaction.Commit()
	if err != nil {
		it.logger("Failed to commit transaction: %v", err)
		return err
	}
	
	it.logger("Transaction committed successfully")
	return nil
}

// Rollback rolls back the transaction with logging
func (it *interceptedTransaction) Rollback() error {
	it.logger("Rolling back transaction")
	
	err := it.Transaction.Rollback()
	if err != nil {
		it.logger("Failed to rollback transaction: %v", err)
		return err
	}
	
	it.logger("Transaction rolled back successfully")
	return nil
}

// TransactionRetryManager manages transaction retries
type TransactionRetryManager struct {
	maxRetries    int
	retryInterval time.Duration
	backoffFactor float64
}

// NewTransactionRetryManager creates a new transaction retry manager
func NewTransactionRetryManager(maxRetries int, retryInterval time.Duration, backoffFactor float64) *TransactionRetryManager {
	return &TransactionRetryManager{
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		backoffFactor: backoffFactor,
	}
}

// ExecuteWithRetry executes a transaction with retry logic
func (trm *TransactionRetryManager) ExecuteWithRetry(ctx context.Context, txManager interfaces.TransactionDAO, fn func(tx interfaces.Transaction) error) error {
	var lastErr error
	interval := trm.retryInterval
	
	for i := 0; i <= trm.maxRetries; i++ {
		err := txManager.WithTransaction(ctx, fn)
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Don't retry on the last attempt
		if i == trm.maxRetries {
			break
		}
		
		// Check if error is retryable
		if !isRetryableError(err) {
			break
		}
		
		// Wait before retry
		select {
		case <-time.After(interval):
			interval = time.Duration(float64(interval) * trm.backoffFactor)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	
	return fmt.Errorf("transaction failed after %d retries: %w", trm.maxRetries, lastErr)
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error) bool {
	// Add logic to determine if error is retryable
	// For example, deadlock errors, timeout errors, etc.
	return true // Simplified for demo
}