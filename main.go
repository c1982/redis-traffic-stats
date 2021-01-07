package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		debugmode        = flag.Bool("debug", false, "enabled debug logs")
		devicename       = flag.String("interface", "", "interface for monitoring")
		redisport        = flag.Uint("redisport", 6379, "redis port number")
		exporteraddr     = flag.String("addr", ":9100", "prometheus exporter http listen port")
		exporterusername = flag.String("username", "admin", "prometheus exporter username")
		exporterpassword = flag.String("password", "pass", "prometheus exporter password")
		keyseparator     = flag.String("s", "", "separator of key. If is this empty, can not works this logic")
		keycleanerregex  = flag.String("r", "", "cleans all regex match in the key")
		maxkeysizenumber = flag.Int("max", 50, "")
	)

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

	go monitorRespPackets(*redisport, *keyseparator, *keycleanerregex, *maxkeysizenumber)
	go exportPrometheusMetrics(*exporteraddr, *exporterusername, *exporterpassword)

	if err := StartMonitor(*devicename, uint16(*redisport)); err != nil {
		panic(err)
	}
}
