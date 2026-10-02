# SpecAuditor Gov-Edition (Oʻz DSt 1987:2018 Ekspertizasi)

> **Avtonom davlat normonazorat va ekspertiza platformasi**  
> Axborot tizimlarini yaratish boʻyicha Texnik topshiriqlarni (TT) Oʻz DSt 1987:2018, Oʻz DSt 1985 va Oʻz DSt 1986 davlat standartlari asosida avtomatik tahlil qilish tizimi.

---

## 🚀 Asosiy imkoniyatlar (Features)

* **100% Avtonomlik (On-Premises / Air-Gapped):** Tashqi internetga yoki ochiq bulutli servislarga (OpenAI/Anthropic) ulanmasdan, toʻliq mahalliy tarmoqda ishlaydi.
* **Qwen 2.5:7B lokal LLM integratsiyasi:** Oʻzbek tili terminologiyasi va Oʻz DSt 1987 standartlari boʻyicha qatʼiy JSON-sxema asosida xulosa beradi (`num_ctx: 16384`).
* **Offline Demo Koʻrgazma (Fail-safe Demo Mode):** Agar lokal Ollama serveri oʻchirilgan yoki GPU boʻlmagan noutbukda ishga tushirilsa, tizim avtomatik ravishda namunaviy etalon ekspertiza tahlilini koʻrsatadi.
* **Yagona mustaqil binar fayl (Go Clean Architecture):** Hech qanday Node.js, Python yoki tashqi kutubxonalarsiz, frontendi ichiga tikilgan (`//go:embed`) yagona fayl shaklida tarqatiladi.
* **Rasmiy blankada chop etish / PDF Eksport:** Natijani bir marta bosish orqali rasmiy qogʻoz blank koʻrinishida PDF formatida yuklab olish imkoniyati.

---

## 📋 Tizim arxitekturasi (Clean Architecture)

```text
ai_analyst/
├── cmd/
│   └── api/main.go               # Dasturga kirish nuqtasi, DI va serverni ishga tushirish
├── internal/
│   ├── domain/
│   │   ├── audit.go              # Asosiy ma'lumotlar modellari (Request, Report, Finding)
│   │   └── standards.go          # Oʻz DSt 1987:2018 8 ta majburiy bo'limi va taqiqlangan jumlalar
│   ├── service/
│   │   ├── auditor.go            # Biznes-mantiq: validatsiya, scoring, LLM orkestratsiyasi
│   │   ├── prompt.go             # O'zbek tilidagi qat'iy davlat eksperti tizim prompri
│   │   └── mock_data.go          # Namunaviy xato TT va Offline Demo etalon xulosasi
│   ├── adapter/
│   │   ├── llm/ollama.go         # Ollama HTTP-klienti, JSON sanitarizatsiyasi va timeout
│   │   └── repository/memory.go  # Xotiradagi xavfsiz (sync.RWMutex) hisobotlar ombori
│   └── delivery/http/
│       ├── handler.go            # API kontrollerlari (/api/audit, /api/health, /api/preset)
│       ├── middleware.go         # Logging, CORS, Panic Recovery
│       └── response.go           # Standartlashgan JSON javoblar
├── web/
│   ├── static/app.js             # Front-end reaktiv dvigateli (taymer, grafik, filtrlash)
│   ├── index.html                # Dark Mode UI (Tailwind CSS, Inter/Roboto Mono)
│   └── embed.go                  # Go embed.FS orqali binar faylga resurslarni tikish
├── Makefile                      # Boshqaruv komandalari
└── go.mod
```

---

## ⚡ Tezkor ishga tushirish (Quick Start)

### 0-variant: Docker orqali bir buyruqda toʻliq ishga tushirish (Tavsiya etiladi)
Hech qanday Go yoki Ollamani alohida oʻrnatmasdan, Web-server, Ollama va AI modelini (`qwen2.5:7b`) bir buyruq bilan ishga tushirish:
```bash
docker compose up -d --build
# yoki:
make docker-up
```
> Tizim avtomatik ravishda Ollama konteynerini koʻtaradi, `qwen2.5:7b` modelini yuklab oladi va Web-interfeysni `http://localhost:8080` manzilida ishga tushiradi.
>
> *(Agar serverda NVIDIA GPU boʻlsa: `make docker-gpu`)*

---

### 🌐 Masofadan (boshqa tarmoqdan) ulanish uchun ommaviy havola:
Boshqa tumandagi yoki shahardagi hamkasblar tizimni brauzer orqali test qilishlari uchun xavfsiz Cloudflare HTTPS-havola yaratish:
```bash
make tunnel
# yoki: ./tunnel.sh
```
Terminal ekranga tayyor havolani chiqaradi (masalan: `https://xxxx.trycloudflare.com`).
* Havolani qayta koʻrish: `make tunnel-url`
* Havolani toʻxtatish: `make tunnel-stop`

---

### 1-variant: Go orqali toʻgʻridan-toʻgʻri ishga tushirish
```bash
/usr/local/go/bin/go run cmd/api/main.go
```
yoki Makefile orqali:
```bash
make run
```

Brauzerda oching: **[http://localhost:8080](http://localhost:8080)**

---

### 2-variant: Bitta mustaqil binar faylga yigʻish
```bash
make build
# yoki:
/usr/local/go/bin/go build -o bin/spec-auditor cmd/api/main.go

# Ishga tushirish:
./bin/spec-auditor
```

---

## 🤖 Lokal LLM (Ollama) bilan toʻliq rejimda ishlash

Agar kompyuteringizda yoki serveringizda Ollama oʻrnatilgan boʻlsa:

1. Qwen 2.5 7B modelini yuklab oling:
   ```bash
   ollama run qwen2.5:7b
   ```
2. SpecAuditor serverini ishga tushiring:
   ```bash
   ./bin/spec-auditor
   ```
Server avtomatik ravishda Ollama mavjudligini aniqlaydi va yuqori oʻng burchakdagi koʻrsatkich yashil rangda yonadi:  
`Lokal LLM: Qwen-2.5-7B (Faol)`.

*(Agar Ollama ishga tushmagan boʻlsa, tizim avtomatik ravishda **Offline Demo** rejimida ishlaydi, hech qanday xatolik yuz bermaydi).*

---

## 📡 API Endpoints

| Metod | Manzil | Tavsif |
| :--- | :--- | :--- |
| `POST` | `/api/audit` | Texnik topshiriq matnini audit qilish (JSON: `{"document_name": "...", "content": "..."}`) |
| `GET` | `/api/preset` | Namunaviy xato TT matnini olish |
| `GET` | `/api/health` | Server va lokal Ollama holatini tekshirish |
| `GET` | `/api/audit/{id}` | Muolaja qilingan hisobotni ID orqali olish |
| `GET` | `/api/history` | Soʻnggi oʻtkazilgan ekspertizalar roʻyxati |

---

## 🧪 Sinovlar (Unit Tests)

Barcha testlarni ishga tushirish:
```bash
make test
# yoki
/usr/local/go/bin/go test -v ./...
```
