package tui

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zyx-holdings/go-spec/internal/questions"
)

// featuresConfirmedMsg is sent when the user confirms feature area selections.
type featuresConfirmedMsg struct {
	featureIDs []string
}

// goBackMsg is sent when the user navigates back to the previous step.
type goBackMsg struct{}

// FeaturesModel is the Bubble Tea model for Step 2: Feature Area Selection.
type FeaturesModel struct {
	areas   []questions.FeatureArea
	cursor  int
	checked map[string]bool
}

// NewFeaturesModel creates a FeaturesModel with all feature areas in stable sorted order.
func NewFeaturesModel() FeaturesModel {
	var areas []questions.FeatureArea
	for _, fa := range questions.FeatureAreas {
		areas = append(areas, fa)
	}
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].ID < areas[j].ID
	})
	return FeaturesModel{
		areas:   areas,
		checked: make(map[string]bool),
	}
}

// Cursor returns the current cursor position.
func (m FeaturesModel) Cursor() int { return m.cursor }

// SelectedID returns the feature area ID at the current cursor position.
func (m FeaturesModel) SelectedID() string {
	if len(m.areas) == 0 {
		return ""
	}
	return m.areas[m.cursor].ID
}

// IsChecked reports whether the feature area with the given ID is selected.
func (m FeaturesModel) IsChecked(id string) bool { return m.checked[id] }

// Selected returns the IDs of all checked feature areas in display order.
func (m FeaturesModel) Selected() []string {
	var selected []string
	for _, a := range m.areas {
		if m.checked[a.ID] {
			selected = append(selected, a.ID)
		}
	}
	return selected
}

func (m FeaturesModel) Init() tea.Cmd { return nil }

// Update handles keyboard input for the feature area selection step.
func (m FeaturesModel) Update(msg tea.Msg) (FeaturesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.areas)-1 {
				m.cursor++
			}
		case tea.KeyRunes:
			if string(msg.Runes) == " " {
				id := m.areas[m.cursor].ID
				m.checked[id] = !m.checked[id]
			}
		case tea.KeySpace:
			id := m.areas[m.cursor].ID
			m.checked[id] = !m.checked[id]
		case tea.KeyEnter:
			selected := m.Selected()
			return m, func() tea.Msg { return featuresConfirmedMsg{featureIDs: selected} }
		case tea.KeyEsc, tea.KeyBackspace:
			return m, func() tea.Msg { return goBackMsg{} }
		}
	}
	return m, nil
}

// View renders the feature area multi-select step.
func (m FeaturesModel) View() string {
	s := titleStyle.Render("Which feature areas apply to your CLI?") + "\n\n"
	for i, fa := range m.areas {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
		}

		var checkbox string
		if m.checked[fa.ID] {
			checkbox = checkboxChecked.Render("[✓]")
		} else {
			checkbox = checkboxUnchecked.Render("[ ]")
		}

		icon := ""
		if fa.Icon != "" {
			icon = fa.Icon + " "
		}

		var label string
		if i == m.cursor {
			label = selectedStyle.Render(icon + fa.DisplayName)
		} else {
			label = normalStyle.Render(icon + fa.DisplayName)
		}
		desc := descStyle.Render("       " + fa.Description)
		s += fmt.Sprintf("%s%s %s\n%s\n\n", cursor, checkbox, label, desc)
	}
	s += helpStyle.Render("↑/↓ navigate  •  space toggle  •  enter confirm  •  esc/backspace back  •  ctrl+c quit")
	return s
}
