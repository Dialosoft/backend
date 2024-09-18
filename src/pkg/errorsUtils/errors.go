// Utility package for error management in program execution due to errors that Dialosoft manages internally
package errorsUtils

import "errors"

var (
	// Error variable for when one of the parameters is null or empty
	ErrParameterCannotBeNull = errors.New("no parameter can be null")
	ErrUnauthorizedAcces     = errors.New("user or password is not valid")
)
