package examples

import (
	"fmt"
	"sync"
	"time"
)

// RunRateLimiting demonstrates rate limiting patterns.
func RunRateLimiting() {
	fmt.Println("=== Rate Limiting Pattern Example ===")

	// Example 1: Fixed rate limiting
	fmt.Println("\n1. Fixed rate limiting (2 requests per second):")
	limiter := newFixedRateLimiter(2, time.Second)
	var wg sync.WaitGroup

	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			limiter.Wait()
			fmt.Printf("Request %d processed at %v\n", id, time.Now().Format("15:04:05.000"))
		}(i)
	}

	wg.Wait()

	// Example 2: Token bucket rate limiting
	fmt.Println("\n2. Token bucket rate limiting (3 tokens per second, burst of 5):")
	tokenLimiter := newTokenBucketLimiter(3, 5)
	var wg2 sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			if tokenLimiter.Allow() {
				fmt.Printf("Token request %d granted at %v\n", id, time.Now().Format("15:04:05.000"))
			} else {
				fmt.Printf("Token request %d denied at %v\n", id, time.Now().Format("15:04:05.000"))
			}
		}(i)
	}

	wg2.Wait()

	fmt.Println("\nRate Limiting example completed!")
}

// Fixed rate limiter using time.Ticker
type fixedRateLimiter struct {
	ticker *time.Ticker
	stop   chan struct{}
}

func newFixedRateLimiter(rate int, interval time.Duration) *fixedRateLimiter {
	limiter := &fixedRateLimiter{
		ticker: time.NewTicker(interval / time.Duration(rate)),
		stop:   make(chan struct{}),
	}
	return limiter
}

func (r *fixedRateLimiter) Wait() {
	<-r.ticker.C
}

func (r *fixedRateLimiter) Stop() {
	r.ticker.Stop()
	close(r.stop)
}

// Token bucket rate limiter
type tokenBucketLimiter struct {
	tokens     chan struct{}
	rate       time.Duration
	burst      int
	mu         sync.Mutex
	lastRefill time.Time
}

func newTokenBucketLimiter(rate int, burst int) *tokenBucketLimiter {
	limiter := &tokenBucketLimiter{
		tokens:     make(chan struct{}, burst),
		rate:       time.Second / time.Duration(rate),
		burst:      burst,
		lastRefill: time.Now(),
	}

	// Fill the bucket initially
	for i := 0; i < burst; i++ {
		limiter.tokens <- struct{}{}
	}

	// Start refilling tokens
	go limiter.refill()

	return limiter
}

func (t *tokenBucketLimiter) refill() {
	ticker := time.NewTicker(t.rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case t.tokens <- struct{}{}:
			// Token added successfully
		default:
			// Bucket is full, skip
		}
	}
}

func (t *tokenBucketLimiter) Allow() bool {
	select {
	case <-t.tokens:
		return true
	default:
		return false
	}
}

func (t *tokenBucketLimiter) Wait() {
	<-t.tokens
}
