package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		ctx:    ctx,
		cancel: cancel,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		m.cancel()
	}()

	return m
}

func (m *Manager) Context() context.Context {
	return m.ctx
}

func (m *Manager) Stop() {
	m.cancel()
}

func (m *Manager) AddWorker() {
	m.wg.Add(1)
}

func (m *Manager) Done() {
	m.wg.Done()
}

func (m *Manager) Wait() {
	m.wg.Wait()
}