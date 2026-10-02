package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hextec/spec-auditor/internal/domain"
)

// PipelineAuditor executes multi-stage modular auditing with deterministic slicing and rules
type PipelineAuditor struct{}

// NewPipelineAuditor creates a new instance of PipelineAuditor
func NewPipelineAuditor() *PipelineAuditor {
	return &PipelineAuditor{}
}

var (
	okpoRegex         = regexp.MustCompile(`(?i)(?:OKPO|ОКПО)\s*[:\.\-]?\s*\b\d{8}\b`)
	systemCodeRegex   = regexp.MustCompile(`\b\d{8}\.\d{3,6}\.\d{3}\b`)
	timeRegRegex      = regexp.MustCompile(`(?i)(?:vaqt\s+reglamenti|регламент\s+времени|\b\d+\s*(?:soniya|sekund|минут|секунд|s|ms)\b)`)
	useCaseIDRegex    = regexp.MustCompile(`\b[UCED]{1,2}\d{1,3}\b`)
	gost34603Regex    = regexp.MustCompile(`(?i)(?:GOST\s*34\.603|ГОСТ\s*34\.603)`)
	security2814Regex = regexp.MustCompile(`(?i)(?:O[ʻʼ']?z\s*DSt\s*2814|2814|xavfsizlik|безопасност)`)
)

// AuditPipeline runs automated section slicing, targeted rule checking, and aggregation
func (p *PipelineAuditor) AuditPipeline(docName, content string) *domain.AuditReport {
	sliced := SliceTechnicalSpecification(content)
	isRussian := isRussianText(content)

	var findings []domain.SectionFinding
	findingIdx := 1

	addFinding := func(severity domain.Severity, section, fragment, violation, recommendation string) {
		findings = append(findings, domain.SectionFinding{
			ID:             fmt.Sprintf("FIND-%02d", findingIdx),
			Severity:       severity,
			Section:        section,
			RawFragment:    fragment,
			Violation:      violation,
			Recommendation: recommendation,
		})
		findingIdx++
	}

	// 1. SECTION 1 AUDIT: Organization OKPO and Oʻz DSt 1985 designation code
	sec1 := findSection(sliced.Sections, 1)
	if sec1 != nil {
		if !okpoRegex.MatchString(sec1.Content) && !strings.Contains(strings.ToUpper(sec1.Content), "OKPO") {
			if isRussian {
				addFinding(domain.SeverityCritical, "1. Общие сведения",
					extractSnippet(sec1.Content, 100),
					"Нарушены требования Oʻz DSt 1985 (п. 7.4-7.5) и Oʻz DSt 1987:2018 (п. 6.2): отсутствует 8-значный код ОКПО организации.",
					"В раздел 1 включить реквизиты Заказчика с указанием 8-значного кода ОКПО (например: 'Код ОКПО Заказчика: 12345678').")
			} else {
				addFinding(domain.SeverityCritical, "1. Umumiy maʼlumotlar",
					extractSnippet(sec1.Content, 100),
					"Oʻz DSt 1985 7.4-7.5 bandlari va Oʻz DSt 1987:2018 6.2-bandi buzilgan: Tashkilotning 8 belgili OKPO kodi koʻrsatilmagan.",
					"1-boʻlimga tashkilotning 8 belgili OKPO kodini kiriting (masalan: 'Buyurtmachi OKPO kodi: 12345678').")
			}
		}
		if !systemCodeRegex.MatchString(sec1.Content) && !strings.Contains(sec1.Content, "1985") {
			if isRussian {
				addFinding(domain.SeverityWarning, "1. Общие сведения",
					extractSnippet(sec1.Content, 100),
					"Oʻz DSt 1985 (п. 7.5): отсутствует официальное 15-значное обозначение системы.",
					"Включить обозначение системы по Oʻz DSt 1985: 'Обозначение АИС: 12345678.198701.001'.")
			} else {
				addFinding(domain.SeverityWarning, "1. Umumiy maʼlumotlar",
					extractSnippet(sec1.Content, 100),
					"Oʻz DSt 1985 7.5-bandi: Tizimning 15 belgigacha rasmiy shifri koʻrsatilmagan.",
					"1-boʻlimga Oʻz DSt 1985 boʻyicha tizim belgilanishini kiriting: 'AT shifri: 12345678.198701.001'.")
			}
		}
	}

	// 2. SECTION 2 AUDIT: Measurable technical-economic efficiency criteria
	sec2 := findSection(sliced.Sections, 2)
	if sec2 != nil {
		hasNumbers := regexp.MustCompile(`\b\d+[\.,]?\d*\s*(?:%|soat|minut|sekund|barobar|сокращение|процент|час|минут)\b`).MatchString(sec2.Content)
		if !hasNumbers {
			if isRussian {
				addFinding(domain.SeverityCritical, "2. Назначение и цели создания системы",
					extractSnippet(sec2.Content, 120),
					"Oʻz DSt 1987:2018 п. 6.3.2: не приведены измеримые количественные критерии научно-технической и экономической эффективности.",
					"Дополнить раздел измеримыми показателями: 'Сокращение времени обработки заявок с X до Y минут, снижение ручных ошибок не менее чем на Z%'.")
			} else {
				addFinding(domain.SeverityCritical, "2. AT vazifasi va yaratish maqsadlari",
					extractSnippet(sec2.Content, 120),
					"Oʻz DSt 1987:2018 6.3.2-bandi: erishilishi kutilayotgan texnik-iqtisodiy samaradorlikning oʻlchanadigan miqdoriy koʻrsatkichlari keltirilmagan.",
					"Maqsadlarni oʻlchanadigan metrikalar bilan toʻldiring: 'Hisobot shakllantirish vaqtini X soatdan Y minutgacha qisqartirish, operatsion xatoliklarni Z% ga kamaytirish'.")
			}
		}
	}

	// 3. SECTION 4 & GLOBAL: Banned subjective non-measurable requirements
	checkLower := strings.ToLower(content)
	for _, phrase := range domain.BannedSubjectivePhrases {
		pLower := strings.ToLower(phrase)
		idx := strings.Index(checkLower, pLower)
		if idx != -1 {
			start := idx - 20
			if start < 0 {
				start = 0
			}
			end := idx + len(phrase) + 30
			if end > len(content) {
				end = len(content)
			}
			rawSnippet := strings.TrimSpace(content[start:end])

			if isRussian {
				addFinding(domain.SeverityCritical, "4.1. Общие требования",
					rawSnippet,
					fmt.Sprintf("Oʻz DSt 1987:2018 п. 6.5.1.4: использована недопустимая качественная субъективная формулировка '%s'.", phrase),
					"Заменить субъективное требование на количественные параметры SLA: указать точное время отклика (latency не более X сек при Y RPS), пиковую нагрузку и допустимый % ошибок.")
			} else {
				addFinding(domain.SeverityCritical, "4.1. Umumiy talablar",
					rawSnippet,
					fmt.Sprintf("Oʻz DSt 1987:2018 6.5.1.4-bandi buzilgan: oʻlchanmaydigan subyektiv '%s' iborasi ishlatilgan.", phrase),
					"Talabni SLA parametrlariga almashtiring: sahifa yuklanish vaqti latency < 2.0 soniya, 95%% soʻrovlar uchun barqaror ishlash va RPS koʻrsatkichini belgilang.")
			}
		}
	}

	// 4. SECTION 6 AUDIT: Testing standards (GOST 34.603) and commission
	sec6 := findSection(sliced.Sections, 6)
	if sec6 != nil {
		if !gost34603Regex.MatchString(sec6.Content) {
			if isRussian {
				addFinding(domain.SeverityWarning, "6. Порядок контроля и приемки системы",
					extractSnippet(sec6.Content, 100),
					"Oʻz DSt 1987:2018 п. 6.7: отсутствуют обязательные стадии испытаний по ГОСТ 34.603 (предварительные, опытная эксплуатация, приемочные).",
					"Указать испытания строго по ГОСТ 34.603: предварительные испытания, опытная эксплуатация и приемочные испытания приемочной комиссией.")
			} else {
				addFinding(domain.SeverityWarning, "6. ATni nazorat qilish va qabul qilish tartibi",
					extractSnippet(sec6.Content, 100),
					"Oʻz DSt 1987:2018 6.7-bandi: GOST 34.603 boʻyicha majburiy 3 ta sinov bosqichi (dastlabki, tajriba-sanoat, qabul) va qabul hayʼati koʻrsatilmagan.",
					"GOST 34.603 boʻyicha sinov turlari va idoralararo qabul qilish komissiyasi tasdiqlash tartibini belgilang.")
			}
		}
	}

	// 5. USE CASE SCENARIOS (B ilovasi)
	if sliced.HasScenarios {
		scenarioSnippet := extractSnippet(sliced.ScenariosText, 150)
		if scenarioSnippet == "" {
			scenarioSnippet = "Foydalanish ssenariysi..."
		}
		if !useCaseIDRegex.MatchString(sliced.ScenariosText) || !timeRegRegex.MatchString(sliced.ScenariosText) {
			if isRussian {
				addFinding(domain.SeverityCritical, "Сценарий использования (Приложение B)",
					scenarioSnippet,
					"Нарушены требования Приложения B к Oʻz DSt 1987:2018: отсутствуют код прецедента (U1, C1) и регламент времени операции.",
					"Оформить сценарий по Приложению B: указать уникальный код операции (U1_...), жесткий временной регламент в секундах и альтернативные ветки.")
			} else {
				addFinding(domain.SeverityCritical, "Foydalanish ssenariysi (B ilovasi)",
					scenarioSnippet,
					"Oʻz DSt 1987:2018 B ilovasi talabi buzilgan: ssenariy identifikatori (U1, C1) va qatʼiy vaqt reglamenti (Vaqt reglamenti) mavjud emas.",
					"B ilovasi talabiga muvofiq rasmiylashtiring: unikal ssenariy kodi (masalan: U1_ACTION), aniq vaqt reglamenti (maksimal 3-5 soniya) va muqobil yoʻllarni koʻrsating.")
			}
		}
	} else {
		if isRussian {
			addFinding(domain.SeverityCritical, "Приложение Б: Сценарии использования",
				"",
				"Нарушены требования Oʻz DSt 1987:2018 (Приложение Б): в документе полностью отсутствуют сценарии использования системы.",
				"Дополнить ТЗ Приложением Б с детальным описанием сценариев пользователей (ID прецедента, жесткий регламент времени операции, альтернативные потоки).")
		} else {
			addFinding(domain.SeverityCritical, "B ilovasi: Foydalanish ssenariylari",
				"",
				"Oʻz DSt 1987:2018 B ilovasi talabi buzilgan: hujjatda tizimdan foydalanish ssenariylari umuman mavjud emas.",
				"Texnik topshiriqqa B ilovasini kiriting: unikal ssenariy kodi (masalan: U1_ACTION), aniq vaqt reglamenti va foydalanuvchi qadamlari.")
		}
	}

	// 6. SYNONYM HEADERS CHECK (Non-standard naming)
	for _, sec := range sliced.Sections {
		if sec.IsSynonym {
			for _, p := range sectionPatterns {
				if p.Number == sec.Number {
					stdTitle := p.TitleUz
					if isRussian {
						stdTitle = p.TitleRu
					}
					if isRussian {
						addFinding(domain.SeverityWarning, stdTitle,
							extractSnippet(sec.Title, 80),
							fmt.Sprintf("Oʻz DSt 1987:2018 (Приложение А): заголовок раздела («%s») не стандартизирован. Государственный стандарт требует точного наименования раздела.", sec.Title),
							fmt.Sprintf("Переименовать заголовок блока в «%s» в строгом соответствии со стандартом Oʻz DSt 1987:2018.", stdTitle))
					} else {
						addFinding(domain.SeverityWarning, stdTitle,
							extractSnippet(sec.Title, 80),
							fmt.Sprintf("Oʻz DSt 1987:2018 (A ilovasi): boʻlim sarlavhasi («%s») standartlashtirilmagan. Standart boʻyicha rasmiy nom talab etiladi.", sec.Title),
							fmt.Sprintf("Boʻlim sarlavhasini Oʻz DSt 1987:2018 ga muvofiq «%s» deb rasmiylashtiring.", stdTitle))
					}
				}
			}
		}
	}

	// 7. MISSING MANDATORY SECTIONS AUDIT (Create explicit CRITICAL findings)
	for _, p := range sectionPatterns {
		if findSection(sliced.Sections, p.Number) == nil {
			stdTitle := p.TitleUz
			if isRussian {
				stdTitle = p.TitleRu
			}
			if isRussian {
				addFinding(domain.SeverityCritical, stdTitle,
					"",
					fmt.Sprintf("Грубое нарушение Oʻz DSt 1987:2018 (Приложение А): обязательный раздел «%s» полностью отсутствует в структуре документа.", stdTitle),
					fmt.Sprintf("Включить в ТЗ обязательный раздел «%s» в соответствии с требованиями государственного стандарта.", stdTitle))
			} else {
				addFinding(domain.SeverityCritical, stdTitle,
					"",
					fmt.Sprintf("Oʻz DSt 1987:2018 (A ilovasi) qoʻpol ravishda buzilgan: majburiy «%s» boʻlimi hujjatda umuman mavjud emas.", stdTitle),
					fmt.Sprintf("Texnik topshiriqqa Oʻz DSt 1987:2018 talablariga muvofiq «%s» majburiy boʻlimini kiriting.", stdTitle))
			}
		}
	}

	// 8. SCORING AGGREGATOR
	totalMissing := len(sliced.MissingSections)
	detectedCount := len(sliced.Sections)
	criticalCount := 0
	warningCount := 0
	sec4Violations := 0

	for _, f := range findings {
		if f.Severity == domain.SeverityCritical {
			criticalCount++
		} else if f.Severity == domain.SeverityWarning {
			warningCount++
		}
		if strings.HasPrefix(f.Section, "4.") {
			sec4Violations++
		}
	}

	// Structure completeness: 100 * (detected / 8)
	structScore := float64(detectedCount) / 8.0 * 100.0
	if structScore > 100.0 {
		structScore = 100.0
	}

	// Scenarios score
	scenariosScore := 100.0
	if !sliced.HasScenarios {
		scenariosScore = 0.0
	} else if !useCaseIDRegex.MatchString(sliced.ScenariosText) || !timeRegRegex.MatchString(sliced.ScenariosText) {
		scenariosScore = 45.0
	}

	// Measurability score
	var measurabilityScore float64
	sec4 := findSection(sliced.Sections, 4)
	if sec4 == nil {
		measurabilityScore = 0.0 // Requirements are not formulated
	} else {
		measurabilityScore = 100.0 - float64(sec4Violations)*15.0
		if measurabilityScore < 15.0 {
			measurabilityScore = 15.0
		}
	}

	// Overall compliance score: weighted
	overallScore := (structScore * 0.45) + (scenariosScore * 0.25) + (measurabilityScore * 0.30)
	overallScore = float64(int(overallScore*10)) / 10.0 // Round to 1 decimal place

	status := domain.CalculateStatus(overallScore)
	if totalMissing > 0 {
		status = domain.StatusRejected
	}

	// Build dynamic executive summary
	var execSummary string
	if isRussian {
		if status == domain.StatusApproved {
			execSummary = fmt.Sprintf("Представленный проект ТЗ соответствует государственному стандарту Oʻz DSt 1987:2018 (Индекс: %.1f%%). Структура полная, требования измеримы. Документ РЕКОМЕНДОВАН к государственной экспертизе.", overallScore)
		} else if status == domain.StatusConditionalPass {
			execSummary = fmt.Sprintf("Проект ТЗ частично соответствует Oʻz DSt 1987:2018 (Индекс: %.1f%%). Выявлено %d замечаний, требующих доработки до проведения государственной экспертизы.", overallScore, len(findings))
		} else {
			if totalMissing > 0 {
				execSummary = fmt.Sprintf("Проект ТЗ не соответствует стандарту Oʻz DSt 1987:2018 (Индекс: %.1f%%). В документе отсутствует %d обязательных разделов и выявлено %d критических несоответствий. В текущей редакции к государственной экспертизе ОТКЛОНЕНО (RAD ETILDI).", overallScore, totalMissing, criticalCount)
			} else {
				execSummary = fmt.Sprintf("Проект ТЗ не соответствует стандарту Oʻz DSt 1987:2018 (Индекс: %.1f%%). Выявлено %d критических несоответствий. В текущей редакции к государственной экспертизе ОТКЛОНЕНО (RAD ETILDI).", overallScore, criticalCount)
			}
		}
	} else {
		if status == domain.StatusApproved {
			execSummary = fmt.Sprintf("Taqdim etilgan Texnik topshiriq loyihasi Oʻz DSt 1987:2018 talablariga toʻliq javob beradi (Muvofiqlik indeksi: %.1f%%). Barcha 8 ta majburiy boʻlim toʻliq shakllantirilgan. Davlat ekspertizasiga TAVSIYA ETILADI.", overallScore)
		} else if status == domain.StatusConditionalPass {
			execSummary = fmt.Sprintf("Texnik topshiriq loyihasi Oʻz DSt 1987:2018 standartiga qisman muvofiq (Muvofiqlik indeksi: %.1f%%). Qayd etilgan %d ta kamchilik bartaraf etilgach, takroriy ekspertizaga kiritilishi mumkin.", overallScore, len(findings))
		} else {
			if totalMissing > 0 {
				execSummary = fmt.Sprintf("Texnik topshiriq loyihasi Oʻz DSt 1987:2018 davlat standarti talablariga mos kelmaydi (Muvofiqlik indeksi: %.1f%%). Hujjatda %d ta majburiy boʻlim tushirib qoldirilgan hamda %d ta qatʼiy xatolik aniqlangan. Koʻrsatilgan kamchiliklar bartaraf etilmaguncha davlat ekspertizasidan oʻtkazish RAD ETILADI.", overallScore, totalMissing, criticalCount)
			} else {
				execSummary = fmt.Sprintf("Texnik topshiriq loyihasi Oʻz DSt 1987:2018 davlat standarti talablariga mos kelmaydi (Muvofiqlik indeksi: %.1f%%). Hujjatda %d ta qatʼiy xatolik aniqlangan. Davlat ekspertizasidan oʻtkazish RAD ETILADI.", overallScore, criticalCount)
			}
		}
	}

	// Build SectionMap for visual presentation of auto-slicing
	var sectionMap []domain.SectionOverview
	for _, p := range sectionPatterns {
		sec := findSection(sliced.Sections, p.Number)
		issuesForSec := 0
		secPrefix := fmt.Sprintf("%d.", p.Number)
		for _, f := range findings {
			if strings.HasPrefix(f.Section, secPrefix) {
				issuesForSec++
			}
		}

		title := p.TitleUz
		if isRussian {
			title = p.TitleRu
		}

		if sec != nil {
			snippet := extractSnippet(sec.Content, 80)
			sectionMap = append(sectionMap, domain.SectionOverview{
				Number:      p.Number,
				Title:       title,
				Detected:    true,
				StandardRef: p.Reference,
				Snippet:     snippet,
				CharCount:   len([]rune(sec.Content)),
				IssuesCount: issuesForSec,
			})
		} else {
			sectionMap = append(sectionMap, domain.SectionOverview{
				Number:      p.Number,
				Title:       title,
				Detected:    false,
				StandardRef: p.Reference,
				IssuesCount: issuesForSec,
			})
		}
	}

	// Appendix B: Foydalanish ssenariylari
	scenariosIssues := 0
	for _, f := range findings {
		lowerSec := strings.ToLower(f.Section)
		if strings.Contains(lowerSec, "b ilovasi") || strings.Contains(lowerSec, "приложение б") || strings.Contains(lowerSec, "сценари") || strings.Contains(lowerSec, "ssenariy") {
			scenariosIssues++
		}
	}
	appBTitle := "B ilovasi: Foydalanish ssenariylari"
	if isRussian {
		appBTitle = "Приложение Б: Сценарии использования"
	}
	sectionMap = append(sectionMap, domain.SectionOverview{
		Number:      9,
		Title:       appBTitle,
		Detected:    sliced.HasScenarios,
		StandardRef: "Oʻz DSt 1987:2018 B ilovasi",
		Snippet:     extractSnippet(sliced.ScenariosText, 80),
		CharCount:   len([]rune(sliced.ScenariosText)),
		IssuesCount: scenariosIssues,
	})

	return &domain.AuditReport{
		AuditID:          fmt.Sprintf("pipe-%d", time.Now().Unix()),
		AnalyzedAt:       time.Now(),
		DocumentName:     docName,
		ComplianceScore:  overallScore,
		Status:           status,
		Metrics: domain.MetricScore{
			StructureCompleteness: structScore,
			ScenariosCompliance:   scenariosScore,
			MeasurabilityScore:    measurabilityScore,
		},
		MissingSections:  sliced.MissingSections,
		SectionMap:       sectionMap,
		Findings:         findings,
		ExecutiveSummary: execSummary,
		IsDemoMode:       false,
	}
}

func findSection(sections []ParsedSection, num int) *ParsedSection {
	for i := range sections {
		if sections[i].Number == num {
			return &sections[i]
		}
	}
	return nil
}

func extractSnippet(text string, maxLen int) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	runes := []rune(cleaned)
	if len(runes) <= maxLen {
		return string(runes)
	}
	return string(runes[:maxLen]) + "..."
}
