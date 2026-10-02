package document

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ExtractText extracts textual content from either a .docx file or a plain text file
func ExtractText(filename string, data []byte) (string, error) {
	lowerName := strings.ToLower(filename)
	if strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md") {
		return string(data), nil
	}

	if strings.HasSuffix(lowerName, ".docx") {
		return ExtractDocxText(data)
	}

	return "", fmt.Errorf("qoʻllab-quvvatlanmaydigan fayl formati: %s (faqat .docx, .txt formatlar qabul qilinadi)", filename)
}

// ExtractDocxText parses word/document.xml from a docx zip archive in memory
func ExtractDocxText(data []byte) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("docx arxivini ochishda xatolik: %w", err)
	}

	var documentXMLFile *zip.File
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			documentXMLFile = f
			break
		}
	}

	if documentXMLFile == nil {
		return "", errors.New("docx faylida word/document.xml topilmadi")
	}

	rc, err := documentXMLFile.Open()
	if err != nil {
		return "", fmt.Errorf("document.xml faylini oʻqishda xatolik: %w", err)
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	var sb strings.Builder
	inText := false

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("xml dekodlashda xatolik: %w", err)
		}

		switch el := token.(type) {
		case xml.StartElement:
			// w:t represents text node in OpenXML
			if el.Name.Local == "t" {
				inText = true
			}
		case xml.EndElement:
			if el.Name.Local == "t" {
				inText = false
			} else if el.Name.Local == "p" {
				// Paragraph end -> line break
				sb.WriteString("\n")
			}
		case xml.CharData:
			if inText {
				sb.Write(el)
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
