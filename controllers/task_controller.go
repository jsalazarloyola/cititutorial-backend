package controllers

import (
	"context"
	"go-template/models"
	"go-template/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	ErrorMessage string            `json:"error_message"`
	Task    map[string]interface{} `json:"task"`
}

type TaskListResponse struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	ErrorMessage string              `json:"error_message"`
	Tasks   []models.Task `json:"task_list"`
}

// Retorna un handler, que se configura en routes.go
func CreateTask(service *services.TasksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// timeout de 10 segundos
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var task models.Task
		defer cancel()

		// valida el body del requerimiento
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest,
				TaskResponse{
					Status:  http.StatusBadRequest,
					Message: "error",
					ErrorMessage: err.Error(),
					// Data:    map[string]interface{}{"data": },
				},
			)
			return
		}

		newTask := models.Task{
			Id: primitive.NewObjectID(),  // ID automática
			Title:       task.Title,
			Description: task.Description,
			Completed:   false,
			CreatedAt:   time.Now().Format("2006-01-02"),
		}

		// Inserta D:
		result, err := service.Collection.InsertOne(ctx, newTask)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				TaskResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					ErrorMessage: err.Error(),
					// Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}
		//https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-gin-gonic-version-269m
		// Reporte de estado
		c.JSON(http.StatusCreated,
			TaskResponse{
				Status:  http.StatusCreated,
				Message: "success",
				Task:    map[string]interface{}{"task": result},
			},
		)
	}
}

// Obtiene una lista con las tareas ("lista" en el sentido de respuesta en
// formato JSON donde uno de los campos es la respuesta)
func GetTasks(service *services.TasksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var tasks []models.Task
		// Empuja al stack de funciones a ejecutar tras esta función
		defer cancel()

		// Recupera las tareas como mapa de BSON
		results, err := service.Collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				TaskResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					ErrorMessage: err.Error(),
					// Task:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleTask models.Task
			if err = results.Decode(&singleTask); err != nil {
				c.JSON(http.StatusInternalServerError,
					TaskResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						ErrorMessage: err.Error(),
						// Task:    map[string]interface{}{"data": err.Error()},
					},
				)
			}

			tasks = append(tasks, singleTask)
		}
		// La respuesta final con los datos
		c.JSON(http.StatusOK,
			TaskListResponse{
				Status:  http.StatusOK,
				Message: "success",
				Tasks:    tasks,
			},
		)
	}
}

// Actualiza la tarea que indica el requerimiento
func EditTask(service *services.TasksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var task models.Task

		// Se debe transformar el ID, un string hexadecimal, a ObjectID
		// para poder buscarlo
		taskId, _ := primitive.ObjectIDFromHex(c.Param("id"))
		log.Print(taskId)

		if err := c.BindJSON(&task); err != nil {
			c.JSON(
				http.StatusBadRequest,
				TaskResponse{
					Status: http.StatusBadRequest,
					Message: "error",
					Task:    map[string]interface{}{"data": err.Error()},
				},
			)
		}
		log.Print(task)

		update := bson.M{
			"title":       task.Title,
			"description": task.Description,
			"completed":   task.Completed,
		}
		result, err := service.Collection.UpdateOne(ctx, bson.M{"_id": taskId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				TaskResponse{
					Status: http.StatusInternalServerError,
					Message: "error",
					Task:    map[string]interface{}{"data": err.Error()},
				},
			)
		}

		var updatedTask models.Task
		if result.MatchedCount == 1 {
			err := service.Collection.FindOne(ctx, bson.M{"_id": taskId}).Decode(&updatedTask)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					TaskResponse{
						Status: http.StatusInternalServerError,
						Message: "error",
						Task:    map[string]interface{}{"data": err.Error()},
					},
				)
				return
			}
		}
		c.JSON(
			http.StatusOK,
			TaskResponse{
				Status: http.StatusOK,
				Message: "succes",
				Task: map[string]interface{}{"data": updatedTask},
			},
		)
	}
}
