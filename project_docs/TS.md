---

# ТЕХНИЧЕСКОЕ ЗАДАНИЕ (ENTERPRISE SPECIFICATION)

## Проект: HEXTEC SpecAuditor Gov-Edition (Oʻz DSt 1987:2018 Ekspertizasi)

---

### 1. БИЗНЕС-МОДЕЛЬ И КОНТЕКСТ ВНЕДРЕНИЯ

#### 1.1. Проблематика госсектора (Pain Points)

* **Нормативная обязанность:** По постановлению КМ РУз № 256 государственные информационные системы обязаны проходить государственную экспертизу ТЗ. Экспертиза проводится строго по государственному стандарту **Oʻz DSt 1987:2018** («Axborot tizimini yaratish uchun texnik topshiriq»).


* **Операционная рутина экспертов:** Эксперты ведомств вручную сверяют поступающие ТЗ (по 40–120 страниц) с десятками таблиц и пунктов стандартов Oʻz DSt 1985, 1986, 1987. На первичное чтение и выявление дефектов уходит до 5–7 рабочих дней на один документ.


* **Человеческий фактор:** Пропускаются критические расхождения (отсутствие обязательных разделов, неизмеримые формулировки SLA, игнорирование структуры кодирования документов по Oʻz DSt 1985, отсутствие обязательных этапов по Oʻz DSt 1986).


* **Проблема безопасности данных:** Государственные ТЗ содержат закрытую информацию о ведомственной инфраструктуре. Их категорически запрещено отправлять в публичные облачные API (OpenAI, Anthropic, Google Cloud).

#### 1.2. Решение и ценность (Value Proposition)

**HEXTEC SpecAuditor** — автономная интеллектуальная платформа для первичного государственного нормоконтроля и аудита спецификаций:

* **Автономия (On-Premises / Air-Gapped):** Работает полностью внутри локального контура министерства/ведомства без выхода во внешнюю сеть.
* **Сокращение рутины на 85%:** Анализирует ТЗ за 15–30 секунд и формирует предварительное экспертное заключение (*Ekspertiza xulosasi*) с точными ссылками на пункты стандарта.


* **Стандартизация замечаний:** Исключает предвзятость, формируя для разработчиков исчерпывающий перечень правок, без которых документ не допускается к тендеру или разработке.



#### 1.3. Монетизация и эксплуатационная модель

* **Для закрытого ведомства (On-Premises License):** Поставка лицензии на серверные мощности заказчика с годовой технической поддержкой и обновлениями базы стандартов.
* **Пользовательские роли (RBAC в перспективе):**
* *Ekspert (Davlat xodimi):* загрузка ТЗ, аудит, утверждение и выгрузка официального заключения.
* *Tizim maʼmuri:* управление локальными моделями, просмотр аудиторских логов.



---

### 2. СИСТЕМНАЯ И ПРОГРАММНАЯ АРХИТЕКТУРА

#### 2.1. Архитектурный стиль: Clean Architecture (Hexagonal / Ports & Adapters)

Проект реализуется на **Go 1.22+** и делится на 4 изолированных слоя:

1. **Domain Layer (`internal/domain`):**
* Чистые структуры данных, не зависящие от фреймворков и библиотек.
* Сущности: `AuditRequest`, `AuditReport`, `SectionFinding`, `MetricScore`, `StandardSection`.


2. **Usecase / Service Layer (`internal/service`):**
* Бизнес-логика валидации: предварительный синтаксический анализ структуры ТЗ, генерация промптов, оркестрация запросов к модели, расчет сводного процента соответствия (*Compliance Index*).


3. **Adapter / Infrastructure Layer (`internal/adapter`):**
* `ollama_client.go`: реализация HTTP-порта общения с локальной LLM (`qwen2.5:7b` или `llama3.1:8b`) через endpoint `/api/generate` с механизмом таймаута (`context.WithTimeout`) и очистки ответа от мусорных markdown-тегов.
* `memory_repo.go`: потокобезопасное хранилище истории аудитов (`sync.RWMutex`) для возможности мгновенного просмотра последних проверок без необходимости сразу разворачивать PostgreSQL.


4. **Delivery / Transport Layer (`internal/delivery/http`):**
* Маршрутизатор (`net/http` или легковесный роутер).
* Middleware: Request Logging, CORS, Panic Recovery.
* Статический сервер для встроенного UI через директиву `//go:embed`.



#### 2.2. Файловая структура репозитория

```text
spec-auditor/
├── cmd/
│   └── api/
│       └── main.go                 # Dependency Injection, парсинг флагов/env, запуск HTTP-сервера
├── internal/
│   ├── domain/
│   │   ├── audit.go               # Модели запросов, отчетов, замечаний
│   │   └── standards.go           # Константы разделов Oʻz DSt 1987:2018
│   ├── service/
│   │   ├── auditor.go             # Основной пайплайн аудита
│   │   └── prompt.go              # Системный промпт на узбекском языке со знанием стандарта
│   ├── adapter/
│   │   ├── llm/
│   │   │   └── ollama.go          # Клиент к локальному инференсу
│   │   └── repository/
│   │       └── memory.go          # Хранилище отчетов в оперативной памяти
│   └── delivery/
│       └── http/
│           ├── handler.go         # Контроллеры API эндпоинтов
│           └── response.go        # Стандартизированный JSON-ответ
├── web/
│   ├── static/
│   │   └── app.js                 # Клиентский движок рендеринга и работы с API
│   └── index.html                 # Полнофункциональный интерфейс (Tailwind CSS)
├── Makefile
├── go.mod
└── go.sum

```

---

### 3. НОРМАТИВНЫЙ БАЗИС И СИСТЕМНЫЙ ПРОМПТ (DOMAIN KNOWLEDGE)

Системный промпт для LLM компилируется динамически и накладывает на модель роль строгого государственного эксперта технической документации.

#### 3.1. Нормативный чеклист по Oʻz DSt 1987:2018:

1. 8 обязательных разделов (A ilovasi):


* *1. Umumiy maʼlumotlar:* наличие кода организации (8 символов по ОКПО) и обозначения ИС (до 15 символов) по Oʻz DSt 1985 (п. 7.4–7.5), сроки, реквизиты сторон.


* *2. AT vazifasi va yaratish maqsadlari:* обязательное наличие измеримых качественных и количественных критериев эффективности (п. 6.3.2).


* *3. Axborotlashtirish obyektining tavsifi:* условия эксплуатации и режимы работы.


* *4. ATga qoʻyiladigan talablar:* подсистемы, режимы (нормальный, аварийный), требования безопасности по Oʻz DSt 2814, эргономика.


* *5. ATni yaratish ishlarining tarkibi va mazmuni:* обязательное соответствие этапам создания по стандарту Oʻz DSt 1986:2018 и комплектность документов по Oʻz DSt 1985:2018.


* *6. ATni nazorat qilish va qabul qilish tartibi:* порядок испытаний по ГОСТ 34.603 и состав приемной комиссии.


* *7. ATni ishga tushirishga tayyorlash talablari:* обучение персонала, ввод классификаторов.


* *8. Hujjatlashtirishga qoʻyiladigan talablar:* комплектность по Oʻz DSt 1985.




2. Анализ сценариев использования (B ilovasi):


* Идентификаторы сценариев (U1, C1, ED1).


* Наличие строгого временного регламента (Vaqt reglamenti) в секундах/минутах для каждой операции.


* Наличие альтернативных путей и расширений (`extend`, `include`).




3. **Фильтрация размытых требований (Quality of Requirements):**
* Запрет на субъективные фразы: *«tizim tez ishlashi kerak»*, *«qulay dizayn boʻlishi lozim»*, *«yuqori yuklamalarga chidamli boʻlishi lozim»*. Выявление таких мест и требование переписать их в количественные параметры (SLA, RPS, latency, % ошибок).





#### 3.2. Спецификация контракта ответа (Strict JSON Schema)

```json
{
  "audit_id": "hex-uuid-string",
  "analyzed_at": "2026-09-08T10:00:00Z",
  "document_name": "Elektron Toʻlovlar AT Texnik Topshirigʻi",
  "compliance_score": 68.5,
  "status": "RAD ETILDI" // "TASDIQLANDI" (>=85%), "SHARTLI TASDIQLANDI" (70-84%), "RAD ETILDI" (<70%)
  "metrics": {
    "structure_completeness": 75.0,
    "scenarios_compliance": 50.0,
    "measurability_score": 60.0
  },
  "missing_sections": [
    "5. ATni yaratish ishlarining tarkibi va mazmuni (Oʻz DSt 1986 bosqichlari koʻrsatilmagan)",
    "8. Hujjatlashtirishga qoʻyiladigan talablar (Oʻz DSt 1985 komplektligi yoʻq)"
  ],
  "findings": [
    {
      "id": "FIND-01",
      "severity": "CRITICAL", // "CRITICAL", "WARNING", "NOTE"
      "section": "4.1.4 Vazifasi koʻrsatkichlari",
      "raw_fragment": "Tizim foydalanuvchilar soni koʻpayganda sekinlashmasdan tez ishlashi kerak.",
      "violation": "Oʻz DSt 1987:2018 6.5.1.4-bandi talabi buzilgan: aniq oʻlchanadigan parametrlar keltirilmagan.",
      "recommendation": "Talabni quyidagi shaklda aniqlashtiring: 'Tizim bir vaqtning oʻzida 100 ta foydalanuvchi ishlaganda ekran navigatsiyasi uchun javob vaqti 5 sekunddan oshmasligini taʼminlashi kerak'."
    }
  ],
  "executive_summary": "Taqdim etilgan Texnik topshiriq loyihasi Oʻz DSt 1987:2018 talablariga qisman mos keladi. Hujjatda 2 ta majburiy boʻlim toʻliq tushirib qoldirilgan hamda funksional talablarda oʻlchanadigan metrikalar koʻrsatilmagan. Koʻrsatilgan kamchiliklar bartaraf etilgunga qadar davlat ekspertizasidan oʻtkazish tavsiya etilmaydi."
}

```

---

### 4. ПОЛЬЗОВАТЕЛЬСКИЙ ИНТЕРФЕЙС И ОПЫТ (UI/UX DESIGN)

#### 4.1. Визуальный стиль

* **Тема:** Профессиональный Dark Mode (Tailwind slate-900 `#0f172a`, панели slate-800 `#1e293b`, акценты indigo-500 `#6366f1` и emerald-500 `#10b981`).
* **Типографика:** Межбуквенные шрифты (Inter / Roboto Mono для кода и фрагментов ТЗ).
* **Сетка:** Двухколоночный полноэкранный интерфейс (Dashboard Layout) без вертикальной перезагрузки страницы.

#### 4.2. Детализация компонентов экрана

**1. Верхняя панель (Header):**

* Слева: Герб/символика и текстовый логотип: `HEXTEC SpecAuditor // Oʻz DSt 1987:2018 Ekspertizasi`.
* По центру: Бейдж текущей версии стандартов: `Normativ baza: Oʻz DSt 1985, 1986, 1987:2018`.


* Справа: Индикатор статуса локального бэкенда и модели: зеленый пульсирующий светодиод `Lokal LLM: Qwen-2.5-7B (Faol)`.

**2. Левая рабочая панель (Панель ввода — 45% ширины):**

* **Блок быстрых действий (Top Bar):**
* Кнопка: `[ Namunaviy xato TT yuklash ]` — мгновенно подставляет в текстовую область реальный пример слабого ведомственного ТЗ с нарушением разделов и размытыми требованиями.
* Кнопка: `[ Maydonni tozalash ]`.


* **Текстовое поле ввода:**
* Увеличенная область ввода (`textarea`) с темным фоном, моноширинным шрифтом, подсветкой строк и автоскроллом.
* Placeholder: *«Texnik topshiriq matnini shu yerga kiriting yoki namunaviy faylni yuklang...»*.
* Счетчик символов и слов в реальном времени.


* **Главная кнопка запуска:**
* Широкая кнопка: `[ EKSPERTIZADAN OʻTKAZISH ]` (Indigo градиент).
* Состояние загрузки: при клике кнопка блокируется, появляется анимированный спиннер и статус: *«Oʻz DSt 1987:2018 qoidalari boʻyicha tahlil qilinmoqda (15-20 soniya)...»*.


**3. Правая панель (Интерактивное экспертное заключение — 55% ширины):**

* **Виджет итогового скоринга (Verdict Card):**
* Большой круговой индикатор или крупный цифровой прогресс-бар: `68% Muvofiqlik indeksi`.
* Статусный бейдж: Красный `RAD ETILDI` (если `<70%`), Желтый `SHARTLI TASDIQLANDI` (`70-84%`), Зеленый `DAVLAT EKSPERTIZASIGA TAYYOR` (`>=85%`).
* Кнопка экспорта: `[ PDF Xulosani yuklab olish ]` (генерирует форматированный вид для печати).


* **Сводка по метрикам (3 карточки):**
* *Tuzilma toʻliqligi (A ilovasi):* `75%`

* *Ssenariylar talabi (B ilovasi):* `50%`

* *Metrikalar aniqligi:* `60%`

* **Блок отсутствующих обязательных разделов:**
* Красные предупреждающие карточки с иконкой тревоги, перечисляющие пропущенные разделы с точным указанием пункта стандарта.

* **Лента выявленных дефектов (Findings Feed):**
* Карточка каждого замечания содержит:
* Бейдж серьезности: `CRITICAL` (красный), `WARNING` (желтый).
* Раздел ТЗ, где найден дефект (например, *4.1.4 Vazifasi koʻrsatkichlari*).


* Цитата из текста ТЗ с ошибкой (на сером фоне с красной рамкой).
* Пункт Oʻz DSt 1987:2018, который нарушен.


* Зеленый блок «Tavsiya»: готовый текст правильной формулировки для вставки в ТЗ.

---

### 5. КРИТЕРИИ ГОТОВНОСТИ И ДЕМОНСТРАЦИОННЫЙ СЦЕНАРИЙ

1. **Автономность запуска:** Репозиторий должен содержать всё необходимое. Запуск одной командой:
```bash
go run cmd/api/main.go

```


Сервер поднимается на `http://localhost:8080` и автоматически готов к работе с локальной Ollama.
2. **Предустановленный тестовый кейс (Demo Preset):** В коде клиента зашит реальный фрагмент некорректного ТЗ. При нажатии одной кнопки интерфейс заполняется, отправляет запрос и через 10–15 секунд выводит структурированный аудит на чистом узбекском языке.
