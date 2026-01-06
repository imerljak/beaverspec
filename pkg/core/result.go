package core

type GenerationResult struct {
	Files    []GeneratedFile
	Warnings []GenerationError
}

type GeneratedFile struct {
	Path     string                 // Relative path from output dir
	Content  []byte                 // File content
	Metadata map[string]interface{} // Optional metadata
}
