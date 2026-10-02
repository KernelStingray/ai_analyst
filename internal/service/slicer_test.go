package service

import (
	"testing"
)

func TestSliceTechnicalSpecification_UzbekPreset(t *testing.T) {
	sliced := SliceTechnicalSpecification(DemoPresetText)

	if sliced.IsUnstructured {
		t.Fatalf("expected structured document, but was marked unstructured")
	}

	// In DemoPresetText, sections 1, 2, 3, 4, 6, 7 are present (6 sections)
	if len(sliced.Sections) != 6 {
		t.Errorf("expected 6 sections, got %d", len(sliced.Sections))
	}

	// Missing sections should be 5 and 8
	if len(sliced.MissingSections) != 2 {
		t.Fatalf("expected 2 missing sections, got %d: %v", len(sliced.MissingSections), sliced.MissingSections)
	}

	hasSec5 := false
	hasSec8 := false
	for _, m := range sliced.MissingSections {
		if len(m) > 0 && m[0] == '5' {
			hasSec5 = true
		}
		if len(m) > 0 && m[0] == '8' {
			hasSec8 = true
		}
	}

	if !hasSec5 || !hasSec8 {
		t.Errorf("expected missing sections 5 and 8, got %v", sliced.MissingSections)
	}

	if !sliced.HasScenarios {
		t.Errorf("expected scenarios to be detected")
	}
}

func TestSliceTechnicalSpecification_RussianPreset(t *testing.T) {
	sliced := SliceTechnicalSpecification(DemoPresetTextRU)

	if sliced.IsUnstructured {
		t.Fatalf("expected structured document, but was marked unstructured")
	}

	if len(sliced.Sections) != 6 {
		t.Errorf("expected 6 sections, got %d", len(sliced.Sections))
	}

	if len(sliced.MissingSections) != 2 {
		t.Errorf("expected 2 missing sections, got %d: %v", len(sliced.MissingSections), sliced.MissingSections)
	}
}
