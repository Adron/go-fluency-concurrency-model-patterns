# Supervisor/Restart Pattern

## Overview

The Supervisor/Restart pattern involves a supervisor goroutine that monitors one or more worker goroutines. If a worker fails (panics or exits unexpectedly), the supervisor restarts it. This pattern is essential for building resilient and fault-tolerant systems, automatically recovering from transient errors, and ensuring critical tasks are always running.

## Implementation Details

### Structure

The supervisor/restart implementation in `examples/supervisor.go` consists of three main components:

1. **Supervisor** - A goroutine that monitors and manages workers
2. **Worker** - A goroutine that performs work and may fail
3. **Coordination** - Channels and synchronization for monitoring and restart

### Code Analysis

```go
func RunSupervisor() {
    var restarts int32
    stop := make(chan struct{})
    done := make(chan struct{})

    // Supervisor goroutine
    go func() {
        for {
            workerDone := make(chan struct{})
            go workerWithFailure(workerDone, stop)
            select {
            case <-workerDone:
                atomic.AddInt32(&restarts, 1)
                fmt.Println("Supervisor: Worker failed, restarting...")
                // Restart after a short delay
                time.Sleep(500 * time.Millisecond)
            case <-stop:
                fmt.Println("Supervisor: Stopping worker supervision.")
                close(done)
                return
            }
        }
    }()

    // Let the supervisor run for a while
    time.Sleep(4 * time.Second)
    close(stop)
    <-done

    fmt.Printf("Supervisor example completed! Worker was restarted %d times.\n", restarts-1)
}
```

### Worker Implementation

```go
func workerWithFailure(done chan<- struct{}, stop <-chan struct{}) {
    fmt.Println("Worker: Started")
    workTime := time.Duration(rand.Intn(1200)+400) * time.Millisecond
    select {
    case <-time.After(workTime):
        // Simulate random failure
        if rand.Float32() < 0.6 {
            fmt.Println("Worker: Simulated failure!")
            done <- struct{}{}
            return
        }
        fmt.Println("Worker: Completed work successfully.")
    case <-stop:
        fmt.Println("Worker: Received stop signal.")
    }
    // Signal normal exit
    done <- struct{}{}
}
```

## How It Works

1. **Supervisor Initialization**: The supervisor starts and launches the first worker
2. **Worker Execution**: The worker performs its task and may fail randomly
3. **Failure Detection**: When a worker fails, it signals the supervisor via a channel
4. **Restart Logic**: The supervisor waits briefly, then launches a new worker
5. **Continuous Monitoring**: This cycle continues until the supervisor receives a stop signal
6. **Graceful Shutdown**: The supervisor stops launching new workers and waits for current workers to finish

The pattern ensures that critical work continues even when individual workers fail.

## Why This Implementation?

### Channel-based Communication
- **Non-blocking Monitoring**: Supervisor can monitor multiple workers without blocking
- **Immediate Failure Detection**: Workers can signal failure immediately
- **Clean Coordination**: Channels provide natural synchronization

### Atomic Counter for Restarts
- **Thread Safety**: Atomic operations ensure accurate restart counting
- **Performance**: Atomic operations are more efficient than mutex-based counting
- **Simplicity**: No need for complex synchronization around the counter

### Separate Stop Channel
- **Graceful Shutdown**: Supervisor can stop monitoring without waiting for failures
- **Clean Termination**: Workers receive stop signals and can clean up properly
- **Resource Management**: Prevents infinite restart loops

### Worker Done Channel
- **Failure Signaling**: Workers signal completion or failure to supervisor
- **Immediate Response**: Supervisor can react immediately to worker state changes
- **Flexible Communication**: Can distinguish between normal completion and failure

### Restart Delay
- **Prevent Thundering Herd**: Brief delay prevents rapid restart loops
- **Resource Recovery**: Gives system time to recover from transient issues
- **Stability**: Prevents overwhelming the system with restart attempts

## Key Design Decisions

1. **Simulated Failures**: Random failures (60% probability) demonstrate the restart mechanism
2. **Variable Work Time**: Random work duration simulates real-world variability
3. **Immediate Restart**: Supervisor restarts workers immediately after failure detection
4. **Stop Signal Handling**: Both supervisor and workers respond to stop signals
5. **Restart Counting**: Tracks the number of restarts for monitoring purposes

## Failure Handling Strategies

### Immediate Restart
- **Pros**: Minimal downtime, immediate recovery
- **Cons**: May restart into the same failure condition
- **Use Case**: Transient failures, network issues

### Exponential Backoff
- **Pros**: Prevents overwhelming the system, allows for recovery
- **Cons**: Longer downtime for legitimate failures
- **Use Case**: Persistent failures, resource exhaustion

### Circuit Breaker
- **Pros**: Prevents cascading failures, allows system recovery
- **Cons**: More complex implementation
- **Use Case**: External service failures, resource constraints

### Health Checks
- **Pros**: Proactive failure detection, better monitoring
- **Cons**: Additional complexity, false positives
- **Use Case**: Long-running processes, critical systems

## Common Use Cases

### System Services
- **Background Workers**: Ensure critical background tasks continue running
- **Data Processing**: Restart failed data processing jobs
- **Scheduled Tasks**: Ensure scheduled operations complete successfully

### Microservices
- **Service Health**: Monitor and restart unhealthy service instances
- **API Endpoints**: Ensure critical API endpoints remain available
- **Message Processors**: Restart failed message processing workers

### Monitoring Systems
- **Health Monitors**: Ensure monitoring agents continue running
- **Alert Processors**: Restart failed alert processing workers
- **Metric Collectors**: Ensure metric collection continues

### Database Operations
- **Connection Managers**: Restart failed database connection pools
- **Migration Workers**: Ensure database migrations complete
- **Backup Jobs**: Restart failed backup operations

### File Processing
- **Upload Handlers**: Restart failed file upload processors
- **Download Workers**: Ensure file downloads complete
- **Processing Jobs**: Restart failed file processing tasks

### Network Services
- **Connection Handlers**: Restart failed network connection handlers
- **Protocol Processors**: Ensure protocol processing continues
- **Load Balancers**: Restart failed load balancing workers

### IoT Applications
- **Device Managers**: Restart failed device communication workers
- **Sensor Processors**: Ensure sensor data processing continues
- **Control Systems**: Restart failed control system workers

### Financial Systems
- **Trading Engines**: Restart failed trading algorithm workers
- **Risk Calculators**: Ensure risk calculations continue
- **Settlement Processors**: Restart failed settlement workers

## Best Practices

### Failure Detection
- **Immediate Signaling**: Workers should signal failure immediately
- **Graceful Degradation**: Workers should handle partial failures gracefully
- **Health Reporting**: Workers should report their health status

### Restart Strategy
- **Appropriate Delays**: Use delays to prevent restart storms
- **Backoff Policies**: Implement exponential backoff for persistent failures
- **Maximum Restarts**: Limit the number of restarts to prevent infinite loops

### Resource Management
- **Cleanup**: Ensure workers clean up resources before exiting
- **Memory Management**: Monitor memory usage to prevent leaks
- **Connection Management**: Properly close connections and file handles

### Monitoring and Logging
- **Restart Tracking**: Log restart events for monitoring
- **Performance Metrics**: Track worker performance and restart frequency
- **Alerting**: Alert on excessive restart rates

The supervisor/restart pattern is particularly effective when you have:
- **Critical Processes**: Tasks that must continue running
- **Unreliable Workers**: Workers that may fail due to external factors
- **Long-running Systems**: Systems that need to operate continuously
- **Fault Tolerance**: Requirements for high availability and reliability
- **Automatic Recovery**: Need for self-healing systems

This pattern provides a robust foundation for building resilient systems that can automatically recover from failures and maintain continuous operation. 