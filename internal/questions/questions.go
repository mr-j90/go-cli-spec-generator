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

// InputType is the kind of input expected for a question.
type InputType = string

const (
	InputText        InputType = "text"
	InputSelect      InputType = "select"
	InputMultiSelect InputType = "multiselect"
	InputConfirm     InputType = "confirm"
)

// Question defines a single question in the spec flow.
type Question struct {
	ID          string
	FeatureArea string    // "_core" or a feature area ID
	Text        string
	InputType   InputType
	Placeholder string
	Options     []string // valid choices for select/multiselect
	Required    bool
	Order       int // relative order within the feature area (1-based)
}

// Profiles is the registry of valid CLI profile IDs.
// Profiles represent the interaction pattern of the CLI tool (spec Section 4.1.1).
var Profiles = map[string]Profile{
	"oneshot": {
		ID:          "oneshot",
		DisplayName: "One-shot",
		Description: "Runs once, performs a task, and exits (e.g. ls, curl)",
	},
	"daemon": {
		ID:          "daemon",
		DisplayName: "Daemon",
		Description: "Runs continuously as a background service or long-lived process",
	},
	"subcommand": {
		ID:          "subcommand",
		DisplayName: "Subcommand",
		Description: "Provides multiple subcommands dispatched by the first argument (e.g. git, docker)",
	},
	"hybrid": {
		ID:          "hybrid",
		DisplayName: "Hybrid",
		Description: "Supports both one-shot invocation and daemon/service mode",
	},
}

// FeatureAreas is the registry of valid feature area IDs.
// "_core" is always implicitly included and is not listed here (spec Section 4.1.2).
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

// All is the ordered list of all questions in the system (spec Sections 5.2 and 5.3).
var All = []Question{
	// ── _core questions — always included (spec Section 5.2) ──────────────────
	{
		ID:          "tool_name",
		FeatureArea: "_core",
		Text:        "What is the name of your CLI tool?",
		InputType:   InputText,
		Placeholder: "mytool",
		Required:    true,
		Order:       1,
	},
	{
		ID:          "tool_description",
		FeatureArea: "_core",
		Text:        "Describe what your CLI tool does in one or two sentences.",
		InputType:   InputText,
		Placeholder: "A tool that...",
		Required:    true,
		Order:       2,
	},
	{
		ID:          "go_module",
		FeatureArea: "_core",
		Text:        "What is the Go module path?",
		InputType:   InputText,
		Placeholder: "github.com/org/mytool",
		Required:    true,
		Order:       3,
	},
	{
		ID:          "binary_name",
		FeatureArea: "_core",
		Text:        "What is the binary executable name?",
		InputType:   InputText,
		Placeholder: "mytool",
		Required:    false,
		Order:       4,
	},
	{
		ID:          "target_os",
		FeatureArea: "_core",
		Text:        "Which OS platforms will this tool support?",
		InputType:   InputMultiSelect,
		Options:     []string{"linux", "darwin", "windows"},
		Required:    false,
		Order:       5,
	},
	{
		ID:          "min_go_version",
		FeatureArea: "_core",
		Text:        "What is the minimum Go version required?",
		InputType:   InputSelect,
		Options:     []string{"1.21", "1.22", "1.23", "1.24"},
		Required:    false,
		Order:       6,
	},
	{
		ID:          "license",
		FeatureArea: "_core",
		Text:        "Which open-source license will this project use?",
		InputType:   InputSelect,
		Options:     []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "none"},
		Required:    false,
		Order:       7,
	},
	{
		ID:          "feature_areas",
		FeatureArea: "_core",
		Text:        "Which feature areas does this CLI require?",
		InputType:   InputMultiSelect,
		Options:     []string{"api", "database", "filesystem", "concurrency", "retry", "config", "hooks", "auth", "logging", "http", "templates", "statemachine"},
		Required:    false,
		Order:       8,
	},

	// ── api (spec Section 5.3.1) ─────────────────────────────────────────────
	{
		ID:          "api_base_url",
		FeatureArea: "api",
		Text:        "What is the base URL of the external API?",
		InputType:   InputText,
		Placeholder: "https://api.example.com/v1",
		Required:    true,
		Order:       1,
	},
	{
		ID:          "api_auth_type",
		FeatureArea: "api",
		Text:        "How does the API authenticate requests?",
		InputType:   InputSelect,
		Options:     []string{"none", "api-key", "bearer-token", "basic-auth", "oauth2"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "api_response_format",
		FeatureArea: "api",
		Text:        "What response format does the API use?",
		InputType:   InputSelect,
		Options:     []string{"json", "xml", "text", "protobuf"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "api_pagination",
		FeatureArea: "api",
		Text:        "Does the API use pagination?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "api_timeout_ms",
		FeatureArea: "api",
		Text:        "What is the default request timeout in milliseconds?",
		InputType:   InputText,
		Placeholder: "5000",
		Required:    false,
		Order:       5,
	},

	// ── database (spec Section 5.3.2) ────────────────────────────────────────
	{
		ID:          "db_driver",
		FeatureArea: "database",
		Text:        "Which database driver will you use?",
		InputType:   InputSelect,
		Options:     []string{"postgres", "mysql", "sqlite", "mongodb"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "db_orm",
		FeatureArea: "database",
		Text:        "Which ORM or query builder will you use?",
		InputType:   InputSelect,
		Options:     []string{"sqlx", "gorm", "ent", "raw-sql"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "db_schema",
		FeatureArea: "database",
		Text:        "Describe the primary schema or data model.",
		InputType:   InputText,
		Placeholder: "users, posts, ...",
		Required:    true,
		Order:       3,
	},
	{
		ID:          "db_migrations",
		FeatureArea: "database",
		Text:        "Which migration tool will you use?",
		InputType:   InputSelect,
		Options:     []string{"goose", "golang-migrate", "atlas", "none"},
		Required:    false,
		Order:       4,
	},
	{
		ID:          "db_connection_pool_size",
		FeatureArea: "database",
		Text:        "What is the maximum connection pool size?",
		InputType:   InputText,
		Placeholder: "10",
		Required:    false,
		Order:       5,
	},

	// ── filesystem (spec Section 5.3.3) ──────────────────────────────────────
	{
		ID:          "fs_operations",
		FeatureArea: "filesystem",
		Text:        "What filesystem operations will this CLI perform?",
		InputType:   InputMultiSelect,
		Options:     []string{"read", "write", "watch", "traverse", "delete"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "fs_base_path",
		FeatureArea: "filesystem",
		Text:        "What is the default base path for filesystem operations?",
		InputType:   InputText,
		Placeholder: "/var/data",
		Required:    false,
		Order:       2,
	},
	{
		ID:          "fs_permissions",
		FeatureArea: "filesystem",
		Text:        "What file permission mode should created files use?",
		InputType:   InputSelect,
		Options:     []string{"0600", "0644", "0755", "0777"},
		Required:    false,
		Order:       3,
	},
	{
		ID:          "fs_atomic_writes",
		FeatureArea: "filesystem",
		Text:        "Use atomic writes (write-then-rename pattern)?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "fs_watch_events",
		FeatureArea: "filesystem",
		Text:        "Which filesystem events should trigger actions?",
		InputType:   InputMultiSelect,
		Options:     []string{"create", "modify", "delete", "rename"},
		Required:    false,
		Order:       5,
	},

	// ── concurrency (spec Section 5.3.4) ─────────────────────────────────────
	{
		ID:          "concurrency_model",
		FeatureArea: "concurrency",
		Text:        "What concurrency model will be used?",
		InputType:   InputSelect,
		Options:     []string{"goroutines", "worker-pool", "pipeline", "fan-out"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "concurrency_max_workers",
		FeatureArea: "concurrency",
		Text:        "What is the maximum number of concurrent workers?",
		InputType:   InputText,
		Placeholder: "8",
		Required:    true,
		Order:       2,
	},
	{
		ID:          "concurrency_cancellation",
		FeatureArea: "concurrency",
		Text:        "How will cancellation be handled?",
		InputType:   InputSelect,
		Options:     []string{"context", "signal", "timeout", "none"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "concurrency_sync_primitives",
		FeatureArea: "concurrency",
		Text:        "Which synchronization primitives will be used?",
		InputType:   InputMultiSelect,
		Options:     []string{"mutex", "rwmutex", "channel", "atomic", "waitgroup"},
		Required:    false,
		Order:       4,
	},
	{
		ID:          "concurrency_error_strategy",
		FeatureArea: "concurrency",
		Text:        "How will concurrent errors be aggregated?",
		InputType:   InputSelect,
		Options:     []string{"first-error", "all-errors", "log-and-continue"},
		Required:    false,
		Order:       5,
	},

	// ── retry (spec Section 5.3.5) ───────────────────────────────────────────
	{
		ID:          "retry_strategy",
		FeatureArea: "retry",
		Text:        "What backoff strategy will be used?",
		InputType:   InputSelect,
		Options:     []string{"fixed", "linear", "exponential", "exponential-with-jitter"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "retry_max_attempts",
		FeatureArea: "retry",
		Text:        "What is the maximum number of retry attempts?",
		InputType:   InputText,
		Placeholder: "3",
		Required:    true,
		Order:       2,
	},
	{
		ID:          "retry_retryable_errors",
		FeatureArea: "retry",
		Text:        "Which error types should trigger a retry?",
		InputType:   InputMultiSelect,
		Options:     []string{"network", "timeout", "5xx", "rate-limit", "all"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "retry_initial_delay_ms",
		FeatureArea: "retry",
		Text:        "What is the initial retry delay in milliseconds?",
		InputType:   InputText,
		Placeholder: "100",
		Required:    false,
		Order:       4,
	},
	{
		ID:          "retry_max_delay_ms",
		FeatureArea: "retry",
		Text:        "What is the maximum retry delay cap in milliseconds?",
		InputType:   InputText,
		Placeholder: "30000",
		Required:    false,
		Order:       5,
	},

	// ── config (spec Section 5.3.6) ──────────────────────────────────────────
	{
		ID:          "config_sources",
		FeatureArea: "config",
		Text:        "Where will configuration come from?",
		InputType:   InputMultiSelect,
		Options:     []string{"config-file", "env-vars", "cli-flags", "defaults"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "config_format",
		FeatureArea: "config",
		Text:        "What file format will be used for configuration files?",
		InputType:   InputSelect,
		Options:     []string{"yaml", "toml", "json", "ini"},
		Required:    false,
		Order:       2,
	},
	{
		ID:          "config_library",
		FeatureArea: "config",
		Text:        "Which configuration library will you use?",
		InputType:   InputSelect,
		Options:     []string{"viper", "koanf", "envconfig", "cobra-flags", "none"},
		Required:    false,
		Order:       3,
	},
	{
		ID:          "config_validation",
		FeatureArea: "config",
		Text:        "How will configuration be validated?",
		InputType:   InputSelect,
		Options:     []string{"struct-tags", "manual", "cue", "none"},
		Required:    false,
		Order:       4,
	},
	{
		ID:          "config_hot_reload",
		FeatureArea: "config",
		Text:        "Should configuration support hot-reload without restart?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       5,
	},

	// ── hooks (spec Section 5.3.7) ───────────────────────────────────────────
	{
		ID:          "hooks_lifecycle",
		FeatureArea: "hooks",
		Text:        "At which lifecycle points will hooks run?",
		InputType:   InputMultiSelect,
		Options:     []string{"pre-run", "post-run", "pre-command", "post-command"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "hooks_mechanism",
		FeatureArea: "hooks",
		Text:        "How are hooks implemented?",
		InputType:   InputSelect,
		Options:     []string{"function-callbacks", "plugin-files", "shell-scripts"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "hooks_error_handling",
		FeatureArea: "hooks",
		Text:        "How should hook errors be handled?",
		InputType:   InputSelect,
		Options:     []string{"abort", "warn-and-continue", "ignore"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "hooks_async",
		FeatureArea: "hooks",
		Text:        "Should hooks run asynchronously?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "hooks_registration",
		FeatureArea: "hooks",
		Text:        "How are hooks registered?",
		InputType:   InputSelect,
		Options:     []string{"config-file", "code-registration", "both"},
		Required:    false,
		Order:       5,
	},

	// ── auth (spec Section 5.3.8) ────────────────────────────────────────────
	{
		ID:          "auth_type",
		FeatureArea: "auth",
		Text:        "What authentication type is used?",
		InputType:   InputSelect,
		Options:     []string{"none", "api-key", "bearer-token", "basic-auth", "oauth2", "certificate"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "auth_credential_storage",
		FeatureArea: "auth",
		Text:        "Where are credentials stored?",
		InputType:   InputSelect,
		Options:     []string{"os-keychain", "config-file", "env-vars", "memory-only"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "auth_token_refresh",
		FeatureArea: "auth",
		Text:        "How is token refresh handled?",
		InputType:   InputSelect,
		Options:     []string{"automatic", "manual", "not-applicable"},
		Required:    false,
		Order:       3,
	},
	{
		ID:          "auth_multi_account",
		FeatureArea: "auth",
		Text:        "Does the CLI support multiple accounts or profiles?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "auth_command_name",
		FeatureArea: "auth",
		Text:        "What is the name of the authentication subcommand?",
		InputType:   InputText,
		Placeholder: "login",
		Required:    false,
		Order:       5,
	},

	// ── logging (spec Section 5.3.9) ─────────────────────────────────────────
	{
		ID:          "logging_library",
		FeatureArea: "logging",
		Text:        "Which logging library will you use?",
		InputType:   InputSelect,
		Options:     []string{"zerolog", "zap", "slog", "logrus", "log"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "logging_format",
		FeatureArea: "logging",
		Text:        "What log output format will be used?",
		InputType:   InputSelect,
		Options:     []string{"json", "text", "logfmt"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "logging_output",
		FeatureArea: "logging",
		Text:        "Where will logs be written?",
		InputType:   InputSelect,
		Options:     []string{"stdout", "stderr", "file", "both"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "logging_levels",
		FeatureArea: "logging",
		Text:        "Which log levels will be used?",
		InputType:   InputMultiSelect,
		Options:     []string{"debug", "info", "warn", "error", "fatal"},
		Required:    false,
		Order:       4,
	},
	{
		ID:          "logging_correlation_id",
		FeatureArea: "logging",
		Text:        "Include correlation IDs in log entries?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       5,
	},

	// ── http (spec Section 5.3.10) ───────────────────────────────────────────
	{
		ID:          "http_role",
		FeatureArea: "http",
		Text:        "What role does this CLI play in HTTP?",
		InputType:   InputSelect,
		Options:     []string{"client-only", "server-only", "both"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "http_client_library",
		FeatureArea: "http",
		Text:        "Which HTTP client library will you use?",
		InputType:   InputSelect,
		Options:     []string{"net/http", "resty", "gentleman"},
		Required:    false,
		Order:       2,
	},
	{
		ID:          "http_server_framework",
		FeatureArea: "http",
		Text:        "Which HTTP server framework will you use?",
		InputType:   InputSelect,
		Options:     []string{"net/http", "chi", "gin", "echo", "fiber"},
		Required:    false,
		Order:       3,
	},
	{
		ID:          "http_middleware",
		FeatureArea: "http",
		Text:        "Which HTTP middleware will be used?",
		InputType:   InputMultiSelect,
		Options:     []string{"logging", "auth", "rate-limiting", "cors", "compression"},
		Required:    false,
		Order:       4,
	},
	{
		ID:          "http_tls",
		FeatureArea: "http",
		Text:        "Will TLS be enabled?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       5,
	},

	// ── templates (spec Section 5.3.11) ──────────────────────────────────────
	{
		ID:          "template_engine",
		FeatureArea: "templates",
		Text:        "Which template engine will be used?",
		InputType:   InputSelect,
		Options:     []string{"text/template", "html/template", "pongo2", "jet"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "template_source",
		FeatureArea: "templates",
		Text:        "Where are templates stored?",
		InputType:   InputSelect,
		Options:     []string{"embedded", "files", "string-literals"},
		Required:    true,
		Order:       2,
	},
	{
		ID:          "template_output_format",
		FeatureArea: "templates",
		Text:        "What format will templates render to?",
		InputType:   InputSelect,
		Options:     []string{"text", "html", "markdown", "yaml", "json"},
		Required:    true,
		Order:       3,
	},
	{
		ID:          "template_helpers",
		FeatureArea: "templates",
		Text:        "Will custom template helper functions be defined?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "template_caching",
		FeatureArea: "templates",
		Text:        "Should parsed templates be cached after first load?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       5,
	},

	// ── statemachine (spec Section 5.3.12) ───────────────────────────────────
	{
		ID:          "sm_library",
		FeatureArea: "statemachine",
		Text:        "Which state machine library will you use?",
		InputType:   InputSelect,
		Options:     []string{"looplab/fsm", "stateless", "custom"},
		Required:    true,
		Order:       1,
	},
	{
		ID:          "sm_initial_state",
		FeatureArea: "statemachine",
		Text:        "What is the name of the initial state?",
		InputType:   InputText,
		Placeholder: "idle",
		Required:    true,
		Order:       2,
	},
	{
		ID:          "sm_persistence",
		FeatureArea: "statemachine",
		Text:        "How will state be persisted between runs?",
		InputType:   InputSelect,
		Options:     []string{"in-memory", "file", "database", "none"},
		Required:    false,
		Order:       3,
	},
	{
		ID:          "sm_visualization",
		FeatureArea: "statemachine",
		Text:        "Generate a state diagram in documentation?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       4,
	},
	{
		ID:          "sm_concurrent_transitions",
		FeatureArea: "statemachine",
		Text:        "Allow concurrent state transitions?",
		InputType:   InputConfirm,
		Required:    false,
		Order:       5,
	},
}

// ByFeatureArea returns all questions belonging to the given feature area,
// preserving their Order within the section.
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

// FilterByFeatures returns the _core questions followed by the questions for
// each of the selected feature areas, in the order provided.
func FilterByFeatures(selected []string) []Question {
	result := ByFeatureArea("_core")
	for _, area := range selected {
		result = append(result, ByFeatureArea(area)...)
	}
	return result
}
