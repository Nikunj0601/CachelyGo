# CachelyGo - Distributed In-Memory Cache with Consistent Reads

A Go implementation of a distributed in-memory cache system with load balancing and consistency guarantees for read operations. The system implements a distributed architecture where multiple cache nodes work together to provide scalable and reliable caching services.

## Features

- Distributed caching across multiple nodes
- Load balancing for distributing requests
- Consistent reads during synchronization
- TTL (Time-To-Live) support for cache entries
- Automatic synchronization between nodes
- Round-robin load balancing
- HTTP API interface

## Installation

1. Clone the repository
2. Ensure Go is installed on your system


## Quick Start

### Start Load Balancer
```bash
go run main.go -mode=load_balancer -nodes=localhost:8081,localhost:8082,localhost:8083

```

### Start Cache Nodes
```bash
go run main.go -mode=node -port=8081 -nodes=localhost:8081,localhost:8082,localhost:8083
go run main.go -mode=node -port=8082 -nodes=localhost:8081,localhost:8082,localhost:8083
go run main.go -mode=node -port=8083 -nodes=localhost:8081,localhost:8082,localhost:8083

```

## API Usage

### Set Cache Entry
```http
POST /cache/set?key=mykey&value=myvalue&ttl=3600
```

### Get Cache Entry
```http
GET /cache/get?key=mykey
```

## Configuration

### Command Line Arguments
- `-mode`: Operating mode ("node" or "load_balancer")
- `-port`: Port number for the server
- `-nodes`: Comma-separated list of cache node addresses

### Default Settings
- Load Balancer Port: 8080
- Default Node Ports: 8081, 8082, 8083
- Default Node List: localhost:8081,localhost:8082,localhost:8083

## Architecture

### Cache Node
- Stores key-value pairs in memory
- Handles TTL expiration
- Maintains synchronization with other nodes
- Tracks last writer information for consistency

### Load Balancer
- Distributes requests across cache nodes
- Uses round-robin algorithm
- Provides single entry point for clients

## Data Consistency

The system ensures consistency through:
- Write synchronization across nodes
- Last writer tracking
- Read redirects during synchronization
- Atomic operations with mutex locks

## Limitations

- In-memory storage (data is lost on restart)
- Basic round-robin load balancing
- No persistent storage
- No automatic node discovery
- No failover handling
