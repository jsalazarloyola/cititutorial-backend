package routes

import (
	"go-template/controllers"
	"go-template/services"

	"github.com/gin-gonic/gin"
)

// Configura las rutas de los requerimientos D:
func RegisterRoutes(router *gin.Engine, service *services.TasksService) {
	router.GET("/ping", controllers.PingController)

	router.GET("/task", controllers.GetTasks(service))

	router.POST("/task", controllers.CreateTask(service))

	router.PUT("/task/:id", controllers.EditTask(service))
}
