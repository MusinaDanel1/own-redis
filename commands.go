package main

import (
	"strconv"
	"strings"
	"time"
)

func handlePing() string {
	return "PONG"
}

func handleSet(args []string) string {
	if len(args) < 3 {
		return "(Error) ERR wrong number of arguments for 'SET' command"
	}
	key := args[1]

	var value string
	var expiration time.Time
	pxIndex := -1

	for i := 2; i < len(args); i++ {
		if strings.EqualFold(args[i], "PX") {
			pxIndex = i
			break
		}
	}

	if pxIndex != -1 {
		if pxIndex+1 >= len(args) {
			return "(Error) ERR wrong number of arguments for 'SET' command"
		}
		value = strings.Join(args[2:pxIndex], " ")
		pxValue, err := strconv.Atoi(args[pxIndex+1])
		if err != nil {
			return "(Error) ERR invalid PX value"
		}
		expiration = time.Now().Add(time.Duration(pxValue) * time.Millisecond)
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
		return "(nil)"
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

func dispatchCommand(args []string) string {
	if len(args) == 0 {
		return "(error) ERR empty command"
	}

	command := strings.ToUpper(args[0])
	switch command {
	case "PING":
		return handlePing()
	case "SET":
		return handleSet(args)
	case "GET":
		return handleGet(args)
	default:
		return "(Error) ERR unknown command"
	}
}
