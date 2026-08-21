package dispatcher

import "sync"

type Dispatcher struct {
	mu   sync.RWMutex
	subs map[chan string]struct{}
}

func New() *Dispatcher {
	return &Dispatcher{
		subs: make(map[chan string]struct{}),
	}
}

func (d *Dispatcher) Subscribe() chan string {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan string, 10)
	d.subs[ch] = struct{}{}
	return ch
}

func (d *Dispatcher) Unsubscribe(ch chan string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.subs[ch]; ok {
		delete(d.subs, ch)
		close(ch)
	}
}

func (d *Dispatcher) Broadcast(msg string) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for ch := range d.subs {
		select {
		case ch <- msg:
		default:
		}
	}
}