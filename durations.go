package main

import (
	"sync"
	"time"
)

//Durations this struct helps for latency
/*
Durations
duratios.m.Lock()
if tcp.SrcPort == redisPort {
	duration, ok := duratios.list[tcp.Seq]
	if ok {
		current := time.Now().UnixNano()
		latency := current - duration
		fmt.Printf("seq: %d, latency %s  len: %d\n", tcp.Seq, time.Nanosecond*time.Duration(latency), len(duratios.list))
		delete(duratios.list, tcp.Seq)
	}
}
duratios.list[tcp.Ack] = time.Now().UnixNano()
duratios.m.Unlock()
*/
type Durations struct {
	m    sync.Mutex
	list map[uint32]int64
}

func (d *Durations) Set(k uint32) {
	d.m.Lock()
	defer d.m.Unlock()
	d.list[k] = time.Now().UnixNano()
}

func (d *Durations) Get(k uint32) time.Duration {
	d.m.Lock()
	defer d.m.Unlock()

	v, ok := d.list[k]
	if !ok {
		return -1
	}

	delete(d.list, k)
	return time.Nanosecond * time.Duration(time.Now().UnixNano()-v)
}
