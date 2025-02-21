#!/bin/bash

toket="token_bot"
chat_id="chat_id"

ngirim_duit() {
    local message="$1"
    curl -s -X POST "https://api.telegram.org/bot$toket/sendMessage" \
        -d "chat_id=$chat_id" \
        -d "text=$message"
}

check_status() {
    local service=$1
    if systemctl is-active --quiet "$service"; then
        return 0
    else
        systemctl restart "$service"
        if systemctl is-active --quiet "$service"; then
            ngirim_duit "⚠️ Layanan $service telah di-restart setelah mengalami kesalahan."
        else
            ngirim_duit "❌ Layanan $service gagal di-restart setelah mengalami kesalahan."
            return 1
        fi
    fi
}

services=("ssh" "dropbear" "ws" "openvpn" "nginx" "haproxy")

while true; do
    for s in "${services[@]}"; do
        declare "$s"="$(check_status "$s")"
        if ! check_status "$s"; then
            ngirim_duit "🚨 Layanan $s tidak aktif setelah di-restart."
        fi
    done
    if ! check_status "vmess@config" || ! check_status "vless@config" || ! check_status "trojan@config" || ! check_status "shadowsocks@config"; then
        ngirim_duit "🚨 VPS akan di-reboot karena ada layanan yang tidak aktif."
        reboot
    fi
    sleep 300
done
