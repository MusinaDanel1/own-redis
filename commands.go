package main

import (
	"strconv"
	"strings"
	"time"
)

func handlePing(args []string) string {
	if len(args) != 1 {
		return "(Error) ERR wrong number of arguments for 'PING' command"
	}
	return "PONG"
}

func handleSet(args []string) string {
	if len(args) < 3 {
		return "(Error) ERR wrong number of arguments for 'SET' command"
	}
	key := strings.TrimSpace(args[1])
	if key == "" {
		return "(Error) ERR key cannot be empty"
	}

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
		if strings.TrimSpace(value) == "" {
			return "(Error) ERR value cannot be empty"
		}
		pxValue, err := strconv.Atoi(args[pxIndex+1])
		if err != nil {
			return "(Error) ERR invalid PX value"
		}
		expiration = time.Now().Add(time.Duration(pxValue) * time.Millisecond)
	} else {
		value = strings.Join(args[2:], " ")
		if strings.TrimSpace(value) == "" {
			return "(Error) ERR value cannot be empty"
		}
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
	key := strings.TrimSpace(args[1])
	if key == "" {
		return "(Error) ERR key cannot be empty"
	}

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
		return handlePing(args)
	case "SET":
		return handleSet(args)
	case "GET":
		return handleGet(args)
	default:
		return "(Error) ERR unknown command"
	}
}
