package main

import (
	"sync"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
)

var (
	duratios *Durations
	tcpchan  chan *layers.TCP
	results  *RespResult
)

func monitorRespPackets(redisport uint) {
	results = NewRespResult()
	tcpchan = make(chan *layers.TCP, 100)
	duratios = &Durations{
		m:    sync.Mutex{},
		list: map[uint32]int64{},
	}

	go func() {
		for range time.Tick(time.Second * 15) {
			log.Debug().Msg("collect reulsts")
			collectResults()
		}
	}()

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

	//TODO: implement slow response and traffics
	results.Add(rsp.Command(), rsp.Args(), rsp.Size())
}

func collectResults() {
	for _, command := range results.Commands() {
		stats := results.calculateCommandStats(command)
		commandCount.WithLabelValues(command).Add(float64(stats.Count))

		for _, arg := range stats.Arguments {
			commandCountDetail.WithLabelValues(command, arg.Argument).Add(float64(arg.Count))
		}

		results.Clear(command)
	}
}
