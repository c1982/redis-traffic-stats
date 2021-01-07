package main

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//StartMonitor monitor redis packet destination or source port
func StartMonitor(devicename string, redisport uint16) error {
	bpffilter := fmt.Sprintf("dst port %d and tcp[((tcp[12:1] & 0xf0) >> 2):1] = 0x2A", redisport)
	handle, err := pcap.OpenLive(devicename, 65535, false, -1*time.Second)
	if err != nil {
		return fmt.Errorf("error opening device %s: %v", devicename, err)
	}
	defer handle.Close()
	handle.SetBPFFilter(bpffilter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			continue
		}

		if len(tcp.Payload) < 1 {
			continue
		}

		tcpchan <- tcp
	}

	return nil
}
