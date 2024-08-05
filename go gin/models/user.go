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

// type User struct {
// 	ID          *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	Username    string              `json:"username" bson:"username"`
// 	Email       string              `json:"email" bson:"email"`
// 	PhoneNumber string `json:"phone_number" bson:"phone_number"`
// 	Address     string `json:"address" bson:"address"`
// 	 // Exclude from JSON serialization
// }

// type Account struct {
// 	ID   *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	User primitive.ObjectID  `json:"user_id,omitempty" bson:"user_id,omitempty"`
// 	Username    string              `json:"username" bson:"username"`
// 	Email       string              `json:"email" bson:"email"`
// 	Password    string
// 	Salt        string `json:"-" bson:"salt"`
// }

type UserRegistrationResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	Message     string `json:"message"`
	UserID      string `json:"user"`
	AccessToken string `json:"token"`
	PublicKey   string `json:"public"`
	PrivateKey  string `json:"private"`
}
