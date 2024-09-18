// Utility package for error management in program execution due to errors that Dialosoft manages internally
package errorsUtils

import "errors"

var (
	// Error variable for when one of the parameters is null or empty
	ErrParameterCannotBeNull = errors.New("no parameter can be null")

	ErrUnauthorizedAcces            = errors.New("user or password is not valid")
	ErrInvalidUUID                  = errors.New("uuid is invalid")
	ErrRefreshTokenExpiredOrInvalid = errors.New("refreshToken expired or invalid")
	ErrRoleIDInRefreshToken         = errors.New("the roleID is empty in the refreshToken when trying to refresh")
)
