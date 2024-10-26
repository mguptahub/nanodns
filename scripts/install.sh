#!/bin/bash

# Define colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Define variables
SERVICE_NAME="nanodns"
BINARY_PATH="/usr/local/bin/nanodns"
SERVICE_PATH="/etc/systemd/system/nanodns.service"
CONFIG_PATH="/etc/nanodns"
ENV_FILE="${CONFIG_PATH}/nanodns.env"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if script is run as root
if [ "$EUID" -ne 0 ]; then
    print_status $RED "Please run as root"
    exit 1
fi

# Create directories
print_status $YELLOW "Creating directories..."
mkdir -p "$CONFIG_PATH"

# Copy binary
print_status $YELLOW "Installing NanoDNS binary..."
if [ -f "./nanodns" ]; then
    cp ./nanodns "$BINARY_PATH"
    chmod +x "$BINARY_PATH"
else
    print_status $RED "nanodns binary not found in current directory"
    exit 1
fi

# Create environment file if it doesn't exist
if [ ! -f "$ENV_FILE" ]; then
    print_status $YELLOW "Creating default environment file..."
    cat > "$ENV_FILE" << EOF
# NanoDNS Environment Configuration

# DNS server port (default: 53)
DNS_PORT=53

# DNS Records
# Format: domain|value|ttl
# Examples:
# A_REC1=app.local|192.168.1.10|300
# A_REC2=api.local|service:myservice
# CNAME_REC1=www.local|app.local
# MX_REC1=local|10|mail.local
# TXT_REC1=local|v=spf1 include:_spf.google.com ~all

# Add your records below:
A_REC1=app.local|127.0.0.1|300
EOF
fi

# Create systemd service file
print_status $YELLOW "Creating systemd service..."
cat > "$SERVICE_PATH" << EOF
[Unit]
Description=NanoDNS Server
After=network.target

[Service]
Type=simple
User=root
EnvironmentFile=${ENV_FILE}
ExecStart=${BINARY_PATH}
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

[Install]
WantedBy=multi-user.target
EOF

# Set permissions
print_status $YELLOW "Setting permissions..."
chmod 644 "$SERVICE_PATH"
chmod 600 "$ENV_FILE"

# Reload systemd
print_status $YELLOW "Reloading systemd..."
systemctl daemon-reload

# Enable and start service
print_status $YELLOW "Enabling and starting NanoDNS service..."
systemctl enable nanodns
systemctl start nanodns

# Check service status
if systemctl is-active --quiet nanodns; then
    print_status $GREEN "NanoDNS service has been installed and started successfully!"
    print_status $GREEN "\nUseful commands:"
    echo "  Check status: systemctl status nanodns"
    echo "  View logs: journalctl -u nanodns"
    echo "  Edit configuration: nano ${ENV_FILE}"
    echo "  Restart service: systemctl restart nanodns"
else
    print_status $RED "Failed to start NanoDNS service. Please check the logs:"
    echo "  journalctl -u nanodns"
fi

print_status $YELLOW "\nConfiguration file location: ${ENV_FILE}"
print_status $YELLOW "Please edit this file to add your DNS records"