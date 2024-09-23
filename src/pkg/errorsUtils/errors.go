// Utility package for error management in program execution due to errors that Dialosoft manages internally
package errorsUtils

import "errors"

var (
	// ErrParameterCannotBeNull is an error for when one of the parameters is null or empty.
	// Used to indicate that a required parameter is missing or not provided.
	ErrParameterCannotBeNull = errors.New("no parameter can be null")

	// ErrUnauthorizedAcces is an error for unauthorized access.
	// Returned when the user credentials (username or password) are incorrect during authentication.
	ErrUnauthorizedAcces = errors.New("user or password is not valid")

	// ErrInvalidUUID is an error for an invalid UUID.
	// Indicates that the provided UUID does not conform to the expected format.
	ErrInvalidUUID = errors.New("uuid is invalid")

	// ErrRefreshTokenExpiredOrInvalid is an error for an expired or invalid refresh token.
	// Used when the refresh token is either expired or has an invalid format.
	ErrRefreshTokenExpiredOrInvalid = errors.New("refreshToken expired or invalid")

	// ErrRoleIDInRefreshToken is an error indicating the roleID is missing in the refresh token.
	// Occurs when attempting to refresh and the refresh token does not contain a valid roleID.
	ErrRoleIDInRefreshToken = errors.New("the roleID is empty in the refreshToken when trying to refresh")

	// ErrNotFound is an error indicating that a requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrInternalServer is a generic error for an internal server issue.
	// Used when an unexpected error occurs on the server side.
	ErrInternalServer = errors.New("internal server error")

	// ErrTokenBlacklisted is an error for a blacklisted token.
	// Returned when a token that has been blacklisted is used.
	ErrTokenBlacklisted = errors.New("token is blacklisted")
)
