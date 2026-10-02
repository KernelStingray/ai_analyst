package domain

import (
	"testing"
)

func TestCalculateStatus(t *testing.T) {
	tests := []struct {
		score    float64
		expected Status
	}{
		{95.0, StatusApproved},
		{85.0, StatusApproved},
		{84.9, StatusConditionalPass},
		{70.0, StatusConditionalPass},
		{69.9, StatusRejected},
		{50.0, StatusRejected},
		{0.0, StatusRejected},
	}

	for _, tt := range tests {
		status := CalculateStatus(tt.score)
		if status != tt.expected {
			t.Errorf("CalculateStatus(%v) = %v, expected %v", tt.score, status, tt.expected)
		}
	}
}

func TestMandatorySectionsCount(t *testing.T) {
	if len(MandatorySections) != 8 {
		t.Errorf("MandatorySections length = %d, expected 8 according to Oʻz DSt 1987:2018 A ilovasi", len(MandatorySections))
	}
}
