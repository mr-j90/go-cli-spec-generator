package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateCommand(t *testing.T) {
	buf := new(bytes.Buffer)
func TestGenerateCmd(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"generate"})

	if err := rootCmd.Execute(); err != nil {
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
func TestGenerateCmdHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"generate", "--help"})

	if err := rootCmd.Execute(); err != nil {
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "generate") {
		t.Errorf("expected help output to contain 'generate', got: %s", output)
		t.Errorf("expected 'generate' in help output, got: %s", output)
	}
}
