# Fan-out/Fan-in Pattern

## Overview

The Fan-out/Fan-in pattern is a concurrency pattern that distributes work across multiple workers (fan-out) and then collects results from all workers (fan-in). This pattern is essential for parallel processing of independent tasks, load balancing, and improving throughput for CPU-intensive operations.

## Implementation Details

### Structure

The fan-out/fan-in implementation in `examples/fan.go` consists of three main components:

1. **Work Generation** - Creates work items to be processed
2. **Fan-out** - Distributes work across multiple workers
3. **Fan-in** - Collects results from all workers

### Code Analysis

```go
func RunFan() {
    // Generate work items
    workItems := generateWorkItems(20)
    
    // Fan out: Distribute work across multiple workers
    numWorkers := 4
    results := fanOut(workItems, numWorkers)
    
    // Fan in: Collect results from all workers
    finalResults := fanIn(results)
    
    // Collect and display results
    for result := range finalResults {
        fmt.Printf("Processed: Item %d -> %s (by Worker %d)\n", 
                   result.OriginalID, result.Processed, result.WorkerID)
    }
}
```

### Work Item Structure

```go
type WorkItem struct {
    ID   int
    Data string
}

type Result struct {
    OriginalID int
    Processed  string
    WorkerID   int
}
```

### Fan-out Implementation

```go
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
```

### Worker Implementation

```go
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
```

### Fan-in Implementation

```go
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
```

## How It Works

1. **Work Generation**: Creates a stream of work items to be processed
2. **Fan-out**: Launches multiple worker goroutines, each consuming from the same job channel
3. **Worker Processing**: Each worker processes jobs independently and sends results to its own result channel
4. **Fan-in**: Multiple goroutines forward results from worker channels to a single output channel
5. **Result Collection**: The main function consumes all results from the fan-in channel

The pattern enables true parallel processing - multiple workers can process different jobs simultaneously, and results are collected as they complete.

## Why This Implementation?

### Channel-based Distribution
- **Automatic Load Balancing**: Workers naturally consume jobs as they become available
- **Backpressure**: If workers are slow, the job channel provides natural backpressure
- **Fair Distribution**: All workers have equal access to jobs

### Individual Result Channels
- **Isolation**: Each worker has its own result channel, preventing blocking
- **Order Independence**: Results can be collected in any order
- **Scalability**: Easy to add or remove workers

### WaitGroup Synchronization
- **Proper Cleanup**: Ensures all workers complete before closing channels
- **Resource Management**: Prevents goroutine leaks
- **Coordinated Shutdown**: All workers finish before fan-in completes

### Goroutine-per-Forward
- **Non-blocking Collection**: Each worker's results are forwarded independently
- **Concurrent Collection**: Results are collected as soon as they're available
- **Efficient Resource Usage**: No worker blocks waiting for others

## Key Design Decisions

1. **Shared Job Channel**: All workers read from the same job channel, providing automatic load balancing
2. **Individual Result Channels**: Each worker has its own result channel to prevent blocking
3. **Read-only Channel Returns**: The fan-out function returns read-only channels for safety
4. **Simulated Processing Time**: Random delays simulate real work and make concurrency visible
5. **Structured Results**: Results include metadata about which worker processed each item

## Common Use Cases

### Parallel Data Processing
- **Image Processing**: Process multiple images simultaneously
- **Document Processing**: Parse, analyze, or transform multiple documents
- **Data Validation**: Validate large datasets across multiple workers

### API Request Handling
- **Microservice Calls**: Make concurrent API calls to multiple services
- **Data Aggregation**: Fetch data from multiple sources simultaneously
- **Load Testing**: Simulate multiple concurrent users

### Batch Processing
- **File Processing**: Process multiple files in parallel
- **Database Operations**: Execute multiple queries concurrently
- **Report Generation**: Generate multiple reports simultaneously

### Real-time Systems
- **Sensor Data Processing**: Process data from multiple sensors
- **Log Analysis**: Analyze logs from multiple sources
- **Monitoring**: Collect metrics from multiple systems

### Machine Learning
- **Model Training**: Train multiple models in parallel
- **Feature Processing**: Process features across multiple workers
- **Hyperparameter Tuning**: Test multiple parameter combinations

The fan-out/fan-in pattern is particularly effective when you have independent tasks that can be processed in parallel and when you want to maximize resource utilization while maintaining result collection order. 