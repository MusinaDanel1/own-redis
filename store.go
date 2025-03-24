package main

import (
	"sync"
	"time"
)

type Entry struct {
	value     string
	expiresAt time.Time
}

var (
	store = make(map[string]Entry)
	mu    sync.RWMutex
)
