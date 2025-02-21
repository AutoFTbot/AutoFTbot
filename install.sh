#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}
print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

if [ "$EUID" -ne 0 ]; then 
    print_error "Script harus dijalankan sebagai root"
    exit 1
fi

print_success "Mengunduh script restart.sh..."
wget -O /root/restart.sh https://raw.githubusercontent.com/AutoFTbot/AutoFTbot/main/restart.sh

chmod +x /root/restart.sh

cat > /etc/systemd/system/restart.service << 'EOL'
[Unit]
Description=Auto Restart Service Monitor
After=network.target

[Service]
Type=simple
User=root
ExecStart=/root/restart.sh
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOL

systemctl daemon-reload
systemctl enable restart
systemctl start restart

if systemctl is-active --quiet restart; then
    print_success "Service restart berhasil diinstall dan dijalankan"
    print_success "Silakan edit token bot dan chat ID di /root/restart.sh"
    print_success "Gunakan: nano /root/restart.sh"
else
    print_error "Gagal menjalankan service"
    exit 1
fi 
