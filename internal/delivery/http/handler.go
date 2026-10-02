package http

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hextec/spec-auditor/internal/adapter/document"
	"github.com/hextec/spec-auditor/internal/domain"
	"github.com/hextec/spec-auditor/internal/service"
)

// Handler manages API endpoints for the audit platform
type Handler struct {
	auditor service.Auditor
}

// NewHandler constructs a Handler instance
func NewHandler(auditor service.Auditor) *Handler {
	return &Handler{
		auditor: auditor,
	}
}

// RegisterRoutes sets up API endpoints on the given ServeMux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/audit", h.HandleAudit)
	mux.HandleFunc("POST /api/upload", h.HandleUpload)
	mux.HandleFunc("GET /api/preset", h.HandlePreset)
	mux.HandleFunc("GET /api/health", h.HandleHealth)
	mux.HandleFunc("GET /api/audit/", h.HandleGetAudit)
	mux.HandleFunc("GET /api/history", h.HandleHistory)
}

// HandleAudit processes an incoming specification for compliance
func (h *Handler) HandleAudit(w http.ResponseWriter, r *http.Request) {
	var req domain.AuditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Yaroqsiz JSON formati: "+err.Error())
		return
	}

	report, err := h.auditor.Audit(r.Context(), req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, report)
}

// HandleUpload accepts a multipart .docx or .txt file and returns extracted text
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// 10 MB limit for file upload
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "Faylni yuklashda xatolik (maksimal 10 MB): "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Fayl topilmadi: "+err.Error())
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Faylni oʻqishda xatolik: "+err.Error())
		return
	}

	text, err := document.ExtractText(header.Filename, fileBytes)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	docName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"document_name": docName,
		"content":       text,
		"char_count":    len([]rune(text)),
	})
}

// HandlePreset returns sample flawed specification text in UZ or RU
func (h *Handler) HandlePreset(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	WriteJSON(w, http.StatusOK, map[string]string{
		"preset": h.auditor.GetPreset(lang),
		"lang":   lang,
	})
}

// HandleHealth returns backend and local LLM operational status
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	llmOnline := h.auditor.IsLLMAvailable(r.Context())
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "ok",
		"llm_available": llmOnline,
		"llm_model":     "qwen2.5:7b",
		"standards":     []string{"Oʻz DSt 1987:2018", "Oʻz DSt 1985:2018", "Oʻz DSt 1986:2018"},
		"demo_mode":     !llmOnline,
	})
}

// HandleGetAudit retrieves an existing report by ID
func (h *Handler) HandleGetAudit(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/audit/")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "Hisobot ID koʻrsatilmadi")
		return
	}

	report, err := h.auditor.GetReport(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "Hisobot topilmadi")
		return
	}

	WriteJSON(w, http.StatusOK, report)
}

// HandleHistory returns recent audit summaries
func (h *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	reports := h.auditor.ListReports(limit)
	WriteJSON(w, http.StatusOK, reports)
}
