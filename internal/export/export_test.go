package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zyx-holdings/go-spec/internal/export"
)

// ---- Format constant tests ----

func TestFormat_Constants(t *testing.T) {
	cases := []struct {
		name string
		f    export.Format
		want string
	}{
		{"PDF", export.FormatPDF, "pdf"},
		{"DOCX", export.FormatDOCX, "docx"},
		{"Markdown", export.FormatMarkdown, "markdown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if string(tc.f) != tc.want {
				t.Errorf("Format%s = %q, want %q", tc.name, tc.f, tc.want)
			}
		})
	}
}

// ---- ExportPDF tests ----

func TestExportPDF_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.pdf")

	if err := export.ExportPDF("hello world", path); err != nil {
		t.Fatalf("ExportPDF: unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("ExportPDF: output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("ExportPDF: output file is empty")
	}
}

func TestExportPDF_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.pdf")

	if err := export.ExportPDF("", path); err != nil {
		t.Fatalf("ExportPDF(empty): unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("ExportPDF(empty): output file not found: %v", err)
	}
}

func TestExportPDF_InvalidPath_ReturnsError(t *testing.T) {
	err := export.ExportPDF("content", "/nonexistent/dir/output.pdf")
	if err == nil {
		t.Error("ExportPDF(invalid path): expected error, got nil")
	}
}

// docxAvailable probes whether the UniOffice license is present by attempting
// a throw-away export. Tests that require a successful DOCX write call this
// helper and skip themselves when the library is unlicensed.
func docxAvailable(t *testing.T) bool {
	t.Helper()
	dir := t.TempDir()
	err := export.ExportDOCX("probe", filepath.Join(dir, "probe.docx"))
	if err != nil && strings.Contains(err.Error(), "license") {
		t.Skip("skipping DOCX test: UniOffice license not available in this environment")
		return false
	}
	return true
}

// ---- ExportDOCX tests ----

func TestExportDOCX_CreatesFile(t *testing.T) {
	if !docxAvailable(t) {
		return
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "output.docx")

	if err := export.ExportDOCX("hello world", path); err != nil {
		t.Fatalf("ExportDOCX: unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("ExportDOCX: output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("ExportDOCX: output file is empty")
	}
}

func TestExportDOCX_EmptyContent(t *testing.T) {
	if !docxAvailable(t) {
		return
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.docx")

	if err := export.ExportDOCX("", path); err != nil {
		t.Fatalf("ExportDOCX(empty): unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("ExportDOCX(empty): output file not found: %v", err)
	}
}

func TestExportDOCX_ReturnsError_OnFailure(t *testing.T) {
	// ExportDOCX must return an error on an unwritable path regardless of
	// license status. The invalid-path case fails before any license check.
	err := export.ExportDOCX("content", "/nonexistent/dir/output.docx")
	if err == nil {
		t.Error("ExportDOCX(invalid path): expected error, got nil")
	}
}
