# Rate Limiting Pattern

## Overview

The Rate Limiting pattern controls the frequency of operations to prevent resource exhaustion and ensure fair usage. This pattern is essential for API request throttling, resource protection, preventing DoS attacks, and ensuring system stability under load.

## Implementation Details

### Structure

The rate limiting implementation in `examples/rate_limiting.go` demonstrates two main techniques:

1. **Fixed Rate Limiting** - Using `time.Ticker` for consistent rate control
2. **Token Bucket Rate Limiting** - Using a token bucket algorithm for burst handling

### Code Analysis

Let's break down the main function and understand how each component works:

```go
func RunRateLimiting() {
    // Example 1: Fixed rate limiting
    fmt.Println("\n1. Fixed rate limiting (2 requests per second):")
    limiter := newFixedRateLimiter(2, time.Second)
    var wg sync.WaitGroup

    for i := 1; i <= 6; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            limiter.Wait()
            fmt.Printf("Request %d processed at %v\n", id, time.Now().Format("15:04:05.000"))
        }(i)
    }

    wg.Wait()

    // Example 2: Token bucket rate limiting
    fmt.Println("\n2. Token bucket rate limiting (3 tokens per second, burst of 5):")
    tokenLimiter := newTokenBucketLimiter(3, 5)
    var wg2 sync.WaitGroup

    for i := 1; i <= 10; i++ {
        wg2.Add(1)
        go func(id int) {
            defer wg2.Done()
            if tokenLimiter.Allow() {
                fmt.Printf("Token request %d granted at %v\n", id, time.Now().Format("15:04:05.000"))
            } else {
                fmt.Printf("Token request %d denied at %v\n", id, time.Now().Format("15:04:05.000"))
            }
        }(i)
    }

    wg2.Wait()
}
```

**Step-by-step breakdown:**

1. **Example 1: Fixed Rate Limiter Setup**:
   - `limiter := newFixedRateLimiter(2, time.Second)` creates a limiter that allows 2 requests per second
   - Rate of 2 requests per second means each request must wait 500ms between requests
   - This demonstrates a simple, predictable rate limiting approach

2. **Fixed Rate Limiter Request Launch**:
   - Launches 6 concurrent goroutines to simulate multiple simultaneous requests
   - `for i := 1; i <= 6; i++` creates requests numbered 1-6
   - Uses closure `func(id int) { ... }(i)` to capture the request ID
   - `wg.Add(1)` tracks each request goroutine

3. **Fixed Rate Limiter Request Processing**:
   - `defer wg.Done()` ensures the request signals completion when it exits
   - `limiter.Wait()` blocks until the rate limiter allows the request to proceed
   - All requests are synchronized to the same ticker, ensuring exactly 2 requests per second
   - `time.Now().Format("15:04:05.000")` shows precise timing to demonstrate rate control

4. **Fixed Rate Limiter Coordination**:
   - `wg.Wait()` waits for all 6 requests to complete
   - Requests will be processed in batches of 2 per second due to the rate limiting

5. **Example 2: Token Bucket Limiter Setup**:
   - `tokenLimiter := newTokenBucketLimiter(3, 5)` creates a token bucket with 3 tokens per second and burst capacity of 5
   - Rate of 3 tokens per second means tokens are refilled every 333ms
   - Burst capacity of 5 means up to 5 requests can be processed immediately

6. **Token Bucket Limiter Request Launch**:
   - Launches 10 concurrent goroutines to test burst handling
   - `for i := 1; i <= 10; i++` creates requests numbered 1-10
   - Uses closure to capture request ID
   - `wg2.Add(1)` tracks each request goroutine

7. **Token Bucket Limiter Request Processing**:
   - `defer wg2.Done()` ensures proper cleanup
   - `if tokenLimiter.Allow()` checks if a token is available without blocking
   - If token available: request is granted and processed immediately
   - If no token available: request is denied immediately (non-blocking)
   - Shows precise timing to demonstrate burst handling vs rate limiting

8. **Token Bucket Limiter Coordination**:
   - `wg2.Wait()` waits for all 10 requests to complete
   - First 5 requests will likely be granted immediately (burst capacity)
   - Remaining requests will be granted as tokens are refilled (3 per second)

### Fixed Rate Limiter Implementation

```go
type fixedRateLimiter struct {
    ticker *time.Ticker
    stop   chan struct{}
}

func newFixedRateLimiter(rate int, interval time.Duration) *fixedRateLimiter {
    limiter := &fixedRateLimiter{
        ticker: time.NewTicker(interval / time.Duration(rate)),
        stop:   make(chan struct{}),
    }
    return limiter
}

func (r *fixedRateLimiter) Wait() {
    <-r.ticker.C
}

func (r *fixedRateLimiter) Stop() {
    r.ticker.Stop()
    close(r.stop)
}
```

**Fixed rate limiter breakdown:**

1. **Data Structure Design**:
   - `ticker *time.Ticker`: Go's built-in ticker that fires at regular intervals
   - `stop chan struct{}`: Signal channel for graceful shutdown
   - Simple structure with minimal memory footprint

2. **Constructor Function**:
   - `newFixedRateLimiter(rate int, interval time.Duration)` takes rate and interval parameters
   - `interval / time.Duration(rate)` calculates the tick interval
   - Example: 1 second / 2 requests = 500ms between requests
   - `time.NewTicker()` creates a ticker that fires every 500ms

3. **Wait Method**:
   - `Wait()` blocks until the next tick occurs
   - `<-r.ticker.C` waits for the next tick from the ticker
   - All requests calling `Wait()` are synchronized to the same ticker
   - Ensures exactly the specified rate of operations

4. **Stop Method**:
   - `Stop()` provides cleanup functionality
   - `r.ticker.Stop()` stops the ticker and frees resources
   - `close(r.stop)` signals any waiting goroutines to stop
   - Prevents resource leaks

### Token Bucket Limiter Implementation

```go
type tokenBucketLimiter struct {
    tokens    chan struct{}
    rate      time.Duration
    burst     int
    mu        sync.Mutex
    lastRefill time.Time
}

func newTokenBucketLimiter(rate int, burst int) *tokenBucketLimiter {
    limiter := &tokenBucketLimiter{
        tokens:     make(chan struct{}, burst),
        rate:       time.Second / time.Duration(rate),
        burst:      burst,
        lastRefill: time.Now(),
    }

    // Fill the bucket initially
    for i := 0; i < burst; i++ {
        limiter.tokens <- struct{}{}
    }

    // Start refilling tokens
    go limiter.refill()

    return limiter
}

func (t *tokenBucketLimiter) refill() {
    ticker := time.NewTicker(t.rate)
    defer ticker.Stop()

    for range ticker.C {
        select {
        case t.tokens <- struct{}{}:
            // Token added successfully
        default:
            // Bucket is full, skip
        }
    }
}

func (t *tokenBucketLimiter) Allow() bool {
    select {
    case <-t.tokens:
        return true
    default:
        return false
    }
}

func (t *tokenBucketLimiter) Wait() {
    <-t.tokens
}
```

**Token bucket limiter breakdown:**

1. **Data Structure Design**:
   - `tokens chan struct{}`: Channel-based bucket storing tokens
   - `rate time.Duration`: Time interval between token refills
   - `burst int`: Maximum number of tokens (burst capacity)
   - `mu sync.Mutex`: Protects concurrent access (though not used in this implementation)
   - `lastRefill time.Time`: Tracks last refill time (though not used in this implementation)

2. **Constructor Function**:
   - `newTokenBucketLimiter(rate int, burst int)` creates a token bucket limiter
   - `make(chan struct{}, burst)` creates a buffered channel with burst capacity
   - `time.Second / time.Duration(rate)` calculates refill interval
   - Example: 1 second / 3 tokens = 333ms between refills

3. **Initial Token Filling**:
   - `for i := 0; i < burst; i++` fills the bucket with initial tokens
   - `limiter.tokens <- struct{}{}` adds tokens to the channel
   - Bucket starts full, allowing immediate burst of requests

4. **Background Refilling**:
   - `go limiter.refill()` starts background token refilling
   - Runs independently of request processing
   - Ensures continuous token availability

5. **Refill Method**:
   - `refill()` runs in background goroutine
   - `time.NewTicker(t.rate)` creates ticker for regular refills
   - `for range ticker.C` continuously refills tokens
   - `select` with `default` prevents blocking if bucket is full

6. **Allow Method (Non-blocking)**:
   - `Allow()` checks if token is available without blocking
   - `select` with `default` provides non-blocking token consumption
   - Returns `true` if token available, `false` if bucket empty
   - Immediate response for rate limit checking

7. **Wait Method (Blocking)**:
   - `Wait()` blocks until a token becomes available
   - `<-t.tokens` waits for token from the channel
   - FIFO order ensures fair token distribution
   - Used when blocking behavior is desired

**Key Design Patterns:**

1. **Channel-based Token Storage**: Uses buffered channels for thread-safe token management with FIFO order.

2. **Background Refilling**: Continuous token refilling in background goroutine ensures steady token availability.

3. **Non-blocking Token Check**: `Allow()` method provides immediate response without blocking.

4. **Ticker-based Rate Control**: Both limiters use `time.Ticker` for precise timing control.

5. **Closure Pattern**: `func(id int) { ... }(i)` captures loop variables in request goroutines.

6. **WaitGroup Coordination**: Proper coordination ensures all requests complete before program exits.

7. **Resource Cleanup**: Both limiters provide cleanup methods to prevent resource leaks.

## How It Works

### Fixed Rate Limiting
1. **Ticker Creation**: Creates a ticker that fires at the desired rate
2. **Request Processing**: Each request waits for the next tick
3. **Rate Control**: Requests are processed at exactly the specified rate
4. **Synchronization**: All requests are synchronized to the ticker

### Token Bucket Rate Limiting
1. **Bucket Initialization**: Creates a bucket filled with initial tokens
2. **Token Consumption**: Requests consume tokens from the bucket
3. **Token Refilling**: Tokens are refilled at a constant rate
4. **Burst Handling**: Can handle bursts up to the bucket capacity

## Why This Implementation?

### Fixed Rate Limiter
- **Predictable Rate**: Ensures exactly the specified rate of operations
- **Simple Implementation**: Straightforward using Go's time.Ticker
- **Synchronization**: All operations are naturally synchronized
- **Resource Efficiency**: Minimal memory and CPU overhead

### Token Bucket Limiter
- **Burst Handling**: Can handle bursts of requests up to bucket capacity
- **Flexible Rate**: Allows for variable request patterns
- **Non-blocking**: Requests can be denied immediately if no tokens available
- **Fair Distribution**: Tokens are consumed in FIFO order

### Channel-based Token Storage
- **Thread Safety**: Channels provide natural thread safety
- **FIFO Order**: Tokens are consumed in the order they were added
- **Non-blocking Operations**: Can check for tokens without blocking
- **Efficient Implementation**: Minimal overhead for token management

### Goroutine-based Refilling
- **Continuous Refilling**: Tokens are refilled continuously in the background
- **Independent Operation**: Refilling doesn't block request processing
- **Automatic Management**: No manual intervention required
- **Resource Efficiency**: Single goroutine handles all refilling

## Key Design Decisions

1. **Two Different Approaches**: Demonstrates both fixed-rate and token-bucket techniques
2. **Simulated Requests**: Multiple concurrent requests show rate limiting in action
3. **Timing Information**: Output includes timestamps to show rate control
4. **Non-blocking Token Check**: Token bucket allows immediate denial of requests
5. **Proper Cleanup**: Both limiters provide cleanup methods

## Performance Characteristics

### Fixed Rate Limiter
- **Consistent Rate**: Exactly the specified rate, no variation
- **Low Overhead**: Minimal CPU and memory usage
- **Synchronization**: All operations are synchronized to the ticker
- **Predictable**: Behavior is completely predictable

### Token Bucket Limiter
- **Burst Handling**: Can handle bursts up to bucket capacity
- **Variable Rate**: Average rate is maintained over time
- **Immediate Response**: Requests are processed immediately if tokens available
- **Fair Distribution**: Requests are processed in order

### Memory Usage
- **Fixed Limiter**: Minimal memory usage (just ticker and channel)
- **Token Bucket**: Memory usage proportional to burst capacity
- **Scalability**: Both scale well with multiple concurrent users

## Common Use Cases

### API Rate Limiting
- **REST APIs**: Limit requests per user or IP address
- **GraphQL APIs**: Control query complexity and frequency
- **Webhook Processing**: Limit incoming webhook frequency
- **Third-party APIs**: Respect external API rate limits

### Resource Protection
- **Database Connections**: Limit concurrent database connections
- **File Operations**: Control file system access frequency
- **Network Requests**: Limit outgoing network requests
- **CPU-intensive Operations**: Control CPU usage

### User Experience
- **UI Interactions**: Limit button clicks or form submissions
- **Search Queries**: Control search request frequency
- **File Uploads**: Limit upload frequency and size
- **Real-time Updates**: Control update frequency

### Security
- **Brute Force Protection**: Prevent password guessing attacks
- **DDoS Mitigation**: Limit requests from suspicious sources
- **Account Protection**: Prevent account takeover attempts
- **API Abuse Prevention**: Prevent API abuse and scraping

### System Stability
- **Load Balancing**: Control load on backend services
- **Cache Management**: Limit cache invalidation frequency
- **Background Jobs**: Control job execution frequency
- **Monitoring**: Limit metric collection frequency

### Microservices
- **Service Communication**: Limit inter-service request frequency
- **Circuit Breaker**: Control circuit breaker state changes
- **Health Checks**: Limit health check frequency
- **Configuration Updates**: Control configuration change frequency

## Advanced Patterns

### Sliding Window Rate Limiting
- **Time-based Windows**: Rate limits based on sliding time windows
- **Precise Control**: More precise than fixed windows
- **Memory Efficient**: Uses less memory than token buckets
- **Complex Implementation**: More complex to implement correctly

### Leaky Bucket Rate Limiting
- **Constant Rate**: Processes requests at a constant rate
- **Queue Management**: Queues requests when rate is exceeded
- **Smooth Processing**: Provides smooth, predictable processing
- **Memory Usage**: Can use significant memory during bursts

### Distributed Rate Limiting
- **Multi-node Coordination**: Coordinate rate limits across multiple nodes
- **Shared State**: Use Redis or similar for shared rate limit state
- **Consistency**: Ensure consistent rate limiting across all nodes
- **Network Overhead**: Additional network communication required

### Adaptive Rate Limiting
- **Dynamic Adjustment**: Adjust rate limits based on system load
- **Health-based**: Reduce rates when system is unhealthy
- **User-based**: Different limits for different user types
- **Time-based**: Different limits for different times of day

## Best Practices

### Rate Limit Configuration
- **Appropriate Limits**: Set limits based on system capacity and requirements
- **User Experience**: Consider impact on user experience
- **Monitoring**: Monitor rate limit effectiveness and adjust as needed
- **Documentation**: Clearly document rate limits for API users

### Error Handling
- **Graceful Degradation**: Provide meaningful error messages when limits are exceeded
- **Retry Logic**: Implement appropriate retry logic for rate-limited requests
- **Backoff Strategies**: Use exponential backoff for retries
- **User Feedback**: Inform users when they're approaching limits

### Monitoring and Alerting
- **Rate Limit Metrics**: Track rate limit usage and violations
- **Alerting**: Alert when rate limits are frequently exceeded
- **Trend Analysis**: Analyze rate limit patterns over time
- **Capacity Planning**: Use rate limit data for capacity planning

The rate limiting pattern is particularly effective when you have:
- **Resource Constraints**: Limited system resources that need protection
- **Fair Usage**: Need to ensure fair usage across multiple users
- **Security Requirements**: Need to prevent abuse and attacks
- **API Management**: Need to control API usage and prevent abuse
- **System Stability**: Need to maintain system stability under varying load

This pattern provides essential tools for building robust, scalable systems that can handle varying loads while protecting system resources and ensuring fair usage. 