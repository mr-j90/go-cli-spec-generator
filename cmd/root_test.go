package cmd

import (
	"strings"
	"testing"
)

func TestRootCommand_Help(t *testing.T) {
	output, err := executeRoot([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "go-spec") {
		t.Errorf("expected help output to contain 'go-spec', got: %s", output)
	}
}

func TestRootCommandHasGenerateSubcommand(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "generate [description]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'generate' subcommand to be registered")
	}
}
