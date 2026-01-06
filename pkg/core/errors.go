package core

type ErrorSeverity int

const (
	ErrorSeverityWarning ErrorSeverity = iota
	ErrorSeverityError
	ErrorSeverityFatal
)

type GenerationError struct {
	Severity ErrorSeverity
	Phase    string // "parsing", "validation", "generation"
	Location *Location
	Message  string
	Hint     string // Suggestion to fix
	Code     string // Error code (e.g., "E1001")
}

type Location struct {
	File   string // Spec file or template file
	Line   int
	Column int
	Path   string // JSON path (e.g., "#/paths/users/get")
}

func (e GenerationError) Error() string {
	return e.Message
}

type ValidationError struct {
	GenerationError
}
