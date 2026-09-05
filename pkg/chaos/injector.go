package chaos

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrSimulated = errors.New("simulated chaos error")

type Config struct {
	ErrorProbability float64
	MaxDelay         time.Duration
}

type Injector struct {
	mu     sync.Mutex
	config Config
	r      *rand.Rand
}

func NewInjector(cfg Config) *Injector {
	return &Injector{
		config: cfg,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *Injector) Execute(fn func() error) error {
	i.mu.Lock()
	delay := time.Duration(0)
	if i.config.MaxDelay > 0 {
		delay = time.Duration(i.r.Int63n(int64(i.config.MaxDelay)))
	}
	shouldErr := i.config.ErrorProbability > 0 && i.r.Float64() < i.config.ErrorProbability
	i.mu.Unlock()

	if delay > 0 {
		time.Sleep(delay)
	}

	if shouldErr {
		return ErrSimulated
	}

	return fn()
}