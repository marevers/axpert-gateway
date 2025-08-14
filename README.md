# Axpert Gateway

A Prometheus metrics gateway for Axpert solar inverters, providing real-time monitoring and data collection via USB connectivity.

## Quick Start

### Using Docker (Recommended)

Pull the latest image from GitHub Container Registry:

```bash
docker pull ghcr.io/marevers/axpert-gateway:latest
```

Run the container with USB device access:

```bash
docker run -d \
  --name axpert-gateway \
  --device=/dev/hidraw0 \
  -p 8080:8080 \
  ghcr.io/marevers/axpert-gateway:latest
```

### Building from Source

#### Prerequisites

- Go 1.24.6 or later
- libhidapi-dev and libudev-dev (for USB HID support)
- pkg-config

#### Installation

```bash
# Clone the repository
git clone https://github.com/marevers/axpert-gateway.git
cd axpert-gateway

# Install dependencies (Ubuntu/Debian)
sudo apt-get install pkg-config libhidapi-dev libudev-dev

# Build the application
go build -o axpert-gateway .

# Run the application
./axpert-gateway
```

#### Cross-compile for Raspberry Pi (arm64)

To cross-compile for Raspberry Pi, you can build using the included build script (requires Docker):

```bash
./build.sh
```

## Configuration

The application supports the following command-line flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--log.level` | `info` | Log level for logging (debug, info, warn, error) |
| `--web.listen-address` | `:8080` | The address to listen on for HTTP requests |
| `--web.telemetry-path` | `/metrics` | Path under which to expose metrics |
| `--axpert.interval` | `30` | Interval in seconds for data polling |
| `--axpert.metrics` | `true` | Enable/disable metrics collection |
| `--axpert.control` | `false` | Enable/disable control API (future feature) |

### Example Usage

```bash
# Run with custom settings
./axpert-gateway \
  -log.level=debug \
  -web.listen-address=:9090 \
  -axpert.interval=60 \
  -axpert.control=true
```

## API Endpoints

- **`/`** - Landing page with links to available endpoints
- **`/metrics`** - Prometheus metrics endpoint
- **`/healthz`** - Health check endpoint

## Metrics

The gateway exposes various Axpert inverter metrics in Prometheus format, including:

- Power generation and consumption data
- Battery status and voltage levels
- Inverter operational parameters
- System health indicators

Access metrics at: `http://localhost:8080/metrics`

## Hardware Requirements

- Axpert solar inverter with USB connectivity
- USB cable (typically USB-A to USB-B)
- Linux-based system (Raspberry Pi recommended)
- USB HID support in kernel

## Monitoring Setup

### Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'axpert-gateway'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 30s
    metrics_path: /metrics
```

### Grafana Dashboard

Create dashboards to visualize:
- Solar power generation
- Battery charge levels
- Load consumption
- Inverter efficiency
- System alerts

## Troubleshooting

### Common Issues

1. **USB Permission Denied**
   ```bash
   # Add user to dialout group
   sudo usermod -a -G dialout $USER
   # Or run with appropriate permissions
   sudo ./axpert-gateway
   ```

2. **No Inverters Found**
   - Check USB connection
   - Verify inverter is powered on
   - Ensure USB drivers are installed

3. **Metrics Not Updating**
   - Check inverter connectivity
   - Verify polling interval settings
   - Review application logs

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
