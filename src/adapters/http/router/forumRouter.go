package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type ForumRouter struct {
	ForumController *controller.ForumController
}

func NewForumRouter(forumController *controller.ForumController) *ForumRouter {
	return &ForumRouter{ForumController: forumController}
}

func (r *ForumRouter) SetupForumRoutes(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {
	forumGroup := api.Group("/forums")
	forumProtected := forumGroup.Group("/protected",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManageForums() /* permission middleware for forums */)

	{
		// public

		forumGroup.Get("/get-forums-by-category-id/:categoryID", r.ForumController.GetForumsByCategoryIDAndAllowed)
	}

	{
		// protected routes by authenticated users and with permission

		// forumGroup.Get("/get-all-forums", r.ForumController.GetAllForums)
		// forumGroup.Get("/get-forum-by-id/:id", r.ForumController.GetForumByID)
		// forumGroup.Get("/get-forum-by-id/:id", r.ForumController.GetForumByName)

		forumProtected.Post("/create-new-forum", r.ForumController.CreateForum)
		forumProtected.Put("/update-forum/:id", r.ForumController.UpdateForum)
		forumProtected.Delete("/delete-forum/:id", r.ForumController.DeleteForum)
		forumProtected.Put("/restore-forum/:id", r.ForumController.RestoreForum)
	}
}
