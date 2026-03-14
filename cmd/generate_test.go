package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func executeRoot(args []string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestGenerateCommand_RequiresDescription(t *testing.T) {
	_, err := executeRoot([]string{"generate"})
	if err == nil {
		t.Error("expected error when no description provided")
	}
}

func TestGenerateCommand_WithDescription(t *testing.T) {
	output, err := executeRoot([]string{"generate", "a user authentication system"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "a user authentication system") {
		t.Errorf("expected output to contain description, got: %s", output)
	}
	if !strings.Contains(output, "feature") {
		t.Errorf("expected output to contain default spec type 'feature', got: %s", output)
	}
}

func TestGenerateCommand_WithOutputFlag(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--output", "spec.md", "my feature"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "spec.md") {
		t.Errorf("expected output to mention output file, got: %s", output)
	}
}

func TestGenerateCommand_WithTypeFlag(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--type", "api", "my api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "api") {
		t.Errorf("expected output to contain spec type 'api', got: %s", output)
	}
}
