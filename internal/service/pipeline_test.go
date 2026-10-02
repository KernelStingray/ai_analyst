package service

import (
	"testing"

	"github.com/hextec/spec-auditor/internal/domain"
)

func TestPipelineAuditor_UzbekDemo(t *testing.T) {
	pipe := NewPipelineAuditor()
	report := pipe.AuditPipeline("DXM-Stat TT", DemoPresetText)

	if report.Status != domain.StatusRejected {
		t.Errorf("expected status RAD ETILDI, got %s", report.Status)
	}

	if len(report.MissingSections) != 2 {
		t.Errorf("expected 2 missing sections, got %d: %v", len(report.MissingSections), report.MissingSections)
	}

	if len(report.Findings) < 4 {
		t.Errorf("expected at least 4 findings, got %d", len(report.Findings))
	}

	if report.ComplianceScore < 50 || report.ComplianceScore > 75 {
		t.Errorf("expected compliance score between 50 and 75, got %f", report.ComplianceScore)
	}
}

func TestPipelineAuditor_RussianDemo(t *testing.T) {
	pipe := NewPipelineAuditor()
	report := pipe.AuditPipeline("АИС Мониторинг ТЗ", DemoPresetTextRU)

	if report.Status != domain.StatusRejected {
		t.Errorf("expected status RAD ETILDI, got %s", report.Status)
	}

	if len(report.MissingSections) != 2 {
		t.Errorf("expected 2 missing sections, got %d", len(report.MissingSections))
	}
}
