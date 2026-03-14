// Package session handles saving and restoring spec generation sessions.
package session

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zyx-holdings/go-spec/internal/questions"
)

// AnswerValue holds either a single string or a list of strings.
// It serializes to/from JSON as a plain string or array of strings.
type AnswerValue struct {
	single  string
	multi   []string
	isMulti bool
}

// NewStringValue creates an AnswerValue holding a single string.
func NewStringValue(s string) AnswerValue {
	return AnswerValue{single: s, isMulti: false}
}

// NewMultiValue creates an AnswerValue holding multiple strings.
func NewMultiValue(ss []string) AnswerValue {
	cp := make([]string, len(ss))
	copy(cp, ss)
	return AnswerValue{multi: cp, isMulti: true}
}

// String returns the single-string value. Empty string if multi.
func (v AnswerValue) String() string { return v.single }

// Strings returns the multi-string value. Nil if single.
func (v AnswerValue) Strings() []string { return v.multi }

// IsMulti reports whether the value is a list of strings.
func (v AnswerValue) IsMulti() bool { return v.isMulti }

// IsEmpty reports whether the value contains no meaningful content.
func (v AnswerValue) IsEmpty() bool {
	if v.isMulti {
		return len(v.multi) == 0
	}
	return v.single == ""
}

// MarshalJSON serializes as a JSON string or JSON array of strings.
func (v AnswerValue) MarshalJSON() ([]byte, error) {
	if v.isMulti {
		return json.Marshal(v.multi)
	}
	return json.Marshal(v.single)
}

// UnmarshalJSON deserializes from a JSON string or JSON array of strings.
func (v *AnswerValue) UnmarshalJSON(data []byte) error {
	// Try string first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		v.single = s
		v.isMulti = false
		return nil
	}
	// Try []string.
	var ss []string
	if err := json.Unmarshal(data, &ss); err == nil {
		v.multi = ss
		v.isMulti = true
		return nil
	}
	return fmt.Errorf("answer value must be a string or array of strings, got: %s", data)
}

// Answer represents a single answer to a spec question.
type Answer struct {
	QuestionID string      `json:"question_id"`
	Value      AnswerValue `json:"value"`
	Skipped    bool        `json:"skipped"`
}

// Session holds all state for a spec generation session.
type Session struct {
	ID               string            `json:"id"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	CLIProfile       string            `json:"cli_profile"`
	SelectedFeatures []string          `json:"selected_features"`
	Answers          map[string]Answer `json:"answers"`
	Completed        bool              `json:"completed"`
	ExportFormats    []string          `json:"export_formats"`
	SessionState     string            `json:"_session_state,omitempty"`
}

// Store wraps a Session and provides answer management operations.
type Store struct {
	s *Session
}

// New creates a new Store with an initialized, empty session.
func New() *Store {
	now := time.Now().UTC()
	return &Store{
		s: &Session{
			ID:        newUUID(),
			CreatedAt: now,
			UpdatedAt: now,
			Answers:   make(map[string]Answer),
		},
	}
}

// Session returns the underlying Session (read-only use intended).
func (st *Store) Session() *Session {
	return st.s
}

// SetAnswer records a non-skipped answer for the given question ID.
func (st *Store) SetAnswer(questionID string, value AnswerValue) {
	st.s.Answers[questionID] = Answer{
		QuestionID: questionID,
		Value:      value,
		Skipped:    false,
	}
	st.s.UpdatedAt = time.Now().UTC()
}

// SkipAnswer marks a question as explicitly skipped.
func (st *Store) SkipAnswer(questionID string) {
	st.s.Answers[questionID] = Answer{
		QuestionID: questionID,
		Skipped:    true,
	}
	st.s.UpdatedAt = time.Now().UTC()
}

// GetAnswer returns the Answer for a question and whether it exists.
func (st *Store) GetAnswer(questionID string) (Answer, bool) {
	a, ok := st.s.Answers[questionID]
	return a, ok
}

// IsComplete reports whether all required questions for _core and the
// session's selected feature areas have non-empty, non-skipped answers.
func (st *Store) IsComplete() bool {
	areas := append([]string{"_core"}, st.s.SelectedFeatures...)
	for _, area := range areas {
		for _, q := range questions.ByFeatureArea(area) {
			if !q.Required {
				continue
			}
			ans, ok := st.s.Answers[q.ID]
			if !ok || ans.Skipped || ans.Value.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// Save serializes the session to JSON and writes it to path.
func (st *Store) Save(path string) error {
	data, err := json.MarshalIndent(st.s, "", "  ")
	if err != nil {
		return fmt.Errorf("session: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("session: write %s: %w", path, err)
	}
	return nil
}

// Load deserializes a session from the JSON file at path and returns a Store.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("session: read %s: %w", path, err)
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("session: unmarshal: %w", err)
	}
	if s.Answers == nil {
		s.Answers = make(map[string]Answer)
	}
	return &Store{s: &s}, nil
}

// ValidationError represents a single validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation failures.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	msgs := make([]string, len(ve))
	for i, e := range ve {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// Validate checks the session against the spec rules. It collects all
// validation errors and returns them together rather than failing fast.
//
// Rules:
//   - cli_profile must be a valid profile ID
//   - each selected_features entry must be a valid feature area ID
//   - all required questions for _core + selected features must have
//     non-empty, non-skipped answers
//   - optional questions may be absent (treated as skipped — no error)
func (st *Store) Validate() error {
	var errs ValidationErrors

	// Validate cli_profile.
	if !questions.IsValidProfile(st.s.CLIProfile) {
		errs = append(errs, ValidationError{
			Field:   "cli_profile",
			Message: fmt.Sprintf("invalid profile ID %q", st.s.CLIProfile),
		})
	}

	// Validate each selected_features entry.
	for _, f := range st.s.SelectedFeatures {
		if !questions.IsValidFeatureArea(f) {
			errs = append(errs, ValidationError{
				Field:   "selected_features",
				Message: fmt.Sprintf("invalid feature area ID %q", f),
			})
		}
	}

	// Validate required questions for _core + selected features.
	areas := append([]string{"_core"}, st.s.SelectedFeatures...)
	for _, area := range areas {
		for _, q := range questions.ByFeatureArea(area) {
			if !q.Required {
				// Optional — missing or skipped is fine.
				continue
			}
			ans, ok := st.s.Answers[q.ID]
			if !ok {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("answers.%s", q.ID),
					Message: "required question has no answer",
				})
				continue
			}
			if ans.Skipped {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("answers.%s", q.ID),
					Message: "required question cannot be skipped",
				})
				continue
			}
			if ans.Value.IsEmpty() {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("answers.%s", q.ID),
					Message: "required question has an empty answer",
				})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// newUUID generates a random UUID v4 string.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
