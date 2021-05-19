package errors_helper

import (
	"errors"
	"fmt"
)

var (
	// ErrStylesNotFound errors that styles were not found.
	ErrStylesNotFound = errors.New("styles not found")

	// ErrNot200Ok errors that the returned status code wasnt 200 Ok.
	ErrNot200Ok = errors.New("status code wasn't OK 200")

	// ErrNoAuthURL errors that the the paramater authURL was not set.
	ErrNoAuthURL = errors.New("no authURL was set")

	// ErrNoServiceDetected errors that the the paramater service was not set.
	ErrNoServiceDetected = errors.New("no service detected")

	// ErrNoSubject errors that the email builder didn't set the subject.
	ErrNoSubject = errors.New("subject parameter is missing")

	// ErrNoToParameter errors that the email builder didn't set the to.
	ErrNoToParameter = errors.New("to parameter is missing")

	// ErrNoParts errors that the email builder didn't specify any parts.
	ErrNoParts = errors.New("no parts were detected")

	// ErrNoPartBody errors that the email builder didn't specify the part's body.
	ErrNoPartBody = errors.New("part doesn't contain body")

	// ErrMessageSmall errors the given string is too small.
	ErrMessageSmall = errors.New("message too small")

	// ErrStyleNotFound errors the style hasn't been found.
	ErrStyleNotFound = errors.New("style not found")

	// ErrUserNotFound errors the user hasn't been found.
	ErrUserNotFound = errors.New("user not found")

	errUnexpectedSigningMethod = errors.New("unexpected jwt signing method")

	// ErrStyleNotFromUSO errors that the style being fetched isn't from the uso-archive.
	ErrStyleNotFromUSO = errors.New("style isn't from uso-archive")

	// ErrDuplicateStyle errors that an duplicate style has been found.
	ErrDuplicateStyle = errors.New("duplicate style")

	// ErrNoImportedStyles error that there has not found importeded styles.
	ErrNoImportedStyles = errors.New("no imported styles")

	// ErrFailedFetch errors that the fetch has failed to do it's operation.
	ErrFailedFetch = errors.New("failed to fetch style")

	// ErrFailedProcessData errors that it couldn't process the given style data.
	ErrFailedProcessData = errors.New("failed to process style data")
)

// ErrUnexpectedSigningMethod errors that a unexpected jwt signing method was used.
func ErrUnexpectedSigningMethod(message string) error {
	return fmt.Errorf("%w=%s", errUnexpectedSigningMethod, message)
}