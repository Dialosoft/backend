package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ForumRouter struct {
	ForumController *controller.ForumController
}

func NewForumRouter(forumController *controller.ForumController) *ForumRouter {
	return &ForumRouter{ForumController: forumController}
}

func (r *ForumRouter) SetupForumRoutes(api fiber.Router, middlewares *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {
	forumGroup := api.Group("forums")

	{
		forumGroup.Get("/get-all-forums", r.ForumController.GetAllForums)
		forumGroup.Get("/get-forum-by-id/:id", r.ForumController.GetForumByID)
		forumGroup.Get("/get-forum-by-id/:id", r.ForumController.GetForumByName)
		forumGroup.Post("/create-new-forum", r.ForumController.CreateForum,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		forumGroup.Put("/update-forum/:id", r.ForumController.UpdateForum,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		forumGroup.Delete("/delete-forum/:id", r.ForumController.DeleteForum,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		forumGroup.Put("/restore-forum/:id", r.ForumController.RestoreForum,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
	}
}
