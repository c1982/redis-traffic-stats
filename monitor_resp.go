package main

import (
	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
	"sync"
)

var (
	duratios *Durations
	tcpchan  chan *layers.TCP
)

func monitorRespPackets(redisport uint) {
	tcpchan = make(chan *layers.TCP, 100)
	duratios = &Durations{
		m:    sync.Mutex{},
		list: map[uint32]int64{},
	}

	for {
		select {
		case packet := <-tcpchan:
			if packet.SrcPort == layers.TCPPort(redisport) { //redis response
				//TODO: handle response
			} else if packet.DstPort == layers.TCPPort(redisport) { //redis request
				processRespPacket(packet.Payload)
			}
		}
	}
}

func processRespPacket(payload []byte) {
	rsp, err := NewRespReader(payload)
	if err != nil {
		log.Debug().Caller().Hex("payload", payload).Err(err).Msg("parse error")
		return
	}

	log.Debug().Hex("payload", payload).Msg("payload")
	log.Debug().Str("command", rsp.Command()).
		Str("args", rsp.Args()).
		Float64("size", rsp.Size()).
		Msg("received")

	commandCount.WithLabelValues(rsp.Command()).Inc()
	commandCountDetail.WithLabelValues(rsp.Command(), rsp.Args()).Inc()

	//TODO: implement slow response and traffics
}
