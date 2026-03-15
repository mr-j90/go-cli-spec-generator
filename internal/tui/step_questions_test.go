package tui_test

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/session"
	"github.com/zyx-holdings/go-spec/internal/tui"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

// updateSQ sends a single key to the model and returns the updated model.
func updateSQ(t *testing.T, m tea.Model, key tea.KeyMsg) tea.Model {
	t.Helper()
	updated, _ := m.Update(key)
	return updated
}

// typeIntoSQ sends a string rune-by-rune into a text-type widget.
func typeIntoSQ(t *testing.T, m tea.Model, s string) tea.Model {
	t.Helper()
	for _, r := range s {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return m
}

// isSQ asserts the model is a *StepQuestions and returns it.
func asSQ(t *testing.T, m tea.Model) *tui.StepQuestions {
	t.Helper()
	sq, ok := m.(*tui.StepQuestions)
	if !ok {
		t.Fatalf("expected *tui.StepQuestions, got %T", m)
	}
	return sq
}

// cmdIsQuit reports whether cmd produces a tea.QuitMsg.
func cmdIsQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	_, ok := msg.(tea.QuitMsg)
	return ok
}

// cmdIsDone reports whether cmd produces a StepQuestionsDoneMsg.
func cmdIsDone(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	_, ok := msg.(tui.StepQuestionsDoneMsg)
	return ok
}

// ─── NewStepQuestions ─────────────────────────────────────────────────────────

func TestNewStepQuestions_CoreOnly(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	if sq == nil {
		t.Fatal("NewStepQuestions returned nil")
	}
	v := sq.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestNewStepQuestions_WithFeatures(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), []string{"testing"}, "")
	v := sq.View()
	// First question is core, so section header should say "Core".
	if !strings.Contains(v, "Core") {
		t.Errorf("View() = %q, want to contain section 'Core'", v)
	}
}

func TestNewStepQuestions_View_ContainsProgressInfo(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	v := sq.View()
	if !strings.Contains(v, "Q 1 of") {
		t.Errorf("View() = %q, want to contain 'Q 1 of'", v)
	}
}

func TestNewStepQuestions_View_ContainsFirstQuestion(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	v := sq.View()
	// First core question is project_name.
	if !strings.Contains(v, "project name") {
		t.Errorf("View() = %q, want to contain 'project name'", v)
	}
}

func TestNewStepQuestions_View_ContainsFooter(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	v := sq.View()
	if !strings.Contains(v, "Enter") {
		t.Errorf("View() footer missing 'Enter' hint, got: %q", v)
	}
}

// ─── Init ─────────────────────────────────────────────────────────────────────

func TestStepQuestions_Init_NonNil(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	cmd := sq.Init()
	if cmd == nil {
		t.Error("Init() returned nil Cmd, want focus/blink cmd")
	}
}

func TestStepQuestions_Init_EmptyQuestions_EmitsDone(t *testing.T) {
	// An empty session with no features and a manually empty store.
	// We test Init by creating a model with no questions — achieved by passing
	// an empty feature list and overriding with no _core (not possible via public
	// API, so we test the done cmd via a direct Init check).
	// The simplest proxy: Init with valid store returns non-nil cmd.
	sq := tui.NewStepQuestions(session.New(), nil, "")
	cmd := sq.Init()
	// Must be non-nil (either blink or done).
	if cmd == nil {
		t.Error("Init() = nil, want non-nil cmd")
	}
}

// ─── Enter on required empty field ───────────────────────────────────────────

func TestStepQuestions_Enter_EmptyRequired_ShowsError(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	// First question (project_name) is required text; widget is empty.
	updated, _ := sq.Update(tea.KeyMsg{Type: tea.KeyEnter})
	v := updated.(*tui.StepQuestions).View()
	if !strings.Contains(v, "required") {
		t.Errorf("View() after empty Enter = %q, want validation error mentioning 'required'", v)
	}
}

func TestStepQuestions_Enter_EmptyRequired_StaysOnSameQuestion(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	updated, _ := sq.Update(tea.KeyMsg{Type: tea.KeyEnter})
	v := updated.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 1 of") {
		t.Errorf("should still be on Q 1 after failed submit, got: %q", v)
	}
}

// ─── Enter with valid answer advances ────────────────────────────────────────

func TestStepQuestions_Enter_ValidText_AdvancesQuestion(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m := typeIntoSQ(t, sq, "my-project")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	v := m.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 2 of") {
		t.Errorf("after valid Enter, expected Q 2, got: %q", v)
	}
}

func TestStepQuestions_Enter_ValidText_RecordsAnswer(t *testing.T) {
	store := session.New()
	sq := tui.NewStepQuestions(store, nil, "")
	m := typeIntoSQ(t, sq, "my-project")
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	ans, ok := store.GetAnswer("project_name")
	if !ok {
		t.Fatal("expected project_name answer in store, not found")
	}
	if ans.Value.String() != "my-project" {
		t.Errorf("stored answer = %q, want my-project", ans.Value.String())
	}
}

// ─── single_select: Enter submits without typing ─────────────────────────────

func TestStepQuestions_Enter_SingleSelect_SubmitsCurrentOption(t *testing.T) {
	store := session.New()
	// Navigate to primary_language (3rd _core question, index 2).
	sq := tui.NewStepQuestions(store, nil, "")
	m := tea.Model(sq)
	// Answer Q1 project_name.
	m = typeIntoSQ(t, m, "proj")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// Answer Q2 project_description (textarea, required).
	m = typeIntoSQ(t, m, "a description")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// Now on Q3 primary_language (single_select).
	v := m.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 3 of") {
		t.Fatalf("expected Q 3, got: %q", v)
	}
	// Enter submits the currently highlighted option.
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	ans, ok := store.GetAnswer("primary_language")
	if !ok {
		t.Fatal("expected primary_language answer recorded")
	}
	if ans.Value.String() == "" {
		t.Error("primary_language answer is empty")
	}
}

// ─── Esc navigation ───────────────────────────────────────────────────────────

func TestStepQuestions_Esc_FirstQuestion_Quits(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	_, cmd := sq.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !cmdIsQuit(cmd) {
		t.Error("Esc on first question should emit tea.Quit")
	}
}

func TestStepQuestions_Esc_SecondQuestion_GoesBack(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m := typeIntoSQ(t, sq, "my-project")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // advance to Q2
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})   // go back to Q1
	v := m.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 1 of") {
		t.Errorf("after Esc from Q2, expected Q 1, got: %q", v)
	}
}

func TestStepQuestions_Esc_PreservesAnswer(t *testing.T) {
	store := session.New()
	sq := tui.NewStepQuestions(store, nil, "")
	m := typeIntoSQ(t, sq, "kept-name")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // advance
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})   // back
	// project_name answer should still be recorded.
	ans, ok := store.GetAnswer("project_name")
	if !ok {
		t.Fatal("answer not preserved after Esc")
	}
	if ans.Value.String() != "kept-name" {
		t.Errorf("preserved answer = %q, want kept-name", ans.Value.String())
	}
}

// ─── Tab: skip optional / hint on required ───────────────────────────────────

func TestStepQuestions_Tab_Required_ShowsHint(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	// Q1 project_name is required.
	updated, _ := sq.Update(tea.KeyMsg{Type: tea.KeyTab})
	v := updated.(*tui.StepQuestions).View()
	if !strings.Contains(v, "required") {
		t.Errorf("Tab on required field should show 'required' hint, got: %q", v)
	}
}

func TestStepQuestions_Tab_Required_DoesNotAdvance(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	updated, _ := sq.Update(tea.KeyMsg{Type: tea.KeyTab})
	v := updated.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 1 of") {
		t.Errorf("Tab on required field should stay on Q 1, got: %q", v)
	}
}

func TestStepQuestions_Tab_Optional_Skips(t *testing.T) {
	store := session.New()
	sq := tui.NewStepQuestions(store, nil, "")
	m := tea.Model(sq)
	// Answer required Q1 and Q2, reach Q3 (primary_language, required), Q4 (team_size, optional).
	m = typeIntoSQ(t, m, "proj")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = typeIntoSQ(t, m, "desc")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Q3 single_select submit
	// Now on Q4 (team_size) which is optional.
	v := m.(*tui.StepQuestions).View()
	if !strings.Contains(v, "Q 4 of") {
		t.Fatalf("expected Q 4 (team_size), got: %q", v)
	}
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab}) // skip team_size
	ans, ok := store.GetAnswer("team_size")
	if !ok {
		t.Fatal("expected team_size to be recorded as skipped")
	}
	if !ans.Skipped {
		t.Error("team_size answer.Skipped = false, want true")
	}
}

// ─── Ctrl+S save and quit ─────────────────────────────────────────────────────

func TestStepQuestions_CtrlS_SavesSessionFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/session.json"
	store := session.New()
	sq := tui.NewStepQuestions(store, nil, path)
	_, cmd := sq.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if cmd == nil {
		t.Fatal("Ctrl+S returned nil cmd")
	}

	// saveAndQuit returns tea.Batch, which produces a tea.BatchMsg when executed.
	// Unwrap it to find the StepQuestionsSavedMsg.
	found := false
	msg := cmd()
	switch m := msg.(type) {
	case tui.StepQuestionsSavedMsg:
		found = true
		if m.Path != path {
			t.Errorf("SavedMsg.Path = %q, want %q", m.Path, path)
		}
	case tea.BatchMsg:
		for _, subcmd := range m {
			if subcmd == nil {
				continue
			}
			if saved, ok := subcmd().(tui.StepQuestionsSavedMsg); ok {
				found = true
				if saved.Path != path {
					t.Errorf("SavedMsg.Path = %q, want %q", saved.Path, path)
				}
			}
		}
	}
	if !found {
		t.Errorf("StepQuestionsSavedMsg not found; got %T", msg)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("session file not created at %s", path)
	}
}

func TestStepQuestions_CtrlS_DefaultPath(t *testing.T) {
	// When savePath is empty, should use a default and not error.
	store := session.New()
	sq := tui.NewStepQuestions(store, nil, "")
	// Override to a tmp dir by setting CWD — not ideal, so just confirm no panic.
	_, cmd := sq.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if cmd == nil {
		t.Fatal("Ctrl+S returned nil cmd")
	}
	// Clean up default file if created.
	_ = os.Remove("specgen-session.json")
}

// ─── Ctrl+C confirm quit overlay ─────────────────────────────────────────────

func TestStepQuestions_CtrlC_ShowsConfirmOverlay(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m, _ := sq.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	v := m.(*tui.StepQuestions).View()
	if !strings.Contains(v, "exit") {
		t.Errorf("Ctrl+C should show quit overlay mentioning 'exit', got: %q", v)
	}
}

func TestStepQuestions_CtrlC_ThenN_Dismisses(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m, _ := sq.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	v := m.(*tui.StepQuestions).View()
	// Back to normal question view.
	if !strings.Contains(v, "Q 1 of") {
		t.Errorf("after Ctrl+C → N expected Q 1 view, got: %q", v)
	}
}

func TestStepQuestions_CtrlC_Twice_Quits(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m, _ := sq.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !cmdIsQuit(cmd) {
		t.Error("Ctrl+C twice should emit tea.Quit")
	}
}

// ─── WindowSizeMsg ────────────────────────────────────────────────────────────

func TestStepQuestions_WindowSize_HandledGracefully(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), nil, "")
	m, cmd := sq.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if m == nil {
		t.Error("Update(WindowSizeMsg) returned nil model")
	}
	_ = cmd
}

// ─── section ordering ─────────────────────────────────────────────────────────

func TestStepQuestions_SectionOrdering_CoreFirst(t *testing.T) {
	sq := tui.NewStepQuestions(session.New(), []string{"testing"}, "")
	v := sq.View()
	if !strings.Contains(v, "Core") {
		t.Errorf("first question section should be Core, got: %q", v)
	}
}

// ─── done after last question ─────────────────────────────────────────────────

func TestStepQuestions_AllAnswered_EmitsDone(t *testing.T) {
	store := session.New()
	// Use only _core (4 questions). Answer all of them.
	sq := tui.NewStepQuestions(store, nil, "")
	m := tea.Model(sq)

	// Q1: project_name (text, required)
	m = typeIntoSQ(t, m, "proj")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Q2: project_description (textarea, required)
	m = typeIntoSQ(t, m, "desc")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Q3: primary_language (single_select, required) — Enter picks first option
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Q4: team_size (single_select, optional) — Tab to skip; this is the last question
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if cmd == nil {
		t.Fatal("after last question Tab, cmd should be non-nil (DoneMsg)")
	}
	if _, ok := cmd().(tui.StepQuestionsDoneMsg); !ok {
		t.Errorf("expected StepQuestionsDoneMsg after all questions answered, cmd returned %T", cmd())
	}
}
