package examples

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Pools demonstrates the worker pools pattern
func RunPools() {
	fmt.Println("=== Worker Pools Pattern Example ===")

	// Configuration
	numWorkers := 3
	numJobs := 15

	// Create job channel
	jobs := make(chan int, numJobs)

	// Create result channel
	results := make(chan string, numJobs)

	// Start the worker pool
	var wg sync.WaitGroup

	// Launch workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go workerPool(i, jobs, results, &wg)
	}

	// Send jobs to the pool
	go func() {
		defer close(jobs)
		for i := 1; i <= numJobs; i++ {
			fmt.Printf("Sending job %d to pool\n", i)
			jobs <- i
			time.Sleep(100 * time.Millisecond) // Simulate job generation time
		}
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Printf("\nWorker pool with %d workers processing %d jobs:\n", numWorkers, numJobs)
	fmt.Println()

	count := 0
	for result := range results {
		fmt.Printf("Result: %s\n", result)
		count++
	}

	fmt.Printf("\nWorker pool completed! Processed %d jobs.\n", count)
}

// Worker function for the pool
func workerPool(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d started\n", id)

	for job := range jobs {
		// Simulate work processing
		processingTime := time.Duration(rand.Intn(300)+200) * time.Millisecond
		fmt.Printf("Worker %d processing job %d (will take %v)\n", id, job, processingTime)

		time.Sleep(processingTime)

		result := fmt.Sprintf("Job %d completed by worker %d in %v", job, id, processingTime)
		results <- result
	}

	fmt.Printf("Worker %d finished\n", id)
}
