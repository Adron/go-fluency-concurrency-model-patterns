package main

import (
	"flag"
	"fmt"
	"os"

	"concurrency-model-patterns/examples"
)

func main() {
	// Define command line flags
	pipeline := flag.Bool("pipeline", false, "Run pipeline pattern example")
	fan := flag.Bool("fan", false, "Run fan-out/fan-in pattern example")
	pools := flag.Bool("pools", false, "Run worker pools pattern example")

	// Parse command line flags
	flag.Parse()

	// Check if any flag was provided
	if !*pipeline && !*fan && !*pools {
		fmt.Println("Concurrency Model Patterns Examples")
		fmt.Println("===================================")
		fmt.Println("Usage:")
		fmt.Println("  cmp-pattern --pipeline  - Run pipeline pattern example")
		fmt.Println("  cmp-pattern --fan       - Run fan-out/fan-in pattern example")
		fmt.Println("  cmp-pattern --pools     - Run worker pools pattern example")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  ./cmp-pattern --pipeline")
		fmt.Println("  ./cmp-pattern --fan")
		fmt.Println("  ./cmp-pattern --pools")
		os.Exit(1)
	}

	// Run the selected example
	switch {
	case *pipeline:
		fmt.Println("Running Pipeline Pattern Example...")
		examples.RunPipeline()
	case *fan:
		fmt.Println("Running Fan-out/Fan-in Pattern Example...")
		examples.RunFan()
	case *pools:
		fmt.Println("Running Worker Pools Pattern Example...")
		examples.RunPools()
	}
}
