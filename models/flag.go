package models

// Flags represents feature flags.
type Flags struct {
	Welcome         bool `json:"welcome,omitempty"`
	Sidebar         bool `json:"sidebar,omitempty"`
	SearchAutofocus bool `json:"search_autofocus,omitempty"`
	ViewRedesign    bool `json:"view_redesign,omitempty"`
}
