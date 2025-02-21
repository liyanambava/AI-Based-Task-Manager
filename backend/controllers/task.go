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

	// Get task ID from URL and force ObjectID conversion
	taskID := c.Param("id")
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

	// Update task in MongoDB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": updatedTask.Status}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Fetch updated task
	var newTask models.Task
	err = collection.FindOne(context.TODO(), filter).Decode(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	// WebSocket Broadcast
	go websocket.BroadcastMessage("Task Updated: " + newTask.Title)

	c.JSON(http.StatusOK, newTask)
}

// Delete a task
func DeleteTask(c *gin.Context) {
	collection := GetTaskCollection()
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	// Get task ID from URL and convert it to ObjectID
	taskID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// Delete the task from MongoDB
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// WebSocket Broadcast
	go websocket.BroadcastMessage("Task Deleted: " + taskID)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task_id": taskID})
}
