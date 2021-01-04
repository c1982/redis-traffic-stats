package main

func main() {
	go monitorRespPackets()
	go exportPrometheusMetrics(":9100", "admin", "pass")

	if err := StartMonitor("lo0", 6379); err != nil {
		panic(err)
	}
}
