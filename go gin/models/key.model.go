package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Key struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User             primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PrivateKey       string
	PublicKey        string
	RefreshToken     string
	RefreshTokenUsed []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
