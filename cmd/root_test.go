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
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	buf := new(bytes.Buffer)
func TestRootCmd(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected help output, got empty string")
	}
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
