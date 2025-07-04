package main

import (
	"github.com/devazizi/net-radar/pkg/dns_check"
	"github.com/devazizi/net-radar/pkg/http_check"
	"github.com/devazizi/net-radar/pkg/icmp_check"
	"github.com/devazizi/net-radar/pkg/metrics"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DNSServers      []string `yaml:"dns_servers"`
	HTTPTargets     []string `yaml:"http_targets"`
	Domain          string   `yaml:"domain"`
	IntervalSeconds int      `yaml:"interval_seconds"`
	ICMPTargets     []string `yaml:"icmp_targets"`
}

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func runChecks(cfg *Config) {
	for {
		var wg sync.WaitGroup
		for _, server := range cfg.DNSServers {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				dns_check.SendQueryDNS(s, cfg.Domain)
			}(server)
		}

		for _, target := range cfg.HTTPTargets {
			t := target
			wg.Add(2)

			go func() {
				defer wg.Done()
				http_check.CheckHTTP("http", t, metrics.HttpSuccess, metrics.HttpLatency, metrics.HttpFailures)
			}()

			go func() {
				defer wg.Done()
				http_check.CheckHTTP("https", t, metrics.HttpSuccess, metrics.HttpLatency, metrics.HttpFailures)
			}()
		}

		for _, target := range cfg.ICMPTargets {
			wg.Add(1)
			go func(t string) {
				defer wg.Done()
				icmp_check.CheckICMP(t)
			}(target)
		}

		wg.Wait()
		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func main() {
	cfg, err := loadConfig("../config.yaml")
	if err != nil {
		log.Fatalf("Error loading config.yaml: %v", err)
	}

	metrics.InitMetrics()

	go runChecks(cfg)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Fatalf("Error writing health response: %v", err)
		}
	}))

	log.Println("server available at http://0.0.0.0:2112")
	log.Fatal(http.ListenAndServe("0.0.0.0:2112", nil))
}
