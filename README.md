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
- Node.js 22 or later (for frontend development)
- libhidapi-dev and libudev-dev (for USB HID support)
- pkg-config

#### Installation

```bash
# Clone the repository
git clone https://github.com/marevers/axpert-gateway.git
cd axpert-gateway

# Install Go dependencies (Ubuntu/Debian)
sudo apt-get install pkg-config libhidapi-dev libudev-dev

# Build frontend (required for control interface)
cd frontend
npm install
npx tsc app.ts --target es2017 --lib es2017,dom --outDir .
cd ..

# Build the application
go build -o axpert-gateway .

# Run the application with control API enabled
./axpert-gateway --axpert.control=true
```

#### Cross-compile for Raspberry Pi (arm64)

To cross-compile for Raspberry Pi, you can build using the included build script (requires Docker):

```bash
./build.sh
```

#### Development

The project includes a TypeScript-based frontend located in the `frontend/` directory:

```bash
# Frontend development
cd frontend
npm install                    # Install dependencies
npx tsc app.ts --target es2017 --lib es2017,dom --outDir .  # Compile TypeScript
```

**Frontend Files:**
- `frontend/app.ts` - TypeScript source code
- `frontend/index.html` - Main HTML interface
- `frontend/styles.css` - CSS styling
- `frontend/package.json` - npm configuration

## Configuration

The application supports the following command-line flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--log.level` | `info` | Log level for logging (debug, info, warn, error) |
| `--web.listen-address` | `:8080` | The address to listen on for HTTP requests |
| `--web.telemetry-path` | `/metrics` | Path under which to expose metrics |
| `--axpert.interval` | `30` | Interval in seconds for data polling |
| `--axpert.metrics` | `true` | Enable/disable metrics collection |
| `--axpert.control` | `false` | Enable/disable control API |

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
- **`/control/`** - Web-based control interface (when control API is enabled)
- **`/api/inverters`** - List available inverters (JSON API)
- **`/api/command/:command`** - Execute inverter commands (JSON API)

## Control API & Web Interface

The gateway now includes a comprehensive control API and web interface for managing your Axpert inverters remotely.

### 🔋 Web Control Interface

Access the modern web interface at: `http://localhost:8080/control/`

**Features:**
- **🎯 Inverter Selection** - Dynamic dropdown populated from connected inverters
- **⚡ Output Priority Control** - Set utility/solar/SBU priority with one click
- **🔌 Charger Priority Control** - Configure charging source preferences
- **⚡ Max Charge Current** - Set AC charging current limits (1-100A)
- **🛡️ Safety Confirmations** - Confirmation dialogs prevent accidental changes
- **📱 Responsive Design** - Works on desktop, tablet, and mobile devices
- **🔋 Battery Favicon** - Professional branding with battery emoji

### 🚀 Control API Endpoints

Enable the control API with the `--axpert.control=true` flag.

#### List Inverters
```bash
GET /api/inverters
```

**Response:**
```json
{
  "inverters": [
    {"serialno": "INV001234"},
    {"serialno": "INV005678"}
  ],
  "count": 2
}
```

#### Execute Commands
```bash
POST /api/command/:command
Content-Type: application/json

{
  "value": "solar",
  "serialno": "INV001234"
}
```

**Available Commands:**
- `setOutputPriority` - Values: `utility`, `solar`, `sbu`
- `setChargerPriority` - Values: `utilityfirst`, `solarfirst`, `solarandutility`, `solaronly`
- `setMaxChargeCurrent` - Values: `1-100` (amperes)

**Response:**
```json
{
  "command": "setOutputPriority",
  "value": "solar",
  "status": "success",
  "message": "Command executed successfully"
}
```

### 🔒 Safety Features

- **Targeted Control** - Commands target specific inverters by serial number
- **Confirmation Dialogs** - Web interface requires user confirmation
- **Error Handling** - Comprehensive error messages and status feedback
- **Disabled by Default** - Control API must be explicitly enabled

### 📋 Example Usage

```bash
# Enable control API and start gateway
./axpert-gateway --axpert.control=true

# Set output priority via API
curl -X POST http://localhost:8080/api/command/setOutputPriority \
  -H "Content-Type: application/json" \
  -d '{"value": "solar", "serialno": "INV001234"}'

# Or use the web interface at http://localhost:8080/control/
```

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

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE.txt) file for details.
