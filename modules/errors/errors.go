package errors

import (
	"errors"
	"fmt"
)

var (
	// StylesNotFound errors that styles were not found.
	StylesNotFound = errors.New("styles not found")

	// Not200Ok errors that the returned status code wasnt 200 Ok.
	Not200Ok = errors.New("status code wasn't OK 200")

	// NoAuthURL errors that the the paramater authURL was not set.
	NoAuthURL = errors.New("no authURL was set")

	// NoServiceDetected errors that the the paramater service was not set.
	NoServiceDetected = errors.New("no service detected")

	// NoSubject errors that the email builder didn't set the subject.
	NoSubject = errors.New("subject parameter is missing")

	// NoToParameter errors that the email builder didn't set the to.
	NoToParameter = errors.New("to parameter is missing")

	// NoParts errors that the email builder didn't specify any parts.
	NoParts = errors.New("no parts were detected")

	// NoPartBody errors that the email builder didn't specify the part's body.
	NoPartBody = errors.New("part doesn't contain body")

	// MessageSmall errors the given string is too small.
	MessageSmall = errors.New("message too small")

	// StyleNotFound errors the style hasn't been found.
	StyleNotFound = errors.New("style not found")

	// UserNotFound errors the user hasn't been found.
	UserNotFound = errors.New("user not found")

	// StyleNotFromUSO errors that the style being fetched isn't from the uso-archive.
	StyleNotFromUSO = errors.New("style isn't from uso-archive")

	// DuplicateStyle errors that an duplicate style has been found.
	DuplicateStyle = errors.New("duplicate style")

	// NoImportedStyles error that there has not found importeded styles.
	NoImportedStyles = errors.New("no imported styles")

	// FailedFetch errors that the fetch has failed to do it's operation.
	FailedFetch = errors.New("failed to fetch style")

	// FailedProcessData errors that it couldn't process the given style data.
	FailedProcessData = errors.New("failed to process style data")

	unexpectedSigningMethod = errors.New("unexpected jwt signing method")
)

// UnexpectedSigningMethod errors that a unexpected jwt signing method was used.
func UnexpectedSigningMethod(message string) error {
	return fmt.Errorf("%w=%s", unexpectedSigningMethod, message)
}
