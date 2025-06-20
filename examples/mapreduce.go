package examples

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// RunMapReduce demonstrates the MapReduce pattern.
func RunMapReduce() {
	fmt.Println("=== MapReduce Pattern Example ===")

	// Sample data: words to count
	data := []string{
		"hello world",
		"hello go",
		"world of concurrency",
		"go programming",
		"concurrency patterns",
		"hello concurrency",
		"go world",
		"patterns in go",
	}

	fmt.Printf("Input data: %v\n", data)

	// Map phase: split words and emit (word, 1) pairs
	mapped := mapPhase(data)

	// Shuffle phase: group by key
	grouped := shufflePhase(mapped)

	// Reduce phase: count occurrences
	result := reducePhase(grouped)

	// Display results
	fmt.Println("\nWord count results:")
	for word, count := range result {
		fmt.Printf("  %s: %d\n", word, count)
	}

	fmt.Println("\nMapReduce example completed!")
}

// MapPhase splits text into words and emits (word, 1) pairs
func mapPhase(data []string) <-chan KeyValue {
	out := make(chan KeyValue, len(data)*10) // Buffer for multiple words per line

	var wg sync.WaitGroup
	for _, line := range data {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			words := strings.Fields(strings.ToLower(text))
			for _, word := range words {
				// Simulate some processing time
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				out <- KeyValue{Key: word, Value: 1}
				fmt.Printf("Map: emitted (%s, 1)\n", word)
			}
		}(line)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ShufflePhase groups key-value pairs by key
func shufflePhase(mapped <-chan KeyValue) map[string][]int {
	grouped := make(map[string][]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	for kv := range mapped {
		wg.Add(1)
		go func(kv KeyValue) {
			defer wg.Done()
			mu.Lock()
			grouped[kv.Key] = append(grouped[kv.Key], kv.Value)
			mu.Unlock()
			fmt.Printf("Shuffle: grouped %s -> %v\n", kv.Key, grouped[kv.Key])
		}(kv)
	}

	wg.Wait()
	return grouped
}

// ReducePhase counts occurrences of each word
func reducePhase(grouped map[string][]int) map[string]int {
	result := make(map[string]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	for word, counts := range grouped {
		wg.Add(1)
		go func(word string, counts []int) {
			defer wg.Done()
			// Simulate some processing time
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			total := 0
			for _, count := range counts {
				total += count
			}

			mu.Lock()
			result[word] = total
			mu.Unlock()
			fmt.Printf("Reduce: %s -> %d\n", word, total)
		}(word, counts)
	}

	wg.Wait()
	return result
}

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value int
}
