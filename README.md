# dumb Load Balancer

This project implements a basic load balancer in Go. The load balancer distributes incoming HTTP requests among a pool of servers using a round-robin scheduling algorithm.

## Features

- **Round-Robin Scheduling:** Distributes requests evenly across all available servers.
- **Server Health Check:** Periodically checks if the servers are alive and ready to handle requests.
- **Easy Configuration:** Easily add or remove servers from the pool.
- **Reverse Proxy:** Uses `httputil.ReverseProxy` to forward requests to backend servers.

## Getting Started

### Prerequisites

- Go 1.x or later installed

### Installation

1. Clone the repository:
   ```bash
   git clone [repository URL]
   ```
2. Navigate to the cloned directory.

### Usage

1. Define the servers you want to load balance across in the `main` function. For example:
   ```go
   servers := []Server{
       newServer("http://localhost:3001/"),
       newServer("http://localhost:3002/"),
       // Add more servers as needed
   }
   ```
2. Start the load balancer:
   ```bash
   go run main.go
   ```

## Components

### Server

- Represents a backend server.
- Each server has an address and a reverse proxy associated with it.
- Implements the `Server` interface.

### Load Balancer

- Manages a list of servers and distributes incoming requests.
- Implements round-robin scheduling to select servers for incoming requests.

### Request Handling

- Requests are received and forwarded to the next available server in the pool.
- The load balancer checks server health before forwarding requests.

### Error Handling

- Includes basic error handling for network and server issues.
