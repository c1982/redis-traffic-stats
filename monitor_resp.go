package main

import (
	"regexp"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
)

var (
	durations *Durations
	tcpchan   chan *layers.TCP
)

func monitorRespPackets(redisport uint, sep, cleaner string, maxkeysize int, slowresponsethresold time.Duration, bigresponsethreshold int) {
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
			if packet.SrcPort == layers.TCPPort(redisport) { //response
				ditem, ok := durations.Get(packet.Seq)
				if !ok {
					break
				}
				if l := ditem.ToLatency(); l > slowresponsethresold {
					slowCommands.WithLabelValues(ditem.Command, ditem.Args).Observe(float64(l))
				}
				if size := len(packet.Payload); size > bigresponsethreshold {
					bigCommands.WithLabelValues(ditem.Command, ditem.Args).Observe(float64(size))
				}
			} else if packet.DstPort == layers.TCPPort(redisport) { //request
				rsp, err := parseRespPacket(packet.Payload, separator, cleanerxp, maxkeysize)
				if err != nil {
					log.Debug().Caller().Hex("payload", packet.Payload).Err(err).Msg("request parse error")
					break
				}
				durations.Set(packet.Ack, rsp.Command(), rsp.Args())
				commandCount.WithLabelValues(rsp.Command()).Inc()
				commandCountDetail.WithLabelValues(rsp.Command(), rsp.Args()).Inc()
			}
		}
	}
}

func parseRespPacket(payload []byte, sep []byte, cleaner *regexp.Regexp, maxkeysize int) (rsp *RespReader, err error) {
	rsp, err = NewRespReader(payload, sep, cleaner, maxkeysize)
	if err != nil {
		return rsp, err
	}

	log.Debug().Hex("payload", payload).Msg("payload")
	log.Debug().Str("command", rsp.Command()).
		Str("args", rsp.Args()).
		Float64("size", rsp.Size()).
		Msg("received")

	return rsp, err
}
