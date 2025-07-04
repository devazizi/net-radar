package icmp_check

import (
	"github.com/devazizi/net-radar/pkg/metrics"
	"github.com/go-ping/ping"
	"log"
	"time"
)

func CheckICMP(target string) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		log.Printf("❌ ICMP init error for %s: %v", target, err)
		metrics.IcmpSuccess.WithLabelValues(target).Set(0)
		metrics.IcmpFailures.WithLabelValues(target).Inc()
		return
	}

	pinger.Count = 3
	pinger.Timeout = 900 * time.Millisecond
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		log.Printf("❌ ICMP ping to %s failed: %v", target, err)
		metrics.IcmpSuccess.WithLabelValues(target).Set(0)
		metrics.IcmpFailures.WithLabelValues(target).Inc()
		return
	}

	stats := pinger.Statistics()
	log.Printf("✅ ICMP ping to %s success in %.2fs", target, stats.AvgRtt.Seconds())
	metrics.IcmpSuccess.WithLabelValues(target).Set(1)
	metrics.IcmpLatency.WithLabelValues(target).Set(stats.AvgRtt.Seconds())
}
