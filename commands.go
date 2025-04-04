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

	var pxValue int
	var err error
	filtered := []string{args[0]}

	for i := 1; i < len(args); i++ {
		if strings.EqualFold(args[i], "PX") {
			if i+1 >= len(args) {
				return "(Error) ERR wrong number of arguments for 'SET' command"
			}
			pxValue, err = strconv.Atoi(args[i+1])
			if err != nil || pxValue <= 0 {
				return "(Error) ERR invalid PX value, must be greater than zero"
			}
			i++
		} else {
			filtered = append(filtered, args[i])
		}
	}

	if len(filtered) < 3 {
		return "(Error) ERR wrong number of arguments for 'SET' command: key and value must be provided"
	}

	key := strings.TrimSpace(filtered[1])
	if key == "" {
		return "(Error) ERR key cannot be empty"
	}

	value := strings.TrimSpace(strings.Join(filtered[2:], " "))
	if value == "" {
		return "(Error) ERR value cannot be empty"
	}

	var expiration time.Time
	if pxValue > 0 {
		expiration = time.Now().Add(time.Duration(pxValue) * time.Millisecond)
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
