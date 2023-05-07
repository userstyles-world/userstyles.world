package models

// Flags represents feature flags.
type Flags struct {
	Welcome         bool `json:"welcome"`
	Sidebar         bool `json:"sidebar"`
	SearchAutofocus bool `json:"searchautofocus"`
}
