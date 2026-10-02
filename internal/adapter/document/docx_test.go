package document

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestExtractDocxText(t *testing.T) {
	// Build a minimal synthetic docx zip in memory
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("failed to create entry: %v", err)
	}

	sampleXML := `<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Oʻzbekiston Respublikasi</w:t></w:r></w:p>
    <w:p><w:r><w:t>Texnik topshiriq loyihasi</w:t></w:r></w:p>
  </w:body>
</w:document>`

	if _, err := w.Write([]byte(sampleXML)); err != nil {
		t.Fatalf("failed to write xml: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip: %v", err)
	}

	text, err := ExtractDocxText(buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(text, "Oʻzbekiston Respublikasi") || !strings.Contains(text, "Texnik topshiriq loyihasi") {
		t.Errorf("extracted text missing expected phrases, got:\n%s", text)
	}
}

func TestExtractText_Plain(t *testing.T) {
	content := "Oddiy matn fayli"
	res, err := ExtractText("spec.txt", []byte(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != content {
		t.Errorf("expected %q, got %q", content, res)
	}
}
