package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zyx-holdings/go-spec/internal/export"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// ExportDoneMsg is sent when all selected formats have been exported.
type ExportDoneMsg struct {
	// OutputPaths lists the file paths of all successfully generated files.
	OutputPaths []string
}

// ExportProgressMsg is sent when a single format export finishes (success or error).
type ExportProgressMsg struct {
	Format string
	Path   string
	Err    error
}

type exportFormatItem struct {
	id       string
	label    string
	selected bool
}

type exportPhase int

const (
	exportPhaseSelecting exportPhase = iota
	exportPhaseExporting
	exportPhaseDone
)

type exportFormatProgress struct {
	path string
	err  error
	done bool
}

// ExportStep is the Bubble Tea model for Step 5: Export Format Selection.
type ExportStep struct {
	formats         []exportFormatItem
	cursor          int
	phase           exportPhase
	progress        map[string]exportFormatProgress
	outputDir       string
	sessionID       string
	sess            *session.Session
	specgenVersion  string
	skipStep        bool
	width           int
	height          int
}

var (
	exportTitleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	exportHintStyle        = lipgloss.NewStyle().Faint(true)
	exportCursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	exportCheckedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	exportUncheckedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	exportDimStyle         = lipgloss.NewStyle().Faint(true)
	exportErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	exportSuccessStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	exportValidationStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

// NewExportStep creates a new ExportStep.
//
// preselectedFormats: if non-empty, those formats are pre-checked and the
// selection UI is skipped (simulating the --format CLI flag). Valid IDs are
// "json", "markdown", "docx", "pdf".
//
// outputDir: directory to write exported files; "" uses the current directory.
//
// sessionID: used as the base filename for each exported file.
//
// sess: the active session; may be nil (falls back to a stub JSON export).
//
// specgenVersion: embedded in JSON exports as the "specgen_version" field.
func NewExportStep(preselectedFormats []string, outputDir, sessionID string, width, height int, sess *session.Session, specgenVersion string) ExportStep {
	allFormats := []exportFormatItem{
		{id: "json", label: "JSON"},
		{id: "markdown", label: "Markdown"},
		{id: "docx", label: "Word Document (.docx)"},
		{id: "pdf", label: "PDF"},
	}

	preselected := make(map[string]bool, len(preselectedFormats))
	for _, f := range preselectedFormats {
		preselected[f] = true
	}
	for i, f := range allFormats {
		if preselected[f.id] {
			allFormats[i].selected = true
		}
	}

	return ExportStep{
		formats:        allFormats,
		phase:          exportPhaseSelecting,
		progress:       make(map[string]exportFormatProgress),
		outputDir:      outputDir,
		sessionID:      sessionID,
		sess:           sess,
		specgenVersion: specgenVersion,
		skipStep:       len(preselectedFormats) > 0,
		width:          width,
		height:         height,
	}
}

// Init implements tea.Model.
// When the step was pre-configured via --format, it immediately starts exporting.
func (e ExportStep) Init() tea.Cmd {
	if e.skipStep {
		return func() tea.Msg { return startExportMsg{} }
	}
	return nil
}

// startExportMsg triggers the export phase.
type startExportMsg struct{}

// Update implements tea.Model.
func (e ExportStep) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if e.phase == exportPhaseSelecting {
			return e.handleSelectKey(msg)
		}
	case startExportMsg:
		return e.beginExport()
	case ExportProgressMsg:
		return e.handleProgress(msg)
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height
	}
	return e, nil
}

func (e ExportStep) handleSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if e.cursor > 0 {
			e.cursor--
		}
	case tea.KeyDown:
		if e.cursor < len(e.formats)-1 {
			e.cursor++
		}
	case tea.KeySpace:
		e.formats[e.cursor].selected = !e.formats[e.cursor].selected
	case tea.KeyEnter:
		if e.hasSelection() {
			return e.beginExport()
		}
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			if r == ' ' {
				e.formats[e.cursor].selected = !e.formats[e.cursor].selected
			}
		}
	}
	return e, nil
}

func (e ExportStep) hasSelection() bool {
	for _, f := range e.formats {
		if f.selected {
			return true
		}
	}
	return false
}

func (e ExportStep) beginExport() (tea.Model, tea.Cmd) {
	e.phase = exportPhaseExporting

	outputDir := e.outputDir
	sessionID := e.sessionID
	sess := e.sess
	specgenVersion := e.specgenVersion

	var cmds []tea.Cmd
	for _, f := range e.formats {
		if !f.selected {
			continue
		}
		fID := f.id // capture loop variable by value
		cmds = append(cmds, func() tea.Msg {
			return runExportFormat(fID, outputDir, sessionID, sess, specgenVersion)
		})
	}
	if len(cmds) == 0 {
		return e, nil
	}
	return e, tea.Batch(cmds...)
}

func (e ExportStep) handleProgress(msg ExportProgressMsg) (tea.Model, tea.Cmd) {
	e.progress[msg.Format] = exportFormatProgress{
		path: msg.Path,
		err:  msg.Err,
		done: true,
	}

	// Check whether all selected formats have finished.
	allDone := true
	for _, f := range e.formats {
		if !f.selected {
			continue
		}
		p, ok := e.progress[f.id]
		if !ok || !p.done {
			allDone = false
			break
		}
	}

	if allDone {
		e.phase = exportPhaseDone
		var paths []string
		for _, f := range e.formats {
			if !f.selected {
				continue
			}
			if p, ok := e.progress[f.id]; ok && p.err == nil {
				paths = append(paths, p.path)
			}
		}
		return e, func() tea.Msg { return ExportDoneMsg{OutputPaths: paths} }
	}

	return e, nil
}

// View implements tea.Model.
func (e ExportStep) View() string {
	var sb strings.Builder
	sb.WriteString(exportTitleStyle.Render("Step 5: Choose Export Formats") + "\n\n")

	switch e.phase {
	case exportPhaseSelecting:
		sb.WriteString(exportHintStyle.Render("Space to toggle · Enter to confirm") + "\n\n")
		for i, f := range e.formats {
			cursor := "  "
			if i == e.cursor {
				cursor = exportCursorStyle.Render("> ")
			}
			var check string
			if f.selected {
				check = exportCheckedStyle.Render("[x] " + f.label)
			} else {
				check = exportUncheckedStyle.Render("[ ] " + f.label)
			}
			sb.WriteString(cursor + check + "\n")
		}
		if !e.hasSelection() {
			sb.WriteString("\n" + exportValidationStyle.Render("⚠  Select at least one format to continue") + "\n")
		}

	case exportPhaseExporting:
		sb.WriteString("Exporting...\n\n")
		for _, f := range e.formats {
			if !f.selected {
				continue
			}
			p, ok := e.progress[f.id]
			switch {
			case !ok:
				sb.WriteString(exportDimStyle.Render(fmt.Sprintf("  ⋯  %s", f.label)) + "\n")
			case p.err != nil:
				sb.WriteString(exportErrorStyle.Render(fmt.Sprintf("  ✗  %s: %s", f.label, p.err.Error())) + "\n")
			default:
				sb.WriteString(exportSuccessStyle.Render(fmt.Sprintf("  ✓  %s → %s", f.label, p.path)) + "\n")
			}
		}

	case exportPhaseDone:
		sb.WriteString(exportSuccessStyle.Render("Export complete!") + "\n\n")
		for _, f := range e.formats {
			if !f.selected {
				continue
			}
			p, ok := e.progress[f.id]
			if !ok {
				continue
			}
			if p.err != nil {
				sb.WriteString(exportErrorStyle.Render(fmt.Sprintf("  ✗  %s: %s", f.label, p.err.Error())) + "\n")
			} else {
				sb.WriteString(exportSuccessStyle.Render(fmt.Sprintf("  ✓  %s", p.path)) + "\n")
			}
		}
	}

	return sb.String()
}

// SelectedFormats returns the IDs of currently-selected formats.
func (e ExportStep) SelectedFormats() []string {
	var result []string
	for _, f := range e.formats {
		if f.selected {
			result = append(result, f.id)
		}
	}
	return result
}

// IsDone reports whether all exports have completed.
func (e ExportStep) IsDone() bool { return e.phase == exportPhaseDone }

// IsSkipped reports whether the selection UI is skipped due to --format flag.
func (e ExportStep) IsSkipped() bool { return e.skipStep }

// Cursor returns the current cursor index (exported for testing).
func (e ExportStep) ExportCursor() int { return e.cursor }

// --- package-level export helpers ----------------------------------------

// runExportFormat performs the actual file export for a single format and
// returns an ExportProgressMsg. Designed to be called as a tea.Cmd closure.
func runExportFormat(formatID, outputDir, sessionID string, sess *session.Session, specgenVersion string) tea.Msg {
	ext := formatID // "pdf", "docx", "markdown", "json" all match their extensions
	if ext == "markdown" {
		ext = "md"
	}

	prefix := sessionID
	if outputDir != "" {
		prefix = filepath.Join(outputDir, sessionID)
	}
	filename := prefix + "." + ext

	var err error
	switch formatID {
	case "pdf":
		err = export.ExportPDF("", filename)
	case "docx":
		err = export.ExportDOCX("", filename)
	case "markdown":
		err = writeTextFile(filename, "")
	case "json":
		if sess != nil {
			err = export.ExportJSON(sess, specgenVersion, prefix)
		} else {
			err = writeTextFile(filename, "{}")
		}
	default:
		err = fmt.Errorf("unknown format: %s", formatID)
	}

	return ExportProgressMsg{
		Format: formatID,
		Path:   filename,
		Err:    err,
	}
}

// writeTextFile writes content to path (mode 0600).
func writeTextFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}
