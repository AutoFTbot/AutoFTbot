#!/bin/bash

# Function to display progress bar
show_progress() {
    local progress=$1
    local total=$2
    local percentage=$((progress * 100 / total))
    printf "\rSetup process : [%-50s] %d%%" $(printf '#%.0s' $(seq 1 $((percentage / 2)))) $percentage
}

total_steps=7
current_step=0

update_progress() {
    current_step=$((current_step + 1))
    show_progress $current_step $total_steps
    sleep 1
}

# 1. Uninstall all conflicting Docker packages.
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do 
    sudo apt-get remove -y $pkg > /dev/null 2>&1
done
update_progress

# 2. Update package list dan install dependencies
sudo apt-get update -qq > /dev/null 2>&1
sudo apt-get install -y ca-certificates curl > /dev/null 2>&1
update_progress

# 3. Setup Docker GPG key
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc > /dev/null 2>&1
update_progress

# 4. Set permission untuk keyring
sudo chmod a+r /etc/apt/keyrings/docker.asc
update_progress

# 5. Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
update_progress

# 6. Update package list dengan repositori Docker baru
sudo apt-get update -qq > /dev/null 2>&1
update_progress

# 7. Install Docker packages
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin > /dev/null 2>&1
update_progress

echo -e "\nSetup complete!"

# 8. Pilih VERSION
echo -e "\n+----+----------------------------+----+----------------------------------+"
echo    "| No | Version                      | No | Version                          |"
echo -e "+----+--------------------------------+--+----------------------------------+"
echo    "| 1  | Windows 11 Pro (5.4 GB)      | 6  | Windows 8.1 Pro (4.0 GB)         |"
echo    "| 2  | Windows 11 Enterprise(5.8 GB)| 7  | Windows 8.1 Enterprise (3.7 GB)  |"
echo    "| 3  | Windows 10 Pro (5.7 GB)      | 8  | Windows 7 Enterprise (3.0 GB)    |"
echo    "| 4  | Windows 10 LTSC (5.6 GB)     | 9  | Windows Vista Enterprise (3.0 GB)|"
echo    "| 5  | Windows 10 Enterprise(5.2 GB)| 10 | Windows XP Professional (0.6 GB) |"
echo    "| 11 | Windows Server 2022 (4.7 GB) | 14 | Windows Server 2012 (4.3 GB)     |"
echo    "| 12 | Windows Server 2019 (5.3 GB) | 15 | Windows Server 2008 (3.0 GB)     |"
echo    "| 13 | Windows Server 2016 (6.5 GB) | 16 | Windows Server 2025 (5.0 GB)     |"
echo    "| 17 | Windows 11 LTSC (4.2 GB)     | 18 | Windows Server 2003 (0.6 GB)     |"
echo -e "+----+-------------------------+----+---------------------------------------+"

while true; do
  read -p "Pilih versi windows yang anda inginkan : " version_choice

  case $version_choice in
    1) VERSION="win11" ;;
    2) VERSION="win11e" ;;
    3) VERSION="win10" ;;
    4) VERSION="ltsc10" ;;
    5) VERSION="win10e" ;;
    6) VERSION="win8" ;;
    7) VERSION="win8e" ;;
    8) VERSION="win7" ;;
    9) VERSION="vista" ;;
    10) VERSION="winxp" ;;
    11) VERSION="2022" ;;
    12) VERSION="2019" ;;
    13) VERSION="2016" ;;
    14) VERSION="2012" ;;
    15) VERSION="2008" ;;
    16) VERSION="2025" ;;
    17) VERSION="ltsc11" ;;
    18) VERSION="2003" ;;
    *) echo "Pilihan tidak valid. Silakan coba lagi."; continue ;;
  esac
  break
done

# 9. Input RAM size
read -p "Size RAM : " ram_size

# 10. Input CPU cores
read -p "CPU Core : " cpu_cores

# 11. Input disk size
read -p "Size Storage : " disk_size

# 12. Input password
read -p "Password RDP : " password

# 13. Write compose.yaml file
cat <<EOL > compose.yaml
version: "3"
services:
  windows:
    image: dockurr/windows
    container_name: windows
    environment:
      VERSION: "$VERSION"
      RAM_SIZE: "${ram_size}G"
      CPU_CORES: "${cpu_cores}"
      DISK_SIZE: "${disk_size}G"
      USERNAME: "admin"
      PASSWORD: "${password}"
    devices:
      - /dev/kvm
    cap_add:
      - NET_ADMIN
    ports:
      - 8006:8006
      - 3389:3389/tcp
      - 3389:3389/udp
    stop_grace_period: 2m
    restart: on-failure
    volumes:
      - /var/win10:/storage
EOL

# 14. Run
sudo docker compose -f compose.yaml up -d
