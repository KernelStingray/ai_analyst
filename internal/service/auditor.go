package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hextec/spec-auditor/internal/adapter/llm"
	"github.com/hextec/spec-auditor/internal/adapter/repository"
	"github.com/hextec/spec-auditor/internal/domain"
)

const (
	MaxContentLength = 500000 // Support full 40-150 page technical specifications
	MinContentLength = 50
)

// Auditor defines business logic for specification compliance audits
type Auditor interface {
	Audit(ctx context.Context, req domain.AuditRequest) (*domain.AuditReport, error)
	GetPreset(lang string) string
	IsLLMAvailable(ctx context.Context) bool
	GetReport(id string) (*domain.AuditReport, error)
	ListReports(limit int) []domain.AuditReport
}

// AuditorService implements Auditor
type AuditorService struct {
	repo      repository.AuditRepository
	llmClient *llm.OllamaClient
	pipeline  *PipelineAuditor
}

// NewAuditorService constructs a new AuditorService
func NewAuditorService(repo repository.AuditRepository, llmClient *llm.OllamaClient) *AuditorService {
	return &AuditorService{
		repo:      repo,
		llmClient: llmClient,
		pipeline:  NewPipelineAuditor(),
	}
}

// GetPreset returns the benchmark demo flawed specification text based on language
func (s *AuditorService) GetPreset(lang string) string {
	if strings.ToLower(lang) == "ru" {
		return DemoPresetTextRU
	}
	return DemoPresetText
}

// IsLLMAvailable checks if the local Ollama daemon is reachable
func (s *AuditorService) IsLLMAvailable(ctx context.Context) bool {
	if s.llmClient == nil {
		return false
	}
	return s.llmClient.Ping(ctx)
}

// GetReport fetches an audit report by its ID
func (s *AuditorService) GetReport(id string) (*domain.AuditReport, error) {
	report, found := s.repo.GetByID(id)
	if !found {
		return nil, errors.New("audit report not found")
	}
	return &report, nil
}

// ListReports returns recent audit reports
func (s *AuditorService) ListReports(limit int) []domain.AuditReport {
	return s.repo.ListRecent(limit)
}

// Audit validates and processes the technical specification using the automated multi-stage pipeline
func (s *AuditorService) Audit(ctx context.Context, req domain.AuditRequest) (*domain.AuditReport, error) {
	trimmed := strings.TrimSpace(req.Content)
	length := len([]rune(trimmed))

	if length < MinContentLength {
		return nil, fmt.Errorf("matn hajmi juda qisqa (kamida %d ta belgi kiritilishi shart, sizda: %d)", MinContentLength, length)
	}

	if length > MaxContentLength {
		return nil, fmt.Errorf("matn hajmi belgilangan chegaradan oshdi (maksimal %d ta belgi, sizda: %d)", MaxContentLength, length)
	}

	if req.DocumentName == "" {
		req.DocumentName = "Nomsiz TT Loyihasi"
	}

	// 1. Run Automated Section Slicer & Pipeline Audit
	pipeReport := s.pipeline.AuditPipeline(req.DocumentName, trimmed)

	// Check whether Ollama is reachable
	llmOnline := s.llmClient != nil && s.llmClient.Ping(ctx)

	if !llmOnline {
		pipeReport.IsDemoMode = true
		_ = s.repo.Save(*pipeReport)
		return pipeReport, nil
	}

	// 2. If Ollama is active, augment with AI executive summary or findings
	aiCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// For large technical specifications (> 6,000 chars), synthesize a concise
	// executive digest of detected sections to avoid CPU context saturation
	llmInput := trimmed
	if length > 6000 {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Hujjat nomi: %s\nJami hajmi: %d belgi\n\n", req.DocumentName, length))
		if len(pipeReport.MissingSections) > 0 {
			b.WriteString("Tushirib qoldirilgan boʻlimlar:\n")
			for _, ms := range pipeReport.MissingSections {
				b.WriteString(fmt.Sprintf("- %s\n", ms))
			}
			b.WriteString("\n")
		}
		b.WriteString("Aniqlangan boʻlimlar mazmuni:\n")
		for _, sec := range pipeReport.SectionMap {
			if sec.Detected && sec.Snippet != "" {
				b.WriteString(fmt.Sprintf("[%s]: %s\n", sec.Title, sec.Snippet))
			}
		}
		runes := []rune(trimmed)
		limit := 3000
		if len(runes) < limit {
			limit = len(runes)
		}
		b.WriteString(fmt.Sprintf("\nHujjat matnidan parcha:\n%s\n", string(runes[:limit])))
		llmInput = b.String()
	}

	systemPrompt := BuildSystemPrompt()
	rawResult, err := s.llmClient.GenerateAudit(aiCtx, systemPrompt, llmInput)
	if err == nil && rawResult != nil {
		if rawResult.ExecutiveSummary != "" {
			pipeReport.ExecutiveSummary = rawResult.ExecutiveSummary
		}
		if rawResult.ComplianceScore > 0 {
			// Weighted blend of deterministic pipeline score and LLM perception
			blendedScore := (pipeReport.ComplianceScore * 0.7) + (rawResult.ComplianceScore * 0.3)
			pipeReport.ComplianceScore = float64(int(blendedScore*10)) / 10.0
			pipeReport.Status = domain.CalculateStatus(pipeReport.ComplianceScore)
		}
		// Merge any unique findings discovered by the LLM
		for _, f := range rawResult.Findings {
			if !hasDuplicateFinding(pipeReport.Findings, f) {
				pipeReport.Findings = append(pipeReport.Findings, f)
			}
		}
	}

	_ = s.repo.Save(*pipeReport)
	return pipeReport, nil
}

func hasDuplicateFinding(existing []domain.SectionFinding, candidate domain.SectionFinding) bool {
	for _, e := range existing {
		if strings.EqualFold(e.Section, candidate.Section) && strings.Contains(candidate.RawFragment, e.RawFragment) {
			return true
		}
	}
	return false
}

func isDemoLike(text string) bool {
	return strings.Contains(text, "DXM-Stat") ||
		strings.Contains(text, "Markaz rahbariyati") ||
		strings.Contains(text, "Oʻzbekiston Respublikasi Vazirligi huzuridagi") ||
		strings.Contains(text, "АИС Мониторинг") ||
		strings.Contains(text, "МИНИСТЕРСТВО ЦИФРОВЫХ ТЕХНОЛОГИЙ")
}

func isRussianText(text string) bool {
	cyrillicCount := 0
	for _, r := range text {
		if (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'ё' || r == 'Ё' {
			cyrillicCount++
			if cyrillicCount > 20 {
				return true
			}
		}
	}
	return false
}
