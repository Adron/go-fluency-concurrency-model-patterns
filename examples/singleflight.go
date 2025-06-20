package examples

import (
	"fmt"
	"sync"
	"time"
)

// RunSingleflight demonstrates the singleflight (spaceflight) pattern.
func RunSingleflight() {
	fmt.Println("=== Singleflight (Spaceflight) Pattern Example ===")

	// Create a singleflight group
	sf := newSingleflight()

	// Simulate multiple concurrent requests for the same key
	key := "user:123"
	numRequests := 5

	var wg sync.WaitGroup
	results := make([]string, numRequests)

	fmt.Printf("Making %d concurrent requests for key: %s\n", numRequests, key)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Request %d: Starting...\n", id)

			result := sf.Do(key, func() (interface{}, error) {
				// Simulate expensive operation (e.g., database query, API call)
				fmt.Printf("Request %d: Executing expensive operation...\n", id)
				time.Sleep(2 * time.Second)
				return fmt.Sprintf("Data for %s (processed by request %d)", key, id), nil
			})

			results[id] = result.(string)
			fmt.Printf("Request %d: Completed with result: %s\n", id, result)
		}(i)
	}

	wg.Wait()

	// Show that all results are the same (same execution)
	fmt.Println("\nAll results should be identical:")
	for i, result := range results {
		fmt.Printf("  Request %d: %s\n", i, result)
	}

	// Test with different keys
	fmt.Println("\nTesting with different keys:")
	keys := []string{"user:123", "user:456", "user:123"}

	for i, key := range keys {
		wg.Add(1)
		go func(id int, k string) {
			defer wg.Done()
			result := sf.Do(k, func() (interface{}, error) {
				fmt.Printf("Request %d: Executing for key %s...\n", id, k)
				time.Sleep(1 * time.Second)
				return fmt.Sprintf("Data for %s", k), nil
			})
			fmt.Printf("Request %d: Key %s -> %s\n", id, k, result)
		}(i, key)
	}

	wg.Wait()
	fmt.Println("\nSingleflight example completed!")
}

// Singleflight ensures only one execution per key
type singleflight struct {
	mu    sync.Mutex
	calls map[string]*call
}

type call struct {
	wg   sync.WaitGroup
	val  interface{}
	err  error
	dups int
}

func newSingleflight() *singleflight {
	return &singleflight{
		calls: make(map[string]*call),
	}
}

func (sf *singleflight) Do(key string, fn func() (interface{}, error)) interface{} {
	sf.mu.Lock()

	if c, exists := sf.calls[key]; exists {
		// Another call is in progress for this key
		c.dups++
		sf.mu.Unlock()
		fmt.Printf("Duplicate call for key %s, waiting for result...\n", key)
		c.wg.Wait()
		return c.val
	}

	// Create new call
	c := &call{}
	c.wg.Add(1)
	sf.calls[key] = c
	sf.mu.Unlock()

	// Execute the function
	c.val, c.err = fn()
	c.wg.Done()

	// Clean up
	sf.mu.Lock()
	delete(sf.calls, key)
	sf.mu.Unlock()

	return c.val
}
