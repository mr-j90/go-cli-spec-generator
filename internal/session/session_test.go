package session_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zyx-holdings/go-spec/internal/session"
)

// ---- AnswerValue tests ----

func TestAnswerValue_StringRoundTrip(t *testing.T) {
	v := session.NewStringValue("hello")
	if v.String() != "hello" {
		t.Errorf("String() = %q, want %q", v.String(), "hello")
	}
	if v.IsMulti() {
		t.Error("IsMulti() = true, want false")
	}
	if v.IsEmpty() {
		t.Error("IsEmpty() = true for non-empty string")
	}
}

func TestAnswerValue_MultiRoundTrip(t *testing.T) {
	v := session.NewMultiValue([]string{"a", "b", "c"})
	if !v.IsMulti() {
		t.Error("IsMulti() = false, want true")
	}
	got := v.Strings()
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("Strings() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Strings()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestAnswerValue_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		v     session.AnswerValue
		empty bool
	}{
		{"empty string", session.NewStringValue(""), true},
		{"non-empty string", session.NewStringValue("x"), false},
		{"empty multi", session.NewMultiValue([]string{}), true},
		{"non-empty multi", session.NewMultiValue([]string{"x"}), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.v.IsEmpty() != tc.empty {
				t.Errorf("IsEmpty() = %v, want %v", tc.v.IsEmpty(), tc.empty)
			}
		})
	}
}

// ---- Store method tests ----

func TestStore_SetAndGetAnswer(t *testing.T) {
	st := session.New()
	st.SetAnswer("project_name", session.NewStringValue("my-project"))

	ans, ok := st.GetAnswer("project_name")
	if !ok {
		t.Fatal("GetAnswer: not found")
	}
	if ans.QuestionID != "project_name" {
		t.Errorf("QuestionID = %q, want %q", ans.QuestionID, "project_name")
	}
	if ans.Value.String() != "my-project" {
		t.Errorf("Value = %q, want %q", ans.Value.String(), "my-project")
	}
	if ans.Skipped {
		t.Error("Skipped = true, want false")
	}
}

func TestStore_SkipAnswer(t *testing.T) {
	st := session.New()
	st.SkipAnswer("team_size")

	ans, ok := st.GetAnswer("team_size")
	if !ok {
		t.Fatal("GetAnswer: not found after skip")
	}
	if !ans.Skipped {
		t.Error("Skipped = false, want true")
	}
}

func TestStore_GetAnswer_Missing(t *testing.T) {
	st := session.New()
	_, ok := st.GetAnswer("nonexistent")
	if ok {
		t.Error("GetAnswer returned ok=true for missing question")
	}
}

func TestStore_SetAnswer_UpdatesTimestamp(t *testing.T) {
	st := session.New()
	before := st.Session().UpdatedAt
	time.Sleep(time.Millisecond)
	st.SetAnswer("project_name", session.NewStringValue("x"))
	after := st.Session().UpdatedAt
	if !after.After(before) {
		t.Error("UpdatedAt was not advanced after SetAnswer")
	}
}

func TestStore_SkipAnswer_UpdatesTimestamp(t *testing.T) {
	st := session.New()
	before := st.Session().UpdatedAt
	time.Sleep(time.Millisecond)
	st.SkipAnswer("team_size")
	after := st.Session().UpdatedAt
	if !after.After(before) {
		t.Error("UpdatedAt was not advanced after SkipAnswer")
	}
}

// ---- IsComplete tests ----

func answerCoreRequired(st *session.Store) {
	st.SetAnswer("project_name", session.NewStringValue("proj"))
	st.SetAnswer("project_description", session.NewStringValue("desc"))
	st.SetAnswer("primary_language", session.NewStringValue("Go"))
}

func TestStore_IsComplete_AllCoreAnswered(t *testing.T) {
	st := session.New()
	answerCoreRequired(st)
	if !st.IsComplete() {
		t.Error("IsComplete() = false, want true when all _core required questions answered")
	}
}

func TestStore_IsComplete_MissingRequired(t *testing.T) {
	st := session.New()
	st.SetAnswer("project_name", session.NewStringValue("proj"))
	// project_description and primary_language missing
	if st.IsComplete() {
		t.Error("IsComplete() = true, want false when required questions missing")
	}
}

func TestStore_IsComplete_SkippedRequired(t *testing.T) {
	st := session.New()
	answerCoreRequired(st)
	st.SkipAnswer("project_name") // override with skip
	if st.IsComplete() {
		t.Error("IsComplete() = true, want false when required question is skipped")
	}
}

func TestStore_IsComplete_OptionalSkipped(t *testing.T) {
	st := session.New()
	answerCoreRequired(st)
	st.SkipAnswer("team_size") // optional — should not affect completeness
	if !st.IsComplete() {
		t.Error("IsComplete() = false, want true when only optional questions are skipped")
	}
}

func TestStore_IsComplete_WithFeatures(t *testing.T) {
	st := session.New()
	answerCoreRequired(st)
	st.Session().SelectedFeatures = []string{"testing"}

	// Required testing questions not answered yet.
	if st.IsComplete() {
		t.Error("IsComplete() = true, want false when feature required questions unanswered")
	}

	st.SetAnswer("test_framework", session.NewStringValue("go test"))
	st.SetAnswer("test_types", session.NewMultiValue([]string{"unit", "integration"}))
	if !st.IsComplete() {
		t.Error("IsComplete() = false, want true after answering all required testing questions")
	}
}

// ---- Save / Load round-trip test ----

func TestStore_SaveLoad_RoundTrip(t *testing.T) {
	st := session.New()
	s := st.Session()
	s.CLIProfile = "api_service"
	s.SelectedFeatures = []string{"storage", "testing"}
	s.ExportFormats = []string{"markdown", "pdf"}
	s.Completed = true
	s.SessionState = "some-resume-state"

	st.SetAnswer("project_name", session.NewStringValue("my-api"))
	st.SetAnswer("project_description", session.NewStringValue("A REST API"))
	st.SetAnswer("primary_language", session.NewStringValue("Go"))
	st.SkipAnswer("team_size")
	st.SetAnswer("test_types", session.NewMultiValue([]string{"unit", "e2e"}))

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")

	if err := st.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := session.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	orig := st.Session()
	got := loaded.Session()

	if orig.ID != got.ID {
		t.Errorf("ID: got %q, want %q", got.ID, orig.ID)
	}
	if !orig.CreatedAt.Equal(got.CreatedAt) {
		t.Errorf("CreatedAt: got %v, want %v", got.CreatedAt, orig.CreatedAt)
	}
	if !orig.UpdatedAt.Equal(got.UpdatedAt) {
		t.Errorf("UpdatedAt: got %v, want %v", got.UpdatedAt, orig.UpdatedAt)
	}
	if orig.CLIProfile != got.CLIProfile {
		t.Errorf("CLIProfile: got %q, want %q", got.CLIProfile, orig.CLIProfile)
	}
	if orig.Completed != got.Completed {
		t.Errorf("Completed: got %v, want %v", got.Completed, orig.Completed)
	}
	if orig.SessionState != got.SessionState {
		t.Errorf("SessionState: got %q, want %q", got.SessionState, orig.SessionState)
	}

	// SelectedFeatures
	if len(orig.SelectedFeatures) != len(got.SelectedFeatures) {
		t.Errorf("SelectedFeatures len: got %d, want %d", len(got.SelectedFeatures), len(orig.SelectedFeatures))
	} else {
		for i, f := range orig.SelectedFeatures {
			if got.SelectedFeatures[i] != f {
				t.Errorf("SelectedFeatures[%d]: got %q, want %q", i, got.SelectedFeatures[i], f)
			}
		}
	}

	// ExportFormats
	if len(orig.ExportFormats) != len(got.ExportFormats) {
		t.Errorf("ExportFormats len: got %d, want %d", len(got.ExportFormats), len(orig.ExportFormats))
	}

	// Answers
	for id, origAns := range orig.Answers {
		gotAns, ok := got.Answers[id]
		if !ok {
			t.Errorf("Answers[%q]: missing after load", id)
			continue
		}
		if origAns.Skipped != gotAns.Skipped {
			t.Errorf("Answers[%q].Skipped: got %v, want %v", id, gotAns.Skipped, origAns.Skipped)
		}
		if origAns.Value.IsMulti() != gotAns.Value.IsMulti() {
			t.Errorf("Answers[%q].Value.IsMulti: got %v, want %v", id, gotAns.Value.IsMulti(), origAns.Value.IsMulti())
		}
		if !origAns.Value.IsMulti() && origAns.Value.String() != gotAns.Value.String() {
			t.Errorf("Answers[%q].Value.String: got %q, want %q", id, gotAns.Value.String(), origAns.Value.String())
		}
		if origAns.Value.IsMulti() {
			os_, gs_ := origAns.Value.Strings(), gotAns.Value.Strings()
			if len(os_) != len(gs_) {
				t.Errorf("Answers[%q].Value.Strings len: got %d, want %d", id, len(gs_), len(os_))
			} else {
				for i, v := range os_ {
					if gs_[i] != v {
						t.Errorf("Answers[%q].Value.Strings[%d]: got %q, want %q", id, i, gs_[i], v)
					}
				}
			}
		}
	}
}

func TestStore_Load_FileNotFound(t *testing.T) {
	_, err := session.Load("/nonexistent/path/session.json")
	if err == nil {
		t.Error("Load: expected error for missing file, got nil")
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{invalid json}"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := session.Load(path)
	if err == nil {
		t.Error("Load: expected error for invalid JSON, got nil")
	}
}

// ---- Validate tests (spec Section 13.2) ----

// validStore returns a Store that passes all validation checks.
func validStore() *session.Store {
	st := session.New()
	s := st.Session()
	s.CLIProfile = "api_service"
	s.SelectedFeatures = []string{"testing"}
	st.SetAnswer("project_name", session.NewStringValue("proj"))
	st.SetAnswer("project_description", session.NewStringValue("desc"))
	st.SetAnswer("primary_language", session.NewStringValue("Go"))
	st.SetAnswer("test_framework", session.NewStringValue("go test"))
	st.SetAnswer("test_types", session.NewMultiValue([]string{"unit"}))
	return st
}

func TestValidate_ValidSession(t *testing.T) {
	st := validStore()
	if err := st.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestValidate_InvalidProfile(t *testing.T) {
	st := validStore()
	st.Session().CLIProfile = "not_a_real_profile"
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for invalid profile")
	}
	errs, ok := err.(session.ValidationErrors)
	if !ok {
		t.Fatalf("error is not ValidationErrors: %T", err)
	}
	found := false
	for _, e := range errs {
		if e.Field == "cli_profile" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for field cli_profile, got: %v", err)
	}
}

func TestValidate_EmptyProfile(t *testing.T) {
	st := validStore()
	st.Session().CLIProfile = ""
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for empty profile")
	}
}

func TestValidate_InvalidFeatureArea(t *testing.T) {
	st := validStore()
	st.Session().SelectedFeatures = append(st.Session().SelectedFeatures, "not_a_feature")
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for invalid feature area")
	}
	errs := err.(session.ValidationErrors)
	found := false
	for _, e := range errs {
		if e.Field == "selected_features" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for field selected_features, got: %v", err)
	}
}

func TestValidate_RequiredQuestionMissing(t *testing.T) {
	st := validStore()
	delete(st.Session().Answers, "project_name")
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for missing required question")
	}
	errs := err.(session.ValidationErrors)
	found := false
	for _, e := range errs {
		if e.Field == "answers.project_name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for answers.project_name, got: %v", err)
	}
}

func TestValidate_RequiredQuestionSkipped(t *testing.T) {
	st := validStore()
	st.SkipAnswer("project_name") // override required answer with a skip
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error when required question is skipped")
	}
	errs := err.(session.ValidationErrors)
	found := false
	for _, e := range errs {
		if e.Field == "answers.project_name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for answers.project_name, got: %v", err)
	}
}

func TestValidate_RequiredQuestionEmptyAnswer(t *testing.T) {
	st := validStore()
	st.SetAnswer("project_name", session.NewStringValue(""))
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for empty required answer")
	}
	errs := err.(session.ValidationErrors)
	found := false
	for _, e := range errs {
		if e.Field == "answers.project_name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for answers.project_name, got: %v", err)
	}
}

func TestValidate_OptionalQuestionMissing_IsAllowed(t *testing.T) {
	st := validStore()
	// team_size is optional — not having it should not produce an error
	delete(st.Session().Answers, "team_size")
	if err := st.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil when optional question is absent", err)
	}
}

func TestValidate_OptionalQuestionSkipped_IsAllowed(t *testing.T) {
	st := validStore()
	st.SkipAnswer("team_size") // optional, skipped — should be fine
	if err := st.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil when optional question is skipped", err)
	}
}

func TestValidate_CollectsAllErrors(t *testing.T) {
	// Set up a session with multiple problems.
	st := session.New()
	st.Session().CLIProfile = "bad_profile"
	st.Session().SelectedFeatures = []string{"bad_feature"}
	// No answers at all.

	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want multiple errors")
	}
	errs := err.(session.ValidationErrors)
	// Expect: cli_profile error + selected_features error + 3 missing _core required answers
	// (project_name, project_description, primary_language)
	if len(errs) < 3 {
		t.Errorf("expected at least 3 validation errors, got %d: %v", len(errs), errs)
	}
}

func TestValidate_FeatureRequiredQuestions(t *testing.T) {
	// A session with "deployment" feature selected but deployment questions unanswered.
	st := session.New()
	st.Session().CLIProfile = "api_service"
	st.Session().SelectedFeatures = []string{"deployment"}
	answerCoreRequired(st)

	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want errors for unanswered deployment questions")
	}
	errs := err.(session.ValidationErrors)
	wantFields := map[string]bool{
		"answers.deployment_target": true,
		"answers.containerized":     true,
	}
	for _, e := range errs {
		delete(wantFields, e.Field)
	}
	if len(wantFields) > 0 {
		t.Errorf("missing expected error fields: %v", wantFields)
	}
}

func TestValidate_MultiValueEmpty(t *testing.T) {
	st := validStore()
	st.SetAnswer("test_types", session.NewMultiValue([]string{})) // empty multi
	err := st.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for empty multi-value on required question")
	}
}

func TestValidate_ValidationErrors_Error(t *testing.T) {
	errs := session.ValidationErrors{
		{Field: "cli_profile", Message: "invalid"},
		{Field: "answers.foo", Message: "required question has no answer"},
	}
	got := errs.Error()
	if got == "" {
		t.Error("ValidationErrors.Error() returned empty string")
	}
}
