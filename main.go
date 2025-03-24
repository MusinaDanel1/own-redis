package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	store = make(map[string]Entry)
	mu    sync.RWMutex
)

type Entry struct {
	value     string
	expiresAt time.Time
}

func main() {
	port := flag.String("port", "8080", "Port number for UDP server")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		fmt.Println("Own Redis")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  own-redis [--port <N>]")
		fmt.Println("  own-redis --help")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --help       Show this screen.")
		fmt.Println(" --port N      Port number (default 8080).")
		os.Exit(0)
	}

	addr, err := net.ResolveUDPAddr("udp", ":"+*port)
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error starting UDP server: %v", err)
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on port: %s\n", *port)

	buffer := make([]byte, 1024)
	var data []string
	var command string
	var command1 string
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		data = strings.Fields(string(buffer[:n]))
		fmt.Printf("Received '%s' from %s\n", data, remoteAddr)

		if len(data) != 0 {
			command = strings.ToUpper(data[0])
		}
		switch command {
		case "PING":
			sendResponse(conn, remoteAddr, "PONG")
			command = command1
		case "SET":
			message := handleSet(data)
			sendResponse(conn, remoteAddr, message)
			command = command1
		case "GET":
			message1 := handleGet(data)
			sendResponse(conn, remoteAddr, message1)
			command = command1
			//Получить значение по ключу и вернуть его,  либо (nil) если не найдено или просрочено
		default:
			fmt.Println("This program has only PING, SET, GET command")
		}

	}

}

func sendResponse(conn *net.UDPConn, remoteAddr *net.UDPAddr, response string) {
	_, err := conn.WriteToUDP([]byte(response+"\n"), remoteAddr)
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func handleSet(args []string) string {
	if len(args) < 3 {
		return "(Error) ERR  wrong number of arguments for 'SET' command"
	}
	key := args[1]
	var value string
	var expiration time.Time
	var expireMillis int
	var position int = -1

	for i := 2; i < len(args); i++ {
		if strings.EqualFold(args[i], "PX") {
			position = i
			break
		}
	}
	if position != -1 {
		if position+1 >= len(args) {
			return "(Error) ERR wrong number of arguments for 'SET' command"
		}
		value = strings.Join(args[2:position], " ")
		var err error
		expireMillis, err = strconv.Atoi(args[position+1])
		if err != nil {
			return "(Error) ERR invalid PX value"
		}
		expiration = time.Now().Add(time.Duration(expireMillis) * time.Millisecond)
	} else {
		value = strings.Join(args[2:], " ")
	}

	mu.Lock()
	store[key] = Entry{
		value:     value,
		expiresAt: expiration,
	}
	mu.Unlock()

	return "OK"
}

func handleGet(args []string) string {
	if len(args) < 2 {
		return "(Error) ERR wrong number of arguments for 'GET' command"
	}
	key := args[1]
	mu.RLock()
	entry, exists := store[key]
	mu.RUnlock()

	if !exists {
		return "(nil)"
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		mu.Lock()
		entry, exists = store[key]
		if exists && !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
			delete(store, key)
		}
		mu.Unlock()
		return "(nil)"
	}
	return entry.value
}
