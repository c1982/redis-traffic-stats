package main

import "flag"

func main() {
	devicename := flag.String("interface", "eth0", "Interface for monitoring")
	redisport := flag.Uint("redisport", 6379, "redis port number")
	exporteraddr := flag.String("addr", ":9100", "prometheus exporter http port")
	exporterusername := flag.String("username", "", "prometheus exporter username")
	exporterpassword := flag.String("password", "", "prometheus exporter password")
	flag.Parse()

	go monitorRespPackets()
	go exportPrometheusMetrics(*exporteraddr, *exporterusername, *exporterpassword)

	if err := StartMonitor(*devicename, uint16(*redisport)); err != nil {
		panic(err)
	}
}
