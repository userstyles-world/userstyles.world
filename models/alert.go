package models

type alertKind int

const (
	alertSuccess alertKind = iota
)

// Alert represents an alert message and its kind.
type Alert struct {
	kind    alertKind
	Message string
}

// Success checks if alert is of success kind.
func (b *Alert) Success() bool { return b.kind == alertSuccess }

// NewSuccessAlert is a helper for creating alerts with success kind.
func NewSuccessAlert(message string) *Alert {
	return &Alert{
		kind:    alertSuccess,
		Message: message,
	}
}
