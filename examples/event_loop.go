package examples

import (
	"fmt"
	"math/rand"
	"time"
)

// RunEventLoop demonstrates the event loop pattern.
func RunEventLoop() {
	fmt.Println("=== Event Loop Pattern Example ===")

	// Create event channels
	userEvents := make(chan string, 10)
	systemEvents := make(chan string, 10)
	timerEvents := make(chan string, 10)
	shutdown := make(chan struct{})

	// Start event producers
	go userEventProducer(userEvents)
	go systemEventProducer(systemEvents)
	go timerEventProducer(timerEvents)

	// Start the event loop
	go eventLoop(userEvents, systemEvents, timerEvents, shutdown)

	// Let the system run for a while
	time.Sleep(5 * time.Second)

	// Shutdown
	fmt.Println("Shutting down event loop...")
	close(shutdown)

	// Wait a bit for cleanup
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Event loop example completed!")
}

// Event loop that processes events from multiple sources
func eventLoop(userEvents, systemEvents, timerEvents <-chan string, shutdown <-chan struct{}) {
	fmt.Println("Event loop started...")

	for {
		select {
		case event := <-userEvents:
			fmt.Printf("Event Loop: Processing user event: %s\n", event)
			processUserEvent(event)

		case event := <-systemEvents:
			fmt.Printf("Event Loop: Processing system event: %s\n", event)
			processSystemEvent(event)

		case event := <-timerEvents:
			fmt.Printf("Event Loop: Processing timer event: %s\n", event)
			processTimerEvent(event)

		case <-shutdown:
			fmt.Println("Event Loop: Shutdown signal received, cleaning up...")
			return
		}
	}
}

// Event producers
func userEventProducer(events chan<- string) {
	userActions := []string{"login", "logout", "click", "scroll", "submit"}
	for i := 0; i < 8; i++ {
		time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)
		action := userActions[rand.Intn(len(userActions))]
		events <- fmt.Sprintf("%s (user_%d)", action, i+1)
	}
}

func systemEventProducer(events chan<- string) {
	systemEvents := []string{"backup", "update", "maintenance", "alert", "sync"}
	for i := 0; i < 6; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		event := systemEvents[rand.Intn(len(systemEvents))]
		events <- fmt.Sprintf("%s (system_%d)", event, i+1)
	}
}

func timerEventProducer(events chan<- string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		if count >= 5 {
			break
		}
		events <- fmt.Sprintf("heartbeat (timer_%d)", count+1)
		count++
	}
}

// Event processors
func processUserEvent(event string) {
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("  -> User event processed: %s\n", event)
}

func processSystemEvent(event string) {
	// Simulate processing time
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("  -> System event processed: %s\n", event)
}

func processTimerEvent(event string) {
	// Simulate processing time
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("  -> Timer event processed: %s\n", event)
}
