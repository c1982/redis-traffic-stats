package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		debugmode        = flag.Bool("debug", false, "Enable debug logs")
		devicename       = flag.String("interface", "", "Ethernet infreface name. eth0, ens5")
		redisport        = flag.Uint("redisport", 6379, "Redis server port number")
		exporteraddr     = flag.String("addr", ":9100", "HTTP listener port for prometheus metrics")
		exporterusername = flag.String("username", "admin", "Prometheus metrics username")
		exporterpassword = flag.String("password", "pass", "Prometheus metrics password")
		keyseparator     = flag.String("s", "", "Separator of keys (for split). If it empty does not split keys.")
		keycleanerregex  = flag.String("r", "", "Regex pattern for cleaner in keys")
		maxkeysizenumber = flag.Int("max", 120, "Key size to be lookup")
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
