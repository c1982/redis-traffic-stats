package main

import (
	"regexp"

	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
)

var (
	duratios *Durations
	tcpchan  chan *layers.TCP
)

func monitorRespPackets(redisport uint, sep, cleaner string) {
	var (
		separator []byte
		cleanerxp *regexp.Regexp
	)

	if sep != "" {
		separator = []byte(sep)
	}

	if cleaner != "" {
		cleanerxp = regexp.MustCompile(cleaner)
	}

	tcpchan = make(chan *layers.TCP, 100)

	for {
		select {
		case packet := <-tcpchan:
			if packet.SrcPort == layers.TCPPort(redisport) { //redis response
				//TODO: handle response
			} else if packet.DstPort == layers.TCPPort(redisport) { //redis request
				processRespPacket(packet.Payload, separator, cleanerxp)
			}
		}
	}
}

func processRespPacket(payload []byte, sep []byte, cleaner *regexp.Regexp) {
	rsp, err := NewRespReader(payload, sep, cleaner)
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
