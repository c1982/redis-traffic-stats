package main

import (
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
	commandCountDetail.WithLabelValues(rsp.Command(), rsp.Args()).Inc()
	commandTraffic.WithLabelValues(rsp.Command()).Observe(rsp.Size())
	commandTrafficDetail.WithLabelValues(rsp.Command(), rsp.Args()).Observe(rsp.Size())
	//TODO: slowCommands
}
