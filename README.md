# portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes with configurable rules.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
  path: /var/log/portwatch.log
rules:
  - port: 22
    protocol: tcp
    expected: open
  - port: 8080
    protocol: tcp
    expected: closed
```

portwatch will scan open ports at the defined interval and emit an alert whenever a port's state deviates from the expected value defined in your rules.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to configuration file |
| `--interval` | `60s` | Override scan interval |
| `--verbose` | `false` | Enable verbose logging |

## License

MIT © 2024 Your Name