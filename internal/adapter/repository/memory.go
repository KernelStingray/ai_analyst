package repository

import (
	"sync"

	"github.com/hextec/spec-auditor/internal/domain"
)

// AuditRepository defines the interface for persisting and retrieving audit reports
type AuditRepository interface {
	Save(report domain.AuditReport) error
	GetByID(id string) (domain.AuditReport, bool)
	ListRecent(limit int) []domain.AuditReport
}

// MemoryRepository implements AuditRepository using in-memory thread-safe storage
type MemoryRepository struct {
	mu      sync.RWMutex
	reports map[string]domain.AuditReport
	order   []string
}

// NewMemoryRepository initializes an empty in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		reports: make(map[string]domain.AuditReport),
		order:   make([]string, 0),
	}
}

// Save stores an audit report in memory
func (r *MemoryRepository) Save(report domain.AuditReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.reports[report.AuditID]; !exists {
		// Store newest first in the ordered slice
		r.order = append([]string{report.AuditID}, r.order...)
	}
	r.reports[report.AuditID] = report
	return nil
}

// GetByID finds an audit report by its unique ID
func (r *MemoryRepository) GetByID(id string) (domain.AuditReport, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	report, ok := r.reports[id]
	return report, ok
}

// ListRecent returns the most recent audit reports up to the given limit
func (r *MemoryRepository) ListRecent(limit int) []domain.AuditReport {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 || limit > len(r.order) {
		limit = len(r.order)
	}

	result := make([]domain.AuditReport, 0, limit)
	for i := 0; i < limit; i++ {
		id := r.order[i]
		if rep, ok := r.reports[id]; ok {
			result = append(result, rep)
		}
	}
	return result
}
