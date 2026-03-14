// Package export handles rendering spec documents to PDF and DOCX formats.
package export

import (
	"github.com/signintech/gopdf"
	"github.com/unidoc/unioffice/document"
)

// Format represents a supported export format.
type Format string

const (
	FormatPDF      Format = "pdf"
	FormatDOCX     Format = "docx"
	FormatMarkdown Format = "markdown"
)

// ExportPDF writes content to a PDF file at the given path.
// TODO: implement full PDF rendering.
func ExportPDF(content, outputPath string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	_ = pdf // TODO: write content
	return pdf.WritePdf(outputPath)
}

// ExportDOCX writes content to a DOCX file at the given path.
// TODO: implement full DOCX rendering.
func ExportDOCX(content, outputPath string) error {
	doc := document.New()
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(content)
	return doc.SaveToFile(outputPath)
}
