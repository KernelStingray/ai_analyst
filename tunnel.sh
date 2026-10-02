#!/usr/bin/env bash
set -e

ACTION="${1:-start}"

if [ "$ACTION" = "stop" ]; then
    echo "[INFO] Остановка туннелей..."
    docker compose --profile tunnel rm -f -s tunnel 2>/dev/null || true
    if [ -f .cloudflared.pid ]; then
        kill $(cat .cloudflared.pid) 2>/dev/null || true
        rm -f .cloudflared.pid
    fi
    pkill -f "cloudflared tunnel" 2>/dev/null || true
    rm -f cloudflared.log
    echo "[DONE] Туннель остановлен."
    exit 0
fi

if [ "$ACTION" = "url" ]; then
    URL=""
    # Check if Docker tunnel container is actively running
    if docker ps --filter "name=spec-tunnel" --filter "status=running" --format '{{.Names}}' 2>/dev/null | grep -q 'spec-tunnel'; then
        URL=$(docker compose logs tunnel 2>&1 | grep -o 'https://[-a-zA-Z0-9]*\.trycloudflare\.com' | tail -n 1 || true)
    elif [ -f .cloudflared.pid ] && kill -0 $(cat .cloudflared.pid 2>/dev/null) 2>/dev/null; then
        if [ -f cloudflared.log ]; then
            URL=$(grep -o 'https://[-a-zA-Z0-9]*\.trycloudflare\.com' cloudflared.log | tail -n 1 || true)
        fi
    fi

    if [ -n "$URL" ]; then
        echo "=================================================================="
        echo "  🌐 АКТУАЛЬНАЯ ПУБЛИЧНАЯ ССЫЛКА:"
        echo "     👉 $URL"
        echo "=================================================================="
    else
        echo "[INFO] Активный туннель не запущен."
        echo "       Для создания нового туннеля выполните: make tunnel"
    fi
    exit 0
fi

echo "=================================================================="
echo "  SpecAuditor - Запуск публичного HTTPS-туннеля (Cloudflare)     "
echo "=================================================================="

# Check if Docker compose containers are running
if docker ps --format '{{.Names}}' | grep -q 'spec-auditor'; then
    echo "[INFO] Обнаружен запущенный контейнер Docker spec-auditor."
    echo "[INFO] Пересоздаем туннель для генерации нового адреса..."
    docker compose --profile tunnel rm -f -s tunnel 2>/dev/null || true
    docker compose --profile tunnel up -d --no-deps tunnel

    echo -n "[INFO] Ожидание генерации защищенной ссылки"
    URL=""
    for i in {1..30}; do
        URL=$(docker compose logs tunnel 2>&1 | grep -o 'https://[-a-zA-Z0-9]*\.trycloudflare\.com' | tail -n 1 || true)
        if [ -n "$URL" ]; then
            echo " [OK]"
            break
        fi
        echo -n "."
        sleep 1
    done
    if [ -z "$URL" ]; then
        echo " [TIMEOUT]"
    fi
else
    # Local fallback
    echo "[INFO] Docker контейнер spec-auditor не найден. Переход в локальный режим..."
    if [ ! -f "bin/cloudflared" ]; then
        echo "[INFO] Скачивание официального бинарника cloudflared..."
        mkdir -p bin
        curl -sL -o bin/cloudflared https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
        chmod +x bin/cloudflared
    fi

    # Stop any running local cloudflared instance
    if [ -f .cloudflared.pid ]; then
        kill $(cat .cloudflared.pid) 2>/dev/null || true
        rm -f .cloudflared.pid
    fi
    pkill -f "cloudflared tunnel" 2>/dev/null || true
    rm -f cloudflared.log

    echo "[INFO] Запуск туннеля для http://localhost:8080..."
    nohup ./bin/cloudflared tunnel --no-autoupdate --protocol http2 --url http://localhost:8080 > cloudflared.log 2>&1 &
    echo $! > .cloudflared.pid

    echo -n "[INFO] Ожидание генерации защищенной ссылки"
    URL=""
    for i in {1..30}; do
        URL=$(grep -o 'https://[-a-zA-Z0-9]*\.trycloudflare\.com' cloudflared.log 2>/dev/null | tail -n 1 || true)
        if [ -n "$URL" ]; then
            echo " [OK]"
            break
        fi
        echo -n "."
        sleep 1
    done
    if [ -z "$URL" ]; then
        echo " [TIMEOUT]"
    fi
fi

if [ -n "$URL" ]; then
    echo ""
    echo "=================================================================="
    echo "  ✅ ТУННЕЛЬ УСПЕШНО СГЕНЕРИРОВАН И ЗАПУЩЕН!"
    echo "  "
    echo "  🌐 ПУБЛИЧНАЯ ССЫЛКА ДЛЯ ТЕСТИРОВАНИЯ:"
    echo "     👉 $URL"
    echo "  "
    echo "  Отправьте эту ссылку коллегам. Доступ работает из любой сети!"
    echo "  Посмотреть ссылку снова: make tunnel-url"
    echo "  Для остановки выполните: make tunnel-stop"
    echo "=================================================================="
else
    echo ""
    echo "[WARN] Туннель запускается, но ссылка пока не получена."
    echo "       Проверьте логи командой: docker compose logs tunnel"
    echo "       Или запросите ссылку снова: make tunnel-url"
fi
