// Tests for the Markdown renderer (spec Section 13.3).
package render_test

import (
	"strings"
	"testing"

	"github.com/zyx-holdings/go-spec/internal/render"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// ---- helpers ----

// coreSession returns a Store with the minimum required core answers and the
// given profile set. Use sess.Session() to access the underlying session.
func coreSession(profile string) *session.Store {
	st := session.New()
	s := st.Session()
	s.CLIProfile = profile
	st.SetAnswer("project_name", session.NewStringValue("test-tool"))
	st.SetAnswer("project_description", session.NewStringValue("A test CLI tool for spec rendering."))
	st.SetAnswer("primary_language", session.NewStringValue("Go"))
	return st
}

// ---- Markdown output structure tests ----

func TestMarkdown_ContainsProjectName(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	if !strings.Contains(out, "test-tool") {
		t.Errorf("output does not contain project name %q:\n%s", "test-tool", out)
	}
}

func TestMarkdown_HasATXHeaders(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	// Must have level-1 and level-2 ATX headers
	if !strings.Contains(out, "# ") {
		t.Error("output lacks ATX level-1 header")
	}
	if !strings.Contains(out, "## ") {
		t.Error("output lacks ATX level-2 header")
	}
}

func TestMarkdown_CoreSectionsPresent(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())

	want := []string{
		"## 1. Problem Statement",
		"## 2. Goals and Non-Goals",
		"## 3. System Overview",
	}
	for _, h := range want {
		if !strings.Contains(out, h) {
			t.Errorf("output missing expected section header %q", h)
		}
	}
}

func TestMarkdown_TestingMatrixPresent(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	if !strings.Contains(out, "Testing and Validation Matrix") {
		t.Error("output missing Testing and Validation Matrix section")
	}
}

func TestMarkdown_ImplementationChecklistPresent(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	if !strings.Contains(out, "Implementation Checklist") {
		t.Error("output missing Implementation Checklist section")
	}
}

func TestMarkdown_EmptyFeatureSelection_OnlyCoreSections(t *testing.T) {
	st := coreSession("oneshot")
	// No features selected
	out := render.Markdown(st.Session())

	// Core sections must exist
	if !strings.Contains(out, "## 1. Problem Statement") {
		t.Error("missing Problem Statement with empty features")
	}
	if !strings.Contains(out, "## 2. Goals and Non-Goals") {
		t.Error("missing Goals and Non-Goals with empty features")
	}
	if !strings.Contains(out, "## 3. System Overview") {
		t.Error("missing System Overview with empty features")
	}

	// Feature sections must NOT appear (section 4 would be a feature section,
	// but with no features the next fixed section is Testing at 4).
	if strings.Contains(out, "## 4. Authentication") {
		t.Error("unexpected Authentication section with no selected features")
	}
}

func TestMarkdown_FeatureSectionsOmittedWhenNotSelected(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"storage"}
	// Set required storage answers
	st.SetAnswer("storage_type", session.NewStringValue("Relational DB"))
	st.SetAnswer("database_name", session.NewStringValue("PostgreSQL"))

	out := render.Markdown(st.Session())

	// Storage section must be present
	if !strings.Contains(out, "Storage") {
		t.Error("output missing Storage section despite being selected")
	}

	// Authentication section must NOT appear (not selected)
	if strings.Contains(out, "## 4. Authentication") || strings.Contains(out, "## 5. Authentication") {
		t.Error("unexpected Authentication section when not selected")
	}
}

func TestMarkdown_FeatureSectionNumbering(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"authentication", "storage"}
	st.SetAnswer("auth_method", session.NewStringValue("JWT"))
	st.SetAnswer("session_management", session.NewStringValue("Redis-backed"))
	st.SetAnswer("storage_type", session.NewStringValue("Relational DB"))
	st.SetAnswer("database_name", session.NewStringValue("PostgreSQL"))

	out := render.Markdown(st.Session())

	if !strings.Contains(out, "## 4. Authentication") {
		t.Error("Authentication should be section 4")
	}
	if !strings.Contains(out, "## 5. Storage") {
		t.Error("Storage should be section 5")
	}
	if !strings.Contains(out, "## 6. Testing and Validation Matrix") {
		t.Error("Testing and Validation Matrix should be section 6")
	}
	if !strings.Contains(out, "## 7. Implementation Checklist") {
		t.Error("Implementation Checklist should be section 7")
	}
}

func TestMarkdown_SkippedOptionalQuestion_TODONote(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"storage"}
	st.SetAnswer("storage_type", session.NewStringValue("Relational DB"))
	st.SetAnswer("database_name", session.NewStringValue("PostgreSQL"))
	st.SkipAnswer("caching_strategy") // optional

	out := render.Markdown(st.Session())
	if !strings.Contains(out, "TODO:") {
		t.Error("output should contain TODO note for skipped optional question")
	}
}

func TestMarkdown_AbsentOptionalQuestion_TODONote(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"storage"}
	st.SetAnswer("storage_type", session.NewStringValue("Relational DB"))
	st.SetAnswer("database_name", session.NewStringValue("PostgreSQL"))
	// caching_strategy not answered at all

	out := render.Markdown(st.Session())
	if !strings.Contains(out, "TODO:") {
		t.Error("output should contain TODO note for absent optional question")
	}
}

func TestMarkdown_MultiSelectAnswer_BulletList(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"testing"}
	st.SetAnswer("test_framework", session.NewStringValue("testify"))
	st.SetAnswer("test_types", session.NewMultiValue([]string{"Unit", "Integration"}))

	out := render.Markdown(st.Session())

	if !strings.Contains(out, "- Unit") {
		t.Error("output should contain bullet item 'Unit'")
	}
	if !strings.Contains(out, "- Integration") {
		t.Error("output should contain bullet item 'Integration'")
	}
}

func TestMarkdown_ConfirmQuestion_YesNo(t *testing.T) {
	st := coreSession("oneshot")
	st.Session().SelectedFeatures = []string{"observability"}
	st.SetAnswer("logging_framework", session.NewStringValue("zap"))
	st.SetAnswer("tracing_enabled", session.NewStringValue("true"))

	out := render.Markdown(st.Session())
	if !strings.Contains(out, "Yes") {
		t.Errorf("confirm answer 'true' should render as 'Yes' in output:\n%s", out)
	}
}

func TestMarkdown_SpecialCharactersEscaped(t *testing.T) {
	st := coreSession("oneshot")
	// Project name with special inline markdown characters
	st.SetAnswer("project_name", session.NewStringValue("my_tool*name"))

	out := render.Markdown(st.Session())
	// The asterisk and underscore should be escaped in the output
	if !strings.Contains(out, `my\_tool\*name`) {
		t.Errorf("special markdown characters in project name should be escaped, got output:\n%s", out)
	}
}

func TestMarkdown_SectionsSeparatedByHR(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	// Sections must be separated by horizontal rules
	if !strings.Contains(out, "\n---\n") {
		t.Error("sections should be separated by horizontal rules (---)")
	}
}

func TestMarkdown_SessionIDInHeader(t *testing.T) {
	st := coreSession("oneshot")
	id := st.Session().ID
	out := render.Markdown(st.Session())
	if !strings.Contains(out, id) {
		t.Errorf("output should contain session ID %q", id)
	}
}

func TestMarkdown_SystemOverview_ListsSelectedComponents(t *testing.T) {
	st := coreSession("daemon")
	st.Session().SelectedFeatures = []string{"api", "observability"}
	st.SetAnswer("api_style", session.NewStringValue("REST"))
	st.SetAnswer("api_versioning", session.NewStringValue("/v1/"))
	st.SetAnswer("logging_framework", session.NewStringValue("zap"))
	st.SetAnswer("tracing_enabled", session.NewStringValue("true"))

	out := render.Markdown(st.Session())
	if !strings.Contains(out, "API Component") {
		t.Error("System Overview should list API Component")
	}
	if !strings.Contains(out, "Observability Component") {
		t.Error("System Overview should list Observability Component")
	}
}

func TestMarkdown_TestingMatrix_HasCheckboxes(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	if !strings.Contains(out, "- [ ]") {
		t.Error("Testing and Validation Matrix should contain checkbox items")
	}
}

func TestMarkdown_ImplementationChecklist_HasCheckboxes(t *testing.T) {
	st := coreSession("oneshot")
	out := render.Markdown(st.Session())
	if !strings.Contains(out, "- [ ]") {
		t.Error("Implementation Checklist should contain checkbox items")
	}
}

func TestMarkdown_ValidCommonMark_NoPanicOnAllFeatures(t *testing.T) {
	st := coreSession("hybrid")
	st.Session().SelectedFeatures = []string{
		"authentication", "storage", "api", "testing", "observability", "deployment",
	}
	// Provide required answers for all feature areas
	st.SetAnswer("auth_method", session.NewStringValue("JWT"))
	st.SetAnswer("session_management", session.NewStringValue("stateless JWT, 1h TTL"))
	st.SetAnswer("storage_type", session.NewStringValue("Relational DB"))
	st.SetAnswer("database_name", session.NewStringValue("PostgreSQL"))
	st.SetAnswer("api_style", session.NewStringValue("REST"))
	st.SetAnswer("api_versioning", session.NewStringValue("URL path (/v1/)"))
	st.SetAnswer("test_framework", session.NewStringValue("testify"))
	st.SetAnswer("test_types", session.NewMultiValue([]string{"Unit", "Integration"}))
	st.SetAnswer("logging_framework", session.NewStringValue("zap"))
	st.SetAnswer("tracing_enabled", session.NewStringValue("true"))
	st.SetAnswer("deployment_target", session.NewStringValue("Cloud (AWS)"))
	st.SetAnswer("containerized", session.NewStringValue("true"))

	// Should not panic
	out := render.Markdown(st.Session())
	if out == "" {
		t.Error("Markdown output should not be empty")
	}
}

// ---- wrapProse tests ----

func TestWrapProse_ShortLine(t *testing.T) {
	lines := render.WrapProse("hello world", 100)
	if len(lines) != 1 || lines[0] != "hello world" {
		t.Errorf("WrapProse short: got %v, want [\"hello world\"]", lines)
	}
}

func TestWrapProse_LongLine(t *testing.T) {
	text := strings.Repeat("word ", 30) // 150 chars
	text = strings.TrimSpace(text)
	lines := render.WrapProse(text, 100)
	if len(lines) < 2 {
		t.Errorf("WrapProse long: expected at least 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if len(l) > 100 {
			t.Errorf("WrapProse: line exceeds 100 chars: %q", l)
		}
	}
}

func TestWrapProse_Empty(t *testing.T) {
	lines := render.WrapProse("", 100)
	if len(lines) != 1 {
		t.Errorf("WrapProse empty: expected 1 line, got %d", len(lines))
	}
}

// ---- escapeMarkdown tests ----

func TestEscapeMarkdown_SpecialChars(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"foo*bar", `foo\*bar`},
		{"foo_bar", `foo\_bar`},
		{"[link]", `\[link\]`},
		{"# header", "# header"},   // # not escaped in inline position
		{"back\\slash", `back\\slash`},
		{"pipe|char", `pipe\|char`},
		{"test-tool", "test-tool"}, // hyphen not special inline
		{"v1.0", "v1.0"},          // dot not special inline
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := render.EscapeMarkdown(tc.input)
			if got != tc.want {
				t.Errorf("EscapeMarkdown(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
