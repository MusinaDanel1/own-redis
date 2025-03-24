package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Port number for UDP server")
	help := flag.Bool("help", false, "Show help message")

	flag.Usage = func() {
		fmt.Println("Own Redis")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  own-redis [--port <N>]")
		fmt.Println("  own-redis --help")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --help       Show this screen.")
		fmt.Println("  --port N     Port number (default 8080).")
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	startUDPServer(*port)
}
