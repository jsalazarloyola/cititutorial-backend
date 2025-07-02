package main

import (
	"go-template/middleware"
	"go-template/routes"
	"go-template/services"
	"go-template/utils"

	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
    utils.LoadEnv()
    time.Local = time.UTC

    log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    log.Println("Start test backend...")
    log.Printf("serverUp, %s", os.Getenv("ADDR"))
    
    // Inicializar conexión Mongo
    client := services.InitMongo()
    taskService := services.NewTasksService(client, "taskdb", "tasks")

    // Servicio de login
    loginService := services.NewLoginService()

    // modo de gin
    gin.SetMode(os.Getenv("GIN_MODE"))

    r := gin.Default()
    r.Use(middleware.CorsMiddleware())

    // Registrar rutas
    routes.RegisterRoutes(r, taskService, loginService)

    r.Run(":8000")
}
