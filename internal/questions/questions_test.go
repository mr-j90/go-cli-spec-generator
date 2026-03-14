package questions_test

import (
	"testing"

	"github.com/zyx-holdings/go-spec/internal/questions"
)

// ---- Registry completeness tests ----

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

func TestAll_HasCoreQuestions(t *testing.T) {
	core := questions.ByFeatureArea("_core")
	if len(core) == 0 {
		t.Error("expected at least one _core question, got none")
	}
}

func TestAll_CoreHasRequiredQuestions(t *testing.T) {
	required := []string{"project_name", "project_description", "primary_language"}
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

// ---- ByFeatureArea tests ----

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
	areas := []string{"authentication", "storage", "api", "testing", "observability", "deployment"}
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

// ---- ByID tests ----

func TestByID_Found(t *testing.T) {
	cases := []string{"project_name", "project_description", "primary_language", "test_framework"}
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

// ---- IsValidProfile tests ----

func TestIsValidProfile_KnownProfiles(t *testing.T) {
	known := []string{"api_service", "web_app", "cli_tool", "library", "data_pipeline"}
	for _, id := range known {
		if !questions.IsValidProfile(id) {
			t.Errorf("IsValidProfile(%q) = false, want true", id)
		}
	}
}

func TestIsValidProfile_Invalid(t *testing.T) {
	invalid := []string{"", "not_a_profile", "API_SERVICE", "api service"}
	for _, id := range invalid {
		if questions.IsValidProfile(id) {
			t.Errorf("IsValidProfile(%q) = true, want false", id)
		}
	}
}

// ---- IsValidFeatureArea tests ----

func TestIsValidFeatureArea_KnownAreas(t *testing.T) {
	known := []string{"authentication", "storage", "api", "testing", "observability", "deployment"}
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
	invalid := []string{"", "not_a_feature", "TESTING", "testing "}
	for _, id := range invalid {
		if questions.IsValidFeatureArea(id) {
			t.Errorf("IsValidFeatureArea(%q) = true, want false", id)
		}
	}
}

// ---- Question field tests ----

func TestQuestions_MultiFieldConsistency(t *testing.T) {
	// Any question with Multi: true must not have a single-value answer expected.
	// Here we just verify that known multi questions are indeed marked Multi.
	multiIDs := []string{"test_types", "auth_providers"}
	for _, id := range multiIDs {
		q, ok := questions.ByID(id)
		if !ok {
			t.Fatalf("ByID(%q): not found", id)
		}
		if !q.Multi {
			t.Errorf("question %q expected Multi=true, got false", id)
		}
	}
}

func TestQuestions_RequiredAndOptionalMix(t *testing.T) {
	// Verify that within _core there is at least one optional question.
	core := questions.ByFeatureArea("_core")
	hasOptional := false
	for _, q := range core {
		if !q.Required {
			hasOptional = true
			break
		}
	}
	if !hasOptional {
		t.Error("expected at least one optional _core question (e.g. team_size)")
	}
}
