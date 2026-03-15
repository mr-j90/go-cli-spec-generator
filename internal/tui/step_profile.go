package tui

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zyx-holdings/go-spec/internal/questions"
)

// profileConfirmedMsg is sent when the user selects a profile and presses Enter.
type profileConfirmedMsg struct {
	profileID string
}

// ProfileModel is the Bubble Tea model for Step 1: Profile Selection.
type ProfileModel struct {
	profiles []questions.Profile
	cursor   int
}

// NewProfileModel creates a ProfileModel with profiles in a stable sorted order.
func NewProfileModel() ProfileModel {
	var profiles []questions.Profile
	for _, p := range questions.Profiles {
		profiles = append(profiles, p)
	}
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].ID < profiles[j].ID
	})
	return ProfileModel{profiles: profiles}
}

// Cursor returns the current cursor position.
func (m ProfileModel) Cursor() int { return m.cursor }

// SelectedID returns the profile ID at the current cursor position.
func (m ProfileModel) SelectedID() string {
	if len(m.profiles) == 0 {
		return ""
	}
	return m.profiles[m.cursor].ID
}

func (m ProfileModel) Init() tea.Cmd { return nil }

// Update handles keyboard input for the profile selection step.
func (m ProfileModel) Update(msg tea.Msg) (ProfileModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			id := m.profiles[m.cursor].ID
			return m, func() tea.Msg { return profileConfirmedMsg{profileID: id} }
		}
	}
	return m, nil
}

// View renders the profile selection step.
func (m ProfileModel) View() string {
	s := titleStyle.Render("What kind of CLI are you building?") + "\n\n"
	for i, p := range m.profiles {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
		}

		var label string
		if i == m.cursor {
			label = selectedStyle.Render(p.DisplayName)
		} else {
			label = normalStyle.Render(p.DisplayName)
		}
		desc := descStyle.Render("    " + p.Description)
		s += fmt.Sprintf("%s%s\n%s\n\n", cursor, label, desc)
	}
	s += helpStyle.Render("↑/↓ navigate  •  enter select  •  ctrl+c quit")
	return s
}
