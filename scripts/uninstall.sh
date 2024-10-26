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

# Stop and disable service
print_status $YELLOW "Stopping and disabling NanoDNS service..."
systemctl stop nanodns
systemctl disable nanodns

# Remove service file
print_status $YELLOW "Removing systemd service..."
rm -f "$SERVICE_PATH"
systemctl daemon-reload

# Remove binary
print_status $YELLOW "Removing NanoDNS binary..."
rm -f "$BINARY_PATH"

# Optionally remove configuration
read -p "Do you want to remove configuration files? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status $YELLOW "Removing configuration files..."
    rm -rf "$CONFIG_PATH"
fi

print_status $GREEN "NanoDNS has been uninstalled successfully!"