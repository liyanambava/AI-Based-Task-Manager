package controllers

import (
	"context"
	"net/http"

	"task-manager/config"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get the tasks collection
func GetTaskCollection() *mongo.Collection {
	if config.DB == nil {
		return nil
	}
	return config.DB.Collection("tasks")
}

// Create a task
func CreateTask(c *gin.Context) {
	collection := GetTaskCollection()
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task_id": result.InsertedID})
}

// GetTasks retrieves all tasks from the database
func GetTasks(c *gin.Context) {
	collection := GetTaskCollection()
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	// Fetch all tasks
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
