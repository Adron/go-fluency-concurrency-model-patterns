# MapReduce Pattern

## Overview

The MapReduce pattern processes large datasets by breaking the work into two phases: Map (process data in parallel) and Reduce (aggregate results). This pattern is essential for processing large datasets in parallel, distributed computing, data analytics and aggregation, and batch processing jobs.

## Implementation Details

### Structure

The MapReduce implementation in `examples/mapreduce.go` consists of three main phases:

1. **Map Phase** - Splits text into words and emits (word, 1) pairs
2. **Shuffle Phase** - Groups key-value pairs by key
3. **Reduce Phase** - Counts occurrences of each word

### Code Analysis

Let's break down the main function and understand how each component works:

```go
func RunMapReduce() {
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
}
```

**Step-by-step breakdown:**

1. **Input Data Preparation**:
   - `data := []string{...}` creates a slice of text lines to process
   - Each line contains multiple words that will be counted
   - The data is designed to demonstrate word frequency (e.g., "hello" appears multiple times)
   - This simulates a real-world scenario where you have text documents to analyze

2. **Map Phase Execution**:
   - `mapped := mapPhase(data)` processes the input data in parallel
   - Returns a channel of `KeyValue` pairs where each word is paired with a count of 1
   - This is the first phase of the MapReduce pattern - transforming raw data into key-value pairs
   - The channel-based approach allows for streaming data between phases

3. **Shuffle Phase Execution**:
   - `grouped := shufflePhase(mapped)` groups the key-value pairs by key
   - Takes the channel output from the map phase
   - Groups all occurrences of the same word together
   - Returns a map where keys are words and values are slices of counts

4. **Reduce Phase Execution**:
   - `result := reducePhase(grouped)` aggregates the grouped data
   - Sums up all the counts for each word
   - Produces the final word count results
   - This is the final phase that produces the desired output

5. **Result Display**:
   - Iterates through the final results map
   - Displays each word and its total count
   - Shows the complete word frequency analysis

### Map Phase Implementation

```go
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
```

**Map phase breakdown:**

1. **Channel Setup**:
   - `out := make(chan KeyValue, len(data)*10)` creates a buffered channel for output
   - Buffer size of `len(data)*10` accounts for multiple words per line
   - Buffering prevents blocking when multiple goroutines emit simultaneously
   - Returns a read-only channel (`<-chan KeyValue`) for encapsulation

2. **Parallel Processing Setup**:
   - `var wg sync.WaitGroup` tracks when all map goroutines complete
   - Each line of text gets its own goroutine for parallel processing
   - This enables true parallelism - all lines can be processed simultaneously

3. **Goroutine Launch**:
   - `for _, line := range data` iterates through each line of input
   - `go func(text string) { ... }(line)` launches a goroutine for each line
   - Uses closure to capture the current line in each goroutine
   - `wg.Add(1)` increments the wait group before each goroutine

4. **Word Processing**:
   - `defer wg.Done()` ensures the goroutine signals completion when it exits
   - `strings.Fields(strings.ToLower(text))` splits the line into words and converts to lowercase
   - `strings.Fields()` splits on whitespace, handling multiple spaces correctly
   - `strings.ToLower()` ensures consistent word matching (case-insensitive)

5. **Key-Value Emission**:
   - `for _, word := range words` processes each word in the line
   - `time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)` simulates processing time
   - `out <- KeyValue{Key: word, Value: 1}` emits a key-value pair for each word
   - Each word gets a value of 1 (will be summed in reduce phase)

6. **Channel Cleanup**:
   - `go func() { wg.Wait(); close(out) }()` runs in background
   - Waits for all map goroutines to complete
   - Closes the output channel when all processing is done
   - This signals to downstream phases that no more data is coming

### Shuffle Phase Implementation

```go
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
```

**Shuffle phase breakdown:**

1. **Data Structure Setup**:
   - `grouped := make(map[string][]int)` creates the output map
   - Keys are words (strings), values are slices of counts ([]int)
   - This structure groups all occurrences of each word together
   - `var mu sync.Mutex` provides thread safety for concurrent map access

2. **WaitGroup Setup**:
   - `var wg sync.WaitGroup` tracks when all shuffle goroutines complete
   - Each key-value pair gets its own goroutine for parallel processing
   - This maintains parallelism from the map phase

3. **Channel Consumption**:
   - `for kv := range mapped` consumes key-value pairs from the map phase
   - The loop continues until the map phase closes the channel
   - Each key-value pair is processed as it arrives

4. **Parallel Grouping**:
   - `go func(kv KeyValue) { ... }(kv)` launches a goroutine for each key-value pair
   - Uses closure to capture the current key-value pair
   - `wg.Add(1)` increments the wait group before each goroutine

5. **Thread-Safe Grouping**:
   - `defer wg.Done()` ensures proper cleanup
   - `mu.Lock()` and `mu.Unlock()` protect the shared map during concurrent access
   - `grouped[kv.Key] = append(grouped[kv.Key], kv.Value)` adds the value to the appropriate group
   - This groups all values for the same key together

6. **Completion and Return**:
   - `wg.Wait()` waits for all shuffle goroutines to complete
   - Returns the grouped map with all words and their associated counts
   - The map is now ready for the reduce phase

### Reduce Phase Implementation

```go
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
```

**Reduce phase breakdown:**

1. **Result Structure Setup**:
   - `result := make(map[string]int)` creates the final output map
   - Keys are words (strings), values are total counts (int)
   - This is the final result structure after aggregation
   - `var mu sync.Mutex` provides thread safety for concurrent result writing

2. **WaitGroup Setup**:
   - `var wg sync.WaitGroup` tracks when all reduce goroutines complete
   - Each word group gets its own goroutine for parallel processing
   - This maintains parallelism for the final aggregation phase

3. **Parallel Processing Setup**:
   - `for word, counts := range grouped` iterates through each word and its counts
   - Each word group is processed independently in parallel
   - This allows multiple words to be aggregated simultaneously

4. **Goroutine Launch**:
   - `go func(word string, counts []int) { ... }(word, counts)` launches a goroutine for each word
   - Uses closure to capture the current word and its counts
   - `wg.Add(1)` increments the wait group before each goroutine

5. **Aggregation Processing**:
   - `defer wg.Done()` ensures proper cleanup
   - `time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)` simulates processing time
   - `total := 0` initializes the sum for this word
   - `for _, count := range counts { total += count }` sums all counts for the word

6. **Thread-Safe Result Storage**:
   - `mu.Lock()` and `mu.Unlock()` protect the shared result map
   - `result[word] = total` stores the final count for the word
   - This ensures thread-safe writing to the final result map

7. **Completion and Return**:
   - `wg.Wait()` waits for all reduce goroutines to complete
   - Returns the final result map with word counts
   - This completes the MapReduce pattern

**Key Design Patterns:**

1. **Channel-based Data Flow**: Channels provide natural streaming between phases with automatic backpressure.

2. **Goroutine-per-Item Processing**: Each data item gets its own goroutine for true parallelism.

3. **Buffered Channels**: Appropriate buffering prevents blocking and improves performance.

4. **Mutex Protection**: Thread-safe access to shared data structures during concurrent operations.

5. **WaitGroup Coordination**: Proper synchronization ensures all goroutines complete before proceeding.

6. **Closure Pattern**: `func(text string) { ... }(line)` captures loop variables in goroutines.

7. **Structured Data Types**: `KeyValue` struct provides type safety and clarity.

## How It Works

### Map Phase
1. **Data Input**: Receives input data (lines of text)
2. **Parallel Processing**: Each line is processed in its own goroutine
3. **Word Extraction**: Splits each line into individual words
4. **Key-Value Emission**: Emits (word, 1) pairs for each word
5. **Channel Output**: Sends key-value pairs to the shuffle phase

### Shuffle Phase
1. **Key-Value Reception**: Receives key-value pairs from map phase
2. **Parallel Grouping**: Each key-value pair is processed in its own goroutine
3. **Grouping by Key**: Groups values by their keys
4. **Thread Safety**: Uses mutex to protect the shared grouped map
5. **Output Preparation**: Prepares grouped data for reduce phase

### Reduce Phase
1. **Grouped Data Input**: Receives grouped data from shuffle phase
2. **Parallel Aggregation**: Each word group is processed in its own goroutine
3. **Value Summation**: Sums all values for each key
4. **Result Storage**: Stores final results in the result map
5. **Final Output**: Returns the final word count results

## Why This Implementation?

### Channel-based Communication
- **Natural Flow**: Channels provide natural data flow between phases
- **Backpressure**: Automatic backpressure when downstream phases are slow
- **Synchronization**: Channels handle synchronization between phases
- **Composability**: Easy to modify or extend individual phases

### Goroutine-per-Item Processing
- **True Parallelism**: Each item is processed independently
- **Scalability**: Can utilize multiple CPU cores effectively
- **Fault Isolation**: Failure of one item doesn't affect others
- **Load Distribution**: Work is naturally distributed across workers

### Buffered Channels
- **Performance**: Buffered channels reduce blocking between phases
- **Memory Management**: Appropriate buffer sizes prevent memory issues
- **Flow Control**: Buffers can handle temporary processing delays
- **Efficiency**: Reduces context switching overhead

### Mutex Protection
- **Thread Safety**: Protects shared data structures during concurrent access
- **Simple Implementation**: Straightforward synchronization for grouped data
- **Performance**: Minimal overhead for the typical MapReduce use case
- **Reliability**: Ensures data consistency during concurrent operations

### Structured Data Types
- **Type Safety**: KeyValue struct provides type safety
- **Clarity**: Clear structure makes the code easy to understand
- **Extensibility**: Easy to extend for different data types
- **Debugging**: Structured data makes debugging easier

## Key Design Decisions

1. **Word Count Example**: Simple example that clearly demonstrates the pattern
2. **Simulated Processing Time**: Random delays make concurrency visible
3. **Detailed Logging**: Output shows the flow through each phase
4. **Error Handling**: Simple implementation focuses on the core pattern
5. **Memory Management**: Appropriate buffer sizes prevent memory issues

## Performance Characteristics

### Throughput
- **Parallel Processing**: Throughput scales with the number of CPU cores
- **Phase Overlap**: Phases can overlap in processing (pipelining)
- **Memory Usage**: Memory usage depends on data size and buffer sizes
- **Network Overhead**: In distributed systems, network communication adds overhead

### Latency
- **Processing Time**: Latency depends on the slowest phase
- **Data Size**: Larger datasets increase processing time
- **Parallelism**: More parallelism reduces latency
- **I/O Operations**: File I/O can be a significant bottleneck

### Scalability
- **Horizontal Scaling**: Can distribute across multiple machines
- **Vertical Scaling**: Can utilize multiple CPU cores on single machine
- **Data Partitioning**: Data can be partitioned for parallel processing
- **Load Balancing**: Work is naturally distributed across workers

## Common Use Cases

### Data Processing
- **Log Analysis**: Process large log files to extract insights
- **Text Processing**: Analyze text documents for patterns and statistics
- **Data Cleaning**: Clean and validate large datasets
- **ETL Operations**: Extract, transform, and load data

### Analytics and Reporting
- **Business Intelligence**: Generate reports from large datasets
- **User Behavior Analysis**: Analyze user activity patterns
- **Performance Metrics**: Calculate performance metrics from logs
- **Trend Analysis**: Identify trends in time-series data

### Machine Learning
- **Feature Engineering**: Extract features from raw data
- **Model Training**: Process training data for machine learning models
- **Data Preprocessing**: Prepare data for machine learning algorithms
- **Model Evaluation**: Evaluate models on large datasets

### Search and Indexing
- **Web Crawling**: Process web pages for search indexing
- **Document Indexing**: Index large document collections
- **Content Analysis**: Analyze content for search relevance
- **Inverted Index**: Build inverted indexes for search engines

### Financial Data Processing
- **Risk Analysis**: Analyze financial data for risk assessment
- **Trading Analysis**: Process trading data for market analysis
- **Fraud Detection**: Analyze transactions for fraud patterns
- **Portfolio Optimization**: Process portfolio data for optimization

### Scientific Computing
- **Simulation Data**: Process simulation results
- **Sensor Data**: Analyze data from scientific instruments
- **Image Processing**: Process large image datasets
- **Genomic Analysis**: Analyze genetic data

### Social Media Analysis
- **Sentiment Analysis**: Analyze social media posts for sentiment
- **Trend Detection**: Identify trending topics and hashtags
- **Network Analysis**: Analyze social network connections
- **Content Recommendation**: Process user behavior for recommendations

## Advanced Patterns

### Distributed MapReduce
- **Multi-node Processing**: Distribute processing across multiple machines
- **Fault Tolerance**: Handle machine failures gracefully
- **Load Balancing**: Balance load across multiple nodes
- **Data Locality**: Process data close to where it's stored

### Streaming MapReduce
- **Real-time Processing**: Process data as it arrives
- **Window-based Processing**: Process data in time windows
- **Incremental Updates**: Update results incrementally
- **Low Latency**: Provide results with minimal delay

### Iterative MapReduce
- **Multiple Passes**: Process data through multiple MapReduce cycles
- **Iterative Algorithms**: Support for iterative algorithms like PageRank
- **Convergence**: Continue until convergence criteria are met
- **State Management**: Maintain state across iterations

The MapReduce pattern is particularly effective when you have:
- **Large Datasets**: Datasets too large to process on a single machine
- **Parallelizable Work**: Work that can be divided into independent tasks
- **Batch Processing**: Operations that can be processed in batches
- **Data Analytics**: Need to extract insights from large datasets
- **Distributed Computing**: Need to utilize multiple machines or cores

This pattern provides a powerful framework for processing large datasets efficiently by leveraging parallel processing and distributed computing capabilities. 