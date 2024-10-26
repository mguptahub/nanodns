[![Build](https://github.com/mguptahub/nanodns/actions/workflows/build.yml/badge.svg)](https://github.com/mguptahub/nanodns/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/mguptahub/nanodns?sort=semver)](https://github.com/mguptahub/nanodns/releases)
[![Issues](https://img.shields.io/github/issues/mguptahub/nanodns)](https://github.com/mguptahub/nanodns/issues)

[![License: AGPL v2](https://img.shields.io/badge/License-AGPL%20v2-blue.svg)](https://www.gnu.org/licenses/agpl-2.0)
[![Security Policy](https://img.shields.io/badge/Security-Policy-blue.svg)](SECURITY.md)

[![Docker Pulls](https://img.shields.io/docker/pulls/mguptahub/nanodns)](https://github.com/mguptahub/nanodns/pkgs/container/nanodns)
[![Docker Image Size](https://img.shields.io/docker/image-size/mguptahub/nanodns/latest)](https://github.com/mguptahub/nanodns/pkgs/container/nanodns)

[![Go Version](https://img.shields.io/github/go-mod/go-version/mguptahub/nanodns)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/mguptahub/nanodns)](https://goreportcard.com/report/github.com/mguptahub/nanodns)


# Nano DNS Server


A lightweight DNS server designed for Docker Compose environments, allowing dynamic resolution of service names and custom DNS records.

## Features

- Environment variable-based configuration
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

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| DNS_PORT | UDP port for DNS server | 53 | 5353 |
| A_xxx | A Record Details | - | - |
| CNAME_xxx | CNAME Record Details | - | - |
| MX_xxx | MX Record Details | - | - |
| TXT_xxx | TXT Record Details | - | - |

### Record Format

All records use the `|` character as a separator. The general format is:
```
RECORD_TYPE_NUMBER=domain|value[|ttl]
```

### A Records
```
A_REC1=domain|ip|ttl
A_REC2=domain|service:servicename|ttl
```
Example:
```
A_REC1=app.example.com|192.168.1.10|300
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


## Docker Usage

### Using Docker Run
```bash
docker run -d \
  --name nanodns \
  -p 5353:5353/udp \
  -e DNS_PORT=5353 \
  -e "A_REC1=app.example.com|192.168.1.10|300" \
  -e "A_REC2=api.example.com|service:webapp" \
  -e "TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all" \
  ghcr.io/mguptahub/nanodns:latest
```

### Using Docker Compose
```yaml
name: nanodns
services:
  dns:
    image: ghcr.io/mguptahub/nanodns:latest
    environment:
      - DNS_PORT=5353  # Optional, defaults to 53
      # A Records
      - A_REC1=app.example.com|service:webapp
      - A_REC2=api.example.com|192.168.1.10|300
      # TXT Records
      - TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all
    ports:
      - "${DNS_PORT:-5353}:${DNS_PORT:-5353}/udp"  # Uses DNS_PORT if set, otherwise 5353
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```

### Kubernetes
For detailed instructions on deploying NanoDNS in Kubernetes, see our [Kubernetes Deployment Guide](kubernetes/README.md).

## Running Without Docker Compose

```bash
# Set environment variables
export DNS_PORT=5353
export A_REC1=app.example.com|192.168.1.10
export TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all

# Run the server
./nanodns
```

## Testing Records

```bash
# Test using custom port
dig @localhost -p 5353 app.example.com A

# Test CNAME record
dig @localhost -p 5353 www.example.com CNAME

# Test MX record
dig @localhost -p 5353 example.com MX

# Test TXT record
dig @localhost -p 5353 example.com TXT
```

## Common Issues and Solutions

1. Port 53 already in use (common on macOS and Linux):
   - Use a different port by setting `DNS_PORT=5353` or another available port
   - Update your client configurations to use the custom port

2. Permission denied when using port 53:
   - Use a port number above 1024 to avoid requiring root privileges
   - Set `DNS_PORT=5353` or another high-numbered port

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

### Join as a Contributor

We welcome contributions! Here's how to get started:

1. Star ⭐ and watch 👀 the repository
2. Check [open issues](https://github.com/mguptahub/nanodns/issues) for tasks labeled `good first issue` or `help wanted`
3. Read our [Contributing Guide](CONTRIBUTING.md) for:
   - Development setup
   - Code style guidelines
   - PR process
   - Release workflow

### Community

- Star the repository to show support
- Watch for updates and new releases
- Join discussions in issues and PRs
- Share your use cases and feedback

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- How to create PRs
- Code style guidelines
- Release process


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

