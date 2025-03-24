package main

import (
	"log"
	"net"
	"strings"
)

func sendResponse(conn *net.UDPConn, remoteAddr *net.UDPAddr, response string) {
	_, err := conn.WriteToUDP([]byte(response+"\n"), remoteAddr)
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func startUDPServer(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error starting UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP server listening on port %s", port)
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		log.Printf("Received: %s from %s", message, remoteAddr.String())

		args := strings.Fields(message)

		response := dispatchCommand(args)

		sendResponse(conn, remoteAddr, response)
	}
}
