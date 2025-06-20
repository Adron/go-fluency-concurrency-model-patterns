# Timeouts and Cancellation Pattern

## Overview

The Timeouts and Cancellation pattern uses `context.Context` or channels to control the lifecycle of concurrent operations. This pattern is essential for preventing goroutines from running indefinitely, graceful shutdown of operations, resource cleanup and management, and building responsive applications.

## Implementation Details

### Structure

The timeouts and cancellation implementation in `examples/timeout_cancellation.go` demonstrates three main techniques:

1. **Context-based Timeout** - Using `context.WithTimeout` for deadline-based cancellation
2. **Channel-based Timeout** - Using `time.After` for simple timeout scenarios
3. **Context Cancellation** - Using `context.WithCancel` for manual cancellation

### Code Analysis

```go
func RunTimeoutCancellation() {
    // Example 1: Context-based timeout
    fmt.Println("\n1. Context-based timeout example:")
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    result := make(chan string, 1)
    go longRunningTask(ctx, result)

    select {
    case res := <-result:
        fmt.Printf("Task completed: %s\n", res)
    case <-ctx.Done():
        fmt.Printf("Task timed out: %v\n", ctx.Err())
    }

    // Example 2: Channel-based timeout
    fmt.Println("\n2. Channel-based timeout example:")
    ch := make(chan string, 1)
    go func() {
        time.Sleep(3 * time.Second)
        ch <- "Channel task completed"
    }()

    select {
    case res := <-ch:
        fmt.Printf("Channel task: %s\n", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Channel task timed out")
    }

    // Example 3: Context cancellation
    fmt.Println("\n3. Context cancellation example:")
    ctx2, cancel2 := context.WithCancel(context.Background())
    defer cancel2()

    go func() {
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Cancelling context...")
        cancel2()
    }()

    select {
    case <-time.After(2 * time.Second):
        fmt.Println("Context cancellation example completed")
    case <-ctx2.Done():
        fmt.Printf("Context cancelled: %v\n", ctx2.Err())
    }
}
```

### Long-running Task Implementation

```go
func longRunningTask(ctx context.Context, result chan<- string) {
    // Simulate work with random duration
    workTime := time.Duration(rand.Intn(3000)+1000) * time.Millisecond
    fmt.Printf("Starting long task (will take %v)...\n", workTime)

    select {
    case <-time.After(workTime):
        // Simulate random failure
        if rand.Float32() < 0.6 {
            fmt.Println("Long task cancelled!")
            return
        }
        result <- "Long task completed successfully"
    case <-ctx.Done():
        fmt.Printf("Long task cancelled: %v\n", ctx.Err())
        return
    }
}
```

## How It Works

### Context-based Timeout
1. **Context Creation**: `context.WithTimeout` creates a context with a deadline
2. **Task Execution**: The long-running task receives the context and monitors it
3. **Timeout Detection**: The task checks `ctx.Done()` to detect timeout
4. **Graceful Cancellation**: When timeout occurs, the task stops work and returns

### Channel-based Timeout
1. **Task Launch**: A goroutine starts the potentially long-running task
2. **Timeout Setup**: `time.After` creates a channel that will receive after the timeout
3. **Race Condition**: `select` waits for either task completion or timeout
4. **Timeout Handling**: If timeout occurs first, the task is abandoned

### Context Cancellation
1. **Context Creation**: `context.WithCancel` creates a cancellable context
2. **Cancellation Trigger**: Another goroutine calls `cancel()` after a delay
3. **Cancellation Detection**: The main goroutine detects cancellation via `ctx.Done()`
4. **Immediate Response**: Cancellation is detected immediately when triggered

## Why This Implementation?

### Context-based Approach
- **Standard Pattern**: Uses Go's standard context package
- **Propagation**: Context can be passed down through call chains
- **Rich Information**: Context provides error information and cancellation reasons
- **Integration**: Works well with HTTP servers, gRPC, and other Go libraries

### Channel-based Approach
- **Simplicity**: Simple and straightforward for basic timeout scenarios
- **Performance**: Minimal overhead for simple cases
- **Flexibility**: Can be combined with other channel operations
- **Learning**: Good for understanding timeout concepts

### Select Statement
- **Non-blocking**: Allows waiting for multiple events without blocking
- **Race Handling**: Naturally handles the race between completion and timeout
- **Clean Code**: Provides clear, readable timeout logic
- **Efficiency**: Efficient event-driven waiting

### Buffered Result Channels
- **Non-blocking Sends**: Tasks can send results without blocking
- **Cleanup Prevention**: Prevents goroutine leaks when timeout occurs
- **Flexibility**: Allows for both synchronous and asynchronous result handling

## Key Design Decisions

1. **Multiple Examples**: Demonstrates different timeout and cancellation techniques
2. **Simulated Work**: Random work duration makes timeout behavior visible
3. **Error Information**: Context provides detailed error information
4. **Graceful Handling**: Tasks check for cancellation and stop work appropriately
5. **Resource Cleanup**: Proper defer statements ensure cleanup

## Performance Characteristics

### Context Overhead
- **Minimal Impact**: Context operations have minimal performance overhead
- **Memory Usage**: Context objects are lightweight
- **Propagation Cost**: Context propagation through call chains is efficient

### Channel Overhead
- **Low Latency**: Channel operations are fast
- **Memory Usage**: Channels use minimal memory
- **Scalability**: Scales well with multiple concurrent operations

### Timeout Accuracy
- **System Dependent**: Accuracy depends on system timer resolution
- **Go Runtime**: Go's runtime provides good timer accuracy
- **Practical Limits**: Sub-millisecond timeouts may not be accurate

## Common Use Cases

### HTTP Request Handling
- **Request Timeouts**: Cancel requests that take too long
- **Client Timeouts**: Set timeouts for outgoing HTTP requests
- **Server Timeouts**: Cancel long-running server operations

### Database Operations
- **Query Timeouts**: Cancel slow database queries
- **Connection Timeouts**: Timeout database connection attempts
- **Transaction Timeouts**: Cancel long-running transactions

### External API Calls
- **API Timeouts**: Cancel slow external API calls
- **Rate Limiting**: Cancel requests when rate limits are exceeded
- **Circuit Breaker**: Cancel requests when external services are down

### File Operations
- **Read Timeouts**: Cancel slow file read operations
- **Write Timeouts**: Cancel slow file write operations
- **Network File Systems**: Handle slow network file system operations

### Background Processing
- **Job Timeouts**: Cancel long-running background jobs
- **Batch Processing**: Timeout batch processing operations
- **Data Processing**: Cancel slow data processing tasks

### User Interface
- **User Input**: Cancel operations when user cancels
- **UI Responsiveness**: Ensure UI remains responsive during long operations
- **Progress Indicators**: Cancel operations based on user feedback

### System Operations
- **Backup Operations**: Cancel slow backup operations
- **System Maintenance**: Timeout maintenance operations
- **Resource Cleanup**: Cancel resource cleanup operations

### Microservices
- **Service Communication**: Timeout inter-service communication
- **Health Checks**: Cancel slow health check operations
- **Configuration Updates**: Timeout configuration update operations

## Best Practices

### Context Usage
- **Always Pass Context**: Pass context to all functions that can be cancelled
- **Check Cancellation**: Regularly check `ctx.Done()` in long-running operations
- **Propagate Context**: Pass context down through call chains
- **Use Timeouts**: Set appropriate timeouts for all operations

### Error Handling
- **Check Context Errors**: Always check `ctx.Err()` for cancellation reasons
- **Log Cancellations**: Log when operations are cancelled for debugging
- **Cleanup Resources**: Ensure resources are cleaned up on cancellation
- **Return Errors**: Return appropriate errors when operations are cancelled

### Timeout Values
- **Appropriate Timeouts**: Set timeouts based on expected operation duration
- **User Experience**: Consider user experience when setting timeouts
- **System Resources**: Consider system resource constraints
- **External Dependencies**: Account for external service response times

### Resource Management
- **Goroutine Cleanup**: Ensure goroutines exit when cancelled
- **Memory Cleanup**: Clean up memory allocations on cancellation
- **Connection Cleanup**: Close connections when operations are cancelled
- **File Handle Cleanup**: Close file handles on cancellation

The timeouts and cancellation pattern is particularly effective when you have:
- **Long-running Operations**: Operations that may take longer than expected
- **External Dependencies**: Operations that depend on external services
- **Resource Constraints**: Need to limit resource usage
- **User Experience**: Need to maintain responsive user interfaces
- **System Stability**: Need to prevent resource exhaustion

This pattern provides essential tools for building robust, responsive applications that can handle failures gracefully and maintain system stability under various conditions. 