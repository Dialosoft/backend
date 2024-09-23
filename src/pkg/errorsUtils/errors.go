// Utility package for error management in program execution due to errors that Dialosoft manages internally
package errorsUtils

import "errors"

var (
	// Error variable for when one of the parameters is null or empty
	// Used to indicate that a required parameter is missing or not provided
	ErrParameterCannotBeNull = errors.New("no parameter can be null")

	// Error for unauthorized access
	// Returned when the user or password is incorrect during authentication
	ErrUnauthorizedAcces = errors.New("user or password is not valid")

	// Error for invalid UUID
	// Indicates that the provided UUID does not conform to the expected format
	ErrInvalidUUID = errors.New("uuid is invalid")

	// Error for expired or invalid refresh token
	// Used when the refresh token is either expired or has an invalid format
	ErrRefreshTokenExpiredOrInvalid = errors.New("refreshToken expired or invalid")

	// Error indicating the roleID is empty in the refresh token
	// Occurs when attempting to refresh and the refresh token does not contain a roleID
	ErrRoleIDInRefreshToken = errors.New("the roleID is empty in the refreshToken when trying to refresh")

	ErrNotFound = errors.New("not found")

	ErrInternalServer = errors.New("internal server error")

	ErrTokenBlacklisted = errors.New("token is blacklisted")
)
