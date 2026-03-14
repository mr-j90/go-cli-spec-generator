// Package questions defines the question flow for interactive spec generation.
package questions

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
	Icon        string
}

// Question defines a single question in the spec flow.
type Question struct {
	ID          string
	FeatureArea string // "_core" or a feature area ID
	Text        string
	Required    bool
	Multi       bool // true if the answer is []string
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
		Icon:        "🔐",
	},
	"storage": {
		ID:          "storage",
		DisplayName: "Storage",
		Description: "Data storage and persistence",
		Icon:        "💾",
	},
	"api": {
		ID:          "api",
		DisplayName: "API",
		Description: "External API endpoints",
		Icon:        "🔌",
	},
	"testing": {
		ID:          "testing",
		DisplayName: "Testing",
		Description: "Testing strategy and coverage",
		Icon:        "🧪",
	},
	"observability": {
		ID:          "observability",
		DisplayName: "Observability",
		Description: "Logging, metrics, and tracing",
		Icon:        "📈",
	},
	"deployment": {
		ID:          "deployment",
		DisplayName: "Deployment",
		Description: "Deployment and infrastructure",
		Icon:        "🚀",
	},
	"security": {
		ID:          "security",
		DisplayName: "Security",
		Description: "Security hardening and vulnerability management",
		Icon:        "🛡",
	},
	"caching": {
		ID:          "caching",
		DisplayName: "Caching",
		Description: "Caching strategies and performance optimization",
		Icon:        "⚡",
	},
	"messaging": {
		ID:          "messaging",
		DisplayName: "Messaging",
		Description: "Message queues, pub/sub, and event streaming",
		Icon:        "📨",
	},
	"search": {
		ID:          "search",
		DisplayName: "Search",
		Description: "Full-text search and indexing",
		Icon:        "🔍",
	},
	"notifications": {
		ID:          "notifications",
		DisplayName: "Notifications",
		Description: "Email, push notifications, and webhooks",
		Icon:        "🔔",
	},
	"configuration": {
		ID:          "configuration",
		DisplayName: "Configuration",
		Description: "App configuration, feature flags, and environment management",
		Icon:        "⚙",
	},
}

// All is the ordered list of all questions in the system.
var All = []Question{
	// _core questions — always included
	{ID: "project_name", FeatureArea: "_core", Text: "What is the project name?", Required: true},
	{ID: "project_description", FeatureArea: "_core", Text: "Describe the project in one or two sentences.", Required: true},
	{ID: "primary_language", FeatureArea: "_core", Text: "What is the primary programming language?", Required: true},
	{ID: "team_size", FeatureArea: "_core", Text: "What is the team size?", Required: false},

	// authentication
	{ID: "auth_method", FeatureArea: "authentication", Text: "What authentication method will be used?", Required: true},
	{ID: "auth_providers", FeatureArea: "authentication", Text: "List the OAuth/OIDC providers (if any).", Required: false, Multi: true},
	{ID: "session_management", FeatureArea: "authentication", Text: "How will sessions be managed?", Required: true},

	// storage
	{ID: "storage_type", FeatureArea: "storage", Text: "What type of storage will be used?", Required: true},
	{ID: "database_name", FeatureArea: "storage", Text: "What database technology will be used?", Required: true},
	{ID: "caching_strategy", FeatureArea: "storage", Text: "Describe the caching strategy.", Required: false},

	// api
	{ID: "api_style", FeatureArea: "api", Text: "What API style will be used (REST, gRPC, GraphQL)?", Required: true},
	{ID: "api_versioning", FeatureArea: "api", Text: "How will the API be versioned?", Required: true},
	{ID: "api_auth", FeatureArea: "api", Text: "How will the API authenticate requests?", Required: false},

	// testing
	{ID: "test_framework", FeatureArea: "testing", Text: "What testing framework will be used?", Required: true},
	{ID: "test_coverage_target", FeatureArea: "testing", Text: "What is the target test coverage percentage?", Required: false},
	{ID: "test_types", FeatureArea: "testing", Text: "What types of tests will be written?", Required: true, Multi: true},

	// observability
	{ID: "logging_framework", FeatureArea: "observability", Text: "What logging framework will be used?", Required: true},
	{ID: "metrics_platform", FeatureArea: "observability", Text: "What metrics platform will be used?", Required: false},
	{ID: "tracing_enabled", FeatureArea: "observability", Text: "Will distributed tracing be enabled?", Required: true},

	// deployment
	{ID: "deployment_target", FeatureArea: "deployment", Text: "Where will the application be deployed?", Required: true},
	{ID: "containerized", FeatureArea: "deployment", Text: "Will the application be containerized?", Required: true},
	{ID: "ci_cd_platform", FeatureArea: "deployment", Text: "What CI/CD platform will be used?", Required: false},
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
