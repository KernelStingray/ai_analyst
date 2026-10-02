package service

import (
	"context"
	"strings"
	"testing"

	"github.com/hextec/spec-auditor/internal/adapter/llm"
	"github.com/hextec/spec-auditor/internal/adapter/repository"
	"github.com/hextec/spec-auditor/internal/domain"
)

func TestAuditorService_ValidationAndDemoMode(t *testing.T) {
	repo := repository.NewMemoryRepository()
	// Create client with invalid port so it is definitely offline
	llmCfg := llm.DefaultOllamaConfig()
	llmCfg.BaseURL = "http://127.0.0.1:59999"
	llmClient := llm.NewOllamaClient(llmCfg)

	svc := NewAuditorService(repo, llmClient)
	ctx := context.Background()

	// 1. Too short content
	_, err := svc.Audit(ctx, domain.AuditRequest{
		DocumentName: "Qisqa TT",
		Content:      "Juda qisqa matn",
	})
	if err == nil || !strings.Contains(err.Error(), "juda qisqa") {
		t.Fatalf("expected too short error, got: %v", err)
	}

	// 2. Too long content (> 25,000 chars)
	tooLong := strings.Repeat("A", MaxContentLength+10)
	_, err = svc.Audit(ctx, domain.AuditRequest{
		DocumentName: "Katta TT",
		Content:      tooLong,
	})
	if err == nil || !strings.Contains(err.Error(), "chegaradan oshdi") {
		t.Fatalf("expected too long error, got: %v", err)
	}

	// 3. Offline Demo Mode execution (preset text)
	report, err := svc.Audit(ctx, domain.AuditRequest{
		DocumentName: "DXM-Stat AT Texnik topshirigʻi",
		Content:      DemoPresetText,
	})
	if err != nil {
		t.Fatalf("unexpected error in offline demo mode: %v", err)
	}

	if !report.IsDemoMode {
		t.Errorf("expected IsDemoMode to be true when LLM is offline")
	}

	if report.Status != domain.StatusRejected {
		t.Errorf("expected status to be %v, got %v", domain.StatusRejected, report.Status)
	}

	if len(report.MissingSections) != 2 {
		t.Errorf("expected 2 missing sections, got %d", len(report.MissingSections))
	}

	if len(report.Findings) < 4 {
		t.Errorf("expected at least 4 findings, got %d", len(report.Findings))
	}

	// Check that report was saved into repository
	saved, found := repo.GetByID(report.AuditID)
	if !found {
		t.Errorf("report not found in repository by ID %s", report.AuditID)
	}
	if saved.ComplianceScore != report.ComplianceScore {
		t.Errorf("saved report score mismatch: %v vs %v", saved.ComplianceScore, report.ComplianceScore)
	}
}
