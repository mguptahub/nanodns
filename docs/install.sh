#!/bin/sh

# install.sh - NanoDNS Installation Script

set -e  # Exit on error
set -u  # Exit on undefined variable

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="mguptahub/nanodns"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/usr/local/share"
BINARY_NAME="nanodns"
ENV_FILE="nanodns.env"
TEMP_DIR="/tmp/nanodns_install"
ENV_SOURCE="https://raw.githubusercontent.com/${GITHUB_REPO}/main/.env.example"

# Get latest version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
    jq -r '.tag_name'
}

# Print banner with ASCII art
print_banner() {
    # Get the latest version
    VERSION=$(get_latest_version)
    
    echo
    echo "${WHITE}"
    echo "  ███╗   ██╗ █████╗ ███╗   ██╗ ██████╗"
    echo "  ████╗  ██║██╔══██╗████╗  ██║██╔═══██╗"
    echo "  ██╔██╗ ██║███████║██╔██╗ ██║██║   ██║"
    echo "  ██║╚██╗██║██╔══██║██║╚██╗██║██║   ██║"
    echo "  ██║ ╚████║██║  ██║██║ ╚████║╚██████╔╝"
    echo "  ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝"
    echo "${CYAN}              DNS SERVER${NC}"
    echo
    echo "  ${CYAN}• Lightweight DNS Server${NC}"
    if [ -n "$VERSION" ]; then
        echo "  ${CYAN}• Version: ${VERSION}${NC}"
    fi
    echo "  ${CYAN}• GitHub: ${GITHUB_REPO}${NC}"
    echo
    echo "=================================================="
    echo
}

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
    arch=$(uname -m)
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
    binary_path="${TEMP_DIR}/${BINARY_NAME}"
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
    use_sudo=""
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
    if [ ! -x "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        print_error "Installation failed. Binary not found in ${INSTALL_DIR}"
    fi    
    print_success "Installation completed"
}

# Verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    # Check version
    version=""
    if ! version=$($BINARY_NAME -v 2>&1); then  
        print_error "Failed to execute $BINARY_NAME. Please check the installation."  
    fi  
    print_success "Successfully installed NanoDNS $version"
}

# Cleanup temporary files
cleanup() {
    print_step "Cleaning up..."
    rm -rf "$TEMP_DIR"
    print_success "Cleanup completed"
}

# Download environment file
download_env_file() {
    print_step "Downloading environment configuration..."
    
    # Create config directory if it doesn't exist
    if [ ! -d "$CONFIG_DIR" ]; then
        use_sudo=""
        if [ ! -w "$(dirname "$CONFIG_DIR")" ]; then
            if command -v sudo >/dev/null 2>&1; then
                use_sudo="sudo"
            else
                print_error "Config directory is not writable and sudo is not available"
            fi
        fi
        $use_sudo mkdir -p "$CONFIG_DIR"
    fi
    
    # Download .env.example file
    if ! curl -fsSL "$ENV_SOURCE" -o "${TEMP_DIR}/${ENV_FILE}"; then
        print_error "Failed to download environment configuration"
    fi
    
    # Install env file
    use_sudo=""
    if [ ! -w "$CONFIG_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            use_sudo="sudo"
        else
            print_error "Config directory is not writable and sudo is not available"
        fi
    fi
    
    $use_sudo install -m 644 "${TEMP_DIR}/${ENV_FILE}" "${CONFIG_DIR}/${ENV_FILE}"
    print_success "Environment configuration installed"
}

# Configure environment variable
configure_environment() {
    print_step "Configuring environment..."
    
    # Detect current shell
    CURRENT_SHELL=$(basename "$SHELL")
    
    # Define shell RC file based on current shell
    case "$CURRENT_SHELL" in
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                SHELL_RC="$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                SHELL_RC="$HOME/.bash_profile"
            fi
            ;;
        zsh)
            SHELL_RC="$HOME/.zshrc"
            ;;
        *)
            if [ -f "$HOME/.profile" ]; then
                SHELL_RC="$HOME/.profile"
            fi
            ;;
    esac
    
    if [ -n "$SHELL_RC" ]; then
        print_step "Detected shell: $CURRENT_SHELL, configuring in $SHELL_RC"
        
        # Check if export already exists
        if ! grep -q "export NANODNS_ENV_FILE=" "$SHELL_RC"; then
            echo >> "$SHELL_RC"  # Add a newline for better formatting
            echo "# NanoDNS Environment Configuration" >> "$SHELL_RC"
            echo "export NANODNS_ENV_FILE=\"${CONFIG_DIR}/${ENV_FILE}\"" >> "$SHELL_RC"
            print_success "Environment variable configured in $SHELL_RC"
            echo "  ${YELLOW}Note: Run 'source $SHELL_RC' to apply changes in current session${NC}"
        else
            print_warning "Environment variable already configured in $SHELL_RC"
        fi
    else
        print_warning "No suitable shell configuration file found for $CURRENT_SHELL"
        print_warning "Please manually add this line to your shell configuration:"
        echo "  export NANODNS_ENV_FILE=\"${CONFIG_DIR}/${ENV_FILE}\""
    fi
    
    # Add to /etc/environment for system-wide configuration
    use_sudo=""
    if [ ! -w "/etc/environment" ]; then
        if command -v sudo >/dev/null 2>&1; then
            use_sudo="sudo"
        else
            print_warning "Cannot write to /etc/environment. Sudo not available."
            return 0
        fi
    fi
    
    if [ -n "$use_sudo" ]; then
        if ! $use_sudo grep -q "NANODNS_ENV_FILE=" "/etc/environment"; then
            echo "NANODNS_ENV_FILE=${CONFIG_DIR}/${ENV_FILE}" | $use_sudo tee -a /etc/environment > /dev/null
            print_success "Environment variable configured system-wide in /etc/environment"
        fi
    fi
}

# Uninstall NanoDNS
do_uninstall() {
    print_banner
    print_step "Uninstalling NanoDNS..."

    # Check if binary exists
    if [ ! -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        print_warning "NanoDNS binary is not installed in ${INSTALL_DIR}"
    else
        # Remove binary
        use_sudo=""
        if [ ! -w "$INSTALL_DIR" ]; then
            if command -v sudo >/dev/null 2>&1; then
                use_sudo="sudo"
            else
                print_error "Install directory is not writable and sudo is not available"
            fi
        fi
        
        if ! $use_sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}"; then
            print_error "Failed to remove NanoDNS binary"
        fi
    fi

    # Remove env file
    if [ -f "${CONFIG_DIR}/${ENV_FILE}" ]; then
        use_sudo=""
        if [ ! -w "$CONFIG_DIR" ]; then
            if command -v sudo >/dev/null 2>&1; then
                use_sudo="sudo"
            else
                print_error "Config directory is not writable and sudo is not available"
            fi
        fi
        
        if ! $use_sudo rm -f "${CONFIG_DIR}/${ENV_FILE}"; then
            print_warning "Failed to remove environment configuration"
        fi
    else
        print_warning "Environment configuration not found in ${CONFIG_DIR}"
    fi

    print_success "NanoDNS has been successfully uninstalled"
    print_warning "Please remove NANODNS_ENV_FILE from your shell configuration and /etc/environment manually"
}

# Perform installation
do_install() {
    print_banner
    check_requirements
    detect_system
    get_download_url
    download_binary
    download_env_file
    install_binary
    configure_environment
    verify_installation
    cleanup
    
    echo
    print_success "NanoDNS has been successfully installed!"
    echo "• Binary location: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "• Config location: ${CONFIG_DIR}/${ENV_FILE}"
    echo "• Run 'nanodns --help' to get started"
}

# Print usage information
usage() {
    print_banner
    echo "Usage: $0 [--install|--uninstall|--help]"
    echo
    echo "Commands:"
    echo "  --install      Install NanoDNS"
    echo "  --uninstall    Uninstall NanoDNS"
    echo "  --help         Show this help message (default)"
    echo
    echo "Examples:"
    echo "  Install:    curl -fsSL https://nanodns.mguptahub.com/install.sh | sh -s -- --install"
    echo "  Uninstall:  curl -fsSL https://nanodns.mguptahub.com/install.sh | sh -s -- --uninstall"
    echo "  Help:       curl -fsSL https://nanodns.mguptahub.com/install.sh | sh -s"
    echo
    echo "Note: The '--' in the examples is used to separate shell arguments from script arguments"
}

# Main process
main() {
    # Handle no arguments case (default to help)
    if [ $# -eq 0 ]; then
        usage
        exit 0
    fi

    # Parse command line arguments
    for arg in "$@"; do
        case $arg in
            --install)
                do_install
                exit 0
                ;;
            --uninstall)
                do_uninstall
                exit 0
                ;;
            --help)
                usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $arg"
                usage
                exit 1
                ;;
        esac
    done
}

# Run main with all arguments
main "$@"