#!/usr/bin/env bash

# install.sh - NanoDNS Installation Script

set -e  # Exit on error
set -u  # Exit on undefined variable

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="mguptahub/nanodns"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="nanodns"
TEMP_DIR="/tmp/nanodns_install"

# Print step information
print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

# Print success message
print_success() {
    echo -e "${GREEN}==>${NC} $1"
}

# Print error message and exit
print_error() {
    echo -e "${RED}Error:${NC} $1" >&2
    exit 1
}

# Print warning message
print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Check if command exists
check_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        print_error "Required command '$1' not found. Please install it first."
    fi
}

# Check system requirements
check_requirements() {
    print_step "Checking system requirements..."
    
    # Check for required commands
    check_command curl
    check_command jq
    check_command grep
}

# Detect system information
detect_system() {
    print_step "Detecting system information..."
    
    # Detect OS
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    # Detect architecture and normalize names
    local arch=$(uname -m)
    case $arch in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        *)       print_error "Unsupported architecture: $arch" ;;
    esac
    
    print_success "Detected: $OS-$ARCH"
}

# Get the latest release URL for current system
get_download_url() {
    print_step "Finding latest release..."
    
    DOWNLOAD_URL=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
                  jq -r '.assets[].browser_download_url' | \
                  grep "${OS}-${ARCH}$" || true)
    
    if [ -z "$DOWNLOAD_URL" ]; then
        print_error "No release found for ${OS}-${ARCH}"
    fi
    
    print_success "Found release: $DOWNLOAD_URL"
}

# Download and verify the binary
download_binary() {
    print_step "Downloading NanoDNS..."
    
    # Create and clean temp directory
    mkdir -p "$TEMP_DIR"
    rm -rf "${TEMP_DIR:?}/*"
    
    # Download binary
    local binary_path="${TEMP_DIR}/${BINARY_NAME}"
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$binary_path"; then
        print_error "Failed to download binary"
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    print_success "Download completed"
}

# Install the binary
install_binary() {
    print_step "Installing NanoDNS..."
    
    # Check if we need sudo
    local use_sudo=""
    if [ ! -w "$INSTALL_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            use_sudo="sudo"
        else
            print_error "Install directory is not writable and sudo is not available"
        fi
    fi
    
    # Install binary
    $use_sudo install -m 755 "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Verify installation
    if ! command -v $BINARY_NAME >/dev/null 2>&1; then
        print_error "Installation failed. Binary not found in PATH"
    fi
    
    print_success "Installation completed"
}

# Verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    # Check version
    local version
    version=$($BINARY_NAME -v 2>&1)
    print_success "Successfully installed NanoDNS $version"
}

# Cleanup temporary files
cleanup() {
    print_step "Cleaning up..."
    rm -rf "$TEMP_DIR"
    print_success "Cleanup completed"
}

# Main installation process
main() {
    echo "NanoDNS Installer"
    echo "----------------"
    
    check_requirements
    detect_system
    get_download_url
    download_binary
    install_binary
    verify_installation
    cleanup
    
    echo
    print_success "NanoDNS has been successfully installed!"
    echo "Run 'nanodns --help' to get started"
}

# Run main if script is executed directly (not sourced)
if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    trap cleanup EXIT
    main "$@"
fi