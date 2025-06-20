package examples

import (
	"fmt"
	"math/rand"
	"time"
)

// Pipeline demonstrates a multi-stage data processing pipeline
func RunPipeline() {
	fmt.Println("=== Pipeline Pattern Example ===")

	// Stage 1: Generate numbers
	numbers := generateNumbers(10)

	// Stage 2: Square the numbers
	squared := square(numbers)

	// Stage 3: Add 10 to each number
	result := addTen(squared)

	// Collect and display results
	fmt.Println("Pipeline stages:")
	fmt.Println("1. Generate numbers")
	fmt.Println("2. Square numbers")
	fmt.Println("3. Add 10")
	fmt.Println()

	for num := range result {
		fmt.Printf("Result: %d\n", num)
	}

	fmt.Println("Pipeline completed!")
}

// Stage 1: Generate random numbers
func generateNumbers(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			num := rand.Intn(10) + 1
			fmt.Printf("Generated: %d\n", num)
			out <- num
			time.Sleep(100 * time.Millisecond) // Simulate work
		}
	}()
	return out
}

// Stage 2: Square the numbers
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			squared := num * num
			fmt.Printf("Squared %d -> %d\n", num, squared)
			out <- squared
			time.Sleep(150 * time.Millisecond) // Simulate work
		}
	}()
	return out
}

// Stage 3: Add 10 to each number
func addTen(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			result := num + 10
			fmt.Printf("Added 10 to %d -> %d\n", num, result)
			out <- result
			time.Sleep(100 * time.Millisecond) // Simulate work
		}
	}()
	return out
}
