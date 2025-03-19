package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "Port number for UDP server")
	flag.Parse()

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
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received '%s' from %s\n", message, remoteAddr)
	}
}
