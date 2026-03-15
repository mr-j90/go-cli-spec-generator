package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zyx-holdings/go-spec/internal/questions"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// StepQuestionsDoneMsg is sent when the user has answered all questions.
type StepQuestionsDoneMsg struct{}

// StepQuestionsSavedMsg is sent when the session is saved via Ctrl+S.
// Path is the file the session was written to.
type StepQuestionsSavedMsg struct {
	Path string
}

var (
	sqHeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	sqSectionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	sqProgressFilled = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render
	sqProgressEmpty  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	sqQuestionStyle  = lipgloss.NewStyle().Bold(true).MarginTop(1)
	sqRequiredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	sqErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	sqHintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	sqFooterStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
)

// StepQuestions is the Bubble Tea model for Step 3: Question Flow.
// It shows one question at a time and handles navigation, validation, and
// session saving. All methods use pointer receivers so that the model returned
// from Update() can be type-asserted back to *StepQuestions in tests and callers.
type StepQuestions struct {
	store    *session.Store
	qs       []questions.Question
	index    int
	widget   InputWidget
	validErr string

	// confirmQuit is set to true when the user presses Ctrl+C,
	// displaying a "save and exit?" overlay.
	confirmQuit bool

	// savePath is the file path used when the user presses Ctrl+S.
	savePath string

	width  int
	height int
}

// NewStepQuestions creates a new Step 3 model.
// selectedFeatures is the ordered list of feature areas to include after _core.
// savePath is where the session JSON is written on Ctrl+S (defaults to "specgen-session.json").
func NewStepQuestions(store *session.Store, selectedFeatures []string, savePath string) *StepQuestions {
	qs := buildOrderedQuestions(selectedFeatures)
	sq := &StepQuestions{
		store:    store,
		qs:       qs,
		index:    0,
		savePath: savePath,
		width:    80,
	}
	if len(qs) > 0 {
		sq.widget = newWidgetForQuestion(qs[0], sq.width)
		sq.prefillWidget()
	}
	return sq
}

// Init returns the focus command for the first widget (satisfies tea.Model).
func (m *StepQuestions) Init() tea.Cmd {
	if len(m.qs) == 0 {
		return func() tea.Msg { return StepQuestionsDoneMsg{} }
	}
	return m.widget.Focus()
}

// Update handles messages and coordinates navigation, validation, and saving.
func (m *StepQuestions) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.confirmQuit {
			return m.handleConfirmQuit(msg)
		}

		switch {
		case msg.Type == tea.KeyCtrlC:
			m.confirmQuit = true
			return m, nil

		case msg.Type == tea.KeyCtrlS:
			return m.saveAndQuit()

		case msg.Type == tea.KeyEsc:
			return m.goBack()

		case msg.Type == tea.KeyTab:
			return m.skipOrHint()

		case msg.Type == tea.KeyEnter && !msg.Alt:
			// Enter submits for all widget types (including textarea).
			return m.submit()

		case msg.Type == tea.KeyEnter && msg.Alt:
			// Alt+Enter inserts a newline in textarea widgets.
			m.widget.InsertNewline()
			return m, nil
		}
	}

	// Forward all other messages to the current widget.
	var cmd tea.Cmd
	m.widget, cmd = m.widget.Update(msg)
	return m, cmd
}

// View renders the current question screen (satisfies tea.Model).
func (m *StepQuestions) View() string {
	if len(m.qs) == 0 {
		return "No questions to display.\n"
	}

	if m.confirmQuit {
		return m.confirmQuitView()
	}

	q := m.qs[m.index]

	var parts []string

	// ── Header ────────────────────────────────────────────────────────────────
	parts = append(parts, m.headerView(q))

	// ── Question text ─────────────────────────────────────────────────────────
	label := sqQuestionStyle.Render(q.Text)
	if q.Required {
		label += " " + sqRequiredStyle.Render("*")
	}
	parts = append(parts, label)

	// ── Input widget ──────────────────────────────────────────────────────────
	parts = append(parts, m.widget.View())

	// ── Validation error / usage hint ─────────────────────────────────────────
	if m.validErr != "" {
		parts = append(parts, sqErrorStyle.Render("✗ "+m.validErr))
	} else if q.InputType == "multi_select" {
		parts = append(parts, sqHintStyle.Render("Space to toggle  ↑/↓ to move"))
	}

	// ── Footer ────────────────────────────────────────────────────────────────
	parts = append(parts, m.footerView(q))

	return strings.Join(parts, "\n")
}

// ─── view helpers ─────────────────────────────────────────────────────────────

func (m *StepQuestions) headerView(q questions.Question) string {
	section := sectionDisplayName(q.FeatureArea)
	total := len(m.qs)
	current := m.index + 1

	const barWidth = 20
	filled := 0
	if total > 0 {
		filled = (current * barWidth) / total
	}
	bar := sqProgressFilled(strings.Repeat("█", filled)) +
		sqProgressEmpty(strings.Repeat("░", barWidth-filled))

	pct := 0
	if total > 0 {
		pct = (current * 100) / total
	}

	left := sqSectionStyle.Render(section)
	right := sqHeaderStyle.Render(fmt.Sprintf("Q %d of %d", current, total))
	progress := fmt.Sprintf("[%s] %d%%", bar, pct)

	return left + "  " + right + "  " + progress
}

func (m *StepQuestions) footerView(q questions.Question) string {
	hints := []string{"Enter: submit", "Esc: back"}
	if !q.Required {
		hints = append(hints, "Tab: skip")
	}
	hints = append(hints, "Ctrl+S: save & quit")
	if q.InputType == "textarea" {
		hints = append(hints, "Alt+Enter: newline")
	}
	return sqFooterStyle.Render(strings.Join(hints, "  ·  "))
}

func (m *StepQuestions) confirmQuitView() string {
	return sqQuestionStyle.Render("Save session and exit?") + "\n" +
		"  " + sqHintStyle.Render("Y: save and quit   N: continue   Ctrl+C: quit without saving") + "\n"
}

// ─── action helpers ───────────────────────────────────────────────────────────

// submit validates and records the current answer, then advances.
func (m *StepQuestions) submit() (tea.Model, tea.Cmd) {
	q := m.qs[m.index]
	if q.Required && m.widget.IsEmpty() {
		m.validErr = "This field is required"
		return m, nil
	}
	m.validErr = ""
	m.recordAnswer(q)
	return m.advance()
}

// goBack moves to the previous question or quits if already on the first.
func (m *StepQuestions) goBack() (tea.Model, tea.Cmd) {
	if m.index == 0 {
		return m, tea.Quit
	}
	m.index--
	m.validErr = ""
	m.widget = newWidgetForQuestion(m.qs[m.index], m.width)
	m.prefillWidget()
	return m, m.widget.Focus()
}

// skipOrHint skips the current optional question or shows a hint for required ones.
func (m *StepQuestions) skipOrHint() (tea.Model, tea.Cmd) {
	q := m.qs[m.index]
	if q.Required {
		m.validErr = "This field is required and cannot be skipped"
		return m, nil
	}
	m.validErr = ""
	m.store.SkipAnswer(q.ID)
	return m.advance()
}

// advance moves to the next question or emits StepQuestionsDoneMsg.
func (m *StepQuestions) advance() (tea.Model, tea.Cmd) {
	m.index++
	if m.index >= len(m.qs) {
		return m, func() tea.Msg { return StepQuestionsDoneMsg{} }
	}
	m.widget = newWidgetForQuestion(m.qs[m.index], m.width)
	m.prefillWidget()
	return m, m.widget.Focus()
}

// saveAndQuit persists the session and emits a quit command.
func (m *StepQuestions) saveAndQuit() (tea.Model, tea.Cmd) {
	path := m.savePath
	if path == "" {
		path = "specgen-session.json"
	}
	if err := m.store.Save(path); err != nil {
		m.validErr = fmt.Sprintf("Save failed: %v", err)
		return m, nil
	}
	savedPath := path
	return m, tea.Batch(
		func() tea.Msg { return StepQuestionsSavedMsg{Path: savedPath} },
		tea.Quit,
	)
}

// handleConfirmQuit processes keypresses inside the "save and exit?" overlay.
func (m *StepQuestions) handleConfirmQuit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyCtrlC:
		return m, tea.Quit
	case msg.Type == tea.KeyEsc:
		m.confirmQuit = false
		return m, nil
	case msg.Type == tea.KeyRunes:
		switch strings.ToLower(string(msg.Runes)) {
		case "y":
			return m.saveAndQuit()
		case "n":
			m.confirmQuit = false
			return m, nil
		}
	}
	return m, nil
}

// ─── state helpers ────────────────────────────────────────────────────────────

func (m *StepQuestions) recordAnswer(q questions.Question) {
	if q.InputType == "multi_select" {
		m.store.SetAnswer(q.ID, session.NewMultiValue(m.widget.Values()))
	} else {
		m.store.SetAnswer(q.ID, session.NewStringValue(m.widget.Value()))
	}
}

func (m *StepQuestions) prefillWidget() {
	if m.index >= len(m.qs) {
		return
	}
	q := m.qs[m.index]
	ans, ok := m.store.GetAnswer(q.ID)
	if !ok || ans.Skipped {
		return
	}
	if q.InputType == "multi_select" {
		m.widget.SetValues(ans.Value.Strings())
	} else {
		m.widget.SetValue(ans.Value.String())
	}
}

// ─── package-level helpers ────────────────────────────────────────────────────

// buildOrderedQuestions returns _core questions followed by each selected feature area.
func buildOrderedQuestions(selectedFeatures []string) []questions.Question {
	var qs []questions.Question
	qs = append(qs, questions.ByFeatureArea("_core")...)
	for _, f := range selectedFeatures {
		qs = append(qs, questions.ByFeatureArea(f)...)
	}
	return qs
}

// newWidgetForQuestion creates an InputWidget configured for the given question.
func newWidgetForQuestion(q questions.Question, width int) InputWidget {
	return NewInputWidget(q.InputType, q.Options, q.Placeholder, width)
}

// sectionDisplayName returns a human-readable name for a feature area key.
func sectionDisplayName(featureArea string) string {
	if featureArea == "_core" {
		return "Core"
	}
	if fa, ok := questions.FeatureAreas[featureArea]; ok {
		return fa.DisplayName
	}
	return featureArea
}
