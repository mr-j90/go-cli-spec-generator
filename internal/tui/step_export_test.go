package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/tui"
)

// ── Construction ────────────────────────────────────────────────────────────

func TestNewExportStep_ReturnsModel(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if m.View() == "" {
		t.Error("View() returned empty string")
	}
}

func TestNewExportStep_ViewContainsTitle(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if !strings.Contains(m.View(), "Step 5") {
		t.Errorf("View() does not contain 'Step 5': %q", m.View())
	}
}

func TestNewExportStep_ViewContainsAllFormats(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	v := m.View()
	for _, label := range []string{"JSON", "Markdown", "Word Document", "PDF"} {
		if !strings.Contains(v, label) {
			t.Errorf("View() missing format label %q", label)
		}
	}
}

func TestNewExportStep_InitialCursorIsZero(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if m.ExportCursor() != 0 {
		t.Errorf("ExportCursor() = %d, want 0", m.ExportCursor())
	}
}

func TestNewExportStep_InitReturnsNilWithoutPreselect(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if m.Init() != nil {
		t.Error("Init() should be nil when no preselected formats")
	}
}

func TestNewExportStep_NotSkippedByDefault(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if m.IsSkipped() {
		t.Error("IsSkipped() should be false when no preselected formats")
	}
}

func TestNewExportStep_NotDoneInitially(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if m.IsDone() {
		t.Error("IsDone() should be false initially")
	}
}

func TestNewExportStep_NoFormatsSelectedInitially(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if len(m.SelectedFormats()) != 0 {
		t.Errorf("SelectedFormats() = %v, want empty", m.SelectedFormats())
	}
}

// ── Preselected formats (--format flag) ──────────────────────────────────────

func TestNewExportStep_PreselectedFormatsAreChecked(t *testing.T) {
	m := tui.NewExportStep([]string{"pdf", "json"}, "", "test-session", 80, 24)
	selected := m.SelectedFormats()
	if len(selected) != 2 {
		t.Fatalf("SelectedFormats() len = %d, want 2", len(selected))
	}
	found := map[string]bool{}
	for _, f := range selected {
		found[f] = true
	}
	if !found["pdf"] || !found["json"] {
		t.Errorf("SelectedFormats() = %v, want [pdf json]", selected)
	}
}

func TestNewExportStep_PreselectedMarksAsSkipped(t *testing.T) {
	m := tui.NewExportStep([]string{"markdown"}, "", "test-session", 80, 24)
	if !m.IsSkipped() {
		t.Error("IsSkipped() should be true when preselected formats provided")
	}
}

func TestNewExportStep_PreselectedInitReturnsCmdNotNil(t *testing.T) {
	m := tui.NewExportStep([]string{"pdf"}, "", "test-session", 80, 24)
	if m.Init() == nil {
		t.Error("Init() should return non-nil cmd when preselected formats are set")
	}
}

// ── Navigation ───────────────────────────────────────────────────────────────

func TestExportStep_DownMovesCursor(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeyDown)
	es := updated.(tui.ExportStep)
	if es.ExportCursor() != 1 {
		t.Errorf("after Down: ExportCursor() = %d, want 1", es.ExportCursor())
	}
}

func TestExportStep_UpDoesNotGoBelowZero(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeyUp)
	es := updated.(tui.ExportStep)
	if es.ExportCursor() != 0 {
		t.Errorf("after Up at 0: ExportCursor() = %d, want 0", es.ExportCursor())
	}
}

func TestExportStep_DownDoesNotExceedLastFormat(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	// 4 formats → max cursor = 3
	for i := 0; i < 10; i++ {
		updated, _ := sendKey(m, tea.KeyDown)
		m = updated.(tui.ExportStep)
	}
	if m.ExportCursor() > 3 {
		t.Errorf("ExportCursor() = %d after 10 Downs, want ≤ 3", m.ExportCursor())
	}
}

// ── Format selection ─────────────────────────────────────────────────────────

func TestExportStep_SpaceTogglesFormatOn(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace)
	es := updated.(tui.ExportStep)
	sel := es.SelectedFormats()
	if len(sel) != 1 || sel[0] != "json" {
		t.Errorf("after Space on first item: SelectedFormats() = %v, want [json]", sel)
	}
}

func TestExportStep_SpaceTogglesFormatOff(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace) // select json
	updated, _ = sendKey(updated, tea.KeySpace) // deselect json
	es := updated.(tui.ExportStep)
	if len(es.SelectedFormats()) != 0 {
		t.Errorf("after two Spaces: SelectedFormats() = %v, want []", es.SelectedFormats())
	}
}

func TestExportStep_SpaceOnSecondItem(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeyDown)   // move to markdown
	updated, _ = sendKey(updated, tea.KeySpace) // select markdown
	es := updated.(tui.ExportStep)
	sel := es.SelectedFormats()
	if len(sel) != 1 || sel[0] != "markdown" {
		t.Errorf("SelectedFormats() = %v, want [markdown]", sel)
	}
}

func TestExportStep_MultipleFormatsCanBeSelected(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace) // select json (cursor 0)
	updated, _ = sendKey(updated, tea.KeyDown)
	updated, _ = sendKey(updated, tea.KeySpace) // select markdown (cursor 1)
	es := updated.(tui.ExportStep)
	if len(es.SelectedFormats()) != 2 {
		t.Errorf("SelectedFormats() len = %d, want 2", len(es.SelectedFormats()))
	}
}

// ── Enter key validation ─────────────────────────────────────────────────────

func TestExportStep_EnterWithNoSelectionDoesNotProceed(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, cmd := sendKey(m, tea.KeyEnter)
	es := updated.(tui.ExportStep)
	if es.IsDone() {
		t.Error("IsDone() should be false when Enter pressed with no selection")
	}
	if cmd != nil {
		t.Error("cmd should be nil when Enter pressed with no selection")
	}
}

func TestExportStep_EnterWithSelectionStartsExporting(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace) // select json
	updated, cmd := sendKey(updated, tea.KeyEnter)
	es := updated.(tui.ExportStep)
	// Should be in exporting phase (not done yet, but cmd returned)
	if cmd == nil {
		t.Error("Enter with selection should return a non-nil cmd to start export")
	}
	if es.IsDone() {
		t.Error("IsDone() should be false immediately after starting export")
	}
}

func TestExportStep_ViewShowsValidationWarningWhenNoSelection(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	if !strings.Contains(m.View(), "least one") {
		t.Error("View() should warn when no format is selected")
	}
}

func TestExportStep_ViewHidesValidationWarningWhenSelected(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace)
	es := updated.(tui.ExportStep)
	if strings.Contains(es.View(), "least one") {
		t.Error("View() should not show validation warning when a format is selected")
	}
}

// ── Export progress ───────────────────────────────────────────────────────────

func TestExportStep_ProgressMsgUpdatesView(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace) // select json
	updated, _ = sendKey(updated, tea.KeyEnter) // start export

	// Inject a successful progress message.
	updated, _ = updated.Update(tui.ExportProgressMsg{
		Format: "json",
		Path:   "/tmp/test-session.json",
		Err:    nil,
	})
	es := updated.(tui.ExportStep)
	if !es.IsDone() {
		t.Error("IsDone() should be true after all selected formats finish")
	}
}

func TestExportStep_ProgressMsgDoneEmitsExportDoneMsg(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace) // select json
	updated, _ = sendKey(updated, tea.KeyEnter)

	_, cmd := updated.Update(tui.ExportProgressMsg{
		Format: "json",
		Path:   "/tmp/test-session.json",
		Err:    nil,
	})
	msg := cmdMsg(cmd)
	dm, ok := msg.(tui.ExportDoneMsg)
	if !ok {
		t.Fatalf("after last progress: got msg %T, want ExportDoneMsg", msg)
	}
	if len(dm.OutputPaths) != 1 || dm.OutputPaths[0] != "/tmp/test-session.json" {
		t.Errorf("ExportDoneMsg.OutputPaths = %v, want [/tmp/test-session.json]", dm.OutputPaths)
	}
}

func TestExportStep_PartialProgressDoesNotFinish(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	// Select two formats
	updated, _ := sendKey(m, tea.KeySpace) // json
	updated, _ = sendKey(updated, tea.KeyDown)
	updated, _ = sendKey(updated, tea.KeySpace) // markdown
	updated, _ = sendKey(updated, tea.KeyEnter)

	// Only one progress message received.
	updated, cmd := updated.Update(tui.ExportProgressMsg{
		Format: "json",
		Path:   "/tmp/test-session.json",
	})
	es := updated.(tui.ExportStep)
	if es.IsDone() {
		t.Error("IsDone() should be false when only one of two formats finished")
	}
	if cmd != nil {
		if msg := cmdMsg(cmd); msg != nil {
			if _, ok := msg.(tui.ExportDoneMsg); ok {
				t.Error("should not emit ExportDoneMsg until all selected formats finish")
			}
		}
	}
}

// ── View content by phase ────────────────────────────────────────────────────

func TestExportStep_ViewShowsExportingState(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace)
	updated, _ = sendKey(updated, tea.KeyEnter)
	es := updated.(tui.ExportStep)
	if !strings.Contains(es.View(), "Exporting") {
		t.Error("View() should contain 'Exporting' during export phase")
	}
}

func TestExportStep_ViewShowsDoneState(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := sendKey(m, tea.KeySpace)
	updated, _ = sendKey(updated, tea.KeyEnter)
	updated, _ = updated.Update(tui.ExportProgressMsg{
		Format: "json",
		Path:   "/tmp/test-session.json",
	})
	es := updated.(tui.ExportStep)
	v := es.View()
	if !strings.Contains(v, "complete") {
		t.Errorf("View() should contain 'complete' in done phase, got: %q", v)
	}
	if !strings.Contains(v, "/tmp/test-session.json") {
		t.Error("View() should show output path in done phase")
	}
}

// ── Window resize ────────────────────────────────────────────────────────────

func TestExportStep_WindowResizeUpdatesModel(t *testing.T) {
	m := tui.NewExportStep(nil, "", "test-session", 80, 24)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	es := updated.(tui.ExportStep)
	if es.View() == "" {
		t.Error("View() empty after WindowSizeMsg")
	}
}
