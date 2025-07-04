package dns_check

import (
	"github.com/devazizi/net-radar/pkg/metrics"
	"github.com/miekg/dns"
	"log"
	"time"
)

func SendQueryDNS(server, domain string) {
	client := &dns.Client{Timeout: 900 * time.Millisecond}
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	start := time.Now()
	resp, _, err := client.Exchange(msg, server+":53")
	duration := time.Since(start).Seconds()

	if err != nil || resp == nil || resp.Rcode != dns.RcodeSuccess {
		log.Printf("❌ [%s] DNS query failed: %v", server, err)
		metrics.DnsSuccess.WithLabelValues(server).Set(0)
		metrics.DnsFailures.WithLabelValues(server).Inc()
		return
	}

	metrics.DnsSuccess.WithLabelValues(server).Set(1)
	log.Printf("✅ [%s] DNS query succeeded in %.2fs", server, duration)
	metrics.DnsLatency.WithLabelValues(server).Set(duration)
}
