# Process Monitor

```
    ___                 
   /   |  _____  ____  __
  / /| | / ___/ / __ \/ /
 / ___ |(__  ) / /_/ / /_
/_/  |_/____/  \____/\__/
                          
```

A lightweight, efficient Go-based system process and container monitoring tool with Prometheus metrics export.

## 🚀 Features

- **Real-time Process Monitoring** - Track CPU and memory usage for all running processes
- **Docker Container Support** - Monitor containerized applications
- **Smart Thresholds** - Configurable alerts when resources exceed limits
- **Prometheus Integration** - Export metrics for visualization and alerting
- **System Aggregation** - View system-wide resource consumption at a glance
- **Cross-Platform** - Works on Linux, macOS, and Windows

## 📦 Installation

### From Source

```bash
git clone https://github.com/yourusername/process-monitor.git
cd process-monitor
go mod download
go build -o process-monitor
```

### Quick Start

```bash
./process-monitor
```

## 🔧 Usage

```bash
./process-monitor [flags]
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-interval` | Monitoring interval in seconds | 5 |
| `-cpu-threshold` | CPU usage threshold percentage | 80.0 |
| `-mem-threshold` | Memory usage threshold percentage | 80.0 |
| `-prometheus-port` | Prometheus metrics HTTP port | 8080 |
| `-enable-prometheus` | Enable Prometheus metrics export | true |

### Examples

**Basic monitoring with default settings:**
```bash
./process-monitor
```

**Custom thresholds and interval:**
```bash
./process-monitor -interval=10 -cpu-threshold=75 -mem-threshold=85
```

**Disable Prometheus export:**
```bash
./process-monitor -enable-prometheus=false
```

**Custom Prometheus port:**
```bash
./process-monitor -prometheus-port=9090
```

## 📊 Prometheus Metrics

Access metrics at `http://localhost:8080/metrics`

### Available Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `process_cpu_percent` | Gauge | CPU usage per process | process_name, pid |
| `process_memory_mb` | Gauge | Memory usage per process (MB) | process_name, pid |
| `system_cpu_percent` | Gauge | System-wide average CPU usage | hostname |
| `system_memory_gb` | Gauge | System-wide total memory usage (GB) | hostname |
| `system_process_count` | Gauge | Total number of running processes | hostname |
| `container_cpu_percent` | Gauge | CPU usage per container | container_id, container_name |
| `container_memory_mb` | Gauge | Memory usage per container (MB) | container_id, container_name |

### Prometheus Configuration Example

```yaml
scrape_configs:
  - job_name: 'process-monitor'
    static_configs:
      - targets: ['localhost:8080']
```

## 🖥️ Output Example

```
2025/11/01 10:30:15 Starting process monitor (interval: 5s, CPU threshold: 80.0%, Memory threshold: 80.0%)
2025/11/01 10:30:15 Starting Prometheus metrics server on :8080
2025/11/01 10:30:20 System Summary - Processes: 342, Avg CPU: 12.45%, Total Memory: 8.24 GB
2025/11/01 10:30:20 Container: nginx (a1b2c3d4e5f6) - CPU: 2.34%, Memory: 45.67 MB
2025/11/01 10:30:20 WARNING: Process chrome (PID: 1234) exceeded CPU threshold: 85.23%
```

## 🛠️ Requirements

- Go 1.21 or higher
- Linux, macOS, or Windows
- Docker (optional, for container monitoring)
- Appropriate system permissions for process inspection

## 📁 Project Structure

```
process-monitor/
├── main.go                 # Entry point and CLI
├── monitor/
│   └── monitor.go         # Core monitoring logic
├── utils/
│   └── metrics.go         # Prometheus metrics and calculations
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is open source and available under the MIT License.

## 📧 Contact

**Asxy**  
Email: contact.amish@yahoo.com

For bug reports and feature requests, please open an issue on GitHub.

## 🌟 Acknowledgments

Built with:
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system and process utilities
- [Prometheus Go client](https://github.com/prometheus/client_golang) - Prometheus instrumentation library

---

⭐ If you find this project useful, please consider giving it a star on GitHub!
