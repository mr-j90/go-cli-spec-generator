package cmd

import (
	"strings"
	"testing"
)

func TestVersionCommand_PrintsVersion(t *testing.T) {
	output, err := executeRoot([]string{"version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, Version) {
		t.Errorf("expected output to contain version %q, got: %s", Version, output)
	}
	if !strings.Contains(output, "specgen") {
		t.Errorf("expected output to contain 'specgen', got: %s", output)
	}
}
