
//code written by amish
package monitor

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"process-monitor/utils"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v3/docker"
)

type Monitor struct {
	metrics *utils.Metrics
}

func NewMonitor(metrics *utils.Metrics) *Monitor {
	return &Monitor{
		metrics: metrics,
	}
}

func (m *Monitor) Monitor() error {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var totalCPU float64
	var totalMem uint64
	var processCount int

	hostname, _ := os.Hostname()

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		pid := p.Pid

		totalCPU += cpuPercent
		totalMem += memInfo.RSS
		processCount++

		m.metrics.UpdateProcessMetrics(name, pid, cpuPercent, float64(memInfo.RSS)/1024/1024)

		if cpuPercent > m.metrics.CPUThreshold {
			log.Printf("WARNING: Process %s (PID: %d) exceeded CPU threshold: %.2f%%\n", name, pid, cpuPercent)
		}

		memPercent := float64(memInfo.RSS) / float64(getTotalMemory()) * 100
		if memPercent > m.metrics.MemThreshold {
			log.Printf("WARNING: Process %s (PID: %d) exceeded memory threshold: %.2f%%\n", name, pid, memPercent)
		}
	}

	avgCPU := totalCPU / float64(runtime.NumCPU())
	m.metrics.UpdateSystemMetrics(avgCPU, float64(totalMem)/1024/1024/1024, processCount, hostname)

	log.Printf("System Summary - Processes: %d, Avg CPU: %.2f%%, Total Memory: %.2f GB\n",
		processCount, avgCPU, float64(totalMem)/1024/1024/1024)

	if err := m.monitorContainers(); err != nil {
		log.Printf("Container monitoring error: %v\n", err)
	}

	return nil
}

func (m *Monitor) monitorContainers() error {
	containers, err := docker.GetDockerStat()
	if err != nil {
		if strings.Contains(err.Error(), "no such file") || 
		   strings.Contains(err.Error(), "permission denied") ||
		   strings.Contains(err.Error(), "cannot connect") {
			return nil
		}
		return err
	}

	for _, container := range containers {
		cpuPercent := utils.CalculateContainerCPU(container)
		memUsageMB := float64(container.MemUsage) / 1024 / 1024

		m.metrics.UpdateContainerMetrics(
			container.ContainerID,
			container.Name,
			cpuPercent,
			memUsageMB,
		)

		if cpuPercent > m.metrics.CPUThreshold {
			log.Printf("WARNING: Container %s (%s) exceeded CPU threshold: %.2f%%\n",
				container.Name, container.ContainerID[:12], cpuPercent)
		}

		memPercent := float64(container.MemUsage) / float64(getTotalMemory()) * 100
		if memPercent > m.metrics.MemThreshold {
			log.Printf("WARNING: Container %s (%s) exceeded memory threshold: %.2f%%\n",
				container.Name, container.ContainerID[:12], memPercent)
		}

		log.Printf("Container: %s (%s) - CPU: %.2f%%, Memory: %.2f MB\n",
			container.Name, container.ContainerID[:12], cpuPercent, memUsageMB)
	}

	return nil
}

func getTotalMemory() uint64 {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/meminfo")
		if err != nil {
			return 8 * 1024 * 1024 * 1024
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				var total uint64
				fmt.Sscanf(line, "MemTotal: %d kB", &total)
				return total * 1024
			}
		}
	}
	return 8 * 1024 * 1024 * 1024
}