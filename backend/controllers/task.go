package controllers

import (
	"context"
	"net/http"
	"time"

	"task-manager/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection *mongo.Collection

func init() {
	taskCollection = config.DB.Collection("tasks")
}

// Create a new task
func CreateTask(c *gin.Context) {
	var task struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = time.Now()

	_, err := taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully"})
}

// Get all tasks
func GetTasks(c *gin.Context) {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	var tasks []bson.M
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
