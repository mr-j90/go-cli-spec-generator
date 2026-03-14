package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zyx-holdings/go-spec/internal/questions"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// ReviewDoneMsg is sent when the user confirms the review and proceeds to export.
type ReviewDoneMsg struct{}

// ReviewBackMsg is sent when the user navigates back to the questions step.
type ReviewBackMsg struct{}

// ReviewJumpEditMsg is sent when the user requests jump-edit for a question.
type ReviewJumpEditMsg struct {
	// QuestionIndex is the 1-based index into the flat question list.
	QuestionIndex int
}

type reviewSection struct {
	featureArea string
	displayName string
	questions   []questions.Question
	collapsed   bool
}

// ReviewStep is the Bubble Tea model for Step 4: Review.
type ReviewStep struct {
	store     *session.Store
	sections  []reviewSection
	viewport  viewport.Model
	// cursor: 0..len(sections)-1 = sections; len(sections) = Export; len(sections)+1 = Back
	cursor    int
	width     int
	height    int
	jumpEdit  bool   // true when 'e' key was pressed and we're awaiting a number
	numBuffer string // accumulated digits for jump-edit
	ready     bool
}

var (
	reviewHeaderStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	reviewCursorMarker     = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	reviewQuestionStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	reviewAnswerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	reviewSkippedStyle     = lipgloss.NewStyle().Faint(true).Italic(true)
	reviewButtonStyle      = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	reviewActiveButtonStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1).
				Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	reviewHintStyle = lipgloss.NewStyle().Faint(true)
)

// NewReviewStep creates a ReviewStep for the given session store.
// width and height are the terminal dimensions (0 values are safe; a minimum
// viewport height is applied automatically).
func NewReviewStep(store *session.Store, width, height int) ReviewStep {
	s := ReviewStep{
		store:  store,
		width:  width,
		height: height,
	}
	s.buildSections()
	vph := s.viewportHeight()
	if vph < 1 {
		vph = 1
	}
	s.viewport = viewport.New(width, vph)
	s.viewport.SetContent(s.renderContent())
	s.ready = true
	return s
}

func (s *ReviewStep) buildSections() {
	sess := s.store.Session()
	s.sections = []reviewSection{
		{
			featureArea: "_core",
			displayName: "Core",
			questions:   questions.ByFeatureArea("_core"),
		},
	}
	for _, f := range sess.SelectedFeatures {
		fa, ok := questions.FeatureAreas[f]
		if !ok {
			continue
		}
		s.sections = append(s.sections, reviewSection{
			featureArea: f,
			displayName: fa.DisplayName,
			questions:   questions.ByFeatureArea(f),
		})
	}
}

// Init implements tea.Model.
func (s ReviewStep) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (s ReviewStep) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKey(msg)
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		if s.ready {
			s.viewport.Width = msg.Width
			vph := s.viewportHeight()
			if vph < 1 {
				vph = 1
			}
			s.viewport.Height = vph
		}
		s.viewport.SetContent(s.renderContent())
		return s, nil
	}
	// Forward other messages (e.g. mouse scroll) to the viewport.
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s ReviewStep) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Jump-edit mode: accumulate digits after 'e', confirm with Enter.
	if s.jumpEdit {
		switch msg.Type {
		case tea.KeyEsc:
			s.jumpEdit = false
			s.numBuffer = ""
		case tea.KeyEnter:
			if s.numBuffer != "" {
				var idx int
				fmt.Sscanf(s.numBuffer, "%d", &idx)
				s.jumpEdit = false
				s.numBuffer = ""
				return s, func() tea.Msg { return ReviewJumpEditMsg{QuestionIndex: idx} }
			}
			s.jumpEdit = false
			s.numBuffer = ""
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if r >= '0' && r <= '9' {
					s.numBuffer += string(r)
				}
			}
		}
		return s, nil
	}

	totalItems := len(s.sections) + 2
	exportIdx := len(s.sections)
	backIdx := len(s.sections) + 1

	switch msg.Type {
	case tea.KeyUp:
		if s.cursor > 0 {
			s.cursor--
			s.viewport.SetContent(s.renderContent())
		}
		return s, nil

	case tea.KeyDown:
		if s.cursor < totalItems-1 {
			s.cursor++
			s.viewport.SetContent(s.renderContent())
		}
		return s, nil

	case tea.KeyEnter:
		switch {
		case s.cursor < len(s.sections):
			s.sections[s.cursor].collapsed = !s.sections[s.cursor].collapsed
			s.viewport.SetContent(s.renderContent())
		case s.cursor == exportIdx:
			return s, func() tea.Msg { return ReviewDoneMsg{} }
		case s.cursor == backIdx:
			return s, func() tea.Msg { return ReviewBackMsg{} }
		}
		return s, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			switch r {
			case 'e', 'E':
				s.jumpEdit = true
				s.numBuffer = ""
				return s, nil
			case 'b', 'B':
				return s, func() tea.Msg { return ReviewBackMsg{} }
			}
		}
	}

	// Forward unhandled keys (PgUp/PgDn, etc.) to viewport.
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

// renderContent builds the viewport content string (sections only; buttons are
// rendered outside the viewport in View).
func (s ReviewStep) renderContent() string {
	sess := s.store.Session()
	var sb strings.Builder

	for i, sec := range s.sections {
		// Cursor indicator.
		marker := "  "
		if s.cursor == i {
			marker = reviewCursorMarker.Render("> ")
		}
		toggle := "▼"
		if sec.collapsed {
			toggle = "▶"
		}
		header := reviewHeaderStyle.Render(fmt.Sprintf("[%s] %s", toggle, sec.displayName))
		sb.WriteString(marker + header + "\n")

		if sec.collapsed {
			continue
		}

		offset := s.sectionOffset(i)
		for j, q := range sec.questions {
			qNum := offset + j + 1
			ans, ok := sess.Answers[q.ID]

			var ansLine string
			switch {
			case !ok || ans.Skipped:
				ansLine = reviewSkippedStyle.Render("(skipped)")
			case ans.Value.IsEmpty():
				ansLine = reviewSkippedStyle.Render("(skipped)")
			case ans.Value.IsMulti():
				vals := ans.Value.Strings()
				if len(vals) == 0 {
					ansLine = reviewSkippedStyle.Render("(skipped)")
				} else {
					ansLine = reviewAnswerStyle.Render(strings.Join(vals, ", "))
				}
			default:
				ansLine = reviewAnswerStyle.Render(ans.Value.String())
			}

			qLine := reviewQuestionStyle.Render(fmt.Sprintf("  %d. %s", qNum, q.Text))
			sb.WriteString("    " + qLine + "\n")
			sb.WriteString("       " + ansLine + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// sectionOffset returns the cumulative number of questions before sectionIdx.
func (s ReviewStep) sectionOffset(sectionIdx int) int {
	total := 0
	for i := 0; i < sectionIdx; i++ {
		total += len(s.sections[i].questions)
	}
	return total
}

// viewportHeight calculates the height available for the viewport.
func (s ReviewStep) viewportHeight() int {
	// Title (1) + hint (1) + blank (2) + buttons (1) + padding (1) = 6 reserved lines.
	if s.height > 6 {
		return s.height - 6
	}
	return 10
}

// View implements tea.Model.
func (s ReviewStep) View() string {
	title := titleStyle.Render("Step 4: Review Your Answers")
	hint := reviewHintStyle.Render("↑/↓ navigate  ·  Enter toggle/select  ·  e<N> jump to question  ·  b back")

	var content string
	if s.ready {
		content = s.viewport.View()
	} else {
		content = "Loading...\n"
	}

	exportIdx := len(s.sections)
	backIdx := len(s.sections) + 1

	var exportBtn, backBtn string
	if s.cursor == exportIdx {
		exportBtn = reviewActiveButtonStyle.Render("[ Export ]")
	} else {
		exportBtn = reviewButtonStyle.Render("[ Export ]")
	}
	if s.cursor == backIdx {
		backBtn = reviewActiveButtonStyle.Render("[ Back ]")
	} else {
		backBtn = reviewButtonStyle.Render("[ Back ]")
	}
	buttons := exportBtn + "  " + backBtn

	if s.jumpEdit {
		buttons += "\n" + reviewCursorMarker.Render(fmt.Sprintf("Jump to question: %s_", s.numBuffer))
	}

	return title + "\n" + hint + "\n\n" + content + "\n" + buttons + "\n"
}

// Cursor returns the current cursor index (exported for testing).
func (s ReviewStep) Cursor() int { return s.cursor }

// SectionCount returns the number of sections (exported for testing).
func (s ReviewStep) SectionCount() int { return len(s.sections) }

// IsSectionCollapsed returns whether the section at idx is collapsed.
func (s ReviewStep) IsSectionCollapsed(idx int) bool {
	if idx < 0 || idx >= len(s.sections) {
		return false
	}
	return s.sections[idx].collapsed
}

// IsJumpEditMode reports whether the model is in jump-edit input mode.
func (s ReviewStep) IsJumpEditMode() bool { return s.jumpEdit }
