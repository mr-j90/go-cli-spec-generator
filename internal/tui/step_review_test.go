package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/session"
	"github.com/zyx-holdings/go-spec/internal/tui"
)

// newTestStore returns a fresh session store with no answers for use in tests.
func newTestStore() *session.Store {
	return session.New()
}

// newTestStoreWithFeatures returns a store with the given feature areas selected.
func newTestStoreWithFeatures(features []string) *session.Store {
	st := session.New()
	st.Session().SelectedFeatures = features
	return st
}

// sendKey sends a single KeyMsg to the model and returns the updated model.
func sendKey(m tea.Model, key tea.KeyType) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: key})
}

// sendRune sends a rune KeyMsg to the model and returns the updated model.
func sendRune(m tea.Model, r rune) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
}

// cmdMsg executes a tea.Cmd and returns the resulting tea.Msg. Returns nil if
// cmd is nil.
func cmdMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

// ── Construction ────────────────────────────────────────────────────────────

func TestNewReviewStep_ReturnsModel(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if m.View() == "" {
		t.Error("View() returned empty string")
	}
}

func TestNewReviewStep_ViewContainsTitle(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if !strings.Contains(m.View(), "Step 4") {
		t.Errorf("View() does not contain 'Step 4': %q", m.View())
	}
}

func TestNewReviewStep_ViewContainsCoreSection(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if !strings.Contains(m.View(), "Core") {
		t.Errorf("View() does not contain 'Core': %q", m.View())
	}
}

func TestNewReviewStep_HasOneSectionForEmptySession(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if m.SectionCount() != 1 {
		t.Errorf("SectionCount() = %d, want 1 (core only)", m.SectionCount())
	}
}

func TestNewReviewStep_SectionsMatchSelectedFeatures(t *testing.T) {
	st := newTestStoreWithFeatures([]string{"authentication", "storage"})
	m := tui.NewReviewStep(st, 80, 24)
	// core + authentication + storage = 3
	if m.SectionCount() != 3 {
		t.Errorf("SectionCount() = %d, want 3", m.SectionCount())
	}
}

func TestNewReviewStep_ViewContainsFeatureSection(t *testing.T) {
	st := newTestStoreWithFeatures([]string{"authentication"})
	m := tui.NewReviewStep(st, 80, 24)
	if !strings.Contains(m.View(), "Authentication") {
		t.Errorf("View() does not contain 'Authentication'")
	}
}

func TestNewReviewStep_InitialCursorIsZero(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if m.Cursor() != 0 {
		t.Errorf("Cursor() = %d, want 0", m.Cursor())
	}
}

func TestNewReviewStep_InitReturnsNil(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if m.Init() != nil {
		t.Error("Init() returned non-nil cmd, want nil")
	}
}

// ── Navigation ───────────────────────────────────────────────────────────────

func TestReviewStep_DownMovesCursor(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendKey(m, tea.KeyDown)
	rs := updated.(tui.ReviewStep)
	if rs.Cursor() != 1 {
		t.Errorf("after Down: Cursor() = %d, want 1", rs.Cursor())
	}
}

func TestReviewStep_UpDoesNotGoBelowZero(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendKey(m, tea.KeyUp)
	rs := updated.(tui.ReviewStep)
	if rs.Cursor() != 0 {
		t.Errorf("after Up at 0: Cursor() = %d, want 0", rs.Cursor())
	}
}

func TestReviewStep_DownDoesNotExceedMax(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// totalItems = 1 section + Export + Back = 3; max cursor = 2
	for i := 0; i < 10; i++ {
		updated, _ := sendKey(m, tea.KeyDown)
		m = updated.(tui.ReviewStep)
	}
	// cursor must not exceed 2 (Back button index)
	if m.Cursor() > 2 {
		t.Errorf("Cursor() = %d after 10 Downs, want ≤ 2", m.Cursor())
	}
}

func TestReviewStep_UpDownRoundtrip(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendKey(m, tea.KeyDown)
	updated, _ = sendKey(updated, tea.KeyUp)
	rs := updated.(tui.ReviewStep)
	if rs.Cursor() != 0 {
		t.Errorf("after Down then Up: Cursor() = %d, want 0", rs.Cursor())
	}
}

// ── Section collapse ─────────────────────────────────────────────────────────

func TestReviewStep_EnterTogglesCollapse(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// cursor is on section 0 (core)
	if m.IsSectionCollapsed(0) {
		t.Fatal("section 0 should not be collapsed initially")
	}
	updated, _ := sendKey(m, tea.KeyEnter)
	rs := updated.(tui.ReviewStep)
	if !rs.IsSectionCollapsed(0) {
		t.Error("section 0 should be collapsed after Enter")
	}
}

func TestReviewStep_EnterTogglesCollapseBackToExpanded(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendKey(m, tea.KeyEnter) // collapse
	updated, _ = sendKey(updated, tea.KeyEnter) // expand
	rs := updated.(tui.ReviewStep)
	if rs.IsSectionCollapsed(0) {
		t.Error("section 0 should be expanded after two Enter presses")
	}
}

func TestReviewStep_CollapsedSectionHidesQuestions(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// Collapse core section
	updated, _ := sendKey(m, tea.KeyEnter)
	rs := updated.(tui.ReviewStep)
	// The core section has question "What is the project name?" — should not appear when collapsed
	if strings.Contains(rs.View(), "project name") {
		t.Error("collapsed section should not show question text")
	}
}

// ── Buttons ──────────────────────────────────────────────────────────────────

func TestReviewStep_ExportButtonSendsReviewDoneMsg(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// Navigate: section 0 → Export button (index 1)
	updated, _ := sendKey(m, tea.KeyDown)
	_, cmd := sendKey(updated, tea.KeyEnter)
	msg := cmdMsg(cmd)
	if _, ok := msg.(tui.ReviewDoneMsg); !ok {
		t.Errorf("Enter on Export button: got msg %T, want ReviewDoneMsg", msg)
	}
}

func TestReviewStep_BackButtonSendsReviewBackMsg(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// Navigate: section 0 → Export (1) → Back (2)
	updated, _ := sendKey(m, tea.KeyDown)
	updated, _ = sendKey(updated, tea.KeyDown)
	_, cmd := sendKey(updated, tea.KeyEnter)
	msg := cmdMsg(cmd)
	if _, ok := msg.(tui.ReviewBackMsg); !ok {
		t.Errorf("Enter on Back button: got msg %T, want ReviewBackMsg", msg)
	}
}

func TestReviewStep_BKeySendsReviewBackMsg(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	_, cmd := sendRune(m, 'b')
	msg := cmdMsg(cmd)
	if _, ok := msg.(tui.ReviewBackMsg); !ok {
		t.Errorf("'b' key: got msg %T, want ReviewBackMsg", msg)
	}
}

func TestReviewStep_EnterOnSectionDoesNotProceed(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	// cursor on section 0 — Enter should collapse, not emit ReviewDoneMsg
	_, cmd := sendKey(m, tea.KeyEnter)
	msg := cmdMsg(cmd)
	if _, ok := msg.(tui.ReviewDoneMsg); ok {
		t.Error("Enter on section should not emit ReviewDoneMsg")
	}
}

// ── Jump-edit mode ───────────────────────────────────────────────────────────

func TestReviewStep_EKeyEntersJumpEditMode(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendRune(m, 'e')
	rs := updated.(tui.ReviewStep)
	if !rs.IsJumpEditMode() {
		t.Error("'e' key should enter jump-edit mode")
	}
}

func TestReviewStep_EscExitsJumpEditMode(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendRune(m, 'e')
	updated, _ = sendKey(updated, tea.KeyEsc)
	rs := updated.(tui.ReviewStep)
	if rs.IsJumpEditMode() {
		t.Error("Esc should exit jump-edit mode")
	}
}

func TestReviewStep_JumpEditEnterSendsJumpMsg(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendRune(m, 'e')
	updated, _ = sendRune(updated, '3')
	_, cmd := sendKey(updated, tea.KeyEnter)
	msg := cmdMsg(cmd)
	jm, ok := msg.(tui.ReviewJumpEditMsg)
	if !ok {
		t.Fatalf("Enter after 'e3': got msg %T, want ReviewJumpEditMsg", msg)
	}
	if jm.QuestionIndex != 3 {
		t.Errorf("ReviewJumpEditMsg.QuestionIndex = %d, want 3", jm.QuestionIndex)
	}
}

func TestReviewStep_JumpEditViewShowsPrompt(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := sendRune(m, 'e')
	rs := updated.(tui.ReviewStep)
	if !strings.Contains(rs.View(), "Jump to question") {
		t.Error("jump-edit mode should show prompt in View()")
	}
}

// ── Skipped answers ──────────────────────────────────────────────────────────

func TestReviewStep_SkippedAnswerShownDimmed(t *testing.T) {
	st := session.New()
	st.SkipAnswer("project_name")
	m := tui.NewReviewStep(st, 80, 24)
	if !strings.Contains(m.View(), "skipped") {
		t.Error("skipped answer should appear as '(skipped)' in View()")
	}
}

func TestReviewStep_AnsweredQuestionShown(t *testing.T) {
	st := session.New()
	st.SetAnswer("project_name", session.NewStringValue("MyProject"))
	m := tui.NewReviewStep(st, 80, 24)
	if !strings.Contains(m.View(), "MyProject") {
		t.Error("answered question value should appear in View()")
	}
}

func TestReviewStep_MultiAnswerShownJoined(t *testing.T) {
	st := session.New()
	st.SetAnswer("auth_providers", session.NewMultiValue([]string{"google", "github"}))
	st.Session().SelectedFeatures = []string{"authentication"}
	m := tui.NewReviewStep(st, 80, 24)
	v := m.View()
	if !strings.Contains(v, "google") || !strings.Contains(v, "github") {
		t.Error("multi-value answer should appear joined in View()")
	}
}

// ── Window resize ────────────────────────────────────────────────────────────

func TestReviewStep_WindowResizeUpdatesModel(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	rs := updated.(tui.ReviewStep)
	// After resize the model should still render valid output.
	if rs.View() == "" {
		t.Error("View() empty after WindowSizeMsg")
	}
}

// ── Buttons visible in view ──────────────────────────────────────────────────

func TestReviewStep_ViewContainsExportButton(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if !strings.Contains(m.View(), "Export") {
		t.Error("View() should contain 'Export' button")
	}
}

func TestReviewStep_ViewContainsBackButton(t *testing.T) {
	m := tui.NewReviewStep(newTestStore(), 80, 24)
	if !strings.Contains(m.View(), "Back") {
		t.Error("View() should contain 'Back' button")
	}
}
