package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

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
			_, err := conn.WriteToUDP([]byte("PONG\n"), remoteAddr)
			var command1 string
			command = command1
			if err != nil {
				log.Printf("Error sending to UDP: %v", err)
			}
		case "SET":
			//проверить количество аргументов, объединить значение, обработать РХ и сохранить данные
		case "GET":
			//Получить значение по ключу и вернуть его,  либо (nil) если не найдено или просрочено
		default:
			fmt.Println("This program has only PING, SET, GET command")
		}

	}

}
