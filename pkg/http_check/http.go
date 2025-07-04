package http_check

import (
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

func CheckHTTP(scheme, target string, successMetric *prometheus.GaugeVec, latencyMetric *prometheus.GaugeVec, failureMetric *prometheus.CounterVec) {
	url := scheme + "://" + target
	start := time.Now()

	client := &http.Client{
		Timeout: 900 * time.Millisecond,
	}
	if scheme == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	resp, err := client.Get(url)
	duration := time.Since(start).Seconds()

	if err != nil || resp.StatusCode >= 400 {
		log.Printf("❌ %s check failed for %s: %v", scheme, target, err)
		successMetric.WithLabelValues(target).Set(0)
		failureMetric.WithLabelValues(target).Inc()
		return
	}

	log.Printf("✅ %s check OK for %s in %.2fs", scheme, target, duration)
	successMetric.WithLabelValues(target).Set(1)
	latencyMetric.WithLabelValues(target).Set(duration)
}
