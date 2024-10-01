package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type PostRouter struct {
	PostController *controller.PostController
}

func NewPostRouter(postController *controller.PostController) *PostRouter {
	return &PostRouter{PostController: postController}
}

func (r *PostRouter) SetupPostRoutes(api fiber.Router, middlewares *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {
	postGroup := api.Group("/posts") // middlewares.GetAndVerifyAccessToken(),
	// middlewares.VerifyRefreshToken(),
	postProtected := postGroup.Group("/protected", middlewares.GetAndVerifyAccessToken(), middlewares.VerifyRefreshToken())

	{
		postGroup.Get("/get-all-posts-by-forum/:forumID", r.PostController.GetAllPostsByForum)
		// postGroup.Get("/get-all-posts", r.PostController.GetAllPosts)
		// postGroup.Get("/get-all-posts-simple", r.PostController.GetAllPostsAndReturnSimpleResponse)
		// postGroup.Get("/get-post-by-id/:id", r.PostController.GetPostByID)
		// postGroup.Get("/get-posts-by-user-id/:userID", r.PostController.GetPostsByUserID)
		// postGroup.Get("/get-like-count/:id", r.PostController.GetPostNumberOfLikes)
		// postGroup.Get("/get-post-likes-by-user-id/:userID", r.PostController.GetPostLikesByUserID)
		postProtected.Post("/create-new-post", r.PostController.CreateNewPost)
		// postGroup.Put("/update-post-title/:id", r.PostController.UpdatePostTitle)
		// postGroup.Put("/update-post-content/:id", r.PostController.UpdatePostContent)
		// postGroup.Delete("/delete-post/:id", r.PostController.DeletePost)
		// postGroup.Put("/restore-post/:id", r.PostController.RestorePost)
		// postGroup.Put("/like-post", r.PostController.LikePost)
		// postGroup.Put("/unlike-post", r.PostController.UnlikePost)
	}
}
