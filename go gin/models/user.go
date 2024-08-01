package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username    string              `json:"username" bson:"username"`
	Email       string              `json:"email" bson:"email"`
	Password    string
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Address     string `json:"address" bson:"address"`
	Salt        string `json:"-" bson:"salt"` // Exclude from JSON serialization
}
