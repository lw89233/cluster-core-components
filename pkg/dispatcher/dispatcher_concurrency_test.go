package dispatcher

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDispatcher_Concurrency(t *testing.T) {
	d := New()
	var wg sync.WaitGroup

	numSubscribers := 500
	numBroadcasters := 500

	wg.Add(numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		go func(id int) {
			defer wg.Done()

			sub := d.Subscribe()

			select {
			case <-sub:
			case <-time.After(200 * time.Millisecond):
			}

			d.Unsubscribe(sub)
		}(i)
	}

	wg.Add(numBroadcasters)
	for i := 0; i < numBroadcasters; i++ {
		go func(id int) {
			defer wg.Done()
			d.Broadcast(fmt.Sprintf("CONCURRENT_MSG_%d", id))
		}(i)
	}

	wg.Wait()
}