package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/zyx-holdings/go-spec/internal/questions"
	"github.com/zyx-holdings/go-spec/internal/session"
)

// jsonAnswerEntry is the per-question entry in the JSON export answers map.
type jsonAnswerEntry struct {
	Value   session.AnswerValue `json:"value"`
	Skipped bool                `json:"skipped"`
}

// jsonMetadata holds summary statistics for the export (spec Section 8.1).
type jsonMetadata struct {
	TotalQuestions  int `json:"total_questions"`
	Answered        int `json:"answered"`
	Skipped         int `json:"skipped"`
	RequiredAnswered int `json:"required_answered"`
	RequiredTotal   int `json:"required_total"`
}

// jsonExport is the top-level JSON export document (spec Section 8.1).
type jsonExport struct {
	SpecgenVersion  string                      `json:"specgen_version"`
	SessionID       string                      `json:"session_id"`
	CreatedAt       string                      `json:"created_at"`
	CLIProfile      string                      `json:"cli_profile"`
	SelectedFeatures []string                   `json:"selected_features"`
	Answers         map[string]jsonAnswerEntry  `json:"answers"`
	Metadata        jsonMetadata                `json:"metadata"`
}

// ExportJSON writes sess as structured, pretty-printed JSON to
// <outputPrefix>.json. specgenVersion is embedded verbatim as the
// "specgen_version" field.
func ExportJSON(sess *session.Session, specgenVersion, outputPrefix string) error {
	answers := make(map[string]jsonAnswerEntry, len(sess.Answers))
	for id, ans := range sess.Answers {
		answers[id] = jsonAnswerEntry{
			Value:   ans.Value,
			Skipped: ans.Skipped,
		}
	}

	// Compute metadata over the questions relevant to this session.
	activeQs := questions.FilterByFeatures(sess.SelectedFeatures)
	var answered, skipped, requiredAnswered, requiredTotal int
	for _, q := range activeQs {
		if q.Required {
			requiredTotal++
		}
		ans, ok := sess.Answers[q.ID]
		if !ok {
			continue
		}
		if ans.Skipped {
			skipped++
		} else if !ans.Value.IsEmpty() {
			answered++
			if q.Required {
				requiredAnswered++
			}
		}
	}

	selectedFeatures := sess.SelectedFeatures
	if selectedFeatures == nil {
		selectedFeatures = []string{}
	}

	doc := jsonExport{
		SpecgenVersion:   specgenVersion,
		SessionID:        sess.ID,
		CreatedAt:        sess.CreatedAt.UTC().Format(time.RFC3339),
		CLIProfile:       sess.CLIProfile,
		SelectedFeatures: selectedFeatures,
		Answers:          answers,
		Metadata: jsonMetadata{
			TotalQuestions:  len(activeQs),
			Answered:        answered,
			Skipped:         skipped,
			RequiredAnswered: requiredAnswered,
			RequiredTotal:   requiredTotal,
		},
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("export json: marshal: %w", err)
	}

	path := outputPrefix + ".json"
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("export json: write %s: %w", path, err)
	}
	return nil
}
