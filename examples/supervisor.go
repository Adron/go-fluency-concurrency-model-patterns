package examples

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// RunSupervisor demonstrates the supervisor/restart pattern.
func RunSupervisor() {
	fmt.Println("=== Supervisor/Restart Pattern Example ===")

	var restarts int32
	stop := make(chan struct{})
	done := make(chan struct{})

	// Supervisor goroutine
	go func() {
		for {
			workerDone := make(chan struct{})
			go workerWithFailure(workerDone, stop)
			select {
			case <-workerDone:
				atomic.AddInt32(&restarts, 1)
				fmt.Println("Supervisor: Worker failed, restarting...")
				// Restart after a short delay
				time.Sleep(500 * time.Millisecond)
			case <-stop:
				fmt.Println("Supervisor: Stopping worker supervision.")
				close(done)
				return
			}
		}
	}()

	// Let the supervisor run for a while
	time.Sleep(4 * time.Second)
	close(stop)
	<-done

	fmt.Printf("Supervisor example completed! Worker was restarted %d times.\n", restarts-1)
}

// workerWithFailure simulates a worker that randomly fails
func workerWithFailure(done chan<- struct{}, stop <-chan struct{}) {
	fmt.Println("Worker: Started")
	workTime := time.Duration(rand.Intn(1200)+400) * time.Millisecond
	select {
	case <-time.After(workTime):
		// Simulate random failure
		if rand.Float32() < 0.6 {
			fmt.Println("Worker: Simulated failure!")
			done <- struct{}{}
			return
		}
		fmt.Println("Worker: Completed work successfully.")
	case <-stop:
		fmt.Println("Worker: Received stop signal.")
	}
	// Signal normal exit
	done <- struct{}{}
}
