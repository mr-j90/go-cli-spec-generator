// Package tui implements the interactive terminal UI using Bubble Tea.
package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// titleStyle is the default title style for the TUI.
// TODO: expand styles as the TUI is implemented.
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))

// Model is the root Bubble Tea model for the spec generator TUI.
// TODO: implement full TUI model with questions flow.
type Model struct {
	input textinput.Model
}

// New creates a new TUI model.
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Describe your spec..."
	ti.Focus()
	return Model{input: ti}
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return titleStyle.Render("specgen") + "\n\n" + m.input.View() + "\n"
}
