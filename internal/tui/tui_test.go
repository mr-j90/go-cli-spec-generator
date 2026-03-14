package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/tui"
)

func TestNew_ReturnsModel(t *testing.T) {
	m := tui.New()
	// View() must return non-empty output after construction.
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string after New()")
	}
}

func TestNew_ViewContainsTitle(t *testing.T) {
	m := tui.New()
	v := m.View()
	if !strings.Contains(v, "specgen") {
		t.Errorf("View() = %q, expected to contain 'specgen'", v)
	}
}

func TestInit_ReturnsCmd(t *testing.T) {
	m := tui.New()
	// Init() must return a non-nil Cmd (textinput.Blink).
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() returned nil Cmd, want non-nil (textinput.Blink)")
	}
}

func TestUpdate_CtrlC_Quits(t *testing.T) {
	m := tui.New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("Update(Ctrl+C) returned nil Cmd, want tea.Quit")
	}
	// tea.Quit returns a tea.QuitMsg when executed; compare the message type.
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("Update(Ctrl+C) cmd() = %T, want tea.QuitMsg", msg)
	}
}

func TestUpdate_Esc_Quits(t *testing.T) {
	m := tui.New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Update(Esc) returned nil Cmd, want tea.Quit")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("Update(Esc) cmd() = %T, want tea.QuitMsg", msg)
	}
}

func TestUpdate_NonQuitKey_DoesNotQuit(t *testing.T) {
	m := tui.New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	// cmd may be nil or non-nil (textinput blink), but must NOT be tea.Quit.
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Error("Update(rune 'a') returned tea.Quit, want no quit")
		}
	}
}

func TestUpdate_ReturnsModel(t *testing.T) {
	m := tui.New()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if updated == nil {
		t.Error("Update() returned nil model")
	}
}
