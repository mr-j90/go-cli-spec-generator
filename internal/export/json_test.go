package export_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zyx-holdings/go-spec/internal/export"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// buildSession constructs a minimal Session for testing purposes.
func buildSession() *session.Session {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	s := &session.Session{
		ID:               "test-session-id",
		CreatedAt:        now,
		CLIProfile:       "oneshot",
		SelectedFeatures: []string{"storage"},
		Answers:          make(map[string]session.Answer),
	}
	// _core required: project_name, project_description, primary_language
	s.Answers["project_name"] = session.Answer{
		QuestionID: "project_name",
		Value:      session.NewStringValue("myapp"),
		Skipped:    false,
	}
	s.Answers["project_description"] = session.Answer{
		QuestionID: "project_description",
		Value:      session.NewStringValue("A test app"),
		Skipped:    false,
	}
	s.Answers["primary_language"] = session.Answer{
		QuestionID: "primary_language",
		Value:      session.NewStringValue("Go"),
		Skipped:    false,
	}
	// storage required: storage_type, database_name
	s.Answers["storage_type"] = session.Answer{
		QuestionID: "storage_type",
		Value:      session.NewStringValue("Relational DB"),
		Skipped:    false,
	}
	s.Answers["database_name"] = session.Answer{
		QuestionID: "database_name",
		Value:      session.NewStringValue("PostgreSQL"),
		Skipped:    false,
	}
	// optional skipped
	s.Answers["team_size"] = session.Answer{
		QuestionID: "team_size",
		Skipped:    true,
	}
	return s
}

// readJSON reads and unmarshals the JSON file produced by ExportJSON.
func readJSON(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readJSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("readJSON: invalid JSON: %v", err)
	}
	return out
}

// ---- ExportJSON tests ----

func TestExportJSON_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: unexpected error: %v", err)
	}

	info, err := os.Stat(prefix + ".json")
	if err != nil {
		t.Fatalf("ExportJSON: output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("ExportJSON: output file is empty")
	}
}

func TestExportJSON_ProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: unexpected error: %v", err)
	}

	data, err := os.ReadFile(prefix + ".json")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !json.Valid(data) {
		t.Errorf("ExportJSON: output is not valid JSON:\n%s", data)
	}
}

func TestExportJSON_TopLevelFields(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")
	sess := buildSession()

	if err := export.ExportJSON(sess, "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")

	checkStringField(t, doc, "specgen_version", "0.1.0")
	checkStringField(t, doc, "session_id", "test-session-id")
	checkStringField(t, doc, "created_at", "2024-06-01T12:00:00Z")
	checkStringField(t, doc, "cli_profile", "oneshot")
}

func TestExportJSON_SelectedFeatures(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	raw, ok := doc["selected_features"]
	if !ok {
		t.Fatal("missing field: selected_features")
	}
	arr, ok := raw.([]interface{})
	if !ok {
		t.Fatalf("selected_features: want array, got %T", raw)
	}
	if len(arr) != 1 || arr[0] != "storage" {
		t.Errorf("selected_features = %v, want [storage]", arr)
	}
}

func TestExportJSON_SelectedFeaturesNilBecomesEmptyArray(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	sess := buildSession()
	sess.SelectedFeatures = nil

	if err := export.ExportJSON(sess, "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	raw, ok := doc["selected_features"]
	if !ok {
		t.Fatal("missing field: selected_features")
	}
	arr, ok := raw.([]interface{})
	if !ok {
		t.Fatalf("selected_features: want array, got %T", raw)
	}
	if len(arr) != 0 {
		t.Errorf("selected_features = %v, want []", arr)
	}
}

func TestExportJSON_AnswersMap(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	rawAnswers, ok := doc["answers"]
	if !ok {
		t.Fatal("missing field: answers")
	}
	answersMap, ok := rawAnswers.(map[string]interface{})
	if !ok {
		t.Fatalf("answers: want object, got %T", rawAnswers)
	}

	// Check a known answered question.
	rawEntry, ok := answersMap["project_name"]
	if !ok {
		t.Fatal("answers: missing key project_name")
	}
	entry, ok := rawEntry.(map[string]interface{})
	if !ok {
		t.Fatalf("answers.project_name: want object, got %T", rawEntry)
	}
	if entry["value"] != "myapp" {
		t.Errorf("answers.project_name.value = %v, want myapp", entry["value"])
	}
	if entry["skipped"] != false {
		t.Errorf("answers.project_name.skipped = %v, want false", entry["skipped"])
	}

	// Check a skipped question.
	rawSkipped, ok := answersMap["team_size"]
	if !ok {
		t.Fatal("answers: missing key team_size")
	}
	skippedEntry, ok := rawSkipped.(map[string]interface{})
	if !ok {
		t.Fatalf("answers.team_size: want object, got %T", rawSkipped)
	}
	if skippedEntry["skipped"] != true {
		t.Errorf("answers.team_size.skipped = %v, want true", skippedEntry["skipped"])
	}
}

func TestExportJSON_Metadata(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	rawMeta, ok := doc["metadata"]
	if !ok {
		t.Fatal("missing field: metadata")
	}
	meta, ok := rawMeta.(map[string]interface{})
	if !ok {
		t.Fatalf("metadata: want object, got %T", rawMeta)
	}

	// _core (4 qs, 3 required) + storage (3 qs, 2 required) = 7 total, 5 required
	checkNumericField(t, meta, "total_questions", 7)
	checkNumericField(t, meta, "required_total", 5)

	// 5 answered (project_name, project_description, primary_language, storage_type, database_name)
	checkNumericField(t, meta, "answered", 5)

	// 1 skipped (team_size)
	checkNumericField(t, meta, "skipped", 1)

	// all 5 required questions are answered
	checkNumericField(t, meta, "required_answered", 5)
}

func TestExportJSON_MetadataKeys(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	meta, ok := doc["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("metadata: missing or wrong type")
	}

	for _, key := range []string{"total_questions", "answered", "skipped", "required_answered", "required_total"} {
		if _, exists := meta[key]; !exists {
			t.Errorf("metadata: missing key %q", key)
		}
	}
}

func TestExportJSON_PrettyPrinted(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "0.1.0", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	data, err := os.ReadFile(prefix + ".json")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	// Pretty-printed JSON contains newlines and indentation.
	content := string(data)
	hasNewlines := false
	hasIndent := false
	for _, line := range []string{"\n", "  "} {
		switch line {
		case "\n":
			for _, c := range content {
				if c == '\n' {
					hasNewlines = true
					break
				}
			}
		case "  ":
			for i := 0; i < len(content)-1; i++ {
				if content[i] == ' ' && content[i+1] == ' ' {
					hasIndent = true
					break
				}
			}
		}
	}
	if !hasNewlines {
		t.Error("ExportJSON: output is not pretty-printed (no newlines)")
	}
	if !hasIndent {
		t.Error("ExportJSON: output is not pretty-printed (no 2-space indentation)")
	}
}

func TestExportJSON_InvalidPath_ReturnsError(t *testing.T) {
	sess := buildSession()
	err := export.ExportJSON(sess, "0.1.0", "/nonexistent/dir/output")
	if err == nil {
		t.Error("ExportJSON(invalid path): expected error, got nil")
	}
}

func TestExportJSON_SpecgenVersion(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "output")

	if err := export.ExportJSON(buildSession(), "9.9.9", prefix); err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}

	doc := readJSON(t, prefix+".json")
	checkStringField(t, doc, "specgen_version", "9.9.9")
}

// ---- helpers ----

func checkStringField(t *testing.T, doc map[string]interface{}, key, want string) {
	t.Helper()
	raw, ok := doc[key]
	if !ok {
		t.Errorf("missing field: %s", key)
		return
	}
	got, ok := raw.(string)
	if !ok {
		t.Errorf("%s: want string, got %T", key, raw)
		return
	}
	if got != want {
		t.Errorf("%s = %q, want %q", key, got, want)
	}
}

func checkNumericField(t *testing.T, doc map[string]interface{}, key string, want int) {
	t.Helper()
	raw, ok := doc[key]
	if !ok {
		t.Errorf("missing field: %s", key)
		return
	}
	// JSON numbers unmarshal to float64.
	got, ok := raw.(float64)
	if !ok {
		t.Errorf("%s: want number, got %T", key, raw)
		return
	}
	if int(got) != want {
		t.Errorf("%s = %d, want %d", key, int(got), want)
	}
}
