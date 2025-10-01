#!/bin/bash

# Docker Installation Script for Ubuntu/Debian (All Versions)
# Compatible with Ubuntu 18.04+, Debian 9+, and derivatives

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root. Consider running as regular user and using sudo."
    fi
}

detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$NAME
        VERSION=$VERSION_ID
        CODENAME=$VERSION_CODENAME
    else
        log_error "Cannot detect OS version. This script requires /etc/os-release"
        exit 1
    fi
    
    log_info "Detected OS: $OS $VERSION ($CODENAME)"
}

check_docker_installed() {
    if command -v docker &> /dev/null; then
        log_warning "Docker is already installed: $(docker --version)"
        read -p "Do you want to reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Skipping Docker installation"
            return 1
        fi
    fi
    return 0
}

update_packages() {
    log_info "Updating package lists..."
    sudo apt-get update
    log_success "Package lists updated"
}

install_prerequisites() {
    log_info "Installing prerequisites..."
    sudo apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        software-properties-common \
        wget
    log_success "Prerequisites installed"
}

add_docker_gpg_key() {
    log_info "Adding Docker's official GPG key..."
    sudo rm -f /usr/share/keyrings/docker-archive-keyring.gpg
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    log_success "Docker GPG key added"
}

add_docker_repo() {
    log_info "Adding Docker repository..."
    ARCH=$(dpkg --print-architecture)
    if [[ "$OS" == *"Ubuntu"* ]]; then
        echo "deb [arch=$ARCH signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $CODENAME stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    elif [[ "$OS" == *"Debian"* ]]; then
        echo "deb [arch=$ARCH signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $CODENAME stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    else
        log_error "Unsupported OS: $OS"
        exit 1
    fi
    
    log_success "Docker repository added"
}

install_docker() {
    log_info "Installing Docker..."
    sudo apt-get update
    sudo apt-get install -y \
        docker-ce \
        docker-ce-cli \
        containerd.io \
        docker-buildx-plugin \
        docker-compose-plugin
    
    log_success "Docker installed successfully"
}

start_docker_service() {
    log_info "Starting Docker service..."
    sudo systemctl start docker
    sudo systemctl enable docker
    sleep 3
    if sudo systemctl is-active --quiet docker; then
        log_success "Docker service is running"
    else
        log_error "Failed to start Docker service"
        exit 1
    fi
}

add_user_to_docker_group() {
    log_info "Adding current user to docker group..."
    sudo usermod -aG docker $USER
    log_success "User added to docker group"
    log_warning "You need to log out and log back in for group changes to take effect"
    log_warning "Or run: newgrp docker"
}

verify_installation() {
    log_info "Verifying installation..."
    if command -v docker &> /dev/null; then
        log_success "Docker version: $(docker --version)"
    else
        log_error "Docker command not found"
        exit 1
    fi
    if command -v docker &> /dev/null && docker compose version &> /dev/null; then
        log_success "Docker Compose version: $(docker compose version)"
    else
        log_warning "Docker Compose not found or not working"
    fi
    log_info "Testing Docker with hello-world container..."
    if sudo docker run --rm hello-world &> /dev/null; then
        log_success "Docker is working correctly"
    else
        log_warning "Docker test failed, but installation might still be correct"
    fi
}

main() {
    log_info "Starting Docker installation for $OS $VERSION"
    
    check_root
    detect_os
    
    if ! check_docker_installed; then
        exit 0
    fi
    
    update_packages
    install_prerequisites
    add_docker_gpg_key
    add_docker_repo
    install_docker
    start_docker_service
    add_user_to_docker_group
    verify_installation
    
    log_success "Docker installation completed successfully!"
    log_info "To use Docker without sudo, log out and log back in, or run: newgrp docker"
    log_info "To start using Docker Compose: docker compose --help"
}

main "$@"
