package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task model represents a task in MongoDB
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"` // e.g., "pending", "completed"
}
