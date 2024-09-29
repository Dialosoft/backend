package errorsUtils

import "errors"

var (
	ErrFailedToRetrievePosts = errors.New("failed to retrieve posts due to a system error")

	ErrFailedToFindPost = errors.New("failed to find the post due to a system error")

	// ErrPostNotFound is returned when the requested post cannot be found or has been deleted.
	ErrPostNotFound = errors.New("the post you are looking for does not exist or has been deleted")

	// ErrNoPostsObtained is returned when no posts match the query criteria.
	ErrNoPostsObtained = errors.New("no posts were found based on the given criteria")

	// ErrUserNotFound is returned when the requested user cannot be found or is deactivated.
	ErrUserNotFound = errors.New("the user does not exist, is deactivated, or banned")

	// ErrPostCreationFailed is returned when there is a failure creating a new post.
	ErrPostCreationFailed = errors.New("failed to create a new post due to system or validation issues")

	// ErrPostUpdateFailed is returned when a post update operation fails.
	ErrPostUpdateFailed = errors.New("unable to update the post, possibly due to invalid data or the post being deleted")

	// ErrPostDeletionFailed is returned when a post cannot be deleted.
	ErrPostDeletionFailed = errors.New("failed to delete the post, it may have already been deleted or an error occurred")

	// ErrUserUnauthorized is returned when a user attempts an unauthorized action.
	ErrUserUnauthorized = errors.New("you are not authorized to perform this action")

	// ErrPostAlreadyLiked is returned when a user tries to like a post they have already liked.
	ErrPostAlreadyLiked = errors.New("you have already liked this post")

	// ErrPostLikeFailed is returned when an error occurs while liking a post.
	ErrPostLikeFailed = errors.New("unable to like the post due to a system error")

	// ErrInvalidPaginationParameters is returned when pagination parameters are invalid.
	ErrInvalidPaginationParameters = errors.New("invalid pagination parameters: limit or offset are out of bounds")

	// ErrPostHasNoLikes is returned when a post has no recorded likes.
	ErrPostHasNoLikes = errors.New("the post has no likes yet")

	// ErrPostContentInvalid is returned when the content of a post is invalid.
	ErrPostContentInvalid = errors.New("the content of the post does not meet the required validation rules")

	// ErrDatabaseConnection is returned when there is a failure connecting to the database.
	ErrDatabaseConnection = errors.New("unable to connect to the database")

	// ErrPostRestorationFailed is returned when a post restoration operation fails.
	ErrPostRestorationFailed = errors.New("failed to restore the post due to a system error")
)
