package main

import (
	"fmt"
	"sync"

	"github.com/google/gopacket/layers"
)

var (
	duratios *Durations
	tcpchan  chan *layers.TCP
)

func monitorRespPackets() {

	tcpchan = make(chan *layers.TCP, 100)

	duratios = &Durations{
		m:    sync.Mutex{},
		list: map[uint32]int64{},
	}

	for {
		select {
		case packet := <-tcpchan:
			processRespPacket(packet.Payload)
		}
	}
}

func processRespPacket(payload []byte) {
	rsp, err := NewRespReader(payload)
	if err != nil {
		return
	}
	commandCount.WithLabelValues(rsp.Command()).Inc()
	fmt.Printf("cmd: %s\n", rsp.Command())
}
