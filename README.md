<div align="center" style="border-bottom:1px solid #444;width:100%;" alt="NanoDNS Logo">

  <img src="./docs/assets/images/nanodns.webp" style="text-align:center;" width="200"/>
  <br>
  <p align="center">
    An ultra-lightweight DNS server that runs anywhere - from Docker containers to Kubernetes pods to Linux services. Perfect for internal networks and ISVs distributing self-hosted applications, it provides custom domain resolution and service discovery without databases or external dependencies.
  </p>
  <br>

  <p align="center">
    <a href="https://github.com/mguptahub/nanodns/actions/workflows/build.yml">
      <img src="https://github.com/mguptahub/nanodns/actions/workflows/build.yml/badge.svg" alt="Build" />
    </a>
    <a href="https://github.com/mguptahub/nanodns/releases">
      <img src="https://img.shields.io/github/v/release/mguptahub/nanodns?sort=semver" alt="Release" />
    </a>
    <a href="https://github.com/mguptahub/nanodns/issues">
      <img src="https://img.shields.io/github/issues/mguptahub/nanodns" alt="Issues" />
    </a>
    <a href="go.mod">
      <img src="https://img.shields.io/github/go-mod/go-version/mguptahub/nanodns" alt="Go Version" />
    </a>
    <a href="go.mod">
      <img src="https://goreportcard.com/badge/github.com/mguptahub/nanodns" alt="Go Report" />
    </a>
    <a href="https://www.gnu.org/licenses/agpl-2.0">
      <img src="https://img.shields.io/badge/License-AGPL%20v2-blue.svg" alt="License: AGPL v2" />
    </a>
    <a href="SECURITY.md">
      <img src="https://img.shields.io/badge/Security-Policy-blue.svg" alt="Security Policy" />
    </a>
  </p>
</div>

## Features

- Environment variable-based configuration (Support .env file)
- Support for A, CNAME, MX, and TXT records
- Docker service name resolution
- Optional TTL configuration (default: 60 seconds)
- Lightweight and fast
- Configurable port

## Installation

### Download

Download the latest release from the [releases page](https://github.com/mguptahub/nanodns/releases).

### Platform-specific Instructions

#### Linux Service Installation

Release assets include scripts for installing/uninstalling NanoDNS as a system service:

```bash
# Make scripts executable
chmod +x install.sh uninstall.sh

# Install service
sudo ./install.sh

# View status and logs
systemctl status nanodns
journalctl -u nanodns -f

# Edit configuration
sudo nano /etc/nanodns/nanodns.env

# Uninstall service
sudo ./uninstall.sh
```

#### macOS

If you see the warning "Apple could not verify this app", run these commands:

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine nanodns-darwin-arm64

# Make executable
chmod +x nanodns-darwin-arm64

# Run the binary
./nanodns-darwin-arm64
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DNS_PORT | UDP port for DNS server | `10053` |
| DNS_RELAY_SERVERS | Comma-separated upstream DNS servers | `8.8.8.8:53,1.1.1.1:53` |
| DNS_DEFAULT_TTL | Default TTL | `60` |
| LOG_DIR | Log file directory path | `/tmp/log/nanodns` |
| SERVICE_LOG | Service log filename | `service.log` |
| ACTION_LOG | Action log filename | `actions.log` |
| MAX_LOG_SIZE | Max log file size before rotation (in bytes) | `1048576` |
| MAX_LOG_BACKUPS | Max log file backups before rotation | `5` |

### DNS Records as Environment Variables

| Variable | Description |
|----------|-------------|
| A_xxx | A Record Details |
| CNAME_xxx | CNAME Record Details |
| MX_xxx | MX Record Details |
| TXT_xxx | TXT Record Details |


### DNS Resolution Strategy

NanoDNS follows this resolution order:

1. Check configured local records first
2. If no local record found and relay is enabled, forward to upstream DNS servers
3. Return first successful response from relay servers

### Record Format

All records use the `|` character as a separator. The general format is:

```txt
RECORD_TYPE_NUMBER=domain|value[|ttl]
```

### A Records

```
A_REC1=domain|ip|ttl
A_REC2=domain|service:servicename|ttl
```
Example:
```
A_REC1=app.example.com|10.10.0.1|300
A_REC2=api.example.com|service:webapp
```

### CNAME Records

```
CNAME_REC1=domain|target|ttl
```
Example:
```
CNAME_REC1=www.example.com|app.example.com|3600
```

### MX Records

```
MX_REC1=domain|priority|mailserver|ttl
```
Example:
```
MX_REC1=example.com|10|mail1.example.com|3600
MX_REC2=example.com|20|mail2.example.com
```

### TXT Records

```
TXT_REC1=domain|"text value"|ttl
```
Example:
```
TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all|3600
TXT_REC2=_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com
```

## Sample `.env` file

```ini

# Basic Configuration
DNS_PORT=10053
DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53
DNS_DEFAULT_TTL=60

# A Records
A_REC1=domain|ip-addr|[ttl]

# CNAME Records
# ...

# MX Records
# ...

# TXT Records
# ...

```

## Docker Usage

### Using Docker Run

```bash
docker run -d \
  --name nanodns \
  -p 10053:10053/udp \
  -e DNS_PORT=10053 \
  -e DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53 \
  -e "A_REC1=app.example.com|10.10.0.1|300" \
  -e "A_REC2=api.example.com|service:webapp" \
  -e "TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all" \
  -v ${PWD}/.env:/app/.env \
  ghcr.io/mguptahub/nanodns:latest
```

### Using Docker Compose

```yaml
name: nanodns
services:
  dns:
    image: ghcr.io/mguptahub/nanodns:latest
    environment:
      # DNS Server Configuration
      - DNS_PORT=10053  # Optional, defaults to 53
      - DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53  # Optional relay servers

      # Local Records
      - A_REC1=app.example.com|service:webapp
      - A_REC2=api.example.com|10.10.0.5|300
      - TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all
    ports:
      - "${DNS_PORT:-10053}:${DNS_PORT:-10053}/udp"
    volumes:
      - ./.env:/app/.env
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```

### Kubernetes

For detailed instructions on deploying NanoDNS in Kubernetes, see our [Kubernetes Deployment Guide](kubernetes/README.md).

## Running Without Docker Compose

Install using the script

```bash
curl -fsSL https://nanodns.mguptahub.com/install.sh | sh -s -- --install
```

Start using the script

```bash
# Check the values in /usr/local/share/nanodns.env before starting
nanodns start
```

Help Command
```
nanodns --help
```

```

Usage: nanodns [command | options]

commands:
  start                              Run the binary as a daemon
  stop                               Stop the running daemon service
  status                             Show service status
  logs                               Show service logs
  logs -a                            Show action logs

options:
  -v | --version                     Show the binary version
  -h | --help                        Show the help information
```

## Testing Records

```bash
# Test local records
dig @localhost -p 10053 app.example.com A

# Test relay resolution (for non-local domains)
dig @localhost -p 10053 google.com A

# Test other record types
dig @localhost -p 10053 www.example.com CNAME
dig @localhost -p 10053 example.com MX
dig @localhost -p 10053 example.com TXT
```

## Common Issues and Solutions

1. Port 53 already in use (common on macOS and Linux):
   - Use a different port by setting `DNS_PORT=10053` or another available port
   - Update your client configurations to use the custom port

2. Permission denied when using port 53:
   - Use a port number above 1024 to avoid requiring root privileges
   - Set `DNS_PORT=10053` or another high-numbered port

3. DNS Relay Issues:
   - Verify upstream DNS servers are accessible
   - Check network connectivity to relay servers
   - Ensure correct format in DNS_RELAY_SERVERS (comma-separated, with ports)
   - Monitor logs for relay errors

## Issues and Support

### Opening New Issues

Before opening a new issue:

1. Check existing issues to avoid duplicates
2. Use issue templates when available
3. Include:
   - NanoDNS version
   - Operating system
   - Clear steps to reproduce
   - Expected vs actual behavior
   - Error messages if any

## Community

- Star the repository to show support
- Watch for updates and new releases
- Join discussions in issues and PRs
- Share your use cases and feedback

## Join as a Contributor

We welcome contributions! Here's how to get started:

1. Star ⭐ and watch 👀 the repository
2. Check [open issues](https://github.com/mguptahub/nanodns/issues) for tasks labeled `good first issue` or `help wanted`
3. Read our [Contributing Guide](CONTRIBUTING.md) for:
   - Development setup
   - Code style guidelines
   - PR process
   - Release workflow


## License and Usage Terms

NanoDNS is open-source software licensed under AGPLv2. This means:

✅ You CAN:
- Use NanoDNS in your development environment
- Use NanoDNS as part of your infrastructure
- Package NanoDNS with your GPL-compatible software (with attribution)
- Modify and distribute NanoDNS (while keeping it open source)

❌ You CANNOT:
- Sell NanoDNS as a standalone product
- Include NanoDNS in proprietary software
- Remove or modify the license and copyright notices

📝 You MUST:
- Include the original license
- State significant changes made
- Include the complete corresponding source code
- Include attribution to this repository
  ```
  This software uses NanoDNS (https://github.com/mguptahub/nanodns)
  ```

### Why AGPLv2?

1. **Simplicity**: 
   - Clearer and more concise terms compared to v3
   - Well-established legal precedents

2. **Compatibility**: 
   - Works well with other GPL v2 software
   - Broader ecosystem compatibility

3. **Core Protection**:
   - Ensures source code remains open
   - Prevents commercial exploitation
   - Requires attribution

### Commercial Usage Notice

While NanoDNS can be used within commercial products as a supporting utility:
1. The complete source code must be available
2. Proper attribution must be included
3. Any modifications must be shared under AGPLv2
4. It cannot be sold as a standalone product or service

### Proper Attribution

Add this to your documentation:
```markdown
This product uses NanoDNS (https://github.com/mguptahub/nanodns), 
an open-source DNS server licensed under AGPL-2.0.
```

### Important Note
If you plan to use NanoDNS in your project, ensure your project's license is compatible with AGPLv2. When in doubt, open an issue for clarification.


---

