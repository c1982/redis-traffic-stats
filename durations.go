package main

import (
	"sync"
	"time"
)

type DurationItem struct {
	Latency int64 //unixnano
	Command string
	Args    string
}

func (d DurationItem) ToLatency() time.Duration {
	return time.Nanosecond * time.Duration(time.Now().UnixNano()-d.Latency)
}

type Durations struct {
	m    sync.Mutex
	list map[uint32]DurationItem
}

func (d *Durations) Set(k uint32, command, args string) {
	d.m.Lock()
	defer d.m.Unlock()
	d.list[k] = DurationItem{time.Now().UnixNano(), command, args}
}

func (d *Durations) Get(k uint32) (item DurationItem, exist bool) {
	d.m.Lock()
	defer d.m.Unlock()

	item, exist = d.list[k]
	if !exist {
		return item, exist
	}
	delete(d.list, k)
	return item, true
}
