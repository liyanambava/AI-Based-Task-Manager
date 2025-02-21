package controllers

import (
	"context"
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/websocket" // Import the websocket package

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Insert task into MongoDB
	result, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Broadcast task creation to all WebSocket clients
	websocket.BroadcastTaskCreated("New task created: " + task.Title)

	// Return response with task ID
	c.JSON(http.StatusCreated, gin.H{"task_id": result.InsertedID})
}

// GetTasks retrieves all tasks from the database
func GetTasks(c *gin.Context) {
	collection := GetTaskCollection()
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	// Fetch all tasks from MongoDB
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []bson.M
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks) // ✅ Always returns `_id`

}

func UpdateTaskStatus(c *gin.Context) {
	collection := GetTaskCollection()
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	// Get task ID from URL
	taskID := c.Param("id")

	// Convert taskID to ObjectID (allow both string and ObjectID formats)
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// Parse request body
	var updatedTask struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// ✅ Check if the task exists before updating
	var existingTask models.Task
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// ✅ Update task status in MongoDB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": updatedTask.Status}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// ✅ Fetch updated task
	err = collection.FindOne(context.TODO(), filter).Decode(&existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	// ✅ WebSocket Broadcast
	go websocket.BroadcastMessage("Task Updated: " + existingTask.Title)

	// ✅ Return updated task
	c.JSON(http.StatusOK, existingTask)
}
