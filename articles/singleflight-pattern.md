# Singleflight (Spaceflight) Pattern

## Overview

The Singleflight pattern ensures that only one execution of a function is in-flight for a given key at a time. Duplicate callers wait for the result instead of executing the function again. This pattern is essential for preventing duplicate expensive operations, caching with concurrent access, API call deduplication, and database query optimization.

## Implementation Details

### Structure

The singleflight implementation in `examples/singleflight.go` consists of three main components:

1. **Singleflight Group** - Manages in-flight calls and their results
2. **Call Tracking** - Tracks active calls for each key
3. **Result Sharing** - Shares results among duplicate callers

### Code Analysis

Let's break down the main function and understand how each component works:

```go
func RunSingleflight() {
    // Create a singleflight group
    sf := newSingleflight()

    // Simulate multiple concurrent requests for the same key
    key := "user:123"
    numRequests := 5

    var wg sync.WaitGroup
    results := make([]string, numRequests)

    fmt.Printf("Making %d concurrent requests for key: %s\n", numRequests, key)

    // Launch concurrent requests
    for i := 0; i < numRequests; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Request %d: Starting...\n", id)
            
            result := sf.Do(key, func() (interface{}, error) {
                // Simulate expensive operation (e.g., database query, API call)
                fmt.Printf("Request %d: Executing expensive operation...\n", id)
                time.Sleep(2 * time.Second)
                return fmt.Sprintf("Data for %s (processed by request %d)", key, id), nil
            })
            
            results[id] = result.(string)
            fmt.Printf("Request %d: Completed with result: %s\n", id, result)
        }(i)
    }

    wg.Wait()

    // Show that all results are the same (same execution)
    fmt.Println("\nAll results should be identical:")
    for i, result := range results {
        fmt.Printf("  Request %d: %s\n", i, result)
    }
}
```

**Step-by-step breakdown:**

1. **Singleflight Group Creation**:
   - `sf := newSingleflight()` creates a new singleflight group instance
   - The group manages all in-flight calls and their results
   - This is the central coordinator that prevents duplicate executions

2. **Request Configuration**:
   - `key := "user:123"` defines the key that all requests will use
   - `numRequests := 5` specifies how many concurrent requests to simulate
   - Using the same key ensures that the singleflight pattern will be triggered

3. **Result Storage Setup**:
   - `var wg sync.WaitGroup` tracks when all request goroutines complete
   - `results := make([]string, numRequests)` creates a slice to store results
   - This allows us to verify that all requests get the same result

4. **Concurrent Request Launch**:
   - `for i := 0; i < numRequests; i++` launches 5 concurrent requests
   - Each request gets a unique ID (0-4) for tracking and debugging
   - Uses closure `func(id int) { ... }(i)` to capture the request ID

5. **Request Goroutine Implementation**:
   - `defer wg.Done()` ensures the request signals completion when it exits
   - `fmt.Printf("Request %d: Starting...\n", id)` shows when each request starts
   - This helps demonstrate the concurrent nature of the requests

6. **Singleflight Call Execution**:
   - `result := sf.Do(key, func() (interface{}, error) { ... })` calls the singleflight group
   - The function parameter defines the expensive operation to deduplicate
   - Only one execution of this function will occur, regardless of how many requests call it

7. **Expensive Operation Simulation**:
   - `fmt.Printf("Request %d: Executing expensive operation...\n", id)` shows which request actually executes
   - `time.Sleep(2 * time.Second)` simulates a 2-second expensive operation
   - Only the first request will see this message - others will wait for the result

8. **Result Processing**:
   - `results[id] = result.(string)` stores the result for this request
   - `fmt.Printf("Request %d: Completed with result: %s\n", id, result)` shows completion
   - All requests should complete with the same result

9. **Verification and Display**:
   - `wg.Wait()` waits for all requests to complete
   - The final loop displays all results to verify they're identical
   - This demonstrates that only one execution occurred, but all requests got the result

### Singleflight Implementation

```go
type singleflight struct {
    mu    sync.Mutex
    calls map[string]*call
}

type call struct {
    wg    sync.WaitGroup
    val   interface{}
    err   error
    dups  int
}

func newSingleflight() *singleflight {
    return &singleflight{
        calls: make(map[string]*call),
    }
}

func (sf *singleflight) Do(key string, fn func() (interface{}, error)) interface{} {
    sf.mu.Lock()
    
    if c, exists := sf.calls[key]; exists {
        // Another call is in progress for this key
        c.dups++
        sf.mu.Unlock()
        fmt.Printf("Duplicate call for key %s, waiting for result...\n", key)
        c.wg.Wait()
        return c.val
    }

    // Create new call
    c := &call{}
    c.wg.Add(1)
    sf.calls[key] = c
    sf.mu.Unlock()

    // Execute the function
    c.val, c.err = fn()
    c.wg.Done()

    // Clean up
    sf.mu.Lock()
    delete(sf.calls, key)
    sf.mu.Unlock()

    return c.val
}
```

**Singleflight implementation breakdown:**

1. **Data Structure Design**:
   - `singleflight struct` contains the main singleflight group
   - `mu sync.Mutex` protects concurrent access to the calls map
   - `calls map[string]*call` tracks active calls by key
   - Simple structure with minimal memory footprint

2. **Call Structure Design**:
   - `call struct` represents an in-flight function call
   - `wg sync.WaitGroup` coordinates between executing and waiting requests
   - `val interface{}` stores the result of the function execution
   - `err error` stores any error from the function execution
   - `dups int` tracks the number of duplicate requests (for monitoring)

3. **Constructor Function**:
   - `newSingleflight()` creates a new singleflight group
   - `calls: make(map[string]*call)` initializes the empty calls map
   - Returns a pointer to the singleflight group for method calls

4. **Do Method - Initial Lock and Duplicate Check**:
   - `sf.mu.Lock()` acquires the mutex to protect the calls map
   - `if c, exists := sf.calls[key]; exists` checks if a call is already in progress
   - If a call exists, this is a duplicate request

5. **Duplicate Request Handling**:
   - `c.dups++` increments the duplicate counter for monitoring
   - `sf.mu.Unlock()` releases the mutex immediately (no need to hold it while waiting)
   - `fmt.Printf("Duplicate call for key %s, waiting for result...\n", key)` shows duplicate detection
   - `c.wg.Wait()` blocks until the original call completes
   - `return c.val` returns the shared result to the duplicate request

6. **New Call Creation**:
   - `c := &call{}` creates a new call structure for this key
   - `c.wg.Add(1)` increments the WaitGroup (will be decremented when execution completes)
   - `sf.calls[key] = c` stores the call in the map to track it
   - `sf.mu.Unlock()` releases the mutex before executing the function

7. **Function Execution**:
   - `c.val, c.err = fn()` executes the provided function
   - This is the actual expensive operation that we want to deduplicate
   - Only one request will reach this point for each key

8. **Result Sharing and Cleanup**:
   - `c.wg.Done()` signals that execution is complete
   - This unblocks all waiting duplicate requests
   - `sf.mu.Lock()` acquires the mutex for cleanup
   - `delete(sf.calls, key)` removes the call from the map to free memory
   - `sf.mu.Unlock()` releases the mutex
   - `return c.val` returns the result to the original request

**Key Design Patterns:**

1. **Mutex Protection**: `sync.Mutex` ensures thread-safe access to the shared calls map during concurrent operations.

2. **WaitGroup Coordination**: `sync.WaitGroup` coordinates between the executing request and all waiting duplicate requests.

3. **Map-based Call Tracking**: Uses a map to efficiently track active calls by key with O(1) lookup time.

4. **Function-based Interface**: Generic interface `func() (interface{}, error)` supports any expensive operation.

5. **Automatic Cleanup**: Call entries are automatically removed from the map after completion to prevent memory leaks.

6. **Duplicate Monitoring**: Tracks the number of duplicate calls for observability and debugging.

7. **Closure Pattern**: `func(id int) { ... }(i)` captures loop variables in request goroutines for proper ID tracking.

## How It Works

1. **Request Arrival**: Multiple requests arrive for the same key
2. **First Request**: The first request creates a call entry and executes the function
3. **Duplicate Detection**: Subsequent requests detect an existing call for the key
4. **Result Waiting**: Duplicate requests wait for the first request to complete
5. **Result Sharing**: All requests receive the same result from the single execution
6. **Cleanup**: The call entry is removed after completion

The pattern ensures that expensive operations are executed only once, even under high concurrency.

## Why This Implementation?

### Mutex-based Synchronization
- **Thread Safety**: Protects the calls map during concurrent access
- **Simple Implementation**: Straightforward synchronization for the call tracking
- **Performance**: Minimal overhead for the typical singleflight use case
- **Reliability**: Ensures data consistency during concurrent operations

### WaitGroup for Result Coordination
- **Synchronization**: WaitGroup coordinates between the executing request and waiting requests
- **Non-blocking**: Executing request can complete without blocking
- **Clean Coordination**: All requests are properly synchronized
- **Resource Efficiency**: Minimal overhead for coordination

### Map-based Call Tracking
- **Key-based Lookup**: Efficient lookup of active calls by key
- **Memory Management**: Calls are removed after completion to prevent memory leaks
- **Scalability**: Can handle multiple different keys simultaneously
- **Flexibility**: Easy to extend for additional metadata

### Function-based Interface
- **Generic Design**: Can work with any function that returns (interface{}, error)
- **Type Safety**: Caller specifies the function to execute
- **Flexibility**: Supports different types of expensive operations
- **Error Handling**: Proper error propagation from the executed function

### Duplicate Counting
- **Monitoring**: Tracks the number of duplicate calls for observability
- **Debugging**: Helps identify when singleflight is being used effectively
- **Metrics**: Can be used for performance monitoring
- **Optimization**: Helps identify opportunities for caching

## Key Design Decisions

1. **Map-based Storage**: Uses a map to track active calls by key
2. **WaitGroup Coordination**: Uses WaitGroup to synchronize result sharing
3. **Mutex Protection**: Protects the calls map during concurrent access
4. **Automatic Cleanup**: Removes call entries after completion
5. **Generic Interface**: Supports any function with (interface{}, error) signature

## Performance Characteristics

### Throughput
- **Deduplication**: Eliminates duplicate expensive operations
- **Resource Efficiency**: Reduces resource usage for duplicate requests
- **Scalability**: Scales well with multiple concurrent users
- **Memory Usage**: Minimal memory overhead for call tracking

### Latency
- **First Request**: Latency is the same as the original operation
- **Duplicate Requests**: Latency is reduced to just coordination overhead
- **Coordination Overhead**: Minimal overhead for result sharing
- **Predictable**: Latency is predictable and consistent

### Memory Usage
- **Call Tracking**: Memory usage proportional to number of active calls
- **Automatic Cleanup**: Memory is freed when calls complete
- **Efficient Storage**: Minimal memory overhead per call
- **Scalability**: Memory usage doesn't grow with duplicate requests

## Common Use Cases

### Caching Systems
- **Cache Miss Handling**: Prevent multiple cache misses for the same key
- **Cache Population**: Ensure only one request populates the cache
- **Cache Warming**: Prevent duplicate cache warming operations
- **Distributed Caching**: Coordinate cache operations across multiple nodes

### Database Operations
- **Query Deduplication**: Prevent duplicate database queries
- **Connection Pooling**: Ensure only one connection is established per key
- **Transaction Management**: Prevent duplicate transaction operations
- **Data Loading**: Prevent duplicate data loading operations

### API Calls
- **External API Calls**: Prevent duplicate calls to external APIs
- **Rate Limiting**: Ensure rate limits are respected across concurrent requests
- **Authentication**: Prevent duplicate authentication requests
- **Data Fetching**: Prevent duplicate data fetching operations

### File Operations
- **File Reading**: Prevent duplicate file read operations
- **File Processing**: Ensure only one process handles a file
- **Configuration Loading**: Prevent duplicate configuration loading
- **Resource Loading**: Prevent duplicate resource loading

### Background Jobs
- **Job Deduplication**: Prevent duplicate background job execution
- **Scheduled Tasks**: Ensure scheduled tasks run only once
- **Data Processing**: Prevent duplicate data processing operations
- **Report Generation**: Prevent duplicate report generation

### System Operations
- **Service Discovery**: Prevent duplicate service discovery requests
- **Health Checks**: Prevent duplicate health check operations
- **Configuration Updates**: Prevent duplicate configuration updates
- **Maintenance Operations**: Prevent duplicate maintenance operations

### User Interface
- **Data Loading**: Prevent duplicate data loading in UI components
- **Form Submission**: Prevent duplicate form submissions
- **Search Operations**: Prevent duplicate search requests
- **Navigation**: Prevent duplicate navigation operations

## Advanced Patterns

### Time-based Expiration
- **Call Expiration**: Expire calls after a certain time
- **Stale Results**: Handle cases where results become stale
- **Refresh Logic**: Automatically refresh expired results
- **TTL Management**: Manage time-to-live for cached results

### Error Handling
- **Error Propagation**: Properly propagate errors to all callers
- **Retry Logic**: Implement retry logic for failed operations
- **Circuit Breaker**: Integrate with circuit breaker patterns
- **Fallback Mechanisms**: Provide fallback mechanisms for failures

### Distributed Singleflight
- **Multi-node Coordination**: Coordinate across multiple nodes
- **Shared State**: Use Redis or similar for shared state
- **Consistency**: Ensure consistent behavior across nodes
- **Network Overhead**: Handle additional network communication

### Metrics and Monitoring
- **Call Tracking**: Track the number of calls and duplicates
- **Performance Metrics**: Monitor performance impact
- **Hit Rates**: Track cache hit rates and effectiveness
- **Alerting**: Alert on unusual patterns or failures

## Best Practices

### Key Design
- **Appropriate Keys**: Use keys that properly identify the operation
- **Key Uniqueness**: Ensure keys are unique for different operations
- **Key Stability**: Use stable keys that don't change frequently
- **Key Size**: Keep keys reasonably sized for performance

### Function Design
- **Idempotent Functions**: Ensure functions are idempotent
- **Error Handling**: Properly handle and propagate errors
- **Resource Management**: Ensure proper resource cleanup
- **Timeout Handling**: Implement appropriate timeouts

### Memory Management
- **Call Cleanup**: Ensure calls are properly cleaned up
- **Memory Monitoring**: Monitor memory usage of call tracking
- **Leak Prevention**: Prevent memory leaks from abandoned calls
- **Size Limits**: Implement limits on the number of active calls

### Performance Optimization
- **Key Hashing**: Use efficient key hashing for large key sets
- **Call Pooling**: Reuse call objects to reduce allocation overhead
- **Lazy Cleanup**: Implement lazy cleanup to reduce lock contention
- **Metrics Collection**: Collect metrics for performance tuning

The singleflight pattern is particularly effective when you have:
- **Expensive Operations**: Operations that are costly to execute
- **High Concurrency**: Multiple concurrent requests for the same data
- **Cache Misses**: Frequent cache misses for the same keys
- **External Dependencies**: Operations that depend on external services
- **Resource Constraints**: Limited resources that need to be conserved

This pattern provides essential tools for building efficient, scalable systems that can handle high concurrency while minimizing resource usage and improving performance. 