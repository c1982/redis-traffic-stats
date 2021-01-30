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

	durations = NewDurations()
	tcpchan = make(chan *layers.TCP, 100)

	for {
		select {
		case packet := <-tcpchan:
			if packet.SrcPort == layers.TCPPort(redisport) { //response
				ditem, ok := durations.Get(packet.Seq)
				if !ok {
					continue
				}
				if l := ditem.ToLatency(); l > slowresponsethresold {
					if ditem.Args != "" {
						slowCommands.WithLabelValues(ditem.Command, ditem.Args).Observe(float64(l))
					}
				}
				if size := len(packet.Payload); size > bigresponsethreshold {
					if ditem.Args != "" {
						bigCommands.WithLabelValues(ditem.Command, ditem.Args).Observe(float64(size))
					}
				}

				log.Debug().Str("command", ditem.Command).Str("args", ditem.Args).Int("size", len(packet.Payload)).Msg("response")
			} else if packet.DstPort == layers.TCPPort(redisport) { //request
				rsp, err := parseRespPacket(packet.Payload, separator, cleanerxp, maxkeysize)
				if err != nil {
					log.Debug().Caller().Hex("payload", packet.Payload).Err(err).Msg("request parse error")
					continue
				}
				if rsp.Args() != "" {
					durations.Set(packet.Ack, rsp.Command(), rsp.Args())
					commandCountDetail.WithLabelValues(rsp.Command(), rsp.Args()).Inc()
				}
				if rsp.Command() != "" {
					commandCount.WithLabelValues(rsp.Command()).Inc()
				}

				log.Debug().Str("command", rsp.Command()).Str("args", rsp.Args()).Float64("size", rsp.Size()).Msg("request")
			}
		}
	}
}

func parseRespPacket(payload []byte, sep []byte, cleaner *regexp.Regexp, maxkeysize int) (rsp *RespReader, err error) {
	rsp, err = NewRespReader(payload, sep, cleaner, maxkeysize)
	if err != nil {
		return rsp, err
	}
	return rsp, err
}
