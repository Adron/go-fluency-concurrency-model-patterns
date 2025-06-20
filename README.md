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
    └── producer_consumer.go # Producer-consumer pattern implementation
    └── supervisor.go       # Supervisor/restart pattern implementation
    └── pubsub.go           # Publish-subscribe (pub/sub) pattern implementation
    └── timeout_cancellation.go # Timeouts and cancellation pattern implementation
    └── rate_limiting.go        # Rate limiting pattern implementation
    └── mapreduce.go            # MapReduce pattern implementation
    └── singleflight.go         # Singleflight (spaceflight) pattern implementation
    └── event_loop.go           # Event loop pattern implementation
    └── resource_pooling.go     # Resource pooling pattern implementation
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

### Producer-Consumer Pattern
```bash
./cmp-pattern --producer-consumer
```
Demonstrates multiple producers generating data and multiple consumers processing data from a shared buffer (channel):
- 2 producers generate random numbers
- 3 consumers process the numbers
- Bounded buffer (channel) for synchronization

### Supervisor/Restart Pattern
```bash
./cmp-pattern --supervisor
```
Demonstrates a supervisor goroutine that monitors a worker goroutine and restarts it if it fails:
- Supervisor launches a worker that may randomly fail
- If the worker fails, the supervisor restarts it
- After a set time, the supervisor stops monitoring

### Publish-Subscribe (Pub/Sub) Pattern
```bash
./cmp-pattern --pubsub
```
Demonstrates a publisher sending messages to multiple subscribers:
- Subscribers register to receive messages
- Publisher broadcasts messages to all subscribers
- All subscribers receive each message

### Timeouts and Cancellation Pattern
```bash
./cmp-pattern --timeout-cancellation
```
Demonstrates handling timeouts and cancellation using context and channels:
- Context-based timeout for long-running tasks
- Channel-based timeout using time.After
- Context cancellation to stop operations

### Rate Limiting Pattern
```bash
./cmp-pattern --rate-limiting
```
Demonstrates different rate limiting techniques:
- Fixed rate limiting using time.Ticker
- Token bucket rate limiting with burst capacity
- Controlling request frequency and resource usage

### MapReduce Pattern
```bash
./cmp-pattern --mapreduce
```
Demonstrates the MapReduce pattern for distributed data processing:
- Map phase: process data in parallel and emit key-value pairs
- Shuffle phase: group data by key
- Reduce phase: aggregate results for each key
- Word count example with concurrent processing

### Singleflight (Spaceflight) Pattern
```bash
./cmp-pattern --singleflight
```
Demonstrates the singleflight pattern that ensures only one execution of a function is in-flight for a given key:
- Multiple concurrent requests for the same key
- Only one execution runs, others wait for the result
- Prevents duplicate expensive operations
- Useful for caching and deduplication

### Event Loop Pattern
```bash
./cmp-pattern --event-loop
```
Demonstrates an event loop that processes events from multiple sources:
- Single goroutine event loop using select
- Multiple event producers (user, system, timer)
- Centralized event processing and dispatching
- Graceful shutdown handling

### Resource Pooling Pattern
```bash
./cmp-pattern --resource-pooling
```
Demonstrates resource pooling to manage expensive resources efficiently:
- Database connection pooling with reuse
- HTTP client pooling for API requests
- Pre-populated pools with maximum size limits
- Automatic resource creation and cleanup

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

### Running the Producer-Consumer Example
```bash
$ ./cmp-pattern --producer-consumer
Running Producer-Consumer Pattern Example...
=== Producer-Consumer Pattern Example ===
Producer 1 produced: 42
Producer 2 produced: 17
Consumer 1 consumed: 42
Producer 1 produced: 88
Consumer 2 consumed: 17
Producer 2 produced: 53
Consumer 3 consumed: 88
...
Producer-Consumer example completed!
```

### Running the Supervisor/Restart Example
```bash
$ ./cmp-pattern --supervisor
Running Supervisor/Restart Pattern Example...
=== Supervisor/Restart Pattern Example ===
Worker: Started
Worker: Simulated failure!
Supervisor: Worker failed, restarting...
Worker: Started
Worker: Completed work successfully.
Supervisor: Worker failed, restarting...
...
Supervisor example completed! Worker was restarted 3 times.
```

### Running the Publish-Subscribe (Pub/Sub) Example
```bash
$ ./cmp-pattern --pubsub
Running Publish-Subscribe (Pub/Sub) Pattern Example...
=== Publish-Subscribe (Pub/Sub) Pattern Example ===
Publisher sending: Message 1
Subscriber 1 received: Message 1
Subscriber 2 received: Message 1
Subscriber 3 received: Message 1
Publisher sending: Message 2
Subscriber 1 received: Message 2
Subscriber 2 received: Message 2
Subscriber 3 received: Message 2
...
Pub/Sub example completed!
```

### Running the Timeouts and Cancellation Example
```bash
$ ./cmp-pattern --timeout-cancellation
Running Timeouts and Cancellation Pattern Example...
=== Timeouts and Cancellation Pattern Example ===

1. Context-based timeout example:
Starting long task (will take 2.3s)...
Task timed out: context deadline exceeded

2. Channel-based timeout example:
Channel task timed out

3. Context cancellation example:
Cancelling context...
Context cancelled: context canceled

Timeouts and Cancellation example completed!
```

### Running the Rate Limiting Example
```bash
$ ./cmp-pattern --rate-limiting
Running Rate Limiting Pattern Example...
=== Rate Limiting Pattern Example ===

1. Fixed rate limiting (2 requests per second):
Request 1 processed at 15:04:05.000
Request 2 processed at 15:04:05.500
Request 3 processed at 15:04:06.000
Request 4 processed at 15:04:06.500
Request 5 processed at 15:04:07.000
Request 6 processed at 15:04:07.500

2. Token bucket rate limiting (3 tokens per second, burst of 5):
Token request 1 granted at 15:04:07.500
Token request 2 granted at 15:04:07.500
Token request 3 granted at 15:04:07.500
Token request 4 granted at 15:04:07.500
Token request 5 granted at 15:04:07.500
Token request 6 denied at 15:04:07.500
...

Rate Limiting example completed!
```

### Running the MapReduce Example
```bash
$ ./cmp-pattern --mapreduce
Running MapReduce Pattern Example...
=== MapReduce Pattern Example ===
Input data: [hello world hello go world of concurrency go programming concurrency patterns hello concurrency go world patterns in go]
Map: emitted (hello, 1)
Map: emitted (world, 1)
Map: emitted (hello, 1)
Map: emitted (go, 1)
...
Shuffle: grouped hello -> [1 1 1]
Shuffle: grouped world -> [1 1 1]
...
Reduce: hello -> 3
Reduce: world -> 3
Reduce: go -> 4
...

Word count results:
  hello: 3
  world: 3
  go: 4
  concurrency: 3
  patterns: 2
  of: 1
  programming: 1
  in: 1

MapReduce example completed!
```

### Running the Singleflight (Spaceflight) Example
```bash
$ ./cmp-pattern --singleflight
Running Singleflight (Spaceflight) Pattern Example...
=== Singleflight (Spaceflight) Pattern Example ===
Making 5 concurrent requests for key: user:123
Request 0: Starting...
Request 1: Starting...
Request 2: Starting...
Request 3: Starting...
Request 4: Starting...
Request 0: Executing expensive operation...
Duplicate call for key user:123, waiting for result...
Duplicate call for key user:123, waiting for result...
Duplicate call for key user:123, waiting for result...
Duplicate call for key user:123, waiting for result...
Request 0: Completed with result: Data for user:123 (processed by request 0)
Request 1: Completed with result: Data for user:123 (processed by request 0)
Request 2: Completed with result: Data for user:123 (processed by request 0)
Request 3: Completed with result: Data for user:123 (processed by request 0)
Request 4: Completed with result: Data for user:123 (processed by request 0)

All results should be identical:
  Request 0: Data for user:123 (processed by request 0)
  Request 1: Data for user:123 (processed by request 0)
  Request 2: Data for user:123 (processed by request 0)
  Request 3: Data for user:123 (processed by request 0)
  Request 4: Data for user:123 (processed by request 0)

Singleflight example completed!
```

### Running the Event Loop Example
```bash
$ ./cmp-pattern --event-loop
Running Event Loop Pattern Example...
=== Event Loop Pattern Example ===
Event loop started...
Event Loop: Processing timer event: heartbeat (timer_1)
  -> Timer event processed: heartbeat (timer_1)
Event Loop: Processing user event: login (user_1)
  -> User event processed: login (user_1)
Event Loop: Processing system event: backup (system_1)
  -> System event processed: backup (system_1)
Event Loop: Processing timer event: heartbeat (timer_2)
  -> Timer event processed: heartbeat (timer_2)
...
Shutting down event loop...
Event Loop: Shutdown signal received, cleaning up...
Event loop example completed!
```

### Running the Resource Pooling Example
```bash
$ ./cmp-pattern --resource-pooling
Running Resource Pooling Pattern Example...
=== Resource Pooling Pattern Example ===

1. Database Connection Pool Example:
Worker 1: Got DB connection 1
Worker 2: Got DB connection 2
Worker 3: Got DB connection 3
Worker 4: Got DB connection 4
Worker 1: Executing query on connection 1
Worker 2: Executing query on connection 2
Worker 1: Released DB connection 1
Worker 2: Released DB connection 2
Worker 5: Got DB connection 1
Worker 5: Executing query on connection 1
...
DB pool closed. Total connections created: 5

2. HTTP Client Pool Example:
Worker 1: Got HTTP client 1
Worker 2: Got HTTP client 2
Worker 3: Got HTTP client 3
Worker 1: Making API request with client 1
Worker 1: Released HTTP client 1
Worker 4: Got HTTP client 1
...
HTTP client pool closed. Total clients created: 4

Resource Pooling example completed!
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

### Producer-Consumer Pattern
The producer-consumer pattern decouples the production of data from its consumption using a shared buffer (channel). Multiple producers can generate data and send it to the buffer, while multiple consumers read from the buffer and process the data. This pattern is useful for:
- Buffering bursts of work
- Decoupling work rates between producers and consumers
- Load leveling and resource management

### Supervisor/Restart Pattern
The supervisor/restart pattern involves a supervisor goroutine that monitors one or more worker goroutines. If a worker fails (panics or exits unexpectedly), the supervisor restarts it. This pattern is useful for:
- Building resilient and fault-tolerant systems
- Automatically recovering from transient errors
- Ensuring critical tasks are always running

### Publish-Subscribe (Pub/Sub) Pattern
The publish-subscribe (pub/sub) pattern decouples message producers (publishers) from consumers (subscribers). Publishers send messages to a topic or channel, and all subscribers to that topic receive the messages. This pattern is useful for:
- Event-driven architectures
- Broadcasting messages to multiple consumers
- Decoupling senders and receivers

### Timeouts and Cancellation Pattern
The timeouts and cancellation pattern uses context.Context or channels to control the lifecycle of concurrent operations. This pattern is essential for:
- Preventing goroutines from running indefinitely
- Graceful shutdown of operations
- Resource cleanup and management
- Building responsive applications

### Rate Limiting Pattern
The rate limiting pattern controls the frequency of operations to prevent resource exhaustion and ensure fair usage. This pattern is useful for:
- API request throttling
- Resource protection
- Preventing DoS attacks
- Ensuring system stability under load

### MapReduce Pattern
The MapReduce pattern processes large datasets by breaking the work into two phases: Map (process data in parallel) and Reduce (aggregate results). This pattern is useful for:
- Processing large datasets in parallel
- Distributed computing
- Data analytics and aggregation
- Batch processing jobs

### Singleflight (Spaceflight) Pattern
The singleflight pattern ensures that only one execution of a function is in-flight for a given key at a time. Duplicate callers wait for the result instead of executing the function again. This pattern is useful for:
- Preventing duplicate expensive operations
- Caching with concurrent access
- API call deduplication
- Database query optimization

### Event Loop Pattern
The event loop pattern uses a single goroutine to wait for and dispatch events from multiple sources using select. This pattern is useful for:
- Network servers and clients
- GUI applications
- State machines
- Centralized event processing
- Non-blocking I/O handling

### Resource Pooling Pattern
The resource pooling pattern maintains a pool of reusable resources to avoid the overhead of creating and destroying them. This pattern is useful for:
- Database connection management
- HTTP client reuse
- File handle pooling
- Memory allocation optimization
- Reducing resource creation overhead

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