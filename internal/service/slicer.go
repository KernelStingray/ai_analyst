package service

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hextec/spec-auditor/internal/domain"
)

// ParsedSection represents an extracted section from the technical specification
type ParsedSection struct {
	Number      int    // 1 to 8, or 9 for use case scenarios
	Title       string // Original header title found in document
	Content     string // Text content belonging to this section
	StandardRef string // Clause reference from Oʻz DSt 1987:2018
	IsSynonym   bool   // True if matched via synonym rather than strict Oʻz DSt title
}

// SliceResult contains partitioned sections and deterministically detected omissions
type SliceResult struct {
	Sections        []ParsedSection
	MissingSections []string
	HasScenarios    bool
	ScenariosText   string
	IsUnstructured  bool
}

// SectionPattern matches section headers in both Uzbek and Russian
type SectionPattern struct {
	Number       int
	StrictRegex  *regexp.Regexp
	SynonymRegex *regexp.Regexp
	TitleUz      string
	TitleRu      string
	Reference    string
}

var sectionPatterns = []SectionPattern{
	{
		Number:       1,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:1[\.\)]|I[\.\)]|1-boʻlim|1-бўлим|Раздел\s*1)\s*[:\.\-]?\s*)(?:Umumiy\s+ma[ʼ'ʻ]lumotlar|Общие\s+сведения)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:1[\.\)]|I[\.\)]|1-boʻlim|1-бўлим|Раздел\s*1)\s*[:\.\-]?\s*)?(?:Umumiy\s+ma[ʼ'ʻ]lumotlar|Общие\s+сведения|Loyiha\s+haqida|О\s+проекте|Введение|Kirish|Паспорт\s+проекта|Преамбула)[^\n]*`),
		TitleUz:      "1. Umumiy maʼlumotlar",
		TitleRu:      "1. Общие сведения",
		Reference:    "Oʻz DSt 1987:2018 6.2-bandi, Oʻz DSt 1985 7.4-7.5",
	},
	{
		Number:       2,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:2[\.\)]|II[\.\)]|2-boʻlim|2-бўлим|Раздел\s*2)\s*[:\.\-]?\s*)(?:AT\s+vazifasi\s+va|Tizim\s+vazifasi|Назначение\s+и\s+цели)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:2[\.\)]|II[\.\)]|2-boʻlim|2-бўлим|Раздел\s*2)\s*[:\.\-]?\s*)?(?:AT\s+vazifasi\s+va|Tizim\s+vazifasi|Назначение\s+и\s+цели|Цели\s+и\s+задачи|Бизнес-цели|Maqsad\s+va\s+vazifalar|Loyiha\s+maqsadi)[^\n]*`),
		TitleUz:      "2. AT vazifasi va yaratish maqsadlari",
		TitleRu:      "2. Назначение и цели создания системы",
		Reference:    "Oʻz DSt 1987:2018 6.3-bandi (6.3.2-band)",
	},
	{
		Number:       3,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:3[\.\)]|III[\.\)]|3-boʻlim|3-бўлим|Раздел\s*3)\s*[:\.\-]?\s*)(?:Axborotlashtirish\s+obyektining|Obyekt\s+tavsifi|Характеристика\s+объекта)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:3[\.\)]|III[\.\)]|3-boʻlim|3-бўлим|Раздел\s*3)\s*[:\.\-]?\s*)?(?:Axborotlashtirish\s+obyektining|Obyekt\s+tavsifi|Характеристика\s+объекта|Бизнес-модель|Контекст\s+внедрения|Проблематика\s+госсектора|Описание\s+объекта|Tizim\s+konteksti)[^\n]*`),
		TitleUz:      "3. Axborotlashtirish obyektining tavsifi",
		TitleRu:      "3. Характеристика объекта автоматизации",
		Reference:    "Oʻz DSt 1987:2018 6.4-bandi",
	},
	{
		Number:       4,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:4[\.\)]|IV[\.\)]|4-boʻlim|4-бўлим|Раздел\s*4)\s*[:\.\-]?\s*)(?:ATga\s+qo[ʻʼ']yiladigan\s+talablar|Tizimga\s+talablar|Требования\s+к\s+системе)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:4[\.\)]|IV[\.\)]|4-boʻlim|4-бўлим|Раздел\s*4)\s*[:\.\-]?\s*)?(?:ATga\s+qo[ʻʼ']yiladigan\s+talablar|Tizimga\s+talablar|Требования\s+к\s+системе|Архитектура|Инфраструктура|Технические\s+требования|Функциональные\s+требования|Стек\s+технологий|Tizim\s+arxitekturasi|Texnik\s+talablar)[^\n]*`),
		TitleUz:      "4. ATga qoʻyiladigan talablar",
		TitleRu:      "4. Требования к системе",
		Reference:    "Oʻz DSt 1987:2018 6.5-bandi, Oʻz DSt 2814",
	},
	{
		Number:       5,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:5[\.\)]|V[\.\)]|5-boʻlim|5-бўлим|Раздел\s*5)\s*[:\.\-]?\s*)(?:ATni\s+yaratish\s+ishlarining|Ishlarning\s+tarkibi|Состав\s+и\s+содержание\s+работ)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:5[\.\)]|V[\.\)]|5-boʻlim|5-бўлим|Раздел\s*5)\s*[:\.\-]?\s*)?(?:ATni\s+yaratish\s+ishlarining|Ishlarning\s+tarkibi|Состав\s+и\s+содержание\s+работ|Этапы\s+разработки|План\s+реализации|Этапы\s+создания|Ishlar\s+bosqichlari|План\s+работ|Этап\s*1\b)[^\n]*`),
		TitleUz:      "5. ATni yaratish ishlarining tarkibi va mazmuni",
		TitleRu:      "5. Состав и содержание работ по созданию системы",
		Reference:    "Oʻz DSt 1987:2018 6.6-bandi, Oʻz DSt 1986:2018",
	},
	{
		Number:       6,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:6[\.\)]|VI[\.\)]|6-boʻlim|6-бўлим|Раздел\s*6)\s*[:\.\-]?\s*)(?:ATni\s+nazorat\s+qilish|Nazorat\s+va\s+qabul|Порядок\s+контроля\s+и\s+приемки)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:6[\.\)]|VI[\.\)]|6-boʻlim|6-бўлим|Раздел\s*6)\s*[:\.\-]?\s*)?(?:ATni\s+nazorat\s+qilish|Nazorat\s+va\s+qabul|Порядок\s+контроля\s+и\s+приемки|Тестирование|Порядок\s+приемки|Приемка\s+системы|Sinovlar\s+va\s+qabul)[^\n]*`),
		TitleUz:      "6. ATni nazorat qilish va qabul qilish tartibi",
		TitleRu:      "6. Порядок контроля и приемки системы",
		Reference:    "Oʻz DSt 1987:2018 6.7-bandi, GOST 34.603",
	},
	{
		Number:       7,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:7[\.\)]|VII[\.\)]|7-boʻlim|7-бўлим|Раздел\s*7)\s*[:\.\-]?\s*)(?:ATni\s+ishga\s+tushirishga|Tayyorlash\s+talablari|Требования\s+к\s+подготовке\s+объекта)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:7[\.\)]|VII[\.\)]|7-boʻlim|7-бўлим|Раздел\s*7)\s*[:\.\-]?\s*)?(?:ATni\s+ishga\s+tushirishga|Tayyorlash\s+talablari|Требования\s+к\s+подготовке\s+объекта|Ввод\s+в\s+действие|Подготовка\s+к\s+вводу|Tayyorgarlik\s+talablari)[^\n]*`),
		TitleUz:      "7. ATni ishga tushirishga tayyorlash talablari",
		TitleRu:      "7. Требования к подготовке объекта к вводу в действие",
		Reference:    "Oʻz DSt 1987:2018 6.8-bandi",
	},
	{
		Number:       8,
		StrictRegex:  regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:8[\.\)]|VIII[\.\)]|8-boʻlim|8-бўлим|Раздел\s*8)\s*[:\.\-]?\s*)(?:Hujjatlashtirishga\s+qo[ʻʼ']yiladigan|Hujjatlarga\s+talablar|Требования\s+к\s+документированию)[^\n]*`),
		SynonymRegex: regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:(?:8[\.\)]|VIII[\.\)]|8-boʻlim|8-бўлим|Раздел\s*8)\s*[:\.\-]?\s*)?(?:Hujjatlashtirishga\s+qo[ʻʼ']yiladigan|Hujjatlarga\s+talablar|Требования\s+к\s+документированию|Требования\s+к\s+документации|Hujjatlashtirish)[^\n]*`),
		TitleUz:      "8. Hujjatlashtirishga qoʻyiladigan talablar",
		TitleRu:      "8. Требования к документированию",
		Reference:    "Oʻz DSt 1987:2018 6.9-bandi, Oʻz DSt 1985:2018",
	},
}

var scenariosRegex = regexp.MustCompile(`(?i)(?:^|\n)\s*(?:#+\s*)?(?:Foydalanish\s+ssenariy|Ssenariylar|Прецедент|Сценарий\s+использования|Сценарии\s+использования|Use[\s\-_]?cases?|Пользовательские\s+сценарии)[^\n]*`)

type sectionMatch struct {
	pattern   SectionPattern
	startIdx  int
	headerEnd int
	header    string
	isSynonym bool
}

// SliceTechnicalSpecification partitions the full text into sections and detects missing mandatory ones
func SliceTechnicalSpecification(fullText string) SliceResult {
	var matches []sectionMatch

	// Find occurrences of each mandatory section header
	for _, p := range sectionPatterns {
		// 1. Try strict regex first
		loc := p.StrictRegex.FindStringIndex(fullText)
		if loc != nil {
			matches = append(matches, sectionMatch{
				pattern:   p,
				startIdx:  loc[0],
				headerEnd: loc[1],
				header:    strings.TrimSpace(fullText[loc[0]:loc[1]]),
				isSynonym: false,
			})
			continue
		}

		// 2. Try synonym regex if strict wasn't found
		locSyn := p.SynonymRegex.FindStringIndex(fullText)
		if locSyn != nil {
			matches = append(matches, sectionMatch{
				pattern:   p,
				startIdx:  locSyn[0],
				headerEnd: locSyn[1],
				header:    strings.TrimSpace(fullText[locSyn[0]:locSyn[1]]),
				isSynonym: true,
			})
		}
	}

	// Sort matches by start position in text
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].startIdx < matches[j].startIdx
	})

	result := SliceResult{
		Sections:        make([]ParsedSection, 0),
		MissingSections: make([]string, 0),
	}

	foundNumbers := make(map[int]bool)

	if len(matches) < 2 {
		result.IsUnstructured = true
	}

	// Extract content between consecutive headers
	for i := 0; i < len(matches); i++ {
		current := matches[i]
		foundNumbers[current.pattern.Number] = true

		endIdx := len(fullText)
		if i+1 < len(matches) {
			endIdx = matches[i+1].startIdx
		}

		content := strings.TrimSpace(fullText[current.headerEnd:endIdx])

		result.Sections = append(result.Sections, ParsedSection{
			Number:      current.pattern.Number,
			Title:       current.header,
			Content:     content,
			StandardRef: current.pattern.Reference,
			IsSynonym:   current.isSynonym,
		})
	}

	// Check use-case scenarios (B ilovasi)
	scenariosLoc := scenariosRegex.FindStringIndex(fullText)
	if scenariosLoc != nil {
		result.HasScenarios = true
		result.ScenariosText = strings.TrimSpace(fullText[scenariosLoc[1]:])
	}

	// Detect which of the 8 mandatory sections were completely omitted
	isRussian := isRussianText(fullText)
	for _, sec := range domain.MandatorySections {
		num := 0
		fmt.Sscanf(sec.Number, "%d", &num)
		if !foundNumbers[num] {
			if isRussian && sec.TitleRu != "" {
				result.MissingSections = append(result.MissingSections,
					fmt.Sprintf("%s. %s (раздел полностью отсутствует)", sec.Number, sec.TitleRu))
			} else {
				result.MissingSections = append(result.MissingSections,
					fmt.Sprintf("%s. %s (Oʻz DSt 1987:2018 boʻyicha majburiy boʻlim kiritilmagan)", sec.Number, sec.TitleUz))
			}
		}
	}

	return result
}
