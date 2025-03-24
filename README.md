# Own Redis

Own Redis is a lightweight in-memory key-value store built as a learning project. It uses the UDP protocol for communication between the client and the storage. The project implements basic commands similar to Redis, such as **PING**, **SET**, and **GET**.

## Overview

This project demonstrates the fundamentals of TCP/IP networking, the UDP protocol, and basic in-memory database concepts. The server processes UDP packets that contain one of the commands and responds accordingly. For safe concurrent access to the in-memory storage, it uses synchronization primitives from the `sync` package.

### Supported Commands

- **PING**  
  Used to check if the server is alive.  
  **Example:**  
  ```bash
  PING
  ```  
  **Response:**  
  ```bash
  PONG
  ```

- **SET**  
  Stores a value by key. The command supports an optional `PX` parameter to set an expiration time (in milliseconds).  
  **Examples:**  
  ```bash
  SET mykey myvalue
  ```  
  **Response:**  
  ```bash
  OK
  ```  
  or with expiration:  
  ```bash
  SET mykey myvalue PX 10000
  ```  
  (The value will expire after 10 seconds.)

- **GET**  
  Retrieves the value by key. If the key is not found or the value has expired, it returns `(nil)`.  
  **Example:**  
  ```bash
  GET mykey
  ```  
  **Response:**  
  ```bash
  myvalue
  ```

## Project Structure

```
own-redis/
├── main.go       // Entry point: processes command-line flags and starts the server.
├── server.go     // Implements the UDP server and response sending function.
├── commands.go   // Contains command handlers for PING, SET, and GET, and routes commands.
└── store.go      // Defines the data structure for storage and the global in-memory store with synchronization.
```

## Requirements

- Go (version 1.13 or later)
- Only built-in Go packages are used (no external dependencies)
- Code should follow the [gofumpt](https://github.com/mvdan/gofumpt) style guidelines

## Installation

1. Clone the repository or copy the project files into a directory.
2. Open a terminal and navigate to the project root directory.

## Building

Run the following command in your terminal:

```bash
go build -o own-redis .
```

This command compiles the project and creates an executable named `own-redis`.

## Running

By default, the server starts on port 8080. You can change the port using the `--port` flag:

```bash
./own-redis --port 8080
```

To display the help screen, run:

```bash
./own-redis --help
```

## Testing

You can test the server using the netcat (or ncat) utility:

1. Start the server.
2. Open a new terminal and run:
   ```bash
   nc -u 127.0.0.1 8080
   ```
3. Type a command, such as `PING`, and press Enter. You should receive the response `PONG`.

You can also test the SET and GET commands:
```bash
SET mykey myvalue
GET mykey
```

## Additional Information

- The server processes each UDP request in a separate goroutine and uses `recover` to prevent a panic in one request from stopping the entire server.
- Access to the global in-memory store is protected by a mutex (`sync.RWMutex`), ensuring safe concurrent access.
- This project is intended for educational purposes to demonstrate basic networking and in-memory storage concepts.

## Contact

Project author: Askaruly Nurislam  
Email: askaruly@hotmail.com  
GitHub: [GitHub profile](https://github.com)
