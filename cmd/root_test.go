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
	if !strings.Contains(output, "specgen") {
		t.Errorf("expected help output to contain 'specgen', got: %s", output)
	}
}

func TestRootCommand_Version(t *testing.T) {
	output, err := executeRoot([]string{"--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, Version) {
		t.Errorf("expected version output to contain %q, got: %s", Version, output)
	}
}

func TestRootCommand_NoColorFlag(t *testing.T) {
	_, err := executeRoot([]string{"--no-color", "--help"})
	if err != nil {
		t.Fatalf("unexpected error with --no-color flag: %v", err)
	}
}

func TestRootCommand_HasExpectedSubcommands(t *testing.T) {
	expected := []string{"new", "generate", "resume", "version"}
	registered := map[string]bool{}
	for _, sub := range rootCmd.Commands() {
		registered[sub.Name()] = true
	}
	for _, name := range expected {
		if !registered[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}
