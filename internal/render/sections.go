// Package render converts spec data into formatted document output (spec Section 8.2).
package render

import (
	"fmt"
	"strings"

	"github.com/zyx-holdings/go-spec/internal/questions"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// Section represents a numbered document section with a title and body lines.
type Section struct {
	Number int
	Title  string
	Lines  []string
}

// sectionBuilder collects helpers for building Section content from a session.
type sectionBuilder struct {
	sess *session.Session
}

// answer returns the string value for a question ID, the skipped flag, and
// whether an answer was recorded at all.
func (b *sectionBuilder) answer(id string) (value string, skipped bool, found bool) {
	ans, ok := b.sess.Answers[id]
	if !ok {
		return "", false, false
	}
	if ans.Skipped {
		return "", true, true
	}
	return ans.Value.String(), false, true
}

// answerMulti returns the []string value for a multi-select question.
func (b *sectionBuilder) answerMulti(id string) (values []string, skipped bool, found bool) {
	ans, ok := b.sess.Answers[id]
	if !ok {
		return nil, false, false
	}
	if ans.Skipped {
		return nil, true, true
	}
	if ans.Value.IsMulti() {
		return ans.Value.Strings(), false, true
	}
	return []string{ans.Value.String()}, false, true
}

// prose returns the answer as prose, or a TODO note when skipped/absent.
func (b *sectionBuilder) prose(id string) string {
	val, skipped, found := b.answer(id)
	if !found || skipped {
		return "TODO: answer not provided."
	}
	return val
}

// bulletList returns the answer as a markdown bullet list, or a TODO note.
func (b *sectionBuilder) bulletList(id string) []string {
	vals, skipped, found := b.answerMulti(id)
	if !found || skipped || len(vals) == 0 {
		return []string{"- TODO: answer not provided."}
	}
	lines := make([]string, len(vals))
	for i, v := range vals {
		lines[i] = "- " + EscapeMarkdown(v)
	}
	return lines
}

// confirmText returns "Yes" or "No" for a confirm question, or TODO.
func (b *sectionBuilder) confirmText(id string) string {
	val, skipped, found := b.answer(id)
	if !found || skipped {
		return "TODO: answer not provided."
	}
	switch strings.ToLower(val) {
	case "true", "yes", "y", "1":
		return "Yes"
	case "false", "no", "n", "0":
		return "No"
	default:
		return EscapeMarkdown(val)
	}
}

// hasFeature reports whether the given feature area was selected.
func (b *sectionBuilder) hasFeature(id string) bool {
	for _, f := range b.sess.SelectedFeatures {
		if f == id {
			return true
		}
	}
	return false
}

// profileDisplay returns the human-readable profile name.
func (b *sectionBuilder) profileDisplay() string {
	p, ok := questions.Profiles[b.sess.CLIProfile]
	if !ok {
		if b.sess.CLIProfile == "" {
			return "TODO: profile not selected."
		}
		return b.sess.CLIProfile
	}
	return p.DisplayName + " — " + p.Description
}

// featureDisplay returns the human-readable feature area name.
func featureDisplay(id string) string {
	fa, ok := questions.FeatureAreas[id]
	if !ok {
		return id
	}
	return fa.DisplayName
}

// ---- Section builders ----

// buildProblemStatement builds section 1: Problem Statement.
func (b *sectionBuilder) buildProblemStatement(num int) Section {
	name := b.prose("project_name")
	desc := b.prose("project_description")

	var lines []string
	lines = append(lines, WrapProse(fmt.Sprintf(
		"**%s** is a software project described as follows: %s",
		EscapeMarkdown(name), EscapeMarkdown(desc),
	), 100)...)

	return Section{Number: num, Title: "Problem Statement", Lines: lines}
}

// buildGoalsAndNonGoals builds section 2: Goals and Non-Goals.
func (b *sectionBuilder) buildGoalsAndNonGoals(num int) Section {
	lang := b.prose("primary_language")
	teamSize, skipped, found := b.answer("team_size")

	var lines []string
	lines = append(lines, "### Goals", "")
	lines = append(lines, fmt.Sprintf("- Implement the system described in the Problem Statement using **%s**.", EscapeMarkdown(lang)))
	lines = append(lines, "- Deliver a working, tested, and documented implementation.")
	lines = append(lines, "- Meet all acceptance criteria for each selected feature area.")

	lines = append(lines, "", "### Non-Goals", "")
	lines = append(lines, "- Features not listed in the selected feature areas are out of scope.")
	if found && !skipped && teamSize != "" {
		lines = append(lines, fmt.Sprintf("- Building for a team significantly larger than **%s** is not a current goal.", EscapeMarkdown(teamSize)))
	}
	lines = append(lines, "- Performance optimization beyond stated targets is deferred to a future iteration.")

	return Section{Number: num, Title: "Goals and Non-Goals", Lines: lines}
}

// buildSystemOverview builds section 3: System Overview.
func (b *sectionBuilder) buildSystemOverview(num int) Section {
	var lines []string

	// Main components based on profile and selected features
	lines = append(lines, "### Main Components", "")
	lines = append(lines, fmt.Sprintf("- **Core Application** — %s", EscapeMarkdown(b.profileDisplay())))

	for _, fid := range b.sess.SelectedFeatures {
		fa, ok := questions.FeatureAreas[fid]
		if !ok {
			continue
		}
		lines = append(lines, fmt.Sprintf("- **%s Component** — %s", fa.DisplayName, fa.Description))
	}

	lines = append(lines, "", "### External Dependencies", "")

	if len(b.sess.SelectedFeatures) == 0 {
		lines = append(lines, "No external feature areas selected. Dependencies are limited to the core application.")
	} else {
		for _, fid := range b.sess.SelectedFeatures {
			dep := externalDependencyNote(fid, b)
			if dep != "" {
				lines = append(lines, fmt.Sprintf("- %s", dep))
			}
		}
	}

	return Section{Number: num, Title: "System Overview", Lines: lines}
}

// externalDependencyNote returns a one-line dependency note for a feature area.
func externalDependencyNote(fid string, b *sectionBuilder) string {
	switch fid {
	case "authentication":
		method := b.prose("auth_method")
		return fmt.Sprintf("Authentication via **%s**", EscapeMarkdown(method))
	case "storage":
		dbName := b.prose("database_name")
		storageType := b.prose("storage_type")
		return fmt.Sprintf("Storage: **%s** (%s)", EscapeMarkdown(dbName), EscapeMarkdown(storageType))
	case "api":
		style := b.prose("api_style")
		return fmt.Sprintf("API layer: **%s**", EscapeMarkdown(style))
	case "testing":
		fw := b.prose("test_framework")
		return fmt.Sprintf("Test framework: **%s**", EscapeMarkdown(fw))
	case "observability":
		logging := b.prose("logging_framework")
		return fmt.Sprintf("Observability: **%s** for logging", EscapeMarkdown(logging))
	case "deployment":
		target := b.prose("deployment_target")
		return fmt.Sprintf("Deployment target: **%s**", EscapeMarkdown(target))
	default:
		fa, ok := questions.FeatureAreas[fid]
		if !ok {
			return ""
		}
		return fmt.Sprintf("%s integration required", fa.DisplayName)
	}
}

// buildFeatureSection builds a numbered section for a given feature area.
func (b *sectionBuilder) buildFeatureSection(num int, fid string) Section {
	fa, ok := questions.FeatureAreas[fid]
	title := fid
	if ok {
		title = fa.DisplayName
	}

	qs := questions.ByFeatureArea(fid)
	var lines []string

	for _, q := range qs {
		lines = append(lines, fmt.Sprintf("### %s", q.Text))
		lines = append(lines, "")

		switch q.InputType {
		case questions.InputTypeMultiSelect:
			bullets := b.bulletList(q.ID)
			lines = append(lines, bullets...)
		case questions.InputTypeConfirm:
			lines = append(lines, b.confirmText(q.ID))
		default:
			lines = append(lines, EscapeMarkdown(b.prose(q.ID)))
		}
		lines = append(lines, "")
	}

	// Trim trailing blank line
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return Section{Number: num, Title: title, Lines: lines}
}

// buildTestingMatrix builds the Testing and Validation Matrix section.
func (b *sectionBuilder) buildTestingMatrix(num int) Section {
	var lines []string
	lines = append(lines, "Auto-generated testing checklist based on selected feature areas.", "")

	// Core checks always present
	lines = append(lines,
		"- [ ] Unit tests for core business logic",
		"- [ ] Integration tests for external dependencies",
		"- [ ] Build passes with no errors",
		"- [ ] All required configuration values are validated at startup",
	)

	for _, fid := range b.sess.SelectedFeatures {
		for _, check := range testingChecksForFeature(fid, b) {
			lines = append(lines, "- [ ] "+check)
		}
	}

	return Section{Number: num, Title: "Testing and Validation Matrix", Lines: lines}
}

// testingChecksForFeature returns feature-specific testing checklist items.
func testingChecksForFeature(fid string, b *sectionBuilder) []string {
	switch fid {
	case "authentication":
		method, _, _ := b.answer("auth_method")
		return []string{
			fmt.Sprintf("%s authentication flow is tested end-to-end", EscapeMarkdown(method)),
			"Unauthorized requests return appropriate error responses",
			"Token expiration and refresh flows are validated",
		}
	case "storage":
		db, _, _ := b.answer("database_name")
		return []string{
			fmt.Sprintf("%s connection and CRUD operations are tested", EscapeMarkdown(db)),
			"Database migrations run cleanly on a fresh schema",
			"Concurrent read/write operations do not cause data corruption",
		}
	case "api":
		style, _, _ := b.answer("api_style")
		return []string{
			fmt.Sprintf("%s endpoints respond with correct status codes and payloads", EscapeMarkdown(style)),
			"API versioning is enforced and backward compatibility is verified",
			"Invalid request payloads return structured error responses",
		}
	case "testing":
		fw, _, _ := b.answer("test_framework")
		target, skipped, found := b.answer("test_coverage_target")
		checks := []string{fmt.Sprintf("%s test suite runs to completion without failures", EscapeMarkdown(fw))}
		if found && !skipped && target != "" {
			checks = append(checks, fmt.Sprintf("Code coverage meets the %s%% target", EscapeMarkdown(target)))
		}
		return checks
	case "observability":
		return []string{
			"Structured log output is validated for required fields",
			"Metrics are emitted for key operations",
			"Tracing spans are correctly propagated across service boundaries",
		}
	case "deployment":
		containerized, _, _ := b.answer("containerized")
		checks := []string{"Deployment scripts execute without errors in staging environment"}
		if strings.ToLower(containerized) == "true" || containerized == "yes" {
			checks = append(checks, "Container image builds successfully and passes health checks")
		}
		return checks
	default:
		fa, ok := questions.FeatureAreas[fid]
		if !ok {
			return nil
		}
		return []string{fmt.Sprintf("%s feature area functionality is tested", fa.DisplayName)}
	}
}

// buildImplementationChecklist builds the Implementation Checklist section.
func (b *sectionBuilder) buildImplementationChecklist(num int) Section {
	var lines []string
	lines = append(lines, "Auto-generated implementation checklist based on selected feature areas.", "")

	// Core items
	lines = append(lines,
		"- [ ] Repository initialized with Go module: `"+EscapeMarkdown(b.prose("project_name"))+"`",
		"- [ ] Core application skeleton created",
		"- [ ] CI pipeline configured",
		"- [ ] README.md with setup and usage instructions written",
	)

	for _, fid := range b.sess.SelectedFeatures {
		for _, item := range implementationItemsForFeature(fid, b) {
			lines = append(lines, "- [ ] "+item)
		}
	}

	return Section{Number: num, Title: "Implementation Checklist", Lines: lines}
}

// implementationItemsForFeature returns feature-specific implementation checklist items.
func implementationItemsForFeature(fid string, b *sectionBuilder) []string {
	switch fid {
	case "authentication":
		method, _, _ := b.answer("auth_method")
		return []string{
			fmt.Sprintf("%s authentication middleware implemented", EscapeMarkdown(method)),
			"User session management implemented",
			"Protected routes enforce authentication",
		}
	case "storage":
		db, _, _ := b.answer("database_name")
		storageType, _, _ := b.answer("storage_type")
		return []string{
			fmt.Sprintf("%s (%s) connection pool configured", EscapeMarkdown(db), EscapeMarkdown(storageType)),
			"Database schema and migrations defined",
			"Repository layer abstracts storage operations",
		}
	case "api":
		style, _, _ := b.answer("api_style")
		versioning, _, _ := b.answer("api_versioning")
		return []string{
			fmt.Sprintf("%s API routes implemented", EscapeMarkdown(style)),
			fmt.Sprintf("API versioning strategy applied: %s", EscapeMarkdown(versioning)),
			"API documentation generated",
		}
	case "testing":
		fw, _, _ := b.answer("test_framework")
		return []string{
			fmt.Sprintf("%s test framework integrated", EscapeMarkdown(fw)),
			"Test helpers and fixtures created",
			"Coverage reporting configured",
		}
	case "observability":
		logging, _, _ := b.answer("logging_framework")
		return []string{
			fmt.Sprintf("%s logger initialized with structured output", EscapeMarkdown(logging)),
			"Request/response logging middleware added",
			"Health check endpoint implemented",
		}
	case "deployment":
		target, _, _ := b.answer("deployment_target")
		cicd, skipped, found := b.answer("ci_cd_platform")
		items := []string{fmt.Sprintf("Deployment manifests created for %s", EscapeMarkdown(target))}
		if found && !skipped && cicd != "" {
			items = append(items, fmt.Sprintf("%s pipeline configured", EscapeMarkdown(cicd)))
		}
		return items
	default:
		fa, ok := questions.FeatureAreas[fid]
		if !ok {
			return nil
		}
		return []string{fmt.Sprintf("%s feature area integrated and configured", fa.DisplayName)}
	}
}
