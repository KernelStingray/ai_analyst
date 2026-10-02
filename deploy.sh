#!/usr/bin/env bash
# SpecAuditor Gov-Edition - On-Premises Deployment Script
set -e

echo "=================================================================="
echo "  SpecAuditor Gov-Edition (Oʻz DSt 1987:2018 Ekspertizasi)  "
echo "  On-Premises / Air-Gapped Deployment                             "
echo "=================================================================="

# Check if Docker is installed
if command -v docker &> /dev/null && (command -v docker-compose &> /dev/null || docker compose version &> /dev/null); then
    echo "[INFO] Docker aniqlandi. Konteynerlar (Web + AI Ollama + Model) ishga tushirilmoqda..."
    docker compose up -d --build
    echo "[SUCCESS] SpecAuditor ishga tushdi: http://localhost:8080"
    echo "[INFO] Loglarni kuzatish uchun: docker compose logs -f"
    exit 0
fi

# Fallback to local binary build
echo "[INFO] Docker topilmadi. Mahalliy binar fayl orqali ishga tushirilmoqda..."
make build
./bin/spec-auditor -port 8080
