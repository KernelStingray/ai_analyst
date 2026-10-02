package domain

// StandardSection represents a mandatory section under Oʻz DSt 1987:2018 A ilovasi
type StandardSection struct {
	Number      string
	TitleUz     string
	TitleRu     string
	Description string
	Reference   string
}

// MandatorySections lists all 8 required sections according to Oʻz DSt 1987:2018 A ilovasi
var MandatorySections = []StandardSection{
	{
		Number:      "1",
		TitleUz:     "Umumiy maʼlumotlar",
		TitleRu:     "Общие сведения",
		Description: "Tizimning toʻliq va qisqartirilgan nomi, buyurtmachi va ishlab chiquvchi rekvizitlari, Oʻz DSt 1985 boʻyicha belgilanishi (OKPO 8 belgi, shifr 15 belgigacha), muddatlar.",
		Reference:   "Oʻz DSt 1987:2018 6.2-bandi, Oʻz DSt 1985 7.4-7.5 bandlari / ГОСТ 34.602-89",
	},
	{
		Number:      "2",
		TitleUz:     "AT vazifasi va yaratish maqsadlari",
		TitleRu:     "Назначение и цели создания системы",
		Description: "Tizimning maqsadi, erishilishi kutilayotgan texnik-iqtisodiy samaradorlikning oʻlchanadigan miqdoriy va sifat koʻrsatkichlari.",
		Reference:   "Oʻz DSt 1987:2018 6.3-bandi (xususan 6.3.2-band) / ГОСТ 34.602",
	},
	{
		Number:      "3",
		TitleUz:     "Axborotlashtirish obyektining tavsifi",
		TitleRu:     "Характеристика объекта автоматизации",
		Description: "Obyektning tuzilishi, faoliyat xususiyatlari, tizimning ishlash shartlari va rejimlari.",
		Reference:   "Oʻz DSt 1987:2018 6.4-bandi",
	},
	{
		Number:      "4",
		TitleUz:     "ATga qoʻyiladigan talablar",
		TitleRu:     "Требования к системе",
		Description: "Butun tizimga qoʻyiladigan talablar (ishonchlilik, xavfsizlik Oʻz DSt 2814, ergonomika), quyi tizimlar va funksiyalarga talablar.",
		Reference:   "Oʻz DSt 1987:2018 6.5-bandi",
	},
	{
		Number:      "5",
		TitleUz:     "ATni yaratish ishlarining tarkibi va mazmuni",
		TitleRu:     "Состав и содержание работ по созданию системы",
		Description: "Tizimni yaratish bosqichlari (Oʻz DSt 1986 / ГОСТ 34.601) va Oʻz DSt 1985 boʻyicha hujjatlar toʻplami.",
		Reference:   "Oʻz DSt 1987:2018 6.6-bandi, Oʻz DSt 1986:2018",
	},
	{
		Number:      "6",
		TitleUz:     "ATni nazorat qilish va qabul qilish tartibi",
		TitleRu:     "Порядок контроля и приемки системы",
		Description: "Sinov turlari (dastlabki, tajriba-sanoat, qabul sinovlari GOST 34.603) va qabul qilish komissiyasi tartibi.",
		Reference:   "Oʻz DSt 1987:2018 6.7-bandi, GOST 34.603",
	},
	{
		Number:      "7",
		TitleUz:     "ATni ishga tushirishga tayyorlash talablari",
		TitleRu:     "Требования к подготовке объекта к вводу в действие",
		Description: "Xodimlarni oʻqitish, axborot bazasini shakllantirish, klassifikatorlar va dastlabki maʼlumotlarni kiritish.",
		Reference:   "Oʻz DSt 1987:2018 6.8-bandi",
	},
	{
		Number:      "8",
		TitleUz:     "Hujjatlashtirishga qoʻyiladigan talablar",
		TitleRu:     "Требования к документированию",
		Description: "Ishlab chiqiladigan texnik va ekspluatatsion hujjatlar roʻyxati, Oʻz DSt 1985:2018 komplektliligi.",
		Reference:   "Oʻz DSt 1987:2018 6.9-bandi, Oʻz DSt 1985:2018",
	},
}

// BannedSubjectivePhrases lists qualitative and subjective expressions forbidden by the standard (UZ + RU)
var BannedSubjectivePhrases = []string{
	// O'zbekcha
	"tizim tez ishlashi kerak",
	"qulay dizayn boʻlishi lozim",
	"yuqori yuklamalarga chidamli boʻlishi lozim",
	"ishonchli ishlashi kerak",
	"zamonaviy texnologiyalar asosida",
	"tezkor qidiruv imkoniyati",
	"sodda va tushunarli interfeys",

	// Русский
	"система должна работать быстро",
	"удобный и интуитивно понятный интерфейс",
	"высокая надежность и отказоустойчивость",
	"на основе современных технологий",
	"быстрый поиск данных",
	"устойчивость к высоким нагрузкам",
	"простой интерфейс",
}

// CalculateStatus calculates decision status based on compliance score
func CalculateStatus(score float64) Status {
	if score >= 85.0 {
		return StatusApproved
	}
	if score >= 70.0 {
		return StatusConditionalPass
	}
	return StatusRejected
}
