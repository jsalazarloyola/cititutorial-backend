package main

import (
	"go-template/middleware"
	"go-template/routes"
	"go-template/services"

	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
    time.Local = time.UTC

    log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    log.Println("Start test backend...")
    log.Printf("serverUp, %s", os.Getenv("ADDR"))
    
    // Inicializar conexión Mongo
    client := services.InitMongo()

    service := services.NewTasksService(client, "taskdb", "tasks")

    r := gin.Default()
    r.Use(middleware.CorsMiddleware())

    // Registrar rutas
    routes.RegisterRoutes(r, service)

    r.Run(":8080")
}
