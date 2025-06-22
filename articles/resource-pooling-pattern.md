# Resource Pooling Pattern

## Overview

The Resource Pooling pattern manages a pool of reusable resources (like database connections, HTTP clients, or file handles) to avoid the overhead of creating and destroying resources frequently. This pattern is essential for improving performance, managing resource limits, reducing overhead, and ensuring efficient resource utilization.

## Implementation Details

### Structure

The resource pooling implementation in `examples/resource_pooling.go` demonstrates two main examples:

1. **Database Connection Pool** - Manages database connections with health checks
2. **HTTP Client Pool** - Manages HTTP clients for making requests

### Code Analysis

Let's break down the main function and understand how each component works:

```go
func RunResourcePooling() {
    // Example 1: Database connection pool
    fmt.Println("1. Database Connection Pool Example:")
    dbPool := newDBConnectionPool(3, 5*time.Second)
    defer dbPool.Close()

    var wg sync.WaitGroup
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            conn := dbPool.Get()
            defer dbPool.Put(conn)
            
            fmt.Printf("Worker %d: Using connection %s\n", id, conn.ID)
            time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
        }(i)
    }
    wg.Wait()

    // Example 2: HTTP client pool
    fmt.Println("\n2. HTTP Client Pool Example:")
    clientPool := newHTTPClientPool(2, 3*time.Second)
    defer clientPool.Close()

    for i := 1; i <= 4; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            client := clientPool.Get()
            defer clientPool.Put(client)
            
            fmt.Printf("Worker %d: Using HTTP client %s\n", id, client.ID)
            time.Sleep(time.Duration(rand.Intn(800)+300) * time.Millisecond)
        }(i)
    }
    wg.Wait()
}
```

**Step-by-step breakdown:**

1. **Database Connection Pool Setup**:
   - `dbPool := newDBConnectionPool(3, 5*time.Second)` creates a pool with 3 connections and 5-second timeout
   - Pool size of 3 means only 3 workers can use connections simultaneously
   - 5-second timeout prevents workers from waiting indefinitely for connections
   - `defer dbPool.Close()` ensures proper cleanup when the function exits

2. **Database Pool Worker Launch**:
   - `var wg sync.WaitGroup` tracks when all database workers complete
   - `for i := 1; i <= 5; i++` launches 5 workers (more than pool size to demonstrate pooling)
   - Uses closure `func(id int) { ... }(i)` to capture the worker ID
   - `wg.Add(1)` increments the wait group before each worker

3. **Database Pool Worker Implementation**:
   - `defer wg.Done()` ensures the worker signals completion when it exits
   - `conn := dbPool.Get()` acquires a connection from the pool (may block if none available)
   - `defer dbPool.Put(conn)` ensures the connection is returned to the pool
   - `time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)` simulates database work (500-1500ms)

4. **Database Pool Coordination**:
   - `wg.Wait()` waits for all database workers to complete
   - Only 3 workers can use connections simultaneously due to pool size
   - Other workers will wait for connections to become available

5. **HTTP Client Pool Setup**:
   - `clientPool := newHTTPClientPool(2, 3*time.Second)` creates a pool with 2 clients and 3-second timeout
   - Smaller pool size (2) demonstrates different resource constraints
   - Shorter timeout (3 seconds) shows different timeout requirements
   - `defer clientPool.Close()` ensures proper cleanup

6. **HTTP Client Pool Worker Launch**:
   - `for i := 1; i <= 4; i++` launches 4 workers (more than pool size)
   - Similar structure to database workers but with different resource type
   - `wg.Add(1)` tracks each HTTP client worker

7. **HTTP Client Pool Worker Implementation**:
   - `defer wg.Done()` ensures proper cleanup
   - `client := clientPool.Get()` acquires an HTTP client from the pool
   - `defer clientPool.Put(client)` ensures the client is returned to the pool
   - `time.Sleep(time.Duration(rand.Intn(800)+300) * time.Millisecond)` simulates HTTP work (300-1100ms)

8. **HTTP Client Pool Coordination**:
   - `wg.Wait()` waits for all HTTP client workers to complete
   - Only 2 workers can use clients simultaneously due to pool size
   - Demonstrates how pooling handles resource constraints

### Database Connection Pool Implementation

```go
type DBConnection struct {
    ID        string
    CreatedAt time.Time
    LastUsed  time.Time
}

type DBConnectionPool struct {
    connections chan *DBConnection
    maxSize     int
    timeout     time.Duration
    closed      bool
    mu          sync.Mutex
}

func newDBConnectionPool(maxSize int, timeout time.Duration) *DBConnectionPool {
    pool := &DBConnectionPool{
        connections: make(chan *DBConnection, maxSize),
        maxSize:     maxSize,
        timeout:     timeout,
    }

    // Pre-populate the pool
    for i := 0; i < maxSize; i++ {
        conn := &DBConnection{
            ID:        fmt.Sprintf("db-conn-%d", i+1),
            CreatedAt: time.Now(),
            LastUsed:  time.Now(),
        }
        pool.connections <- conn
    }

    return pool
}

func (p *DBConnectionPool) Get() *DBConnection {
    select {
    case conn := <-p.connections:
        conn.LastUsed = time.Now()
        fmt.Printf("  Acquired connection: %s\n", conn.ID)
        return conn
    case <-time.After(p.timeout):
        fmt.Println("  Timeout waiting for connection")
        return nil
    }
}

func (p *DBConnectionPool) Put(conn *DBConnection) {
    if conn == nil {
        return
    }

    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        return
    }

    // Check if connection is still healthy
    if time.Since(conn.LastUsed) > p.timeout {
        fmt.Printf("  Discarding stale connection: %s\n", conn.ID)
        return
    }

    select {
    case p.connections <- conn:
        fmt.Printf("  Returned connection: %s\n", conn.ID)
    default:
        fmt.Printf("  Pool full, discarding connection: %s\n", conn.ID)
    }
}

func (p *DBConnectionPool) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.closed = true
    close(p.connections)
}
```

**Database connection pool breakdown:**

1. **Connection Structure Design**:
   - `DBConnection struct` represents a database connection
   - `ID string` provides unique identification for tracking and debugging
   - `CreatedAt time.Time` tracks when the connection was created
   - `LastUsed time.Time` tracks when the connection was last used for health checking

2. **Pool Structure Design**:
   - `DBConnectionPool struct` manages the connection pool
   - `connections chan *DBConnection` uses a buffered channel for thread-safe resource management
   - `maxSize int` defines the maximum number of connections in the pool
   - `timeout time.Duration` defines the timeout for connection acquisition
   - `closed bool` tracks pool state for graceful shutdown
   - `mu sync.Mutex` protects pool state during concurrent operations

3. **Pool Constructor**:
   - `newDBConnectionPool(maxSize int, timeout time.Duration)` creates a new pool
   - `connections: make(chan *DBConnection, maxSize)` creates a buffered channel with pool size
   - The channel buffer size matches the pool size for optimal performance

4. **Pool Pre-population**:
   - `for i := 0; i < maxSize; i++` creates the initial set of connections
   - `conn := &DBConnection{...}` creates each connection with unique ID and timestamps
   - `pool.connections <- conn` adds each connection to the pool
   - Pre-population ensures connections are immediately available

5. **Connection Acquisition (Get Method)**:
   - `select` statement handles both connection acquisition and timeout
   - `case conn := <-p.connections:` acquires a connection from the pool
   - `conn.LastUsed = time.Now()` updates the last used timestamp
   - `case <-time.After(p.timeout):` handles timeout when no connections are available
   - Returns `nil` on timeout to indicate failure

6. **Connection Return (Put Method)**:
   - `if conn == nil { return }` handles nil connections gracefully
   - `p.mu.Lock()` and `defer p.mu.Unlock()` protect pool state during return
   - `if p.closed { return }` prevents returning connections to a closed pool

7. **Health Checking**:
   - `if time.Since(conn.LastUsed) > p.timeout` checks if connection is stale
   - Stale connections are discarded rather than returned to the pool
   - This ensures only healthy connections are reused

8. **Connection Return Logic**:
   - `select { case p.connections <- conn: ... default: ... }` attempts to return connection
   - If pool is full, connection is discarded (default case)
   - This prevents pool overflow and memory leaks

9. **Pool Cleanup (Close Method)**:
   - `p.mu.Lock()` and `defer p.mu.Unlock()` protect pool state during shutdown
   - `p.closed = true` marks the pool as closed
   - `close(p.connections)` closes the channel to signal shutdown

### HTTP Client Pool Implementation

```go
type HTTPClient struct {
    ID        string
    CreatedAt time.Time
    LastUsed  time.Time
}

type HTTPClientPool struct {
    clients chan *HTTPClient
    maxSize int
    timeout time.Duration
    closed  bool
    mu      sync.Mutex
}

func newHTTPClientPool(maxSize int, timeout time.Duration) *HTTPClientPool {
    pool := &HTTPClientPool{
        clients: make(chan *HTTPClient, maxSize),
        maxSize: maxSize,
        timeout: timeout,
    }

    // Pre-populate the pool
    for i := 0; i < maxSize; i++ {
        client := &HTTPClient{
            ID:        fmt.Sprintf("http-client-%d", i+1),
            CreatedAt: time.Now(),
            LastUsed:  time.Now(),
        }
        pool.clients <- client
    }

    return pool
}

func (p *HTTPClientPool) Get() *HTTPClient {
    select {
    case client := <-p.clients:
        client.LastUsed = time.Now()
        fmt.Printf("  Acquired HTTP client: %s\n", client.ID)
        return client
    case <-time.After(p.timeout):
        fmt.Println("  Timeout waiting for HTTP client")
        return nil
    }
}

func (p *HTTPClientPool) Put(client *HTTPClient) {
    if client == nil {
        return
    }

    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        return
    }

    // Check if client is still healthy
    if time.Since(client.LastUsed) > p.timeout {
        fmt.Printf("  Discarding stale HTTP client: %s\n", client.ID)
        return
    }

    select {
    case p.clients <- client:
        fmt.Printf("  Returned HTTP client: %s\n", client.ID)
    default:
        fmt.Printf("  Pool full, discarding HTTP client: %s\n", client.ID)
    }
}

func (p *HTTPClientPool) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.closed = true
    close(p.clients)
}
```

**HTTP client pool breakdown:**

1. **Client Structure Design**:
   - `HTTPClient struct` represents an HTTP client
   - Similar structure to database connections for consistency
   - `ID string` provides unique identification
   - `CreatedAt` and `LastUsed` timestamps for lifecycle tracking

2. **Pool Structure Design**:
   - `HTTPClientPool struct` manages the HTTP client pool
   - Identical structure to database pool for consistency
   - Uses the same channel-based approach for resource management

3. **Pool Constructor and Pre-population**:
   - `newHTTPClientPool(maxSize int, timeout time.Duration)` creates the pool
   - Pre-populates with `maxSize` HTTP clients
   - Each client gets a unique ID and initial timestamps

4. **Client Acquisition and Return**:
   - `Get()` method identical to database pool with timeout handling
   - `Put()` method includes health checking and graceful handling
   - Same mutex protection and closed pool handling

5. **Health Checking and Cleanup**:
   - Same health checking logic as database pool
   - Stale clients are discarded to maintain pool health
   - `Close()` method provides graceful shutdown

**Key Design Patterns:**

1. **Channel-based Resource Management**: Uses buffered channels for thread-safe resource allocation and return.

2. **Pre-populated Pools**: Creates all resources upfront for immediate availability and predictable performance.

3. **Health Checking**: Validates resource health before reuse to ensure reliability.

4. **Timeout Handling**: Prevents indefinite waiting with configurable timeouts.

5. **Graceful Shutdown**: Proper cleanup prevents resource leaks during shutdown.

6. **Mutex Protection**: Thread-safe state management for pool operations.

7. **Resource Tracking**: Comprehensive tracking of resource lifecycle for monitoring and debugging.

## How It Works

### Database Connection Pool
1. **Pool Initialization**: Creates a pool with pre-populated database connections
2. **Connection Acquisition**: Workers request connections from the pool
3. **Connection Usage**: Workers use connections for database operations
4. **Connection Return**: Workers return connections to the pool when done
5. **Health Checking**: Pool checks connection health and discards stale connections
6. **Resource Management**: Pool manages connection lifecycle and cleanup

### HTTP Client Pool
1. **Pool Initialization**: Creates a pool with pre-populated HTTP clients
2. **Client Acquisition**: Workers request HTTP clients from the pool
3. **Client Usage**: Workers use clients for HTTP requests
4. **Client Return**: Workers return clients to the pool when done
5. **Health Checking**: Pool checks client health and discards stale clients
6. **Resource Management**: Pool manages client lifecycle and cleanup

## Why This Implementation?

### Channel-based Pool Management
- **Thread Safety**: Channels provide natural thread safety for pool operations
- **Blocking Semantics**: Workers block when no resources are available
- **FIFO Order**: Resources are distributed in first-in-first-out order
- **Simple Implementation**: Straightforward resource management

### Pre-populated Pools
- **Fast Startup**: Resources are available immediately
- **Predictable Performance**: No delay for resource creation
- **Resource Efficiency**: Optimal resource utilization from start
- **Warm-up Avoidance**: No cold start performance issues

### Health Checking
- **Resource Validation**: Ensures resources are still usable
- **Stale Detection**: Detects and removes stale resources
- **Automatic Cleanup**: Automatically manages resource lifecycle
- **Reliability**: Ensures only healthy resources are used

### Timeout Handling
- **Deadlock Prevention**: Prevents workers from waiting indefinitely
- **Graceful Degradation**: Handles resource exhaustion gracefully
- **User Experience**: Provides feedback when resources are unavailable
- **System Stability**: Prevents system-wide resource starvation

### Mutex Protection
- **Thread Safety**: Protects pool state during concurrent operations
- **Cleanup Coordination**: Ensures proper cleanup during shutdown
- **State Consistency**: Maintains consistent pool state
- **Resource Safety**: Prevents resource leaks during shutdown

## Key Design Decisions

1. **Two Different Pools**: Demonstrates the pattern with different resource types
2. **Health Checking**: Implements resource health validation
3. **Timeout Handling**: Provides timeout mechanisms for resource acquisition
4. **Graceful Shutdown**: Implements proper cleanup and shutdown
5. **Resource Tracking**: Tracks resource usage and health metrics

## Performance Characteristics

### Throughput
- **Resource Reuse**: Eliminates resource creation/destruction overhead
- **Connection Pooling**: Reduces connection establishment time
- **Parallel Processing**: Multiple workers can use resources concurrently
- **Efficient Allocation**: Fast resource allocation from pool

### Latency
- **Immediate Availability**: Pre-populated resources are immediately available
- **Reduced Overhead**: No resource creation overhead for each request
- **Predictable Performance**: Consistent latency regardless of load
- **Connection Reuse**: Reuses existing connections for better performance

### Resource Usage
- **Controlled Memory**: Fixed memory usage regardless of load
- **Efficient Utilization**: Resources are used efficiently across workers
- **Automatic Cleanup**: Stale resources are automatically removed
- **Scalability**: Pool size can be tuned for optimal performance

## Common Use Cases

### Database Operations
- **Connection Pooling**: Manage database connections efficiently
- **Query Optimization**: Reuse connections for multiple queries
- **Connection Limits**: Respect database connection limits
- **Health Monitoring**: Monitor connection health and performance

### HTTP Operations
- **HTTP Client Pooling**: Reuse HTTP clients for multiple requests
- **API Calls**: Efficiently manage API client connections
- **Load Balancing**: Distribute requests across multiple clients
- **Connection Management**: Handle connection limits and timeouts

### File Operations
- **File Handle Pooling**: Manage file handles efficiently
- **I/O Operations**: Reuse file handles for multiple operations
- **Resource Limits**: Respect system file handle limits
- **Performance Optimization**: Reduce file open/close overhead

### Network Operations
- **Socket Pooling**: Manage network socket connections
- **Protocol Handlers**: Reuse protocol handlers for multiple connections
- **Connection Limits**: Handle connection limits and timeouts
- **Load Distribution**: Distribute load across multiple connections

### Memory Management
- **Object Pooling**: Reuse expensive objects
- **Buffer Pooling**: Manage memory buffers efficiently
- **Allocation Optimization**: Reduce memory allocation overhead
- **Garbage Collection**: Minimize garbage collection pressure

### Thread Management
- **Thread Pooling**: Manage worker threads efficiently
- **Task Distribution**: Distribute tasks across worker threads
- **Resource Limits**: Control thread usage and limits
- **Performance Optimization**: Optimize thread creation/destruction

## Advanced Patterns

### Dynamic Pool Sizing
- **Adaptive Sizing**: Adjust pool size based on load
- **Auto-scaling**: Automatically scale pool size up or down
- **Load-based Allocation**: Allocate resources based on current load
- **Performance Monitoring**: Monitor and adjust based on performance metrics

### Resource Health Monitoring
- **Health Checks**: Regular health checks for pooled resources
- **Failure Detection**: Detect and remove failed resources
- **Recovery Mechanisms**: Implement recovery for failed resources
- **Metrics Collection**: Collect health and performance metrics

### Priority-based Pooling
- **Priority Queues**: Prioritize resource allocation
- **Fair Distribution**: Ensure fair resource distribution
- **Resource Reservation**: Reserve resources for high-priority operations
- **Load Balancing**: Balance load across different priority levels

### Distributed Pooling
- **Multi-node Pools**: Coordinate pools across multiple nodes
- **Shared Resources**: Share resources across multiple nodes
- **Load Distribution**: Distribute load across multiple pools
- **Fault Tolerance**: Handle node failures gracefully

## Best Practices

### Pool Configuration
- **Appropriate Size**: Set pool size based on system capacity and requirements
- **Timeout Values**: Set appropriate timeouts for resource acquisition
- **Health Check Intervals**: Regular health checks for resource validation
- **Monitoring**: Monitor pool usage and performance metrics

### Resource Management
- **Proper Cleanup**: Ensure resources are properly cleaned up
- **Leak Prevention**: Prevent resource leaks through proper management
- **Health Monitoring**: Monitor resource health and performance
- **Error Handling**: Handle resource failures gracefully

### Performance Optimization
- **Pool Sizing**: Optimize pool size for your workload
- **Resource Reuse**: Maximize resource reuse for better performance
- **Connection Limits**: Respect system and service connection limits
- **Load Testing**: Test pool performance under various loads

### Monitoring and Alerting
- **Pool Metrics**: Track pool usage, health, and performance
- **Resource Metrics**: Monitor individual resource health and performance
- **Alerting**: Alert on pool exhaustion or resource failures
- **Capacity Planning**: Use metrics for capacity planning

The resource pooling pattern is particularly effective when you have:
- **Expensive Resources**: Resources that are costly to create and destroy
- **Resource Limits**: Limited resources that need efficient management
- **High Concurrency**: Multiple concurrent operations that need resources
- **Performance Requirements**: Need for optimal performance and resource utilization
- **Resource Constraints**: Limited system resources that need careful management

This pattern provides essential tools for building efficient, scalable systems that can handle high concurrency while optimizing resource usage and improving performance. 