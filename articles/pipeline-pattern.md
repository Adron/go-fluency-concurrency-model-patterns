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

Let's break down the main function and understand how each component works:

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

**Step-by-step breakdown:**

1. **Stage 1: Number Generation**:
   - `numbers := generateNumbers(10)` creates the first stage of the pipeline
   - The function returns a channel that will receive 10 random numbers
   - This stage runs in its own goroutine and generates data asynchronously
   - The channel provides the connection to the next stage

2. **Stage 2: Number Squaring**:
   - `squared := square(numbers)` creates the second stage of the pipeline
   - Takes the output channel from Stage 1 as input
   - Squares each number as it receives it from the previous stage
   - Returns a new channel with the squared results

3. **Stage 3: Adding Ten**:
   - `result := addTen(squared)` creates the third stage of the pipeline
   - Takes the output channel from Stage 2 as input
   - Adds 10 to each squared number
   - Returns the final channel with the processed results

4. **Result Collection**:
   - `for num := range result` consumes all results from the final stage
   - The range loop automatically exits when the channel is closed
   - `fmt.Printf("Result: %d\n", num)` displays each final result
   - This is the consumer that drives the entire pipeline

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

**Stage implementation breakdown:**

1. **Function Signature**:
   - `func generateNumbers(count int) <-chan int` defines a function that returns a read-only channel
   - `count int` specifies how many numbers to generate
   - `<-chan int` indicates the function returns a channel that can only be read from
   - This provides encapsulation - callers can only consume, not produce

2. **Channel Creation**:
   - `out := make(chan int)` creates an unbuffered channel for output
   - Unbuffered channels provide natural synchronization between stages
   - Each stage waits for the next stage to be ready before sending data
   - This creates a natural flow control mechanism

3. **Goroutine Launch**:
   - `go func() { ... }()` launches the stage processing in a background goroutine
   - This allows the function to return immediately while processing continues
   - The stage runs independently and can process data concurrently with other stages
   - This is the key to achieving true pipeline concurrency

4. **Resource Management**:
   - `defer close(out)` ensures the output channel is closed when the goroutine exits
   - This signals to downstream stages that no more data is coming
   - Proper cleanup prevents goroutine leaks and channel deadlocks
   - The defer statement guarantees cleanup even if the function panics

5. **Data Generation Loop**:
   - `for i := 0; i < count; i++` generates exactly the specified number of items
   - `num := rand.Intn(10) + 1` creates random numbers between 1 and 10
   - `fmt.Printf("Generated: %d\n", num)` logs the generation for debugging
   - `out <- num` sends the number to the output channel

6. **Work Simulation**:
   - `time.Sleep(100 * time.Millisecond)` simulates processing time
   - This makes the concurrent nature of the pipeline visible
   - Without this delay, the pipeline would run too fast to observe
   - The delay also demonstrates how stages can work at different speeds

### Pipeline Flow

Let's trace through a single number's journey through the pipeline:

1. **Stage 1 (Generate)**: Creates number 7, sends to Stage 2
2. **Stage 2 (Square)**: Receives 7, squares it to 49, sends to Stage 3  
3. **Stage 3 (Add Ten)**: Receives 49, adds 10 to get 59, sends to consumer
4. **Consumer**: Receives 59 and displays "Result: 59"

While this is happening, Stage 1 can already be generating the next number, creating true concurrency.

### Concurrent Processing Benefits

The pipeline pattern enables several key benefits:

1. **Overlapped Processing**: While Stage 2 is squaring number 7, Stage 1 can be generating number 3
2. **Resource Utilization**: Multiple CPU cores can be utilized simultaneously
3. **Throughput**: Overall processing time is reduced compared to sequential processing
4. **Scalability**: Each stage can be optimized or scaled independently

### Channel-based Synchronization

The unbuffered channels provide natural synchronization:

- **Stage 1** blocks when sending if **Stage 2** isn't ready to receive
- **Stage 2** blocks when sending if **Stage 3** isn't ready to receive  
- **Stage 3** blocks when sending if the **consumer** isn't ready to receive
- This creates a natural backpressure mechanism that prevents memory buildup

### Error Handling and Cleanup

The pipeline design includes several safety features:

1. **Deferred Channel Closing**: Each stage properly closes its output channel
2. **Goroutine Management**: Each stage runs in its own goroutine for isolation
3. **Resource Cleanup**: Channels are automatically cleaned up by Go's garbage collector
4. **Flow Control**: Unbuffered channels prevent memory leaks from unbounded buffering

**Key Design Patterns:**

1. **Channel-based Communication**: Each stage communicates through channels, providing natural backpressure and synchronization.

2. **Goroutine-per-Stage**: Each stage runs in its own goroutine, enabling true concurrent processing.

3. **Unbuffered Channels**: Ensures each stage processes one item at a time and provides natural flow control.

4. **Deferred Channel Closing**: Each stage properly closes its output channel when done, signaling completion to downstream stages.

5. **Read-only Channel Returns**: Functions return `<-chan int` for encapsulation and safety.

6. **Simulated Work**: The `time.Sleep()` calls simulate real processing time, making the concurrent nature more visible.

7. **Pipeline Composition**: Stages can be easily composed and modified without affecting other stages.

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
