package questions_test

import (
	"testing"

	"github.com/zyx-holdings/go-spec/internal/questions"
)

// ── Registry completeness ─────────────────────────────────────────────────────

func TestProfiles_Count(t *testing.T) {
	want := 4
	if got := len(questions.Profiles); got != want {
		t.Errorf("len(Profiles) = %d, want %d", got, want)
	}
}

func TestProfiles_ExpectedIDs(t *testing.T) {
	expected := []string{"oneshot", "daemon", "subcommand", "hybrid"}
	for _, id := range expected {
		if _, ok := questions.Profiles[id]; !ok {
			t.Errorf("Profiles missing expected ID %q", id)
		}
	}
}

func TestProfiles_AllHaveRequiredFields(t *testing.T) {
	for id, p := range questions.Profiles {
		if p.ID == "" {
			t.Errorf("Profile[%q].ID is empty", id)
		}
		if p.ID != id {
			t.Errorf("Profile[%q].ID = %q, want key to match ID", id, p.ID)
		}
		if p.DisplayName == "" {
			t.Errorf("Profile[%q].DisplayName is empty", id)
		}
		if p.Description == "" {
			t.Errorf("Profile[%q].Description is empty", id)
		}
	}
}

func TestFeatureAreas_Count(t *testing.T) {
	want := 12
	if got := len(questions.FeatureAreas); got != want {
		t.Errorf("len(FeatureAreas) = %d, want %d", got, want)
	}
}

func TestFeatureAreas_ExpectedIDs(t *testing.T) {
	expected := []string{
		"api", "database", "filesystem", "concurrency",
		"retry", "config", "hooks", "auth",
		"logging", "http", "templates", "statemachine",
	}
	for _, id := range expected {
		if _, ok := questions.FeatureAreas[id]; !ok {
			t.Errorf("FeatureAreas missing expected ID %q", id)
		}
	}
}

func TestFeatureAreas_AllHaveRequiredFields(t *testing.T) {
	for id, fa := range questions.FeatureAreas {
		if fa.ID == "" {
			t.Errorf("FeatureArea[%q].ID is empty", id)
		}
		if fa.ID != id {
			t.Errorf("FeatureArea[%q].ID = %q, want key to match ID", id, fa.ID)
		}
		if fa.DisplayName == "" {
			t.Errorf("FeatureArea[%q].DisplayName is empty", id)
		}
		if fa.Description == "" {
			t.Errorf("FeatureArea[%q].Description is empty", id)
		}
	}
}

// ── Question bank completeness ────────────────────────────────────────────────

func TestAll_QuestionIDs_AreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, q := range questions.All {
		if seen[q.ID] {
			t.Errorf("duplicate question ID: %q", q.ID)
		}
		seen[q.ID] = true
	}
}

func TestAll_QuestionFeatureAreas_AreValid(t *testing.T) {
	for _, q := range questions.All {
		if q.FeatureArea == "_core" {
			continue
		}
		if !questions.IsValidFeatureArea(q.FeatureArea) {
			t.Errorf("question %q has unknown feature area %q", q.ID, q.FeatureArea)
		}
	}
}

func TestAll_CoreQuestions_Count(t *testing.T) {
	core := questions.ByFeatureArea("_core")
	want := 8
	if got := len(core); got != want {
		t.Errorf("_core question count = %d, want %d", got, want)
	}
}

func TestAll_CoreQuestions_RequiredIDs(t *testing.T) {
	required := []string{"tool_name", "tool_description", "go_module"}
	core := questions.ByFeatureArea("_core")
	coreByID := map[string]questions.Question{}
	for _, q := range core {
		coreByID[q.ID] = q
	}
	for _, id := range required {
		q, ok := coreByID[id]
		if !ok {
			t.Errorf("expected _core question %q to exist", id)
			continue
		}
		if !q.Required {
			t.Errorf("expected _core question %q to be required", id)
		}
	}
}

func TestAll_CoreQuestions_HasOptional(t *testing.T) {
	core := questions.ByFeatureArea("_core")
	for _, q := range core {
		if !q.Required {
			return
		}
	}
	t.Error("expected at least one optional _core question (e.g. binary_name)")
}

func TestAll_EachFeatureAreaHasQuestions(t *testing.T) {
	for id := range questions.FeatureAreas {
		t.Run(id, func(t *testing.T) {
			qs := questions.ByFeatureArea(id)
			if len(qs) == 0 {
				t.Errorf("feature area %q has no questions", id)
			}
		})
	}
}

func TestAll_FeatureAreaQuestions_TotalAtLeast60(t *testing.T) {
	total := 0
	for id := range questions.FeatureAreas {
		total += len(questions.ByFeatureArea(id))
	}
	if total < 60 {
		t.Errorf("total feature-area question count = %d, want >= 60", total)
	}
}

func TestAll_OrderWithinSection_IsPositive(t *testing.T) {
	for _, q := range questions.All {
		if q.Order <= 0 {
			t.Errorf("question %q has Order = %d, want > 0", q.ID, q.Order)
		}
	}
}

func TestAll_OrderWithinSection_IsDeterministic(t *testing.T) {
	// Calling ByFeatureArea twice must return the same sequence.
	for id := range questions.FeatureAreas {
		a := questions.ByFeatureArea(id)
		b := questions.ByFeatureArea(id)
		if len(a) != len(b) {
			t.Errorf("ByFeatureArea(%q) returned different lengths on two calls", id)
			continue
		}
		for i := range a {
			if a[i].ID != b[i].ID {
				t.Errorf("ByFeatureArea(%q)[%d] = %q / %q on two calls", id, i, a[i].ID, b[i].ID)
			}
		}
	}
}

func TestAll_Questions_InputTypeNonEmpty(t *testing.T) {
	for _, q := range questions.All {
		if q.InputType == "" {
			t.Errorf("question %q has empty InputType", q.ID)
		}
	}
}

func TestAll_Questions_InputTypeIsKnown(t *testing.T) {
	known := map[string]bool{
		questions.InputText:        true,
		questions.InputSelect:      true,
		questions.InputMultiSelect: true,
		questions.InputConfirm:     true,
	}
	for _, q := range questions.All {
		if !known[q.InputType] {
			t.Errorf("question %q has unknown InputType %q", q.ID, q.InputType)
		}
	}
}

func TestAll_SelectAndMultiSelect_HaveOptions(t *testing.T) {
	for _, q := range questions.All {
		if q.InputType == questions.InputSelect || q.InputType == questions.InputMultiSelect {
			if len(q.Options) == 0 {
				t.Errorf("question %q (InputType=%q) has no Options", q.ID, q.InputType)
			}
		}
	}
}

// ── ByFeatureArea ─────────────────────────────────────────────────────────────

func TestByFeatureArea_Core(t *testing.T) {
	got := questions.ByFeatureArea("_core")
	if len(got) == 0 {
		t.Fatal("ByFeatureArea(_core) returned empty slice")
	}
	for _, q := range got {
		if q.FeatureArea != "_core" {
			t.Errorf("ByFeatureArea(_core): got question with feature area %q", q.FeatureArea)
		}
	}
}

func TestByFeatureArea_KnownArea(t *testing.T) {
	areas := []string{
		"api", "database", "filesystem", "concurrency",
		"retry", "config", "hooks", "auth",
		"logging", "http", "templates", "statemachine",
	}
	for _, area := range areas {
		t.Run(area, func(t *testing.T) {
			got := questions.ByFeatureArea(area)
			if len(got) == 0 {
				t.Errorf("ByFeatureArea(%q) returned empty slice", area)
			}
			for _, q := range got {
				if q.FeatureArea != area {
					t.Errorf("ByFeatureArea(%q): got question with area %q", area, q.FeatureArea)
				}
			}
		})
	}
}

func TestByFeatureArea_Unknown_ReturnsEmpty(t *testing.T) {
	got := questions.ByFeatureArea("not_a_real_area")
	if len(got) != 0 {
		t.Errorf("ByFeatureArea(unknown) = %v, want empty slice", got)
	}
}

// ── ByID ─────────────────────────────────────────────────────────────────────

func TestByID_Found(t *testing.T) {
	cases := []string{"tool_name", "tool_description", "go_module", "logging_library"}
	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			q, ok := questions.ByID(id)
			if !ok {
				t.Fatalf("ByID(%q): not found", id)
			}
			if q.ID != id {
				t.Errorf("ByID(%q).ID = %q, want %q", id, q.ID, id)
			}
		})
	}
}

func TestByID_NotFound(t *testing.T) {
	_, ok := questions.ByID("nonexistent_question")
	if ok {
		t.Error("ByID(nonexistent): ok = true, want false")
	}
}

func TestByID_EmptyID(t *testing.T) {
	_, ok := questions.ByID("")
	if ok {
		t.Error("ByID(empty string): ok = true, want false")
	}
}

// ── IsValidProfile ────────────────────────────────────────────────────────────

func TestIsValidProfile_KnownProfiles(t *testing.T) {
	known := []string{"oneshot", "daemon", "subcommand", "hybrid"}
	for _, id := range known {
		if !questions.IsValidProfile(id) {
			t.Errorf("IsValidProfile(%q) = false, want true", id)
		}
	}
}

func TestIsValidProfile_Invalid(t *testing.T) {
	invalid := []string{"", "not_a_profile", "ONESHOT", "oneshot "}
	for _, id := range invalid {
		if questions.IsValidProfile(id) {
			t.Errorf("IsValidProfile(%q) = true, want false", id)
		}
	}
}

// ── IsValidFeatureArea ────────────────────────────────────────────────────────

func TestIsValidFeatureArea_KnownAreas(t *testing.T) {
	known := []string{
		"api", "database", "filesystem", "concurrency",
		"retry", "config", "hooks", "auth",
		"logging", "http", "templates", "statemachine",
	}
	for _, id := range known {
		if !questions.IsValidFeatureArea(id) {
			t.Errorf("IsValidFeatureArea(%q) = false, want true", id)
		}
	}
}

func TestIsValidFeatureArea_CoreIsNotSelectable(t *testing.T) {
	if questions.IsValidFeatureArea("_core") {
		t.Error("IsValidFeatureArea(_core) = true, want false — _core is not a selectable feature area")
	}
}

func TestIsValidFeatureArea_Invalid(t *testing.T) {
	invalid := []string{"", "not_a_feature", "LOGGING", "logging "}
	for _, id := range invalid {
		if questions.IsValidFeatureArea(id) {
			t.Errorf("IsValidFeatureArea(%q) = true, want false", id)
		}
	}
}

// ── FilterByFeatures ──────────────────────────────────────────────────────────

func TestFilterByFeatures_EmptySelection_ReturnsCoreOnly(t *testing.T) {
	got := questions.FilterByFeatures([]string{})
	core := questions.ByFeatureArea("_core")
	if len(got) != len(core) {
		t.Errorf("FilterByFeatures([]) len = %d, want %d (core only)", len(got), len(core))
	}
	for i, q := range got {
		if q.FeatureArea != "_core" {
			t.Errorf("FilterByFeatures([])[%d].FeatureArea = %q, want _core", i, q.FeatureArea)
		}
	}
}

func TestFilterByFeatures_SingleArea_CorePlusArea(t *testing.T) {
	got := questions.FilterByFeatures([]string{"logging"})
	core := questions.ByFeatureArea("_core")
	logging := questions.ByFeatureArea("logging")
	want := len(core) + len(logging)
	if len(got) != want {
		t.Errorf("FilterByFeatures([logging]) len = %d, want %d", len(got), want)
	}
	// First len(core) should be _core questions.
	for i, q := range got[:len(core)] {
		if q.FeatureArea != "_core" {
			t.Errorf("FilterByFeatures([logging])[%d].FeatureArea = %q, want _core", i, q.FeatureArea)
		}
	}
	// Remainder should be logging questions.
	for i, q := range got[len(core):] {
		if q.FeatureArea != "logging" {
			t.Errorf("FilterByFeatures([logging]) logging section [%d].FeatureArea = %q, want logging", i, q.FeatureArea)
		}
	}
}

func TestFilterByFeatures_MultipleAreas_OrderPreserved(t *testing.T) {
	selected := []string{"auth", "database"}
	got := questions.FilterByFeatures(selected)

	core := questions.ByFeatureArea("_core")
	auth := questions.ByFeatureArea("auth")
	db := questions.ByFeatureArea("database")
	want := len(core) + len(auth) + len(db)

	if len(got) != want {
		t.Errorf("FilterByFeatures([auth,database]) len = %d, want %d", len(got), want)
	}

	// Verify order: core → auth → database
	offset := len(core)
	for i, q := range got[offset : offset+len(auth)] {
		if q.FeatureArea != "auth" {
			t.Errorf("auth section[%d].FeatureArea = %q, want auth", i, q.FeatureArea)
		}
	}
	offset += len(auth)
	for i, q := range got[offset:] {
		if q.FeatureArea != "database" {
			t.Errorf("database section[%d].FeatureArea = %q, want database", i, q.FeatureArea)
		}
	}
}

func TestFilterByFeatures_NoDuplicates(t *testing.T) {
	got := questions.FilterByFeatures([]string{"logging", "auth", "retry"})
	seen := map[string]bool{}
	for _, q := range got {
		if seen[q.ID] {
			t.Errorf("FilterByFeatures returned duplicate question ID: %q", q.ID)
		}
		seen[q.ID] = true
	}
}

func TestFilterByFeatures_AllAreas(t *testing.T) {
	all := []string{
		"api", "database", "filesystem", "concurrency",
		"retry", "config", "hooks", "auth",
		"logging", "http", "templates", "statemachine",
	}
	got := questions.FilterByFeatures(all)
	if len(got) != len(questions.All) {
		t.Errorf("FilterByFeatures(all areas) len = %d, want %d", len(got), len(questions.All))
	}
}
