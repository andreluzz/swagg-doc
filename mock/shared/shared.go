package shared

// Resource swagg-doc:model
// Defines a Resource model
type Resource struct {
	LanguageCode string `json:"language-code"`
	Role         string `json:"role"`
	Department   string `json:"department"`
}
