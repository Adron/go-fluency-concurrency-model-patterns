# Publish-Subscribe (Pub/Sub) Pattern

## Overview

The Publish-Subscribe (Pub/Sub) pattern decouples message producers (publishers) from consumers (subscribers). Publishers send messages to a topic or channel, and all subscribers to that topic receive the messages. This pattern is essential for event-driven architectures, broadcasting messages to multiple consumers, and decoupling senders and receivers.

## Implementation Details

### Structure

The publish-subscribe implementation in `examples/pubsub.go` consists of three main components:

1. **Broadcaster** - Manages subscriptions and message distribution
2. **Publishers** - Send messages to the broadcaster
3. **Subscribers** - Receive messages from the broadcaster

### Code Analysis

```go
func RunPubSub() {
    // Create a broadcaster
    b := newBroadcaster()

    numSubscribers := 3
    var wg sync.WaitGroup

    // Start subscribers
    for i := 1; i <= numSubscribers; i++ {
        ch := b.subscribe()
        wg.Add(1)
        go func(id int, ch <-chan string) {
            defer wg.Done()
            for msg := range ch {
                fmt.Printf("Subscriber %d received: %s\n", id, msg)
            }
            fmt.Printf("Subscriber %d done.\n", id)
        }(i, ch)
    }

    // Start publisher
    go func() {
        for i := 1; i <= 5; i++ {
            msg := fmt.Sprintf("Message %d", i)
            fmt.Printf("Publisher sending: %s\n", msg)
            b.publish(msg)
            time.Sleep(400 * time.Millisecond)
        }
        b.close()
    }()

    wg.Wait()
}
```

### Broadcaster Implementation

```go
type broadcaster struct {
    subscribers []chan string
    closed     bool
    mu         sync.Mutex
}

func newBroadcaster() *broadcaster {
    return &broadcaster{
        subscribers: make([]chan string, 0),
    }
}

func (b *broadcaster) subscribe() <-chan string {
    b.mu.Lock()
    defer b.mu.Unlock()
    ch := make(chan string, 2)
    b.subscribers = append(b.subscribers, ch)
    return ch
}

func (b *broadcaster) publish(msg string) {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.closed {
        return
    }
    for _, ch := range b.subscribers {
        ch <- msg
    }
}

func (b *broadcaster) close() {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.closed {
        return
    }
    for _, ch := range b.subscribers {
        close(ch)
    }
    b.closed = true
}
```

## How It Works

1. **Broadcaster Creation**: A broadcaster is created to manage subscriptions and message distribution
2. **Subscriber Registration**: Subscribers call `subscribe()` to get their own message channel
3. **Message Publishing**: Publishers call `publish()` to send messages to all subscribers
4. **Message Distribution**: The broadcaster sends each message to all subscriber channels
5. **Graceful Shutdown**: When the broadcaster is closed, all subscriber channels are closed

The pattern enables loose coupling between publishers and subscribers, allowing for flexible message distribution.

## Why This Implementation?

### Individual Subscriber Channels
- **Isolation**: Each subscriber has its own channel, preventing blocking
- **Independent Processing**: Subscribers can process messages at their own pace
- **Scalability**: Easy to add or remove subscribers without affecting others

### Mutex-based Synchronization
- **Thread Safety**: Protects subscriber list during concurrent operations
- **Simple Implementation**: Straightforward synchronization for the subscriber list
- **Performance**: Minimal overhead for the typical pub/sub use case

### Buffered Subscriber Channels
- **Non-blocking Publishing**: Publishers don't block if subscribers are slow
- **Message Buffering**: Can handle temporary subscriber slowdowns
- **Flow Control**: Natural backpressure through channel buffering

### Centralized Message Distribution
- **Broadcast Semantics**: All subscribers receive all messages
- **Simple Coordination**: Single point of message distribution
- **Consistent Delivery**: All subscribers receive messages in the same order

### Graceful Shutdown
- **Clean Termination**: Properly closes all subscriber channels
- **Resource Cleanup**: Prevents goroutine leaks
- **Coordinated Exit**: All subscribers exit when broadcaster closes

## Key Design Decisions

1. **Slice-based Subscriber Management**: Simple and efficient for typical use cases
2. **Read-only Channel Returns**: Subscribers receive read-only channels for safety
3. **Buffered Channels**: Prevents blocking during message distribution
4. **Mutex Protection**: Ensures thread-safe subscriber management
5. **Immediate Message Distribution**: Messages are sent to all subscribers immediately

## Performance Characteristics

### Throughput
- **Linear Scaling**: Throughput scales linearly with the number of subscribers
- **Channel Overhead**: Each subscriber adds minimal overhead
- **Publishing Performance**: Publishing is O(n) where n is the number of subscribers

### Latency
- **Immediate Distribution**: Messages are distributed immediately upon publishing
- **Subscriber Processing**: Latency depends on individual subscriber processing speed
- **Channel Buffering**: Buffered channels can absorb temporary processing delays

### Memory Usage
- **Per-Subscriber Overhead**: Each subscriber requires a channel and goroutine
- **Message Buffering**: Buffered channels use additional memory
- **Scalability**: Memory usage grows linearly with subscriber count

## Common Use Cases

### Event-Driven Systems
- **User Interface Events**: Mouse clicks, keyboard input, window events
- **System Events**: File changes, network events, timer events
- **Application Events**: User actions, state changes, error conditions

### Message Broadcasting
- **Chat Applications**: Broadcast messages to all connected users
- **Notification Systems**: Send notifications to multiple recipients
- **Status Updates**: Broadcast system status to all monitoring clients

### Data Distribution
- **Real-time Data**: Distribute sensor data, market data, or metrics
- **Configuration Changes**: Broadcast configuration updates to all services
- **Cache Invalidation**: Notify all cache instances of data changes

### Microservices Communication
- **Service Discovery**: Notify services of new service instances
- **Health Updates**: Broadcast health status changes
- **Load Balancing**: Distribute load information to all instances

### Monitoring and Logging
- **Metric Collection**: Distribute metrics to multiple monitoring systems
- **Log Aggregation**: Send logs to multiple destinations
- **Alert Distribution**: Broadcast alerts to multiple notification channels

### Gaming and Real-time Applications
- **Game State Updates**: Broadcast game state changes to all players
- **Multiplayer Events**: Distribute player actions to all participants
- **Real-time Collaboration**: Share changes across multiple users

### IoT and Sensor Networks
- **Sensor Data**: Distribute sensor readings to multiple processors
- **Device Status**: Broadcast device status changes
- **Control Commands**: Distribute control commands to multiple devices

### Financial Systems
- **Market Data**: Distribute real-time market data to multiple traders
- **Trade Updates**: Broadcast trade execution to all relevant systems
- **Risk Alerts**: Distribute risk alerts to multiple risk managers

## Advanced Patterns

### Topic-based Pub/Sub
- **Multiple Topics**: Support for different message categories
- **Selective Subscription**: Subscribers can choose which topics to receive
- **Hierarchical Topics**: Support for topic hierarchies (e.g., "sports.football.nfl")

### Filtered Pub/Sub
- **Message Filtering**: Subscribers can filter messages based on content
- **Conditional Delivery**: Only deliver messages that match subscriber criteria
- **Performance Optimization**: Reduce unnecessary message processing

### Persistent Pub/Sub
- **Message Persistence**: Store messages for offline subscribers
- **Reliable Delivery**: Ensure messages are delivered even if subscribers are temporarily unavailable
- **Message Ordering**: Maintain message order across subscriber reconnections

### Distributed Pub/Sub
- **Multi-node Distribution**: Distribute messages across multiple nodes
- **Load Balancing**: Balance message distribution across multiple brokers
- **Fault Tolerance**: Continue operation even if some nodes fail

The publish-subscribe pattern is particularly effective when you have:
- **Loose Coupling**: Components that should be independent of each other
- **Multiple Consumers**: Multiple components that need the same information
- **Event-driven Architecture**: Systems that respond to events rather than direct calls
- **Scalable Communication**: Need to add/remove consumers without affecting publishers
- **Broadcast Requirements**: Need to send the same message to multiple recipients

This pattern provides a flexible foundation for building decoupled, event-driven systems that can scale to handle multiple publishers and subscribers efficiently. 