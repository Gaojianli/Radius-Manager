package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"

	"github.com/Gaojianli/raduis_mgnt/controllers"
	"github.com/Gaojianli/raduis_mgnt/middleware"
	"github.com/Gaojianli/raduis_mgnt/static"
)

func SetupRoutes(h *server.Hertz) {
	h.Use(cors.Default())

	userController := &controllers.UserController{}
	radiusController := &controllers.RadiusController{}

	api := h.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/login", middleware.JWTMiddleware.LoginHandler)
				auth.POST("/refresh", middleware.JWTMiddleware.RefreshHandler)
			}

			user := v1.Group("/user")
			user.Use(middleware.JWTMiddleware.MiddlewareFunc())
			{
				user.GET("/profile", userController.GetProfile)
				user.PUT("/change-password", userController.ChangePassword)
				user.GET("/stats", userController.GetStats)
			}

			admin := v1.Group("/admin")
			admin.Use(middleware.JWTMiddleware.MiddlewareFunc(), middleware.RequireAdmin())
			{
				admin.GET("/users", userController.GetUsers)
				admin.POST("/users", userController.AdminCreateUser)
				admin.PUT("/users/:id/password", userController.AdminChangePassword)
				admin.PUT("/users/:id/ban", userController.AdminToggleBanUser)
				admin.DELETE("/users/:id", userController.AdminDeleteUser)
				admin.GET("/auth-logs", userController.GetAuthLogs)
				admin.GET("/stats", userController.GetAdminStats)
			}

			radius := v1.Group("/radius")
			{
				radius.POST("/auth", radiusController.Authenticate)
				radius.POST("/authorize", radiusController.Authorize)
				radius.POST("/accounting", radiusController.Accounting)
			}
		}
	}

	// 设置静态文件服务 (放在所有API路由之后)
	if err := static.SetupStaticRoutes(h); err != nil {
		panic("Failed to setup static routes: " + err.Error())
	}
}
