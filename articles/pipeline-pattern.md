# Pipeline Pattern

## Overview

The Pipeline pattern is a fundamental concurrency pattern that connects multiple stages of processing where each stage processes data and passes it to the next stage. This pattern is particularly useful for data transformation workflows, stream processing, and multi-step computations.

## Implementation Details

### Structure

The pipeline implementation in `examples/pipeline.go` consists of three main stages:

1. **Generate Numbers** - Creates random numbers
2. **Square Numbers** - Squares each number
3. **Add Ten** - Adds 10 to each squared number

### Code Analysis

```go
func RunPipeline() {
    // Stage 1: Generate numbers
    numbers := generateNumbers(10)
    
    // Stage 2: Square the numbers
    squared := square(numbers)
    
    // Stage 3: Add 10 to each number
    result := addTen(squared)
    
    // Collect and display results
    for num := range result {
        fmt.Printf("Result: %d\n", num)
    }
}
```

### Stage Implementation

Each stage follows the same pattern:

```go
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
```

### Key Design Decisions

1. **Channel-based Communication**: Each stage communicates through channels, providing natural backpressure and synchronization.

2. **Goroutine-per-Stage**: Each stage runs in its own goroutine, enabling true concurrent processing.

3. **Buffered vs Unbuffered**: The example uses unbuffered channels, which ensures each stage processes one item at a time and provides natural flow control.

4. **Deferred Channel Closing**: Each stage properly closes its output channel when done, signaling completion to downstream stages.

5. **Simulated Work**: The `time.Sleep()` calls simulate real processing time, making the concurrent nature more visible.

## How It Works

1. **Stage 1** starts generating random numbers and sending them to its output channel
2. **Stage 2** receives numbers from Stage 1, squares them, and sends results to its output channel
3. **Stage 3** receives squared numbers, adds 10, and sends final results
4. **Main function** consumes the final results from Stage 3

The pipeline processes data concurrently - while Stage 2 is processing the first number, Stage 1 can be generating the second number, and so on.

## Why This Implementation?

### Channel-based Design
- **Natural Flow Control**: Channels provide built-in backpressure
- **Thread Safety**: Channels handle synchronization automatically
- **Composability**: Easy to add/remove stages or modify processing logic

### Goroutine-per-Stage
- **True Concurrency**: Each stage can process independently
- **Resource Efficiency**: Better CPU utilization
- **Scalability**: Can easily scale individual stages

### Unbuffered Channels
- **Synchronization**: Ensures proper coordination between stages
- **Memory Efficiency**: No buffering overhead
- **Predictable Flow**: Each stage processes one item at a time

## Common Use Cases

### Data Processing Pipelines
- **ETL (Extract, Transform, Load)**: Extract data from source, transform it, load to destination
- **Image Processing**: Resize → Filter → Compress → Save
- **Text Processing**: Parse → Clean → Analyze → Generate Report

### Stream Processing
- **Log Processing**: Parse → Filter → Aggregate → Store
- **Sensor Data**: Collect → Validate → Process → Alert
- **Financial Data**: Receive → Validate → Calculate → Store

### API Processing
- **Request Pipeline**: Authenticate → Validate → Process → Respond
- **Data Pipeline**: Fetch → Transform → Cache → Return

### Real-time Systems
- **IoT Data**: Collect → Process → Analyze → Act
- **Trading Systems**: Receive → Validate → Execute → Confirm
- **Monitoring**: Collect → Analyze → Alert → Log

The pipeline pattern is particularly effective when you have a series of transformations that can be performed independently and when you want to maximize throughput through concurrent processing. 