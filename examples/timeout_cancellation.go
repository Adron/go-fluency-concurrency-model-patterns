package examples

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RunTimeoutCancellation demonstrates timeouts and cancellation patterns.
func RunTimeoutCancellation() {
	fmt.Println("=== Timeouts and Cancellation Pattern Example ===")

	// Example 1: Context-based timeout
	fmt.Println("\n1. Context-based timeout example:")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := make(chan string, 1)
	go longRunningTask(ctx, result)

	select {
	case res := <-result:
		fmt.Printf("Task completed: %s\n", res)
	case <-ctx.Done():
		fmt.Printf("Task timed out: %v\n", ctx.Err())
	}

	// Example 2: Channel-based timeout
	fmt.Println("\n2. Channel-based timeout example:")
	ch := make(chan string, 1)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "Channel task completed"
	}()

	select {
	case res := <-ch:
		fmt.Printf("Channel task: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Channel task timed out")
	}

	// Example 3: Cancellation with context
	fmt.Println("\n3. Context cancellation example:")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Cancelling context...")
		cancel2()
	}()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Context cancellation example completed")
	case <-ctx2.Done():
		fmt.Printf("Context cancelled: %v\n", ctx2.Err())
	}

	fmt.Println("\nTimeouts and Cancellation example completed!")
}

// longRunningTask simulates a long-running task that respects context cancellation
func longRunningTask(ctx context.Context, result chan<- string) {
	// Simulate work with random duration
	workTime := time.Duration(rand.Intn(3000)+1000) * time.Millisecond
	fmt.Printf("Starting long task (will take %v)...\n", workTime)

	select {
	case <-time.After(workTime):
		result <- "Long task completed successfully"
	case <-ctx.Done():
		fmt.Printf("Long task cancelled: %v\n", ctx.Err())
		return
	}
}
