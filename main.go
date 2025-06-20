package main

import (
	"flag"
	"fmt"
	"os"

	"concurrency-model-patterns/examples"
)

func main() {
	// Define command line flags
	pipeline := flag.Bool("pipeline", false, "Run pipeline pattern example")
	fan := flag.Bool("fan", false, "Run fan-out/fan-in pattern example")
	pools := flag.Bool("pools", false, "Run worker pools pattern example")
	producerConsumer := flag.Bool("producer-consumer", false, "Run producer-consumer pattern example")
	supervisor := flag.Bool("supervisor", false, "Run supervisor/restart pattern example")
	pubsub := flag.Bool("pubsub", false, "Run publish-subscribe (pub/sub) pattern example")
	timeoutCancellation := flag.Bool("timeout-cancellation", false, "Run timeouts and cancellation pattern example")
	rateLimiting := flag.Bool("rate-limiting", false, "Run rate limiting pattern example")
	mapreduce := flag.Bool("mapreduce", false, "Run MapReduce pattern example")
	singleflight := flag.Bool("singleflight", false, "Run singleflight (spaceflight) pattern example")
	eventLoop := flag.Bool("event-loop", false, "Run event loop pattern example")
	resourcePooling := flag.Bool("resource-pooling", false, "Run resource pooling pattern example")

	// Parse command line flags
	flag.Parse()

	// Check if any flag was provided
	if !*pipeline && !*fan && !*pools && !*producerConsumer && !*supervisor && !*pubsub && !*timeoutCancellation && !*rateLimiting && !*mapreduce && !*singleflight && !*eventLoop && !*resourcePooling {
		fmt.Println("Concurrency Model Patterns Examples")
		fmt.Println("===================================")
		fmt.Println("Usage:")
		fmt.Println("  cmp-pattern --pipeline           - Run pipeline pattern example")
		fmt.Println("  cmp-pattern --fan                - Run fan-out/fan-in pattern example")
		fmt.Println("  cmp-pattern --pools              - Run worker pools pattern example")
		fmt.Println("  cmp-pattern --producer-consumer  - Run producer-consumer pattern example")
		fmt.Println("  cmp-pattern --supervisor         - Run supervisor/restart pattern example")
		fmt.Println("  cmp-pattern --pubsub             - Run publish-subscribe (pub/sub) pattern example")
		fmt.Println("  cmp-pattern --timeout-cancellation - Run timeouts and cancellation pattern example")
		fmt.Println("  cmp-pattern --rate-limiting      - Run rate limiting pattern example")
		fmt.Println("  cmp-pattern --mapreduce          - Run MapReduce pattern example")
		fmt.Println("  cmp-pattern --singleflight       - Run singleflight (spaceflight) pattern example")
		fmt.Println("  cmp-pattern --event-loop         - Run event loop pattern example")
		fmt.Println("  cmp-pattern --resource-pooling   - Run resource pooling pattern example")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  ./cmp-pattern --pipeline")
		fmt.Println("  ./cmp-pattern --fan")
		fmt.Println("  ./cmp-pattern --pools")
		fmt.Println("  ./cmp-pattern --producer-consumer")
		fmt.Println("  ./cmp-pattern --supervisor")
		fmt.Println("  ./cmp-pattern --pubsub")
		fmt.Println("  ./cmp-pattern --timeout-cancellation")
		fmt.Println("  ./cmp-pattern --rate-limiting")
		fmt.Println("  ./cmp-pattern --mapreduce")
		fmt.Println("  ./cmp-pattern --singleflight")
		fmt.Println("  ./cmp-pattern --event-loop")
		fmt.Println("  ./cmp-pattern --resource-pooling")
		os.Exit(1)
	}

	// Run the selected example
	switch {
	case *pipeline:
		fmt.Println("Running Pipeline Pattern Example...")
		examples.RunPipeline()
	case *fan:
		fmt.Println("Running Fan-out/Fan-in Pattern Example...")
		examples.RunFan()
	case *pools:
		fmt.Println("Running Worker Pools Pattern Example...")
		examples.RunPools()
	case *producerConsumer:
		fmt.Println("Running Producer-Consumer Pattern Example...")
		examples.RunProducerConsumer()
	case *supervisor:
		fmt.Println("Running Supervisor/Restart Pattern Example...")
		examples.RunSupervisor()
	case *pubsub:
		fmt.Println("Running Publish-Subscribe (Pub/Sub) Pattern Example...")
		examples.RunPubSub()
	case *timeoutCancellation:
		fmt.Println("Running Timeouts and Cancellation Pattern Example...")
		examples.RunTimeoutCancellation()
	case *rateLimiting:
		fmt.Println("Running Rate Limiting Pattern Example...")
		examples.RunRateLimiting()
	case *mapreduce:
		fmt.Println("Running MapReduce Pattern Example...")
		examples.RunMapReduce()
	case *singleflight:
		fmt.Println("Running Singleflight (Spaceflight) Pattern Example...")
		examples.RunSingleflight()
	case *eventLoop:
		fmt.Println("Running Event Loop Pattern Example...")
		examples.RunEventLoop()
	case *resourcePooling:
		fmt.Println("Running Resource Pooling Pattern Example...")
		examples.RunResourcePooling()
	}
}
