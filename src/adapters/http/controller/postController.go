package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostController struct {
	PostService services.PostService
	Layer       string
}

func NewPostController(postService services.PostService, Layer string) *PostController {
	return &PostController{PostService: postService, Layer: Layer}
}

// GetAllPostsByForum retrieves posts for a specific forum.
//
// @Summary Get all posts by forum ID
// @Description Retrieves all posts within a given forum specified by its ID, with pagination options.
// @Tags Posts
// @Accept json
// @Produce json
// @Param forumID path string true "Forum ID"
// @Param limit query int false "Limit of posts (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} response.PostResponse "List of posts"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - No posts found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /posts/get-all-posts-by-forum/{forumID} [get]
func (pc *PostController) GetAllPostsByForum(c fiber.Ctx) error {
	limit := c.Query("limit")
	offset := c.Query("offset")

	forumID := c.Params("id")
	forumUUID, err := uuid.Parse(forumID)
	if err != nil {
		return response.ErrUUIDParse(c, forumID)
	}

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	responses, err := pc.PostService.GetAllPostsByForum(forumUUID, limitInt, offsetInt)
	if err != nil {
		return response.ErrInternalServer(c, err, responses, pc.Layer)
	}

	if responses == nil {
		return response.ErrNotFound(c, pc.Layer)
	}

	return response.Standard(c, "OK", responses)
}

// GetAllPosts retrieves all posts in the system.
//
// @Summary Get all posts
// @Description Retrieves a list of all posts with pagination options.
// @Tags Posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit of posts (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} response.PostResponse "List of posts"
// @Failure 404 {object} response.StandardError "NOT FOUND - No posts found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/get-all-posts [get]
func (pc *PostController) GetAllPosts(c fiber.Ctx) error {

	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	posts, err := pc.PostService.GetAllPosts(limitInt, offsetInt)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, posts, pc.Layer)
	}

	return response.Standard(c, "OK", posts)
}

// GetPostByID retrieves a post by its ID.
//
// @Summary Get post by ID
// @Description Retrieves a post by its unique identifier.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.PostResponse "Post data"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/get-post-by-id/{id} [get]
func (pc *PostController) GetPostByID(c fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c, postID)
	}

	post, err := pc.PostService.GetPostByID(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, post, pc.Layer)
	}

	return response.Standard(c, "OK", post)
}

// GetPostsByUserID retrieves posts by a specific user.
//
// @Summary Get posts by user ID
// @Description Retrieves all posts created by a specific user.
// @Tags Posts
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} response.PostResponse "List of posts by user"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - Posts not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /posts/get-posts-by-user-id/{userID} [get]
func (pc *PostController) GetPostsByUserID(c fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return response.ErrUUIDParse(c, userID)
	}

	posts, err := pc.PostService.GetPostsByUserID(userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, posts, pc.Layer)
	}

	return response.Standard(c, "OK", posts)
}

// GetAllPostsAndReturnSimpleResponse retrieves all posts with minimal data.
//
// @Summary Get all posts (simple response)
// @Description Retrieves a list of all posts with a simplified response format, with pagination options.
// @Tags Posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit of posts (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} response.SimplePostResponse "List of simple post responses"
// @Failure 404 {object} response.StandardError "NOT FOUND - No posts found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /posts/get-all-posts-simple [get]
func (pc *PostController) GetAllPostsAndReturnSimpleResponse(c fiber.Ctx) error {
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	posts, err := pc.PostService.GetAllPostsAndReturnSimpleResponse(limitInt, offsetInt)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, posts, pc.Layer)
	}
	return response.Standard(c, "OK", posts)
}

// GetPostNumberOfLikes retrieves the number of likes for a specific post.
//
// @Summary Get number of likes for a post
// @Description Retrieves the count of likes for a post by its unique identifier.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} int "Number of likes"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid post ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /posts/get-like-count/{id} [get]
func (pc *PostController) GetPostNumberOfLikes(c fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c, postID)
	}

	likes, err := pc.PostService.GetLikeCount(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, likes, pc.Layer)
	}

	return response.Standard(c, "OK", likes)
}

// CreateNewPost creates a new post in the system.
//
// @Summary Create a new post
// @Description Creates a new post with the provided data.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body request.NewPost true "New Post Data"
// @Success 201 {object} response.PostResponse "CREATED - New post created successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/authenticated/create-new-post [post]
func (pc *PostController) CreateNewPost(c fiber.Ctx) error {
	var req request.NewPost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	post, err := pc.PostService.CreateNewPost(uuid.MustParse(req.UserID), req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, post, pc.Layer)
	}

	return response.StandardCreated(c, "CREATED", post)
}

// UpdatePostTitle updates the title of an existing post.
//
// @Summary Update post title
// @Description Updates the title of a specific post by its ID.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body request.UpdatePostTitle true "Updated Post Title Data"
// @Success 200 {object} response.StandardResponse "UPDATED - Post title updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/update-post-title/{id} [put]
func (pc *PostController) UpdatePostTitle(c fiber.Ctx) error {
	var req request.UpdatePostTitle

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	postUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c, req.PostID)
	}

	err = pc.PostService.UpdatePostTitle(postUUID, req.Title)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

// UpdatePostContent updates the content of an existing post.
//
// @Summary Update post content
// @Description Updates the content of a specific post by its ID.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body request.UpdatePostContent true "Updated Post Content Data"
// @Success 200 {object} response.StandardResponse "UPDATED - Post content updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/update-post-content/{id} [put]
func (pc *PostController) UpdatePostContent(c fiber.Ctx) error {
	var req request.UpdatePostContent

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	postUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c, req.PostID)
	}

	err = pc.PostService.UpdatePostContent(postUUID, req.Content)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

// DeletePost deletes an existing post by its ID.
//
// @Summary Delete post
// @Description Deletes a specific post from the system by its ID.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.StandardResponse "DELETED - Post deleted successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid post ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/delete-post/{id} [delete]
func (pc *PostController) DeletePost(c fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c, postID)
	}

	err = pc.PostService.DeletePost(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "DELETED", nil)
}

// RestorePost restores a previously deleted post by its ID.
//
// @Summary Restore post
// @Description Restores a previously deleted post by its unique identifier.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.StandardResponse "RESTORED - Post restored successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid post ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/protected/restore-post/{id} [put]
func (pc *PostController) RestorePost(c fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c, postID)
	}

	err = pc.PostService.RestorePost(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "RESTORED", nil)
}

// LikePost allows a user to like a specific post.
//
// @Summary Like post
// @Description Adds a like to a specific post from a user.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body request.LikeOrUnlikePost true "Post and User ID for liking"
// @Success 200 {object} response.StandardResponse "LIKED - Post liked successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post or user not found"
// @Failure 409 {object} response.StandardError "CONFLICT - Post already liked by the user"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/authenticated/like-post [put]
func (pc *PostController) LikePost(c fiber.Ctx) error {
	var req request.LikeOrUnlikePost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.ErrUUIDParse(c, req.UserID)
	}
	postUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c, req.PostID)
	}

	err = pc.PostService.LikePost(postUUID, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.PersonalizedErr(c, "You already liked this post", fiber.StatusConflict)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "LIKED", nil)
}

// UnlikePost allows a user to remove their like from a specific post.
//
// @Summary Unlike post
// @Description Removes a like from a specific post by a user.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body request.LikeOrUnlikePost true "Post and User ID for unliking"
// @Success 200 {object} response.StandardResponse "UNLIKED - Post unliked successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Post or user not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /posts/authenticated/unlike-post [put]
func (pc *PostController) UnlikePost(c fiber.Ctx) error {
	var req request.LikeOrUnlikePost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, pc.Layer)
	}

	userUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c, req.PostID)
	}
	postUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.ErrUUIDParse(c, req.UserID)
	}

	err = pc.PostService.UnlikePost(postUUID, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, pc.Layer)
	}

	return response.Standard(c, "UNLIKED", nil)
}

// GetPostLikesByUserID retrieves the posts liked by a specific user.
//
// @Summary Get posts liked by user ID
// @Description Retrieves all post IDs that a user has liked.
// @Tags Posts
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} map[string]interface{} "List of post IDs liked by the user"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid user ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - No liked posts found for the user"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /posts/get-post-likes-by-user-id/{userID} [get]
func (pc *PostController) GetPostLikesByUserID(c fiber.Ctx) error {
	postID := c.Params("userID")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c, postID)
	}

	likes, err := pc.PostService.GetPostLikesByUserID(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, pc.Layer)
		}
		return response.ErrInternalServer(c, err, likes, pc.Layer)
	}

	return response.Standard(c, "OK", fiber.Map{
		"postsIDsLikes": likes,
	})
}
