package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	commandCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_traffic_stats_command_count",
		Help: "The total number of redis commands",
	}, []string{"command"})

	commandTraffic = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "redis_traffic_stats_command_traffic",
		Help: "The total number of redis traffic bytes",
	}, []string{"command"})

	commandCountDetail = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_traffic_stats_command_detail_count",
		Help: "The total number of redis commands detail count",
	}, []string{"command", "args"})

	commandTrafficDetail = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "redis_traffic_stats_command_traffic",
		Help: "The total number of redis traffic bytes",
	}, []string{"command", "args"})

	slowCommands = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "redis_traffic_stats_slow_commands",
		Help: "The total number of redis traffic nanosecond",
	}, []string{"command", "args"})
)

func exportPrometheusMetrics(addr, username, password string) {
	basicauth := func(f http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			user, pass, _ := r.BasicAuth()
			if user != username {
				http.Error(w, "Unauthorized.", 401)
				return
			}

			if pass != password {
				http.Error(w, "Unauthorized.", 401)
				return
			}

			f.ServeHTTP(w, r)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", basicauth(promhttp.Handler()))
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
