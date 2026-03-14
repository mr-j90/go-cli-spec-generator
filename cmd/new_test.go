package cmd

import (
	"strings"
	"testing"
)

func TestNewCommand_Help(t *testing.T) {
	output, err := executeRoot([]string{"new", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "new") {
		t.Errorf("expected help output to contain 'new', got: %s", output)
	}
}

func TestNewCommand_Run(t *testing.T) {
	output, err := executeRoot([]string{"new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Launching") {
		t.Errorf("expected output to contain 'Launching', got: %s", output)
	}
}
