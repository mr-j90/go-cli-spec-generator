// Package tui implements the interactive terminal UI using Bubble Tea.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// tuiStep represents the current active step in the wizard.
type tuiStep int

const (
	stepProfile  tuiStep = iota
	stepFeatures         // future steps follow
)

// App is the root Bubble Tea model that manages step navigation across all steps.
type App struct {
	step     tuiStep
	profile  ProfileModel
	features FeaturesModel
	quitting bool
}

// Result holds the final selections produced by the TUI.
type Result struct {
	ProfileID  string
	FeatureIDs []string
}

// New creates a new App. Pass noColor=true to disable ANSI color output.
func New(noColor bool) App {
	initStyles(noColor)
	return App{
		step:     stepProfile,
		profile:  NewProfileModel(),
		features: NewFeaturesModel(),
	}
}

func (a App) Init() tea.Cmd {
	return a.profile.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global quit on Ctrl+C regardless of step.
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyCtrlC {
		a.quitting = true
		return a, tea.Quit
	}

	switch msg := msg.(type) {
	case profileConfirmedMsg:
		_ = msg
		a.step = stepFeatures
		return a, nil

	case featuresConfirmedMsg:
		_ = msg
		// Steps 3-5 are future work; finish for now.
		a.quitting = true
		return a, tea.Quit

	case goBackMsg:
		if a.step > stepProfile {
			a.step--
		}
		return a, nil

	case tea.WindowSizeMsg:
		return a, nil
	}

	// Delegate to the active step model.
	var cmd tea.Cmd
	switch a.step {
	case stepProfile:
		a.profile, cmd = a.profile.Update(msg)
	case stepFeatures:
		a.features, cmd = a.features.Update(msg)
	}
	return a, cmd
}

func (a App) View() string {
	if a.quitting {
		return ""
	}
	switch a.step {
	case stepProfile:
		return a.profile.View()
	case stepFeatures:
		return a.features.View()
	}
	return ""
}

// Run launches the TUI program and returns the user's final selections.
func Run(noColor bool) (Result, error) {
	app := New(noColor)
	p := tea.NewProgram(app)
	finalModel, err := p.Run()
	if err != nil {
		return Result{}, err
	}
	finalApp, ok := finalModel.(App)
	if !ok {
		return Result{}, nil
	}
	return Result{
		ProfileID:  finalApp.profile.SelectedID(),
		FeatureIDs: finalApp.features.Selected(),
	}, nil
}
