package tui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zyx-holdings/go-spec/internal/tui"
)

// ---- App (top-level model) ----

func TestApp_New_ViewContainsProfileTitle(t *testing.T) {
	app := tui.New(true)
	v := app.View()
	if !strings.Contains(v, "What kind of CLI are you building?") {
		t.Errorf("View() = %q, want profile title", v)
	}
}

func TestApp_Init_ReturnsNilCmd(t *testing.T) {
	app := tui.New(true)
	// Profile step Init returns nil; App.Init mirrors it.
	cmd := app.Init()
	if cmd != nil {
		// Execute the cmd — it must not produce a QuitMsg.
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Error("Init() returned tea.Quit, want non-quit cmd")
		}
	}
}

func TestApp_CtrlC_Quits(t *testing.T) {
	app := tui.New(true)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("Update(Ctrl+C) returned nil cmd, want tea.Quit")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("Update(Ctrl+C) cmd() = %T, want tea.QuitMsg", msg)
	}
}

func TestApp_NonQuitKey_DoesNotQuit(t *testing.T) {
	app := tui.New(true)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Error("Update(rune 'a') returned tea.Quit, want no quit")
		}
	}
}

func TestApp_Update_ReturnsModel(t *testing.T) {
	app := tui.New(true)
	updated, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if updated == nil {
		t.Error("Update() returned nil model")
	}
}

func TestApp_NoColor_ViewIsNonEmpty(t *testing.T) {
	app := tui.New(true)
	if app.View() == "" {
		t.Error("View() returned empty string with noColor=true")
	}
}

// ---- ProfileModel (Step 1) ----

func TestProfileModel_InitialView_ContainsTitle(t *testing.T) {
	m := tui.NewProfileModel()
	v := m.View()
	if !strings.Contains(v, "What kind of CLI are you building?") {
		t.Errorf("ProfileModel.View() missing title, got: %q", v)
	}
}

func TestProfileModel_InitialView_ContainsProfiles(t *testing.T) {
	m := tui.NewProfileModel()
	v := m.View()
	for _, name := range []string{"API Service", "CLI Tool", "Library", "Web Application"} {
		if !strings.Contains(v, name) {
			t.Errorf("ProfileModel.View() missing profile %q", name)
		}
	}
}

func TestProfileModel_ArrowDown_MovesCursor(t *testing.T) {
	m := tui.NewProfileModel()
	before := m.Cursor()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	after := m.Cursor()
	if after != before+1 {
		t.Errorf("cursor after KeyDown = %d, want %d", after, before+1)
	}
}

func TestProfileModel_ArrowUp_AtTop_StaysPut(t *testing.T) {
	m := tui.NewProfileModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("cursor after KeyUp at top = %d, want 0", m.Cursor())
	}
}

func TestProfileModel_ArrowDown_AtBottom_StaysPut(t *testing.T) {
	m := tui.NewProfileModel()
	// Drive cursor to the last item.
	for i := 0; i < 20; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	cursorAtEnd := m.Cursor()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.Cursor() != cursorAtEnd {
		t.Errorf("cursor past bottom = %d, want %d", m.Cursor(), cursorAtEnd)
	}
}

func TestProfileModel_Enter_EmitsConfirmedMsg(t *testing.T) {
	m := tui.NewProfileModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter returned nil cmd, want profileConfirmedMsg")
	}
	// The App handles the message; we can't inspect the unexported type, but
	// we can verify the cmd executes without panicking and returns a non-nil msg.
	msg := cmd()
	if msg == nil {
		t.Error("profileConfirmedMsg cmd() returned nil")
	}
}

func TestProfileModel_SelectedID_NonEmpty(t *testing.T) {
	m := tui.NewProfileModel()
	if m.SelectedID() == "" {
		t.Error("SelectedID() returned empty string on fresh model")
	}
}

// ---- FeaturesModel (Step 2) ----

func TestFeaturesModel_InitialView_ContainsTitle(t *testing.T) {
	m := tui.NewFeaturesModel()
	v := m.View()
	if !strings.Contains(v, "Which feature areas apply to your CLI?") {
		t.Errorf("FeaturesModel.View() missing title, got: %q", v)
	}
}

func TestFeaturesModel_InitialView_ContainsAllFeatureAreas(t *testing.T) {
	m := tui.NewFeaturesModel()
	v := m.View()
	for _, name := range []string{
		"Authentication", "Storage", "API", "Testing",
		"Observability", "Deployment", "Security", "Caching",
		"Messaging", "Search", "Notifications", "Configuration",
	} {
		if !strings.Contains(v, name) {
			t.Errorf("FeaturesModel.View() missing feature area %q", name)
		}
	}
}

func TestFeaturesModel_SpaceTogglesItem(t *testing.T) {
	m := tui.NewFeaturesModel()
	id := m.SelectedID()
	if m.IsChecked(id) {
		t.Fatalf("item %q unexpectedly pre-checked", id)
	}
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !m.IsChecked(id) {
		t.Errorf("item %q not checked after Space", id)
	}
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if m.IsChecked(id) {
		t.Errorf("item %q still checked after second Space (should toggle off)", id)
	}
}

func TestFeaturesModel_SpaceRune_TogglesItem(t *testing.T) {
	m := tui.NewFeaturesModel()
	id := m.SelectedID()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	if !m.IsChecked(id) {
		t.Errorf("item %q not checked after space rune", id)
	}
}

func TestFeaturesModel_Enter_WithNoSelections_Confirms(t *testing.T) {
	m := tui.NewFeaturesModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter returned nil cmd")
	}
	msg := cmd()
	if msg == nil {
		t.Error("featuresConfirmedMsg cmd() returned nil")
	}
}

func TestFeaturesModel_Enter_WithSelections_IncludesSelected(t *testing.T) {
	m := tui.NewFeaturesModel()
	// Toggle the first item on.
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	firstID := m.SelectedID() // still on first item
	selected := m.Selected()
	if len(selected) != 1 || selected[0] != firstID {
		t.Errorf("Selected() = %v, want [%q]", selected, firstID)
	}
}

func TestFeaturesModel_Esc_EmitsGoBack(t *testing.T) {
	m := tui.NewFeaturesModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Esc returned nil cmd, want goBackMsg")
	}
	msg := cmd()
	if msg == nil {
		t.Error("goBackMsg cmd() returned nil")
	}
}

func TestFeaturesModel_Backspace_EmitsGoBack(t *testing.T) {
	m := tui.NewFeaturesModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace returned nil cmd, want goBackMsg")
	}
	msg := cmd()
	if msg == nil {
		t.Error("goBackMsg cmd() returned nil")
	}
}

func TestFeaturesModel_ArrowNavigation(t *testing.T) {
	m := tui.NewFeaturesModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.Cursor() != 1 {
		t.Errorf("cursor after KeyDown = %d, want 1", m.Cursor())
	}
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("cursor after KeyUp = %d, want 0", m.Cursor())
	}
}

// ---- App step transitions ----

func TestApp_EnterOnProfile_AdvancesToFeatures(t *testing.T) {
	app := tui.New(true)
	// Simulate Enter on the profile step — produces a profileConfirmedMsg cmd.
	app2, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on profile step returned nil cmd")
	}
	// Execute the cmd to get the profileConfirmedMsg, then feed it back.
	msg := cmd()
	app3, _ := app2.Update(msg)
	v := app3.(interface{ View() string }).View()
	if !strings.Contains(v, "Which feature areas apply to your CLI?") {
		t.Errorf("after profile confirmed, View() = %q, want features title", v)
	}
}

func TestApp_EscOnFeatures_GoesBackToProfile(t *testing.T) {
	app := tui.New(true)
	// Advance to features step.
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	msg := cmd()
	app2, _ := app.Update(msg)

	// Now press Esc on features step.
	app3, cmd2 := app2.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd2 != nil {
		msg2 := cmd2()
		app3, _ = app3.Update(msg2)
	}
	v := app3.(interface{ View() string }).View()
	if !strings.Contains(v, "What kind of CLI are you building?") {
		t.Errorf("after Esc on features, View() = %q, want profile title", v)
	}
}
