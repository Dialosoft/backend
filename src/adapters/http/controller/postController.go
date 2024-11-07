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
