package shared

// Resource swagg-doc:model
// Defines a Resource model
type Resource struct {
	LanguageCode string `json:"language-code"`
	Role         string `json:"role"`
	Department   string `json:"department"`
}

// Component defines a component
type Component struct {
	LanguageCode string `json:"language-code"`
	Role         string `json:"role"`
	Department   string `json:"department"`
}

// Translation defines a translation model
type Translation struct {
	Languages map[string]string `json:"languages"`
}
