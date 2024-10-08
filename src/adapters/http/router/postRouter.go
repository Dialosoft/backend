package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type PostRouter struct {
	PostController *controller.PostController
}

func NewPostRouter(postController *controller.PostController) *PostRouter {
	return &PostRouter{PostController: postController}
}

func (r *PostRouter) SetupPostRoutes(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {

	// public group
	postGroup := api.Group("/posts")

	// protected group
	postProtected := postGroup.Group("/protected",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManagePosts() /* permission middleware for Posts */)

	// authenticated group
	authenticatedGroup := postGroup.Group("/authenticated",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken())

	{
		// public
		postGroup.Get("/get-all-posts-by-forum/:forumID", r.PostController.GetAllPostsByForum)
		postGroup.Get("/get-all-posts-simple", r.PostController.GetAllPostsAndReturnSimpleResponse)
		postGroup.Get("/get-posts-by-user-id/:userID", r.PostController.GetPostsByUserID)
	}

	{
		// protected routes by authenticated users
		authenticatedGroup.Post("/create-new-post", r.PostController.CreateNewPost)
		authenticatedGroup.Put("/like-post", r.PostController.LikePost)
		authenticatedGroup.Put("/unlike-post", r.PostController.UnlikePost)
	}

	{

		// protected routes by authenticated users and with permission (manage posts permission)
		postProtected.Delete("/delete-post/:id", r.PostController.DeletePost)
		postProtected.Get("/get-all-posts", r.PostController.GetAllPosts)
		postProtected.Get("/get-post-by-id/:id", r.PostController.GetPostByID)
		postProtected.Put("/update-post-title/:id", r.PostController.UpdatePostTitle)
		postProtected.Put("/update-post-content/:id", r.PostController.UpdatePostContent)
		postProtected.Put("/restore-post/:id", r.PostController.RestorePost)
		// postGroup.Get("/get-like-count/:id", r.PostController.GetPostNumberOfLikes)
		// postGroup.Get("/get-post-likes-by-user-id/:userID", r.PostController.GetPostLikesByUserID)
	}
}
