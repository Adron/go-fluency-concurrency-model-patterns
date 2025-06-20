# Concurrency Model Patterns

A Go project demonstrating common concurrency patterns and models using Go's goroutines and channels.

## Overview

This project provides practical examples of three fundamental concurrency patterns:

1. **Pipeline Pattern** - Multi-stage data processing pipeline
2. **Fan-out/Fan-in Pattern** - Distributing work across multiple workers and collecting results
3. **Worker Pools Pattern** - Fixed-size pool of workers processing jobs from a queue

## Project Structure

```
concurrency-model-patterns/
├── main.go              # Main application entry point
├── go.mod               # Go module definition
├── .gitignore           # Git ignore patterns
├── README.md            # This file
└── examples/            # Concurrency pattern examples
    ├── pipeline.go      # Pipeline pattern implementation
    ├── fan.go           # Fan-out/fan-in pattern implementation
    └── pools.go         # Worker pools pattern implementation
```

## Building the Application

To build the executable:

```bash
go build -o cmp-pattern
```

This creates an executable named `cmp-pattern` that you can run with different flags.

## Usage

The application supports three command-line flags to run different concurrency pattern examples:

### Pipeline Pattern
```bash
./cmp-pattern --pipeline
```
Demonstrates a 3-stage pipeline:
1. Generate random numbers
2. Square the numbers
3. Add 10 to each result

### Fan-out/Fan-in Pattern
```bash
./cmp-pattern --fan
```
Demonstrates distributing work across multiple workers and collecting results:
- Generates 20 work items
- Distributes them across 4 workers
- Collects and displays processed results

### Worker Pools Pattern
```bash
./cmp-pattern --pools
```
Demonstrates a fixed-size worker pool:
- Creates a pool of 3 workers
- Processes 15 jobs from a queue
- Shows how workers handle jobs concurrently

### Help
If no flag is provided, the application shows usage information:
```bash
./cmp-pattern
```

## Examples

### Running the Pipeline Example
```bash
$ ./cmp-pattern --pipeline
Running Pipeline Pattern Example...
=== Pipeline Pattern Example ===
Pipeline stages:
1. Generate numbers
2. Square numbers
3. Add 10

Generated: 7
Squared 7 -> 49
Added 10 to 49 -> 59
Result: 59
...
```

### Running the Fan-out/Fan-in Example
```bash
$ ./cmp-pattern --fan
Running Fan-out/Fan-in Pattern Example...
=== Fan-out/Fan-in Pattern Example ===
Distributing 20 work items across 4 workers...

Generated work item: 0
Worker 1 processed item 0
Processed: processed-data-0-by-worker-1
...
```

### Running the Worker Pools Example
```bash
$ ./cmp-pattern --pools
Running Worker Pools Pattern Example...
=== Worker Pools Pattern Example ===
Worker 1 started
Worker 2 started
Worker 3 started
Sending job 1 to pool
Worker 1 processing job 1 (will take 234ms)
...
```

## Concurrency Patterns Explained

### Pipeline Pattern
The pipeline pattern connects multiple stages where each stage processes data and passes it to the next stage. This is useful for:
- Data transformation workflows
- Multi-step processing
- Stream processing applications

### Fan-out/Fan-in Pattern
This pattern distributes work across multiple workers (fan-out) and then collects results from all workers (fan-in). Useful for:
- Parallel processing of independent tasks
- Load balancing across multiple workers
- Improving throughput for CPU-intensive tasks

### Worker Pools Pattern
A worker pool maintains a fixed number of workers that process jobs from a shared queue. Benefits include:
- Controlling resource usage
- Handling bursty workloads
- Providing predictable performance

## Requirements

- Go 1.21 or later
- No external dependencies required

## Development

To run the examples during development:

```bash
go run main.go --pipeline
go run main.go --fan
go run main.go --pools
```

## Contributing

Feel free to add more concurrency patterns or improve the existing examples. Each pattern should be implemented in its own file in the `examples/` directory and exposed through a function that can be called from `main.go`. 