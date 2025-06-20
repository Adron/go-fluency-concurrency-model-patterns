package examples

import (
	"fmt"
	"sync"
	"time"
)

// RunPubSub demonstrates the publish-subscribe (pub/sub) pattern.
func RunPubSub() {
	fmt.Println("=== Publish-Subscribe (Pub/Sub) Pattern Example ===")

	// Create a broadcaster
	b := newBroadcaster()

	numSubscribers := 3
	var wg sync.WaitGroup

	// Start subscribers
	for i := 1; i <= numSubscribers; i++ {
		ch := b.subscribe()
		wg.Add(1)
		go func(id int, ch <-chan string) {
			defer wg.Done()
			for msg := range ch {
				fmt.Printf("Subscriber %d received: %s\n", id, msg)
			}
			fmt.Printf("Subscriber %d done.\n", id)
		}(i, ch)
	}

	// Start publisher
	go func() {
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Message %d", i)
			fmt.Printf("Publisher sending: %s\n", msg)
			b.publish(msg)
			time.Sleep(400 * time.Millisecond)
		}
		b.close()
	}()

	wg.Wait()
	fmt.Println("Pub/Sub example completed!")
}

// broadcaster manages subscriptions and publishing
// Not thread-safe for subscribe after close

type broadcaster struct {
	subscribers []chan string
	closed      bool
	mu          sync.Mutex
}

func newBroadcaster() *broadcaster {
	return &broadcaster{
		subscribers: make([]chan string, 0),
	}
}

func (b *broadcaster) subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 2)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *broadcaster) publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	for _, ch := range b.subscribers {
		ch <- msg
	}
}

func (b *broadcaster) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.closed = true
}
