// Package docx menyediakan builder minimal untuk membuat file .docx (OOXML)
// dari scratch tanpa dependency eksternal, menggunakan hanya stdlib Go.
package docx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
)

// ── OOXML boilerplate ────────────────────────────────────────────────────────

const xmlContentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`

const xmlRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

const xmlWordRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`

const xmlDocumentFmt = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
%s
    <w:sectPr/>
  </w:body>
</w:document>`

// ── Builder ──────────────────────────────────────────────────────────────────

// Builder mengumpulkan paragraf OOXML dan membuat file .docx saat Build() dipanggil.
type Builder struct {
	paragraphs []string
}

// New membuat Builder baru.
func New() *Builder {
	return &Builder{}
}

// AddTitle menambahkan paragraf judul: bold, besar (18pt), rata tengah.
func (b *Builder) AddTitle(text string) *Builder {
	p := fmt.Sprintf(
		`    <w:p><w:pPr><w:jc w:val="center"/></w:pPr>`+
			`<w:r><w:rPr><w:b/><w:sz w:val="36"/><w:szCs w:val="36"/></w:rPr>`+
			`<w:t>%s</w:t></w:r></w:p>`,
		xmlEscape(text),
	)
	b.paragraphs = append(b.paragraphs, p)
	return b
}

// AddCentered menambahkan paragraf teks biasa rata tengah.
func (b *Builder) AddCentered(text string) *Builder {
	p := fmt.Sprintf(
		`    <w:p><w:pPr><w:jc w:val="center"/></w:pPr>`+
			`<w:r><w:t xml:space="preserve">%s</w:t></w:r></w:p>`,
		xmlEscape(text),
	)
	b.paragraphs = append(b.paragraphs, p)
	return b
}

// AddBold menambahkan paragraf dengan teks bold (mendukung newline di dalam teks).
func (b *Builder) AddBold(text string) *Builder {
	p := fmt.Sprintf(`    <w:p>%s</w:p>`, textRuns(text, true))
	b.paragraphs = append(b.paragraphs, p)
	return b
}

// AddNormal menambahkan paragraf teks biasa (mendukung newline di dalam teks).
func (b *Builder) AddNormal(text string) *Builder {
	p := fmt.Sprintf(`    <w:p>%s</w:p>`, textRuns(text, false))
	b.paragraphs = append(b.paragraphs, p)
	return b
}

// AddLabelValue menambahkan paragraf dengan label bold diikuti value normal
// pada baris yang sama. Contoh: "What did you complete yesterday? <jawaban>"
func (b *Builder) AddLabelValue(label, value string) *Builder {
	var labelSb strings.Builder
	labelSb.WriteString(`<w:p>`)
	labelSb.WriteString(`<w:r><w:rPr><w:b/></w:rPr>`)
	labelSb.WriteString(fmt.Sprintf(
		`<w:t xml:space="preserve">%s</w:t>`,
		xmlEscape(label),
	))
	labelSb.WriteString(`</w:r>`)
	labelSb.WriteString(`</w:p>`)

	// Paragraph 2: value / jawaban
	var valueSb strings.Builder
	valueSb.WriteString(`<w:p>`)
	valueSb.WriteString(textRuns(value, false))
	valueSb.WriteString(`</w:p>`)

	b.paragraphs = append(b.paragraphs, labelSb.String(), valueSb.String())

	return b
}

// AddSeparator menambahkan garis pemisah (karakter Unicode box-drawing).
func (b *Builder) AddSeparator() *Builder {
	b.paragraphs = append(b.paragraphs,
		`    <w:p><w:r><w:t>─────────────────────────────────────</w:t></w:r></w:p>`)
	return b
}

// AddBlank menambahkan paragraf kosong (baris kosong).
func (b *Builder) AddBlank() *Builder {
	b.paragraphs = append(b.paragraphs, `    <w:p/>`)
	return b
}

// Build menghasilkan file .docx sebagai []byte.
func (b *Builder) Build() ([]byte, error) {
	docXML := fmt.Sprintf(xmlDocumentFmt, strings.Join(b.paragraphs, "\n"))

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	entries := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", xmlContentTypes},
		{"_rels/.rels", xmlRels},
		{"word/_rels/document.xml.rels", xmlWordRels},
		{"word/document.xml", docXML},
	}

	for _, e := range entries {
		f, err := zw.Create(e.name)
		if err != nil {
			return nil, fmt.Errorf("docx: create entry %q: %w", e.name, err)
		}
		if _, err := f.Write([]byte(e.content)); err != nil {
			return nil, fmt.Errorf("docx: write entry %q: %w", e.name, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("docx: close zip: %w", err)
	}

	return buf.Bytes(), nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

// xmlEscape melakukan escaping karakter khusus XML.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// textRuns mengubah teks (termasuk newline) menjadi serangkaian XML <w:r> runs.
// Newline direpresentasikan sebagai <w:br/> di dalam run terpisah.
func textRuns(text string, bold bool) string {
	// Normalise line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")

	var sb strings.Builder
	for i, line := range lines {
		if i > 0 {
			// Line break antar baris
			sb.WriteString(`<w:r><w:br/></w:r>`)
		}
		if line == "" {
			continue
		}
		sb.WriteString(`<w:r>`)
		if bold {
			sb.WriteString(`<w:rPr><w:b/></w:rPr>`)
		}
		sb.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, xmlEscape(line)))
		sb.WriteString(`</w:r>`)
	}
	return sb.String()
}
