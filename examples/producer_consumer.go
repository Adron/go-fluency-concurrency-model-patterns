package examples

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RunProducerConsumer demonstrates the producer-consumer pattern with multiple producers and consumers.
func RunProducerConsumer() {
	fmt.Println("=== Producer-Consumer Pattern Example ===")

	bufferSize := 5
	numProducers := 2
	numConsumers := 3
	numItems := 10

	buffer := make(chan int, bufferSize)
	var wg sync.WaitGroup

	// Start producers
	for p := 1; p <= numProducers; p++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < numItems; i++ {
				item := rand.Intn(100)
				buffer <- item
				fmt.Printf("Producer %d produced: %d\n", id, item)
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
			}
		}(p)
	}

	// Start consumers
	var consumerWg sync.WaitGroup
	for c := 1; c <= numConsumers; c++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			for item := range buffer {
				fmt.Printf("Consumer %d consumed: %d\n", id, item)
				time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
			}
		}(c)
	}

	// Wait for all producers to finish, then close the buffer
	wg.Wait()
	close(buffer)

	// Wait for all consumers to finish
	consumerWg.Wait()

	fmt.Println("Producer-Consumer example completed!")
}
