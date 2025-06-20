package examples

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fan demonstrates the fan-out/fan-in pattern
func RunFan() {
	fmt.Println("=== Fan-out/Fan-in Pattern Example ===")

	// Generate work items
	workItems := generateWorkItems(20)

	// Fan out: Distribute work across multiple workers
	numWorkers := 4
	results := fanOut(workItems, numWorkers)

	// Fan in: Collect results from all workers
	finalResults := fanIn(results)

	fmt.Printf("Distributing %d work items across %d workers...\n", 20, numWorkers)
	fmt.Println()

	// Collect and display results
	count := 0
	for result := range finalResults {
		fmt.Printf("Processed: Item %d -> %s (by Worker %d)\n", result.OriginalID, result.Processed, result.WorkerID)
		count++
	}

	fmt.Printf("\nFan-out/Fan-in completed! Processed %d items.\n", count)
}

// WorkItem represents a unit of work
type WorkItem struct {
	ID   int
	Data string
}

// Result represents the processed work item
type Result struct {
	OriginalID int
	Processed  string
	WorkerID   int
}

// Generate work items
func generateWorkItems(count int) <-chan WorkItem {
	out := make(chan WorkItem)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			item := WorkItem{
				ID:   i,
				Data: fmt.Sprintf("data-%d", i),
			}
			fmt.Printf("Generated work item: %d\n", i)
			out <- item
			time.Sleep(50 * time.Millisecond)
		}
	}()
	return out
}

// Worker function that processes work items
func worker(id int, jobs <-chan WorkItem, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Simulate processing work
		time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)

		result := Result{
			OriginalID: job.ID,
			Processed:  fmt.Sprintf("processed-%s-by-worker-%d", job.Data, id),
			WorkerID:   id,
		}

		fmt.Printf("Worker %d processed item %d\n", id, job.ID)
		results <- result
	}
}

// Fan out: Distribute work across multiple workers
func fanOut(jobs <-chan WorkItem, numWorkers int) []<-chan Result {
	var workers []chan Result
	var wg sync.WaitGroup

	// Create workers
	for i := 0; i < numWorkers; i++ {
		workerResults := make(chan Result)
		workers = append(workers, workerResults)

		wg.Add(1)
		go worker(i+1, jobs, workerResults, &wg)
	}

	// Close worker result channels when all workers are done
	go func() {
		wg.Wait()
		for _, workerChan := range workers {
			close(workerChan)
		}
	}()

	// Convert to read-only channels for return
	var resultChannels []<-chan Result
	for _, ch := range workers {
		resultChannels = append(resultChannels, ch)
	}

	return resultChannels
}

// Fan in: Collect results from multiple channels
func fanIn(inputs []<-chan Result) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	// Function to forward results from one input channel
	forward := func(c <-chan Result) {
		defer wg.Done()
		for result := range c {
			out <- result
		}
	}

	wg.Add(len(inputs))
	for _, input := range inputs {
		go forward(input)
	}

	// Close output channel when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
