package routes

import (
	"go-template/controllers"
	"go-template/middleware"
	"go-template/services"

	"github.com/gin-gonic/gin"
)

// Configura las rutas de los requerimientos D:
func RegisterRoutes(router *gin.Engine, taskService *services.TasksService, loginService *services.LoginService) {
	router.GET("/ping", controllers.PingController)

	authMiddleware := middleware.LoadJWTAuth(loginService)

	// Login no está protegido :'v
	router.POST("/api/login", authMiddleware.LoginHandler)

	protected := router.Group("/api")

	protected.Use(authMiddleware.MiddlewareFunc())
	{
		protected.GET("/task", controllers.GetTasks(taskService))
		protected.GET("/task/:page/:limit", controllers.GetPageTask(taskService))

		protected.POST("/task", controllers.CreateTask(taskService))

		protected.PUT("/task/:id", controllers.EditTask(taskService))
	}

}
