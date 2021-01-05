package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debugmode := flag.Bool("debug", true, "enabled debug logs")
	devicename := flag.String("interface", "lo0", "Interface for monitoring")
	redisport := flag.Uint("redisport", 6379, "redis port number")
	exporteraddr := flag.String("addr", ":9100", "prometheus exporter http port")
	exporterusername := flag.String("username", "admin", "prometheus exporter username")
	exporterpassword := flag.String("password", "pass", "prometheus exporter password")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugmode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().
		Str("devicename", *devicename).
		Uint("redisport", *redisport).
		Str("addr", *exporteraddr).
		Str("username", *exporterusername).
		Msg("redis monitoring started")

	go monitorRespPackets()
	go exportPrometheusMetrics(*exporteraddr, *exporterusername, *exporterpassword)

	if err := StartMonitor(*devicename, uint16(*redisport)); err != nil {
		panic(err)
	}
}
