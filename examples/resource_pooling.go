package examples

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RunResourcePooling demonstrates the resource pooling pattern.
func RunResourcePooling() {
	fmt.Println("=== Resource Pooling Pattern Example ===")

	// Example 1: Database Connection Pool
	fmt.Println("\n1. Database Connection Pool Example:")
	dbPool := newDBConnectionPool(3, 5)

	var wg sync.WaitGroup
	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn := dbPool.getConnection()
			fmt.Printf("Worker %d: Got DB connection %d\n", id, conn.id)

			// Simulate database operation
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
			fmt.Printf("Worker %d: Executing query on connection %d\n", id, conn.id)

			dbPool.releaseConnection(conn)
			fmt.Printf("Worker %d: Released DB connection %d\n", id, conn.id)
		}(i)
	}
	wg.Wait()
	dbPool.close()

	// Example 2: HTTP Client Pool
	fmt.Println("\n2. HTTP Client Pool Example:")
	clientPool := newHTTPClientPool(2, 4)

	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := clientPool.getClient()
			fmt.Printf("Worker %d: Got HTTP client %d\n", id, client.id)

			// Simulate API request
			time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
			fmt.Printf("Worker %d: Making API request with client %d\n", id, client.id)

			clientPool.releaseClient(client)
			fmt.Printf("Worker %d: Released HTTP client %d\n", id, client.id)
		}(i)
	}
	wg.Wait()
	clientPool.close()

	fmt.Println("\nResource Pooling example completed!")
}

// Database Connection Pool
type dbConnection struct {
	id       int
	lastUsed time.Time
}

type dbConnectionPool struct {
	connections chan *dbConnection
	maxSize     int
	created     int
	mu          sync.Mutex
}

func newDBConnectionPool(initial, maxSize int) *dbConnectionPool {
	pool := &dbConnectionPool{
		connections: make(chan *dbConnection, maxSize),
		maxSize:     maxSize,
	}

	// Pre-populate with initial connections
	for i := 0; i < initial; i++ {
		pool.connections <- &dbConnection{
			id:       i + 1,
			lastUsed: time.Now(),
		}
		pool.created++
	}

	return pool
}

func (p *dbConnectionPool) getConnection() *dbConnection {
	select {
	case conn := <-p.connections:
		conn.lastUsed = time.Now()
		return conn
	default:
		// Create new connection if pool is empty and under max size
		p.mu.Lock()
		if p.created < p.maxSize {
			p.created++
			p.mu.Unlock()
			return &dbConnection{
				id:       p.created,
				lastUsed: time.Now(),
			}
		}
		p.mu.Unlock()

		// Wait for available connection
		conn := <-p.connections
		conn.lastUsed = time.Now()
		return conn
	}
}

func (p *dbConnectionPool) releaseConnection(conn *dbConnection) {
	select {
	case p.connections <- conn:
		// Successfully returned to pool
	default:
		// Pool is full, discard connection
		fmt.Printf("Pool full, discarding connection %d\n", conn.id)
	}
}

func (p *dbConnectionPool) close() {
	close(p.connections)
	fmt.Printf("DB pool closed. Total connections created: %d\n", p.created)
}

// HTTP Client Pool
type httpClient struct {
	id       int
	lastUsed time.Time
}

type httpClientPool struct {
	clients chan *httpClient
	maxSize int
	created int
	mu      sync.Mutex
}

func newHTTPClientPool(initial, maxSize int) *httpClientPool {
	pool := &httpClientPool{
		clients: make(chan *httpClient, maxSize),
		maxSize: maxSize,
	}

	// Pre-populate with initial clients
	for i := 0; i < initial; i++ {
		pool.clients <- &httpClient{
			id:       i + 1,
			lastUsed: time.Now(),
		}
		pool.created++
	}

	return pool
}

func (p *httpClientPool) getClient() *httpClient {
	select {
	case client := <-p.clients:
		client.lastUsed = time.Now()
		return client
	default:
		// Create new client if pool is empty and under max size
		p.mu.Lock()
		if p.created < p.maxSize {
			p.created++
			p.mu.Unlock()
			return &httpClient{
				id:       p.created,
				lastUsed: time.Now(),
			}
		}
		p.mu.Unlock()

		// Wait for available client
		client := <-p.clients
		client.lastUsed = time.Now()
		return client
	}
}

func (p *httpClientPool) releaseClient(client *httpClient) {
	select {
	case p.clients <- client:
		// Successfully returned to pool
	default:
		// Pool is full, discard client
		fmt.Printf("Pool full, discarding HTTP client %d\n", client.id)
	}
}

func (p *httpClientPool) close() {
	close(p.clients)
	fmt.Printf("HTTP client pool closed. Total clients created: %d\n", p.created)
}
