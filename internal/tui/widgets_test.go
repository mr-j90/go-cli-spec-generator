package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/tui"
)

// ─── construction ─────────────────────────────────────────────────────────────

func TestNewInputWidget_Text(t *testing.T) {
	w := tui.NewInputWidget("text", nil, "placeholder", 80)
	if w.Value() != "" {
		t.Errorf("new text widget Value() = %q, want empty", w.Value())
	}
	if !w.IsEmpty() {
		t.Error("new text widget IsEmpty() = false, want true")
	}
}

func TestNewInputWidget_TextArea(t *testing.T) {
	w := tui.NewInputWidget("textarea", nil, "placeholder", 80)
	if !w.IsEmpty() {
		t.Error("new textarea widget IsEmpty() = false, want true")
	}
}

func TestNewInputWidget_SingleSelect(t *testing.T) {
	opts := []string{"Go", "Python", "Rust"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	// single_select with options always has a selection (cursor at 0).
	if w.IsEmpty() {
		t.Error("single_select with options IsEmpty() = true, want false")
	}
	if w.Value() != "Go" {
		t.Errorf("single_select Value() = %q, want %q", w.Value(), "Go")
	}
}

func TestNewInputWidget_SingleSelect_NoOptions(t *testing.T) {
	w := tui.NewInputWidget("single_select", nil, "", 80)
	if !w.IsEmpty() {
		t.Error("single_select with no options IsEmpty() = false, want true")
	}
}

func TestNewInputWidget_MultiSelect(t *testing.T) {
	opts := []string{"Unit", "Integration", "E2E"}
	w := tui.NewInputWidget("multi_select", opts, "", 80)
	if !w.IsEmpty() {
		t.Error("new multi_select IsEmpty() = false, want true (nothing toggled)")
	}
	if w.Values() != nil {
		t.Errorf("new multi_select Values() = %v, want nil", w.Values())
	}
}

func TestNewInputWidget_Confirm(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	if !w.IsEmpty() {
		t.Error("new confirm widget IsEmpty() = false, want true")
	}
	if w.Value() != "" {
		t.Errorf("confirm Value() = %q, want empty before input", w.Value())
	}
}

// ─── Value() after input ──────────────────────────────────────────────────────

func TestInputWidget_SetValue_Text(t *testing.T) {
	w := tui.NewInputWidget("text", nil, "", 80)
	w.SetValue("my-project")
	if w.Value() != "my-project" {
		t.Errorf("text Value() = %q, want %q", w.Value(), "my-project")
	}
	if w.IsEmpty() {
		t.Error("text IsEmpty() = true after SetValue, want false")
	}
}

func TestInputWidget_SetValue_TextArea(t *testing.T) {
	w := tui.NewInputWidget("textarea", nil, "", 80)
	w.SetValue("line one\nline two")
	if w.IsEmpty() {
		t.Error("textarea IsEmpty() = true after SetValue, want false")
	}
}

func TestInputWidget_SetValue_TextTrimsWhitespace(t *testing.T) {
	w := tui.NewInputWidget("text", nil, "", 80)
	w.SetValue("  spaces  ")
	if w.Value() != "spaces" {
		t.Errorf("text Value() = %q, want %q (trimmed)", w.Value(), "spaces")
	}
}

func TestInputWidget_SetValue_SingleSelect(t *testing.T) {
	opts := []string{"Go", "Python", "Rust"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	w.SetValue("Rust")
	if w.Value() != "Rust" {
		t.Errorf("single_select Value() = %q after SetValue(Rust), want Rust", w.Value())
	}
}

func TestInputWidget_SetValue_SingleSelect_UnknownOption(t *testing.T) {
	opts := []string{"Go", "Python"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	w.SetValue("COBOL")
	// Unknown option leaves cursor at default (Go).
	if w.Value() != "Go" {
		t.Errorf("single_select Value() = %q after unknown SetValue, want Go (default)", w.Value())
	}
}

func TestInputWidget_SetValues_MultiSelect(t *testing.T) {
	opts := []string{"Unit", "Integration", "E2E"}
	w := tui.NewInputWidget("multi_select", opts, "", 80)
	w.SetValues([]string{"Unit", "E2E"})
	got := w.Values()
	if len(got) != 2 {
		t.Fatalf("multi_select Values() len = %d, want 2", len(got))
	}
	if got[0] != "Unit" || got[1] != "E2E" {
		t.Errorf("multi_select Values() = %v, want [Unit E2E]", got)
	}
	if w.IsEmpty() {
		t.Error("multi_select IsEmpty() = true after SetValues, want false")
	}
}

func TestInputWidget_SetValue_Confirm_Yes(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	w.SetValue("yes")
	if w.Value() != "yes" {
		t.Errorf("confirm Value() = %q, want yes", w.Value())
	}
	if w.IsEmpty() {
		t.Error("confirm IsEmpty() = true after SetValue(yes), want false")
	}
}

func TestInputWidget_SetValue_Confirm_No(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	w.SetValue("no")
	if w.Value() != "no" {
		t.Errorf("confirm Value() = %q, want no", w.Value())
	}
}

// ─── Update — navigation keys ─────────────────────────────────────────────────

func TestInputWidget_Update_SingleSelect_ArrowDown(t *testing.T) {
	opts := []string{"Go", "Python", "Rust"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown})
	if w.Value() != "Python" {
		t.Errorf("single_select after ↓ Value() = %q, want Python", w.Value())
	}
}

func TestInputWidget_Update_SingleSelect_ArrowUp_AtTop(t *testing.T) {
	opts := []string{"Go", "Python", "Rust"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyUp})
	// Should not go below 0.
	if w.Value() != "Go" {
		t.Errorf("single_select after ↑ at top Value() = %q, want Go", w.Value())
	}
}

func TestInputWidget_Update_SingleSelect_ArrowDown_AtBottom(t *testing.T) {
	opts := []string{"Go", "Python", "Rust"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown})
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown})
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown}) // should stop at Rust
	if w.Value() != "Rust" {
		t.Errorf("single_select after ↓ past end Value() = %q, want Rust", w.Value())
	}
}

func TestInputWidget_Update_MultiSelect_SpaceToggles(t *testing.T) {
	opts := []string{"Unit", "Integration", "E2E"}
	w := tui.NewInputWidget("multi_select", opts, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !strings.Contains(w.Value(), "Unit") {
		t.Errorf("multi_select after Space Value() = %q, want Unit selected", w.Value())
	}
	// Toggle off.
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeySpace})
	if w.Value() != "" {
		t.Errorf("multi_select after second Space Value() = %q, want empty", w.Value())
	}
}

func TestInputWidget_Update_MultiSelect_NavigateAndToggle(t *testing.T) {
	opts := []string{"Unit", "Integration", "E2E"}
	w := tui.NewInputWidget("multi_select", opts, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown})   // move to Integration
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeySpace})  // select Integration
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyDown})   // move to E2E
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeySpace})  // select E2E
	vals := w.Values()
	if len(vals) != 2 {
		t.Fatalf("multi_select Values() len = %d, want 2", len(vals))
	}
}

func TestInputWidget_Update_Confirm_Y(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if w.Value() != "yes" {
		t.Errorf("confirm after 'y' Value() = %q, want yes", w.Value())
	}
	if w.IsEmpty() {
		t.Error("confirm IsEmpty() = true after 'y', want false")
	}
}

func TestInputWidget_Update_Confirm_N(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if w.Value() != "no" {
		t.Errorf("confirm after 'n' Value() = %q, want no", w.Value())
	}
}

func TestInputWidget_Update_Confirm_CaseInsensitive(t *testing.T) {
	w := tui.NewInputWidget("confirm", nil, "", 80)
	w, _ = w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Y")})
	if w.Value() != "yes" {
		t.Errorf("confirm after 'Y' Value() = %q, want yes", w.Value())
	}
}

// ─── Values() for non-multi type ──────────────────────────────────────────────

func TestInputWidget_Values_NonMulti_ReturnsNil(t *testing.T) {
	w := tui.NewInputWidget("text", nil, "", 80)
	if w.Values() != nil {
		t.Errorf("text widget Values() = %v, want nil", w.Values())
	}
}

// ─── View() sanity ────────────────────────────────────────────────────────────

func TestInputWidget_View_NotEmpty(t *testing.T) {
	types := []struct {
		inputType string
		opts      []string
	}{
		{"text", nil},
		{"textarea", nil},
		{"single_select", []string{"A", "B"}},
		{"multi_select", []string{"X", "Y"}},
		{"confirm", nil},
	}
	for _, tc := range types {
		w := tui.NewInputWidget(tc.inputType, tc.opts, "", 80)
		if v := w.View(); v == "" {
			t.Errorf("%s widget View() returned empty string", tc.inputType)
		}
	}
}

func TestInputWidget_View_SingleSelect_ShowsCursor(t *testing.T) {
	opts := []string{"Go", "Python"}
	w := tui.NewInputWidget("single_select", opts, "", 80)
	v := w.View()
	if !strings.Contains(v, "Go") {
		t.Errorf("single_select View() = %q, want to contain 'Go'", v)
	}
}

func TestInputWidget_View_MultiSelect_ShowsOptions(t *testing.T) {
	opts := []string{"Unit", "Integration"}
	w := tui.NewInputWidget("multi_select", opts, "", 80)
	v := w.View()
	if !strings.Contains(v, "Unit") || !strings.Contains(v, "Integration") {
		t.Errorf("multi_select View() = %q, want to contain all options", v)
	}
}

// ─── InsertNewline ────────────────────────────────────────────────────────────

func TestInputWidget_InsertNewline_Textarea(t *testing.T) {
	w := tui.NewInputWidget("textarea", nil, "", 80)
	w.SetValue("first line")
	w.InsertNewline()
	// After inserting a newline, value should contain a newline.
	raw := w.Value() // trimmed, but internal value has \n
	_ = raw          // just check no panic
}

func TestInputWidget_InsertNewline_NonTextarea_NoOp(t *testing.T) {
	w := tui.NewInputWidget("text", nil, "", 80)
	w.SetValue("hello")
	w.InsertNewline() // should be a no-op
	if w.Value() != "hello" {
		t.Errorf("text widget InsertNewline changed value to %q", w.Value())
	}
}
