.PHONY: all run build test clean help docker-up docker-down docker-logs docker-gpu tunnel tunnel-url tunnel-stop

GO ?= /usr/local/go/bin/go
BINARY_NAME = spec-auditor

help:
	@echo "SpecAuditor Gov-Edition (Oʻz DSt 1987:2018)"
	@echo "--------------------------------------------------"
	@echo "make docker-up   - Docker orqali barcha servislarni (AI + Web) bir buyruq bilan ishga tushirish"
	@echo "make docker-gpu  - Docker orqali NVIDIA GPU tezlatgich bilan ishga tushirish"
	@echo "make docker-down - Docker konteynerlarini to'xtatish"
	@echo "make docker-logs - Konteyner loglarini kuzatish"
	@echo "make tunnel      - Masofadan ulanish uchun xavfsiz Cloudflare HTTPS havola yaratish"
	@echo "make tunnel-url  - Hozirgi faol ommaviy havolani ko'rish"
	@echo "make tunnel-stop - Ommaviy tunnelni to'xtatish"
	@echo "make run         - Lokal Go orqali ishga tushirish (http://localhost:8080)"
	@echo "make build       - Bitta mustaqil binar faylga yig'ish (bin/$(BINARY_NAME))"
	@echo "make test        - Barcha unit-testlarni ishga tushirish"
	@echo "make clean       - Yig'ilgan binar fayllarni tozalash"

tunnel:
	@./tunnel.sh

tunnel-url:
	@./tunnel.sh url

tunnel-stop:
	@./tunnel.sh stop

docker-up:
	docker compose up -d --build
	@echo "SpecAuditor ishga tushmoqda: http://localhost:8080"
	@echo "AI model holatini kuzatish uchun: make docker-logs"

docker-gpu:
	docker compose -f docker-compose.yml -f docker-compose.gpu.yml up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

run:
	$(GO) run cmd/api/main.go

build:
	mkdir -p bin
	$(GO) build -v -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "Tayyor: bin/$(BINARY_NAME)"

test:
	$(GO) test -v ./...

clean:
	rm -rf bin/
