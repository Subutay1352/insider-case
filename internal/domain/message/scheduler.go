package message

import (
	"context"
	"fmt"
	"insider-case/internal/pkg/logger"
	"sync"
	"time"
)

// Scheduler manages the automatic message sending scheduler
type Scheduler struct {
	processor         MessageProcessor
	ticker            *time.Ticker
	ctx               context.Context
	cancel            context.CancelFunc
	mu                sync.RWMutex
	isRunning         bool
	interval          time.Duration
	processingTimeout time.Duration
}

// NewScheduler creates a new Scheduler
func NewScheduler(processor MessageProcessor, interval time.Duration, processingTimeout time.Duration) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		processor:         processor,
		ctx:               ctx,
		cancel:            cancel,
		interval:          interval,
		processingTimeout: processingTimeout,
		isRunning:         false,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return ErrSchedulerRunning
	}

	// Create new context and ticker
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.ticker = time.NewTicker(s.interval)
	s.isRunning = true

	// Start the scheduler in a goroutine
	go s.run()

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return ErrSchedulerNotRunning
	}

	s.cancel()
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.isRunning = false

	return nil
}

// StopAndWait stops the scheduler and waits for the goroutine to finish
func (s *Scheduler) StopAndWait(ctx context.Context) error {
	s.mu.Lock()
	if !s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("scheduler is not running")
	}

	// Mark as stopping and get the context
	done := make(chan struct{})
	oldCtx := s.ctx
	s.cancel()
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.isRunning = false
	s.mu.Unlock()

	// Wait for the goroutine to finish or timeout
	go func() {
		// Wait for context to be done (goroutine should exit)
		<-oldCtx.Done()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: %v", ErrSchedulerTimeout, ctx.Err())
	}
}

// IsRunning returns whether the scheduler is currently running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// run executes the scheduler loop
func (s *Scheduler) run() {
	// Send immediately on start
	s.sendMessages()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			s.sendMessages()
		}
	}
}

// sendMessages sends queued messages
func (s *Scheduler) sendMessages() {
	ctx, cancel := context.WithTimeout(context.Background(), s.processingTimeout)
	defer cancel()

	if err := s.processor.SendPendingMessages(ctx); err != nil {
		logger.Error("Failed to send queued messages",
			"error", err,
		)
	}
}
