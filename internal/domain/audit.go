package domain

import (
	"time"
)

// Status represents the certification/audit decision
type Status string

const (
	StatusApproved        Status = "TASDIQLANDI"          // >= 85%
	StatusConditionalPass Status = "SHARTLI TASDIQLANDI" // 70-84%
	StatusRejected        Status = "RAD ETILDI"           // < 70%
)

// Severity represents finding criticality level
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityWarning  Severity = "WARNING"
	SeverityNote     Severity = "NOTE"
)

// MetricScore contains sub-scores for specific audit categories
type MetricScore struct {
	StructureCompleteness float64 `json:"structure_completeness"` // A ilovasi: 8 mandatory sections
	ScenariosCompliance   float64 `json:"scenarios_compliance"`   // B ilovasi: Use-case specifications
	MeasurabilityScore    float64 `json:"measurability_score"`    // SLA / quantitative criteria
}

// SectionFinding describes a specific defect or violation in the document
type SectionFinding struct {
	ID             string   `json:"id"`
	Severity       Severity `json:"severity"`
	Section        string   `json:"section"`
	RawFragment    string   `json:"raw_fragment"`
	Violation      string   `json:"violation"`
	Recommendation string   `json:"recommendation"`
}

// AuditRequest is the incoming payload for auditing a specification
type AuditRequest struct {
	DocumentName string `json:"document_name"`
	Content      string `json:"content"`
}

// SectionOverview contains status and metadata of a standard section
type SectionOverview struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Detected    bool   `json:"detected"`
	StandardRef string `json:"standard_ref"`
	Snippet     string `json:"snippet,omitempty"`
	CharCount   int    `json:"char_count,omitempty"`
	IssuesCount int    `json:"issues_count"`
}

// AuditReport is the full structured expert conclusion (Ekspertiza Xulosasi)
type AuditReport struct {
	AuditID          string            `json:"audit_id"`
	AnalyzedAt       time.Time         `json:"analyzed_at"`
	DocumentName     string            `json:"document_name"`
	ComplianceScore  float64           `json:"compliance_score"`
	Status           Status            `json:"status"`
	Metrics          MetricScore       `json:"metrics"`
	MissingSections  []string          `json:"missing_sections"`
	SectionMap       []SectionOverview `json:"section_map,omitempty"`
	Findings         []SectionFinding  `json:"findings"`
	ExecutiveSummary string            `json:"executive_summary"`
	IsDemoMode       bool              `json:"is_demo_mode"`
}
