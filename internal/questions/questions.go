// Package questions defines the question flow for interactive spec generation.
package questions

// Input type constants for question widgets.
const (
	InputTypeText         = "text"
	InputTypeTextArea     = "textarea"
	InputTypeSingleSelect = "single_select"
	InputTypeMultiSelect  = "multi_select"
	InputTypeConfirm      = "confirm"
)

// Profile represents a CLI profile type for a project.
type Profile struct {
	ID          string
	DisplayName string
	Description string
}

// FeatureArea represents a feature area for spec generation.
type FeatureArea struct {
	ID          string
	DisplayName string
	Description string
}

// Question defines a single question in the spec flow.
type Question struct {
	ID          string
	FeatureArea string   // "_core" or a feature area ID
	Text        string
	Required    bool
	Multi       bool     // true if the answer is []string (kept for compatibility)
	InputType   string   // one of the InputType* constants
	Options     []string // non-empty for single_select and multi_select
	Placeholder string   // hint text for text and textarea widgets
}

// Profiles is the registry of valid CLI profile IDs.
var Profiles = map[string]Profile{
	"api_service": {
		ID:          "api_service",
		DisplayName: "API Service",
		Description: "Backend API service with REST or gRPC endpoints",
	},
	"web_app": {
		ID:          "web_app",
		DisplayName: "Web Application",
		Description: "Web application with frontend and backend",
	},
	"cli_tool": {
		ID:          "cli_tool",
		DisplayName: "CLI Tool",
		Description: "Command-line interface tool",
	},
	"library": {
		ID:          "library",
		DisplayName: "Library",
		Description: "Reusable library or package",
	},
	"data_pipeline": {
		ID:          "data_pipeline",
		DisplayName: "Data Pipeline",
		Description: "Data processing or ETL pipeline",
	},
}

// FeatureAreas is the registry of valid feature area IDs.
// "_core" is always implicitly included and is not listed here.
var FeatureAreas = map[string]FeatureArea{
	"authentication": {
		ID:          "authentication",
		DisplayName: "Authentication",
		Description: "User authentication and authorization",
	},
	"storage": {
		ID:          "storage",
		DisplayName: "Storage",
		Description: "Data storage and persistence",
	},
	"api": {
		ID:          "api",
		DisplayName: "API",
		Description: "External API endpoints",
	},
	"testing": {
		ID:          "testing",
		DisplayName: "Testing",
		Description: "Testing strategy and coverage",
	},
	"observability": {
		ID:          "observability",
		DisplayName: "Observability",
		Description: "Logging, metrics, and tracing",
	},
	"deployment": {
		ID:          "deployment",
		DisplayName: "Deployment",
		Description: "Deployment and infrastructure",
	},
}

// All is the ordered list of all questions in the system.
var All = []Question{
	// _core questions — always included
	{
		ID: "project_name", FeatureArea: "_core",
		Text: "What is the project name?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. my-awesome-service",
	},
	{
		ID: "project_description", FeatureArea: "_core",
		Text: "Describe the project in one or two sentences.", Required: true,
		InputType: InputTypeTextArea, Placeholder: "e.g. A REST API that manages user accounts and billing",
	},
	{
		ID: "primary_language", FeatureArea: "_core",
		Text: "What is the primary programming language?", Required: true,
		InputType: InputTypeSingleSelect,
		Options:   []string{"Go", "Python", "JavaScript", "TypeScript", "Rust", "Java", "C#", "C++", "Other"},
	},
	{
		ID: "team_size", FeatureArea: "_core",
		Text: "What is the team size?", Required: false,
		InputType: InputTypeSingleSelect,
		Options:   []string{"1", "2–5", "6–15", "16–50", "50+"},
	},

	// authentication
	{
		ID: "auth_method", FeatureArea: "authentication",
		Text: "What authentication method will be used?", Required: true,
		InputType: InputTypeSingleSelect,
		Options:   []string{"JWT", "OAuth 2.0", "Session-based", "API Keys", "SAML", "Other"},
	},
	{
		ID: "auth_providers", FeatureArea: "authentication",
		Text: "List the OAuth/OIDC providers (if any).", Required: false, Multi: true,
		InputType: InputTypeMultiSelect,
		Options:   []string{"Google", "GitHub", "Microsoft", "Okta", "Auth0", "Other"},
	},
	{
		ID: "session_management", FeatureArea: "authentication",
		Text: "How will sessions be managed?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. Redis-backed sessions with 24h TTL",
	},

	// storage
	{
		ID: "storage_type", FeatureArea: "storage",
		Text: "What type of storage will be used?", Required: true,
		InputType: InputTypeSingleSelect,
		Options:   []string{"Relational DB", "NoSQL DB", "Object Storage", "In-Memory", "Mixed"},
	},
	{
		ID: "database_name", FeatureArea: "storage",
		Text: "What database technology will be used?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. PostgreSQL, MongoDB, SQLite",
	},
	{
		ID: "caching_strategy", FeatureArea: "storage",
		Text: "Describe the caching strategy.", Required: false,
		InputType: InputTypeText, Placeholder: "e.g. Redis with 5-minute TTL on hot paths",
	},

	// api
	{
		ID: "api_style", FeatureArea: "api",
		Text: "What API style will be used?", Required: true,
		InputType: InputTypeSingleSelect,
		Options:   []string{"REST", "gRPC", "GraphQL", "Mixed"},
	},
	{
		ID: "api_versioning", FeatureArea: "api",
		Text: "How will the API be versioned?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. URL path (/v1/), Accept header",
	},
	{
		ID: "api_auth", FeatureArea: "api",
		Text: "How will the API authenticate requests?", Required: false,
		InputType: InputTypeText, Placeholder: "e.g. Bearer token, mTLS, API key",
	},

	// testing
	{
		ID: "test_framework", FeatureArea: "testing",
		Text: "What testing framework will be used?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. Jest, pytest, testify, JUnit",
	},
	{
		ID: "test_coverage_target", FeatureArea: "testing",
		Text: "What is the target test coverage percentage?", Required: false,
		InputType: InputTypeText, Placeholder: "e.g. 80",
	},
	{
		ID: "test_types", FeatureArea: "testing",
		Text: "What types of tests will be written?", Required: true, Multi: true,
		InputType: InputTypeMultiSelect,
		Options:   []string{"Unit", "Integration", "End-to-End", "Contract", "Load", "Smoke"},
	},

	// observability
	{
		ID: "logging_framework", FeatureArea: "observability",
		Text: "What logging framework will be used?", Required: true,
		InputType: InputTypeText, Placeholder: "e.g. zap, logrus, slog, winston",
	},
	{
		ID: "metrics_platform", FeatureArea: "observability",
		Text: "What metrics platform will be used?", Required: false,
		InputType: InputTypeText, Placeholder: "e.g. Prometheus, Datadog, CloudWatch",
	},
	{
		ID: "tracing_enabled", FeatureArea: "observability",
		Text: "Will distributed tracing be enabled?", Required: true,
		InputType: InputTypeConfirm,
	},

	// deployment
	{
		ID: "deployment_target", FeatureArea: "deployment",
		Text: "Where will the application be deployed?", Required: true,
		InputType: InputTypeSingleSelect,
		Options:   []string{"Cloud (AWS)", "Cloud (GCP)", "Cloud (Azure)", "On-Premises", "Hybrid", "Serverless"},
	},
	{
		ID: "containerized", FeatureArea: "deployment",
		Text: "Will the application be containerized?", Required: true,
		InputType: InputTypeConfirm,
	},
	{
		ID: "ci_cd_platform", FeatureArea: "deployment",
		Text: "What CI/CD platform will be used?", Required: false,
		InputType: InputTypeText, Placeholder: "e.g. GitHub Actions, Jenkins, GitLab CI",
	},
}

// ByFeatureArea returns all questions belonging to the given feature area.
func ByFeatureArea(area string) []Question {
	var result []Question
	for _, q := range All {
		if q.FeatureArea == area {
			result = append(result, q)
		}
	}
	return result
}

// ByID returns the question with the given ID, and whether it was found.
func ByID(id string) (Question, bool) {
	for _, q := range All {
		if q.ID == id {
			return q, true
		}
	}
	return Question{}, false
}

// IsValidProfile reports whether id is a known profile ID.
func IsValidProfile(id string) bool {
	_, ok := Profiles[id]
	return ok
}

// IsValidFeatureArea reports whether id is a known feature area ID.
// "_core" is not a selectable feature area, so it returns false.
func IsValidFeatureArea(id string) bool {
	_, ok := FeatureAreas[id]
	return ok
}
