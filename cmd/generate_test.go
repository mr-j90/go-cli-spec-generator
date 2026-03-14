package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// executeRoot is a test helper that runs the root command with the given args
// and returns captured output. It resets all flag state before each run to
// prevent test pollution from Cobra's singleton command tree.
func executeRoot(args []string) (string, error) {
	resetCmdFlags(rootCmd)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

// resetCmdFlags resets all flag values on cmd and its subcommands to defaults.
func resetCmdFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetCmdFlags(sub)
	}
}

func TestGenerateCommand_Help(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "generate") {
		t.Errorf("expected help output to contain 'generate', got: %s", output)
	}
}

func TestGenerateCommand_NoArgs(t *testing.T) {
	output, err := executeRoot([]string{"generate"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Generating") {
		t.Errorf("expected output to contain 'Generating', got: %s", output)
	}
}

func TestGenerateCommand_WithInputFlag(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--input", "session.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "session.json") {
		t.Errorf("expected output to mention input file, got: %s", output)
	}
}

func TestGenerateCommand_WithOutputFlag(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--output", "spec.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "spec.md") {
		t.Errorf("expected output to mention output file, got: %s", output)
	}
}

func TestGenerateCommand_WithFormatFlag(t *testing.T) {
	output, err := executeRoot([]string{"generate", "--format", "pdf"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "pdf") {
		t.Errorf("expected output to mention format, got: %s", output)
	}
}
