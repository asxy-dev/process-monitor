
//code written by amish

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"process-monitor/monitor"
	"process-monitor/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	interval := flag.Int("interval", 5, "Monitoring interval in seconds")
	cpuThreshold := flag.Float64("cpu-threshold", 80.0, "CPU usage threshold percentage")
	memThreshold := flag.Float64("mem-threshold", 80.0, "Memory usage threshold percentage")
	prometheusPort := flag.String("prometheus-port", "8080", "Prometheus metrics port")
	enablePrometheus := flag.Bool("enable-prometheus", true, "Enable Prometheus metrics export")
	flag.Parse()

	metrics := utils.NewMetrics(*cpuThreshold, *memThreshold)
	mon := monitor.NewMonitor(metrics)

	if *enablePrometheus {
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			addr := fmt.Sprintf(":%s", *prometheusPort)
			log.Printf("Starting Prometheus metrics server on %s\n", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Fatalf("Failed to start Prometheus server: %v", err)
			}
		}()
	}

	log.Printf("Starting process monitor (interval: %ds, CPU threshold: %.1f%%, Memory threshold: %.1f%%)\n",
		*interval, *cpuThreshold, *memThreshold)

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := mon.Monitor(); err != nil {
			log.Printf("Error monitoring processes: %v\n", err)
		}
	}
}