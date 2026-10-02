package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hextec/spec-auditor/internal/domain"
)

// OllamaConfig holds connection parameters for the local LLM
type OllamaConfig struct {
	BaseURL string
	Model   string
	Timeout time.Duration
	NumCtx  int
}

// DefaultOllamaConfig returns recommended production defaults for Qwen 2.5 7B
func DefaultOllamaConfig() OllamaConfig {
	return OllamaConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen2.5:7b",
		Timeout: 300 * time.Second,
		NumCtx:  4096,
	}
}

// OllamaClient communicates with local Ollama inference server
type OllamaClient struct {
	config     OllamaConfig
	httpClient *http.Client
}

// NewOllamaClient creates a new Ollama API client
func NewOllamaClient(cfg OllamaConfig) *OllamaClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:11434"
	}
	if cfg.Model == "" {
		cfg.Model = "qwen2.5:7b"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 300 * time.Second
	}
	if cfg.NumCtx <= 0 {
		cfg.NumCtx = 4096
	}

	return &OllamaClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// generateRequest is the payload sent to Ollama /api/generate
type generateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	System  string                 `json:"system,omitempty"`
	Stream  bool                   `json:"stream"`
	Format  string                 `json:"format"`
	Options map[string]interface{} `json:"options"`
}

// generateResponse is the response from Ollama /api/generate
type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// RawLLMAuditResult matches the strict JSON schema expected from LLM
type RawLLMAuditResult struct {
	AuditID          string                  `json:"audit_id"`
	AnalyzedAt       string                  `json:"analyzed_at"`
	DocumentName     string                  `json:"document_name"`
	ComplianceScore  float64                 `json:"compliance_score"`
	Status           string                  `json:"status"`
	Metrics          domain.MetricScore      `json:"metrics"`
	MissingSections  []string                `json:"missing_sections"`
	Findings         []domain.SectionFinding `json:"findings"`
	ExecutiveSummary string                  `json:"executive_summary"`
}

// Ping checks whether the Ollama server is reachable and the target model is downloaded
func (c *OllamaClient) Ping(ctx context.Context) bool {
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/api/tags", strings.TrimRight(c.config.BaseURL, "/"))
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false
	}

	targetModel := strings.ToLower(c.config.Model)
	for _, m := range tags.Models {
		name := strings.ToLower(m.Name)
		if strings.HasPrefix(name, targetModel) || strings.Contains(name, targetModel) {
			return true
		}
	}

	return false
}

// GenerateAudit invokes the local model and parses the returned JSON audit
func (c *OllamaClient) GenerateAudit(ctx context.Context, systemPrompt, documentText string) (*RawLLMAuditResult, error) {
	reqBody := generateRequest{
		Model:  c.config.Model,
		System: systemPrompt,
		Prompt: documentText,
		Stream: false,
		Format: "json",
		Options: map[string]interface{}{
			"num_ctx":     c.config.NumCtx,
			"temperature": 0.1,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", strings.TrimRight(c.config.BaseURL, "/"))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	var genResp generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}

	rawCleaned := SanitizeJSONString(genResp.Response)
	if rawCleaned == "" {
		return nil, errors.New("ollama returned empty response after sanitization")
	}

	var result RawLLMAuditResult
	if err := json.Unmarshal([]byte(rawCleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse structured audit json: %w (raw response: %s)", err, rawCleaned)
	}

	return &result, nil
}

// SanitizeJSONString removes markdown blocks and trims non-JSON characters
func SanitizeJSONString(raw string) string {
	s := strings.TrimSpace(raw)

	// Remove markdown fences: ```json ... ``` or ``` ... ```
	if strings.HasPrefix(s, "```") {
		firstLineEnd := strings.Index(s, "\n")
		if firstLineEnd != -1 {
			s = s[firstLineEnd+1:]
		}
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}

	// Extract everything between the first '{' and the last '}'
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start != -1 && end != -1 && end >= start {
		return s[start : end+1]
	}

	return ""
}
