package llm

import (
	"testing"
)

func TestSanitizeJSONString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Clean JSON passes through",
			input:    `{"audit_id": "123", "status": "RAD ETILDI"}`,
			expected: `{"audit_id": "123", "status": "RAD ETILDI"}`,
		},
		{
			name: "Markdown codefence with json tag",
			input: "```json\n" +
				`{"audit_id": "123", "status": "TASDIQLANDI"}` +
				"\n```",
			expected: `{"audit_id": "123", "status": "TASDIQLANDI"}`,
		},
		{
			name: "Markdown codefence without tag and leading/trailing text",
			input: "Mana ekspertiza xulosasi:\n```\n" +
				`{"compliance_score": 75.5}` +
				"\n```\nUmid qilamanki foydali boʻldi.",
			expected: `{"compliance_score": 75.5}`,
		},
		{
			name:     "No brackets returns empty",
			input:    "Faqat matn hech qanday JSON yoʻq",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SanitizeJSONString(tt.input)
			if actual != tt.expected {
				t.Errorf("SanitizeJSONString() = %q, expected %q", actual, tt.expected)
			}
		})
	}
}
