package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	widgetCursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)
	widgetCheckedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	widgetUncheckedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	widgetHintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

// InputWidget is a unified input component for all question input types.
// It supports text, textarea, single_select, multi_select, and confirm.
type InputWidget struct {
	inputType string

	// text / textarea
	text textinput.Model
	area textarea.Model

	// single_select / multi_select
	options  []string
	cursor   int
	selected map[int]bool

	// confirm
	confirmVal bool
	confirmSet bool
}

// NewInputWidget creates an InputWidget for the given input type, options, and placeholder.
// width is the available terminal width for sizing inner components.
func NewInputWidget(inputType string, options []string, placeholder string, width int) InputWidget {
	w := InputWidget{
		inputType: inputType,
		options:   options,
		selected:  make(map[int]bool),
	}

	innerWidth := width - 4
	if innerWidth < 40 {
		innerWidth = 40
	}

	switch inputType {
	case "text":
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.Width = innerWidth
		_ = ti.Focus() // focus immediately so keystrokes are accepted without Init
		w.text = ti

	case "textarea":
		ta := textarea.New()
		ta.Placeholder = placeholder
		ta.SetWidth(innerWidth)
		ta.SetHeight(4)
		ta.ShowLineNumbers = false
		ta.Focus() // focus immediately so keystrokes are accepted without Init
		w.area = ta
	}

	return w
}

// Focus focuses the widget and returns the initialization command (e.g. cursor blink).
func (w *InputWidget) Focus() tea.Cmd {
	switch w.inputType {
	case "text":
		return w.text.Focus()
	case "textarea":
		w.area.Focus()
		return textarea.Blink
	}
	return nil
}

// Blur removes focus from the widget.
func (w *InputWidget) Blur() {
	switch w.inputType {
	case "text":
		w.text.Blur()
	case "textarea":
		w.area.Blur()
	}
}

// Value returns the current answer as a single string.
// For multi_select it returns a comma-separated list of selected options.
// Returns empty string if no selection has been made.
func (w *InputWidget) Value() string {
	switch w.inputType {
	case "text":
		return strings.TrimSpace(w.text.Value())
	case "textarea":
		return strings.TrimSpace(w.area.Value())
	case "single_select":
		if len(w.options) == 0 {
			return ""
		}
		return w.options[w.cursor]
	case "multi_select":
		var vals []string
		for i, opt := range w.options {
			if w.selected[i] {
				vals = append(vals, opt)
			}
		}
		return strings.Join(vals, ", ")
	case "confirm":
		if !w.confirmSet {
			return ""
		}
		if w.confirmVal {
			return "yes"
		}
		return "no"
	}
	return ""
}

// Values returns the selected options as a slice (multi_select only).
// Returns nil for all other input types.
func (w *InputWidget) Values() []string {
	if w.inputType != "multi_select" {
		return nil
	}
	var vals []string
	for i, opt := range w.options {
		if w.selected[i] {
			vals = append(vals, opt)
		}
	}
	return vals
}

// IsEmpty reports whether the widget has no meaningful value entered.
func (w *InputWidget) IsEmpty() bool {
	switch w.inputType {
	case "text":
		return strings.TrimSpace(w.text.Value()) == ""
	case "textarea":
		return strings.TrimSpace(w.area.Value()) == ""
	case "single_select":
		// A cursor position always represents a selection for non-empty option lists.
		return len(w.options) == 0
	case "multi_select":
		return len(w.Values()) == 0
	case "confirm":
		return !w.confirmSet
	}
	return true
}

// SetValue pre-fills the widget from a previously stored string answer.
// Used when navigating back to a question.
func (w *InputWidget) SetValue(val string) {
	switch w.inputType {
	case "text":
		w.text.SetValue(val)
	case "textarea":
		w.area.SetValue(val)
	case "single_select":
		for i, opt := range w.options {
			if opt == val {
				w.cursor = i
				return
			}
		}
	case "confirm":
		switch strings.ToLower(val) {
		case "yes", "y", "true":
			w.confirmVal = true
			w.confirmSet = true
		case "no", "n", "false":
			w.confirmVal = false
			w.confirmSet = true
		}
	}
}

// SetValues pre-fills a multi_select widget from a previously stored slice.
func (w *InputWidget) SetValues(vals []string) {
	if w.inputType != "multi_select" {
		return
	}
	set := make(map[string]bool, len(vals))
	for _, v := range vals {
		set[v] = true
	}
	for i, opt := range w.options {
		w.selected[i] = set[opt]
	}
}

// InsertNewline appends a newline to a textarea widget.
// This is called by StepQuestions when Alt+Enter is pressed.
func (w *InputWidget) InsertNewline() {
	if w.inputType == "textarea" {
		w.area.SetValue(w.area.Value() + "\n")
	}
}

// Update processes a key message for the widget.
// Enter, Esc, Tab, Ctrl+S, and Ctrl+C are NOT handled here — they are
// intercepted by StepQuestions before being forwarded to the widget.
func (w InputWidget) Update(msg tea.Msg) (InputWidget, tea.Cmd) {
	var cmd tea.Cmd
	switch w.inputType {
	case "text":
		w.text, cmd = w.text.Update(msg)

	case "textarea":
		w.area, cmd = w.area.Update(msg)

	case "single_select":
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.Type {
			case tea.KeyUp:
				if w.cursor > 0 {
					w.cursor--
				}
			case tea.KeyDown:
				if w.cursor < len(w.options)-1 {
					w.cursor++
				}
			}
		}

	case "multi_select":
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.Type {
			case tea.KeyUp:
				if w.cursor > 0 {
					w.cursor--
				}
			case tea.KeyDown:
				if w.cursor < len(w.options)-1 {
					w.cursor++
				}
			case tea.KeySpace:
				w.selected[w.cursor] = !w.selected[w.cursor]
			}
		}

	case "confirm":
		if km, ok := msg.(tea.KeyMsg); ok {
			if km.Type == tea.KeyRunes {
				switch strings.ToLower(string(km.Runes)) {
				case "y":
					w.confirmVal = true
					w.confirmSet = true
				case "n":
					w.confirmVal = false
					w.confirmSet = true
				}
			}
		}
	}
	return w, cmd
}

// View renders the widget to a string.
func (w InputWidget) View() string {
	switch w.inputType {
	case "text":
		return w.text.View()

	case "textarea":
		return w.area.View()

	case "single_select":
		var sb strings.Builder
		for i, opt := range w.options {
			if i == w.cursor {
				sb.WriteString(widgetCursorStyle.Render("▶ " + opt))
			} else {
				sb.WriteString("  " + opt)
			}
			if i < len(w.options)-1 {
				sb.WriteByte('\n')
			}
		}
		return sb.String()

	case "multi_select":
		var sb strings.Builder
		for i, opt := range w.options {
			check := widgetUncheckedStyle.Render("[ ]")
			if w.selected[i] {
				check = widgetCheckedStyle.Render("[✓]")
			}
			prefix := "  "
			if i == w.cursor {
				prefix = widgetCursorStyle.Render("▶ ")
			}
			sb.WriteString(prefix + check + " " + opt)
			if i < len(w.options)-1 {
				sb.WriteByte('\n')
			}
		}
		return sb.String()

	case "confirm":
		yesLabel := "[ Y ]"
		noLabel := "[ N ]"
		if w.confirmSet {
			if w.confirmVal {
				yesLabel = widgetCursorStyle.Render("[Y]")
			} else {
				noLabel = widgetCursorStyle.Render("[N]")
			}
		}
		return yesLabel + "  " + noLabel + "  " + widgetHintStyle.Render("(press Y or N)")
	}
	return ""
}
