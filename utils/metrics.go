//code written by amish


package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/docker"
)

type Metrics struct {
	CPUThreshold float64
	MemThreshold float64

	processCPU    *prometheus.GaugeVec
	processMemory *prometheus.GaugeVec
	systemCPU     *prometheus.GaugeVec
	systemMemory  *prometheus.GaugeVec
	processCount  *prometheus.GaugeVec
	containerCPU  *prometheus.GaugeVec
	containerMem  *prometheus.GaugeVec
}

func NewMetrics(cpuThreshold, memThreshold float64) *Metrics {
	return &Metrics{
		CPUThreshold: cpuThreshold,
		MemThreshold: memThreshold,

		processCPU: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu_percent",
				Help: "CPU usage percentage per process",
			},
			[]string{"process_name", "pid"},
		),

		processMemory: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_memory_mb",
				Help: "Memory usage in MB per process",
			},
			[]string{"process_name", "pid"},
		),

		systemCPU: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_cpu_percent",
				Help: "System-wide average CPU usage percentage",
			},
			[]string{"hostname"},
		),

		systemMemory: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_memory_gb",
				Help: "System-wide total memory usage in GB",
			},
			[]string{"hostname"},
		),

		processCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_process_count",
				Help: "Total number of running processes",
			},
			[]string{"hostname"},
		),

		containerCPU: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "container_cpu_percent",
				Help: "CPU usage percentage per container",
			},
			[]string{"container_id", "container_name"},
		),

		containerMem: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "container_memory_mb",
				Help: "Memory usage in MB per container",
			},
			[]string{"container_id", "container_name"},
		),
	}
}

func (m *Metrics) UpdateProcessMetrics(name string, pid int32, cpu, memMB float64) {
	pidStr := formatPID(pid)
	m.processCPU.WithLabelValues(name, pidStr).Set(cpu)
	m.processMemory.WithLabelValues(name, pidStr).Set(memMB)
}

func (m *Metrics) UpdateSystemMetrics(cpu, memGB float64, count int, hostname string) {
	m.systemCPU.WithLabelValues(hostname).Set(cpu)
	m.systemMemory.WithLabelValues(hostname).Set(memGB)
	m.processCount.WithLabelValues(hostname).Set(float64(count))
}

func (m *Metrics) UpdateContainerMetrics(id, name string, cpu, memMB float64) {
	shortID := id
	if len(id) > 12 {
		shortID = id[:12]
	}
	m.containerCPU.WithLabelValues(shortID, name).Set(cpu)
	m.containerMem.WithLabelValues(shortID, name).Set(memMB)
}

func CalculateContainerCPU(container docker.CgroupDockerStat) float64 {
	if container.CPUStats.TotalUsage == 0 || container.PreCPUStats.TotalUsage == 0 {
		return 0.0
	}

	cpuDelta := float64(container.CPUStats.TotalUsage - container.PreCPUStats.TotalUsage)
	systemDelta := float64(container.CPUStats.SystemUsage - container.PreCPUStats.SystemUsage)

	if systemDelta == 0 {
		return 0.0
	}

	numCPUs := float64(container.CPUStats.OnlineCPUs)
	if numCPUs == 0 {
		numCPUs = float64(len(container.CPUStats.PercpuUsage))
		if numCPUs == 0 {
			numCPUs = 1.0
		}
	}

	cpuPercent := (cpuDelta / systemDelta) * numCPUs * 100.0
	return cpuPercent
}

func formatPID(pid int32) string {
	return string(rune(pid + '0'))
}