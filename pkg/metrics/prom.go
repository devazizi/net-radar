package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	DnsSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "net_radar_dns_query_success",
			Help: "Whether the DNS query to a server succeeded (1) or failed (0).",
		},
		[]string{"server"},
	)

	DnsLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "net_radar_dns_query_latency_seconds",
			Help: "Latency in seconds for DNS query per server.",
		},
		[]string{"server"},
	)

	DnsFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "net_radar_dns_query_failures_total",
			Help: "Total number of DNS query failures per server.",
		},
		[]string{"server"},
	)

	HttpSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "net_radar_http_check_success", Help: "HTTP check success"}, []string{"target"})

	HttpLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "net_radar_http_check_latency_seconds",
		Help: "HTTP latency"},
		[]string{"target"},
	)

	HttpFailures = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "net_radar_http_check_failures_total", Help: "HTTP failures"}, []string{"target"})

	HttpsSuccess  = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "net_radar_https_check_success", Help: "HTTPS check success"}, []string{"target"})
	HttpsLatency  = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "net_radar_https_check_latency_seconds", Help: "HTTPS latency"}, []string{"target"})
	HttpsFailures = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "net_radar_https_check_failures_total", Help: "HTTPS failures"}, []string{"target"})

	IcmpSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "net_radar_icmp_ping_success",
		Help: "Whether ICMP ping to a host succeeded (1) or failed (0).",
	}, []string{"target"})

	IcmpLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "net_radar_icmp_ping_latency_seconds",
		Help: "Round-trip time of ICMP ping.",
	}, []string{"target"})

	IcmpFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "net_radar_icmp_ping_failures_total",
		Help: "Total number of ICMP ping failures.",
	}, []string{"target"})
)

func InitMetrics() {
	prometheus.MustRegister(
		DnsSuccess, DnsLatency, DnsFailures,
		HttpSuccess, HttpLatency, HttpFailures,
		HttpsSuccess, HttpsLatency, HttpsFailures,
		IcmpSuccess, IcmpLatency, IcmpFailures,
	)

	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
