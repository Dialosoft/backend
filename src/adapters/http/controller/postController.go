package controller

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostController struct {
	PostService services.PostService
}

func NewPostController(postService services.PostService) *PostController {
	return &PostController{PostService: postService}
}

func (pc *PostController) GetAllPosts(c fiber.Ctx) error {
	posts, err := pc.PostService.GetAllPosts()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
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
		return response.ErrUUIDParse(c)
	}

	post, err := pc.PostService.GetPostByID(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
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
		return response.ErrUUIDParse(c)
	}

	posts, err := pc.PostService.GetPostsByUserID(userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", posts)
}

func (pc *PostController) CreateNewPost(c fiber.Ctx) error {
	var req request.NewPost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	post, err := pc.PostService.CreateNewPost(uuid.MustParse(req.UserID), req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.StandardCreated(c, "CREATED", post)
}

func (pc *PostController) UpdatePostTitle(c fiber.Ctx) error {
	var req request.UpdatePostTitle

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	postUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.UpdatePostTitle(postUUID, req.Title)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (pc *PostController) UpdatePostContent(c fiber.Ctx) error {
	var req request.UpdatePostContent

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	postUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.UpdatePostContent(postUUID, req.Content)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
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
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.DeletePost(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
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
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.RestorePost(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "RESTORED", nil)
}

func (pc *PostController) LikePost(c fiber.Ctx) error {
	var req request.LikeOrUnlikePost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	userUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}
	postUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.LikePost(postUUID, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "LIKED", nil)
}

func (pc *PostController) UnlikePost(c fiber.Ctx) error {
	var req request.LikeOrUnlikePost
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	userUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}
	postUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	err = pc.PostService.UnlikePost(postUUID, userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "UNLIKED", nil)
}

func (pc *PostController) GetPostLikesByUserID(c fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	likes, err := pc.PostService.GetPostLikesByUserID(postUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", likes)
}
