package main

import (
	"sync"
)

type Durations struct {
	m    sync.Mutex
	list map[uint32]int64
}
