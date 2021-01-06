package main

import (
	"sync"
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
