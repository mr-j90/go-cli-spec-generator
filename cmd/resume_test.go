package cmd

import (
	"strings"
	"testing"
)

func TestResumeCommand_Help(t *testing.T) {
	output, err := executeRoot([]string{"resume", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "resume") {
		t.Errorf("expected help output to contain 'resume', got: %s", output)
	}
}

func TestResumeCommand_Run(t *testing.T) {
	output, err := executeRoot([]string{"resume"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Resuming") {
		t.Errorf("expected output to contain 'Resuming', got: %s", output)
	}
}

func TestResumeCommand_WithSessionFlag(t *testing.T) {
	output, err := executeRoot([]string{"resume", "--session", "my-session.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "my-session.json") {
		t.Errorf("expected output to mention session file, got: %s", output)
	}
}
