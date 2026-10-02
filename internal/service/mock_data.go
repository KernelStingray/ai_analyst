package service

import (
	"time"

	"github.com/hextec/spec-auditor/internal/domain"
)

// DemoPresetText provides a realistic flawed governmental IT specification fragment in Uzbek
const DemoPresetText = `OʻZBEKISTON RESPUBLIKASI VAZIRLIGI HUZURIDAGI MARKAZ
TEXNIK TOPSHIRIQ LOYIHASI

1. Umumiy maʼlumotlar
1.1. Tizimning toʻliq nomi: Davlat xizmatlarini monitoring qilish va statistik tahlil axborot tizimi.
1.2. Qisqartirilgan nomi: "DXM-Stat".
1.3. Buyurtmachi: Markaz rahbariyati.
1.4. Ishlab chiquvchi: Tender gʻolibi boʻlgan kompaniya.
1.5. Ishlarni bajarish muddati: 2026-yil oxirigacha.
(Izoh: Tashkilotning OKPO kodi va Oʻz DSt 1985 boʻyicha belgilanishi koʻrsatilmagan).

2. AT vazifasi va yaratish maqsadlari
2.1. Tizimning asosiy maqsadi — aholiga xizmat koʻrsatish sifatini oshirish va hisobotlarni tezlashtirish.
2.2. Tizim idoralararo maʼlumotlar almashinuvini yaxshilashi va qogʻozbozlikni kamaytirishi lozim.

3. Axborotlashtirish obyektining tavsifi
3.1. Obyekt — Markazning markaziy apparati va 14 ta hududiy boshqarmalari.
3.2. Foydalanuvchilar soni — taxminan 500 nafar xodim.

4. ATga qoʻyiladigan talablar
4.1. Umumiy talablar:
4.1.1. Tizim foydalanuvchilar soni koʻpayganda ham sekinlashmasdan tez ishlashi kerak.
4.1.2. Foydalanuvchilar uchun juda qulay dizayn boʻlishi lozim.
4.1.3. Tizim zamonaviy texnologiyalar asosida yaratilishi va yuqori yuklamalarga chidamli boʻlishi lozim.
4.1.4. Maʼlumotlar bazasi ishonchli ishlashi va buzilmasligi kerak.

6. ATni nazorat qilish va qabul qilish tartibi
6.1. Tizim sinovdan oʻtkaziladi va dalolatnoma tuziladi.
(Izoh: GOST 34.603 talablari va qabul komissiyasi tarkibi yoʻq).

7. ATni ishga tushirishga tayyorlash talablari
7.1. Xodimlar tizimda ishlash boʻyicha yoʻriqnoma bilan tanishishi zarur.

Foydalanish ssenariysi:
Foydalanuvchi hisobotni yuklaydi. Tizim faylni tekshiradi va saqlaydi.`

// GetDemoBenchmarkReport returns a pre-calculated benchmark audit report for offline demonstrations
func GetDemoBenchmarkReport() domain.AuditReport {
	return domain.AuditReport{
		AuditID:         "demo-uzdst-1987-001",
		AnalyzedAt:      time.Now(),
		DocumentName:    "DXM-Stat AT Texnik Topshirigʻi (Namunaviy loyiha)",
		ComplianceScore: 62.5,
		Status:          domain.StatusRejected,
		Metrics: domain.MetricScore{
			StructureCompleteness: 62.5, // 5 out of 8 sections present, 2 completely missing, 1 malformed
			ScenariosCompliance:   35.0, // Scenarios lack IDs, SLA latency, error limits
			MeasurabilityScore:    45.0, // Multiple subjective phrases without SLA
		},
		MissingSections: []string{
			"5. ATni yaratish ishlarining tarkibi va mazmuni (Oʻz DSt 1986:2018 bosqichlari va Oʻz DSt 1985 komplektligi umuman kiritilmagan)",
			"8. Hujjatlashtirishga qoʻyiladigan talablar (Oʻz DSt 1985:2018 boʻyicha majburiy hujjatlar roʻyxati mavjud emas)",
		},
		Findings: []domain.SectionFinding{
			{
				ID:          "FIND-01",
				Severity:    domain.SeverityCritical,
				Section:     "1. Umumiy maʼlumotlar",
				RawFragment: "Tizimning toʻliq nomi: Davlat xizmatlarini monitoring qilish... Qisqartirilgan nomi: 'DXM-Stat'.",
				Violation:   "Oʻz DSt 1985 7.4-7.5 bandlari va Oʻz DSt 1987:2018 6.2-bandi talabi buzilgan: Tashkilotning 8 belgili OKPO kodi va tizimning 15 belgigacha rasmiy shifri koʻrsatilmagan.",
				Recommendation: "1-boʻlimga tashkilotning OKPO kodini hamda Oʻz DSt 1985 boʻyicha tizim belgilanishini kiriting: 'OKPO: 12345678, AT shifri: 12345678.198701.001'.",
			},
			{
				ID:          "FIND-02",
				Severity:    domain.SeverityCritical,
				Section:     "2. AT vazifasi va yaratish maqsadlari",
				RawFragment: "Tizimning asosiy maqsadi — aholiga xizmat koʻrsatish sifatini oshirish va hisobotlarni tezlashtirish.",
				Violation:   "Oʻz DSt 1987:2018 6.3.2-bandi talabi buzilgan: erishilishi kutilayotgan texnik-iqtisodiy samaradorlikning oʻlchanadigan miqdoriy va sifat koʻrsatkichlari keltirilmagan.",
				Recommendation: "Talabni oʻlchanadigan metrikalar bilan toʻldiring: 'Hisobotlarni shakllantirish vaqtini amaldagi 4 soatdan 15 minutgacha qisqartirish, hujjat aylanishidagi xatoliklar ulushini kamida 40% ga kamaytirish'.",
			},
			{
				ID:          "FIND-03",
				Severity:    domain.SeverityCritical,
				Section:     "4.1. Umumiy talablar",
				RawFragment: "Tizim foydalanuvchilar soni koʻpayganda ham sekinlashmasdan tez ishlashi kerak.",
				Violation:   "Oʻz DSt 1987:2018 6.5.1.4-bandi buzilgan: noaniq, subyektiv va oʻlchanmaydigan 'tez ishlashi kerak' iborasi ishlatilgan.",
				Recommendation: "Talabni SLA parametrlariga oʻzgartiring: 'Tizim bir vaqtning oʻzida 500 ta faol sessiyada ishlaganda 95% soʻrovlar uchun sahifa yuklanish vaqti (latency) 2.0 sekunddan oshmasligi lozim'.",
			},
			{
				ID:          "FIND-04",
				Severity:    domain.SeverityWarning,
				Section:     "4.1. Umumiy talablar",
				RawFragment: "Foydalanuvchilar uchun juda qulay dizayn boʻlishi lozim.",
				Violation:   "Oʻz DSt 1987:2018 6.5.2-bandi: ergonomika va texnik estetika talablari noaniq baholovchi sifatlar bilan berilgan.",
				Recommendation: "Quyidagicha tahrirlansin: 'Foydalanuvchi interfeysi davlat tili qoidalariga, WCAG 2.1 AA standarti talablariga va tasdiqlangan yagona idoraviy UI-Kit komponentlariga muvofiq boʻlishi lozim'.",
			},
			{
				ID:          "FIND-05",
				Severity:    domain.SeverityCritical,
				Section:     "4.1. Umumiy talablar",
				RawFragment: "Tizim zamonaviy texnologiyalar asosida yaratilishi va yuqori yuklamalarga chidamli boʻlishi lozim.",
				Violation:   "Oʻz DSt 1987:2018 6.5.1-bandi: subyektiv 'zamonaviy texnologiyalar' va 'yuqori yuklamalar' talablari davlat ekspertizasida qabul qilinmaydi.",
				Recommendation: "Aniq texnik koʻrsatkichlarni belgilang: 'Tizim kamida 200 RPS (Request Per Second) doimiy yuklamaga bardosh berishi, xatoliklar darajasi (HTTP 5xx) 0.05% dan oshmasligi shart'.",
			},
			{
				ID:          "FIND-06",
				Severity:    domain.SeverityCritical,
				Section:     "Foydalanish ssenariysi",
				RawFragment: "Foydalanuvchi hisobotni yuklaydi. Tizim faylni tekshiradi va saqlaydi.",
				Violation:   "Oʻz DSt 1987:2018 B ilovasi talabi buzilgan: identifikator (U1, C1), qatʼiy vaqt reglamenti (Vaqt reglamenti) va muqobil yoʻnalishlar koʻrsatilmagan.",
				Recommendation: "B ilovasiga muvofiq rasmiylashtiring: Ssenariy kodi: U1_REPORT_UPLOAD; Vaqt reglamenti: faylni tekshirish va tasdiqlash maksimal 5 soniya; Istisnolar: Notoʻgʻri formatda 400 xatolik qaytarish.",
			},
			{
				ID:          "FIND-07",
				Severity:    domain.SeverityWarning,
				Section:     "6. ATni nazorat qilish va qabul qilish tartibi",
				RawFragment: "Tizim sinovdan oʻtkaziladi va dalolatnoma tuziladi.",
				Violation:   "Oʻz DSt 1987:2018 6.7-bandi va GOST 34.603 talabi: sinov bosqichlari (dastlabki, tajriba-sanoat, qabul) va qabul hayʼati tarkibi belgilanmagan.",
				Recommendation: "GOST 34.603 boʻyicha 3 ta majburiy sinov bosqichi (Dastlabki sinovlar, Tajriba-sanoat ekspluatatsiyasi, Qabul sinovlari) hamda idoralararo qabul komissiyasi tasdiqlash tartibini kiriting.",
			},
		},
		ExecutiveSummary: "Taqdim etilgan Texnik topshiriq loyihasi Oʻz DSt 1987:2018 talablariga toʻliq mos kelmaydi (Muvofiqlik indeksi: 62.5%). Hujjatda 2 ta majburiy boʻlim (5 va 8-boʻlimlar) toʻliq tushirib qoldirilgan, tizim belgilanishida Oʻz DSt 1985 boʻyicha tashkilot OKPO kodi va shifr mavjud emas, talablarda esa tekshirib boʻlmaydigan subyektiv jumlalar qoʻllanilgan. Koʻrsatilgan 7 ta jiddiy kamchilik toʻliq bartaraf etilmaguncha hujjatni davlat ekspertizasidan oʻtkazish RAD ETILADI.",
		IsDemoMode:       true,
	}
}

// DemoPresetTextRU provides a realistic flawed governmental IT specification fragment in Russian
const DemoPresetTextRU = `МИНИСТЕРСТВО ЦИФРОВЫХ ТЕХНОЛОГИЙ РЕСПУБЛИКИ УЗБЕКИСТАН
ПРОЕКТ ТЕХНИЧЕСКОГО ЗАДАНИЯ

1. Общие сведения
1.1. Полное наименование системы: Автоматизированная информационная система мониторинга государственных услуг.
1.2. Краткое наименование: "АИС Мониторинг".
1.3. Заказчик: Руководство Министерства.
1.4. Разработчик: Победитель конкурсных торгов.
1.5. Сроки выполнения работ: до конца 2026 года.
(Примечание: Код ОКПО организации и обозначение системы по Oʻz DSt 1985 не указаны).

2. Назначение и цели создания системы
2.1. Основная цель системы — повышение качества обслуживания населения и ускорение формирования отчетности.
2.2. Система должна улучшить межведомственное взаимодействие и сократить бумажный документооборот.

3. Характеристика объекта автоматизации
3.1. Объект — Центральный аппарат и 14 региональных подразделений.
3.2. Количество пользователей — около 500 сотрудников.

4. Требования к системе
4.1. Общие требования:
4.1.1. Система должна работать быстро при увеличении количества пользователей без задержек.
4.1.2. Должен быть предусмотрен удобный и интуитивно понятный интерфейс.
4.1.3. Система должна создаваться на основе современных технологий и обладать высокой устойчивостью к нагрузкам.
4.1.4. База данных должна работать надежно и без сбоев.

6. Порядок контроля и приемки системы
6.1. Система проходит испытания с оформлением совместного акта.
(Примечание: Требования ГОСТ 34.603 и состав комиссии отсутствуют).

7. Требования к подготовке объекта к вводу в действие
7.1. Персонал должен изучить инструкцию по эксплуатации системы.

Сценарий использования:
Пользователь загружает отчет. Система проверяет файл и сохраняет его в базу.`

// GetDemoBenchmarkReportRU returns benchmark report for Russian specifications
func GetDemoBenchmarkReportRU() domain.AuditReport {
	return domain.AuditReport{
		AuditID:         "demo-ru-1987-001",
		AnalyzedAt:      time.Now(),
		DocumentName:    "АИС Мониторинг ТЗ (Проект на русском языке)",
		ComplianceScore: 61.0,
		Status:          domain.StatusRejected,
		Metrics: domain.MetricScore{
			StructureCompleteness: 62.5,
			ScenariosCompliance:   35.0,
			MeasurabilityScore:    40.0,
		},
		MissingSections: []string{
			"5. Состав и содержание работ по созданию системы (этапы по Oʻz DSt 1986:2018 и ГОСТ 34.601 полностью отсутствуют)",
			"8. Требования к документированию (комплектность по Oʻz DSt 1985:2018 не определена)",
		},
		Findings: []domain.SectionFinding{
			{
				ID:          "FIND-RU-01",
				Severity:    domain.SeverityCritical,
				Section:     "1. Общие сведения",
				RawFragment: "Полное наименование: Автоматизированная информационная система... Краткое: 'АИС Мониторинг'.",
				Violation:   "Нарушены требования Oʻz DSt 1985 (п. 7.4-7.5) и Oʻz DSt 1987:2018 (п. 6.2): отсутствуют 8-значный код организации по ОКПО и обозначение системы.",
				Recommendation: "В раздел 1 включить реквизиты: 'Код ОКПО Заказчика: 12345678, Обозначение системы по Oʻz DSt 1985: 12345678.198701.001'.",
			},
			{
				ID:          "FIND-RU-02",
				Severity:    domain.SeverityCritical,
				Section:     "2. Назначение и цели создания системы",
				RawFragment: "Основная цель — повышение качества обслуживания населения и ускорение формирования отчетности.",
				Violation:   "Нарушен п. 6.3.2 Oʻz DSt 1987:2018: не приведены измеримые количественные критерии научно-технической и экономической эффективности.",
				Recommendation: "Указать измеримые критерии: 'Сокращение времени формирования сводной отчетности с 4 часов до 15 минут, снижение доли ошибок ручного ввода данных не менее чем на 40%'.",
			},
			{
				ID:          "FIND-RU-03",
				Severity:    domain.SeverityCritical,
				Section:     "4.1. Общие требования",
				RawFragment: "Система должна работать быстро при увеличении количества пользователей без задержек.",
				Violation:   "Нарушен п. 6.5.1.4 Oʻz DSt 1987:2018: использована неизмеримая субъективная формулировка 'работать быстро'.",
				Recommendation: "Заменить на SLA: 'Время отклика интерфейса системы при 500 одновременно работающих пользователях не должно превышать 2.0 секунд для 95% транзакций'.",
			},
			{
				ID:          "FIND-RU-04",
				Severity:    domain.SeverityWarning,
				Section:     "4.1. Общие требования",
				RawFragment: "Должен быть предусмотрен удобный и интуитивно понятный интерфейс.",
				Violation:   "Нарушен п. 6.5.2 Oʻz DSt 1987:2018: требования эргономики заданы качественными оценочными суждениями.",
				Recommendation: "Сформулировать требования: 'Пользовательский интерфейс должен соответствовать стандарту доступности WCAG 2.1 уровня AA и единым государственным рекомендациям UI-Kit'.",
			},
			{
				ID:          "FIND-RU-05",
				Severity:    domain.SeverityCritical,
				Section:     "4.1. Общие требования",
				RawFragment: "Система должна создаваться на основе современных технологий и обладать высокой устойчивостью к нагрузкам.",
				Violation:   "Нарушен п. 6.5.1 Oʻz DSt 1987:2018: запрещено использование размытых фраз без параметров пропускной способности.",
				Recommendation: "Задать нагрузочные параметры: 'Архитектура должна выдерживать нагрузку не менее 200 RPS при коэффициенте ошибок HTTP 5xx не более 0.05%'.",
			},
			{
				ID:          "FIND-RU-06",
				Severity:    domain.SeverityCritical,
				Section:     "Сценарий использования",
				RawFragment: "Пользователь загружает отчет. Система проверяет файл и сохраняет его в базу.",
				Violation:   "Нарушено Приложение B к Oʻz DSt 1987:2018: отсутствуют идентификатор прецедента (U1, C1), регламент времени операции и альтернативные сценарии.",
				Recommendation: "Оформить сценарий по Приложению B: Код: U1_REPORT_UPLOAD, Регламент времени: не более 5 секунд на валидацию и сохранение, Альтернатива: возврат ошибки 400 при неверном формате.",
			},
		},
		ExecutiveSummary: "Представленный проект ТЗ не соответствует государственному стандарту Oʻz DSt 1987:2018 (Индекс соответствия: 61.0%). В документе полностью опущены 2 обязательных раздела (5 и 8), отсутствуют шифры по Oʻz DSt 1985, а функциональные требования содержат неизмеримые субъективные формулировки. Также согласно Закону РУз 'О государственном языке' итоговое ТЗ для госэкспертизы должно быть оформлено на государственном (узбекском) языке. До устранения нарушений ТЗ к государственной экспертизе РЕКОМЕНДОВАНО К ОТКЛОНЕНИЮ (RAD ETILDI).",
		IsDemoMode:       true,
	}
}
