package handlers

import (
	"context"
	"fmt"
	"net/http"

	config "myproject/dbs"
	"myproject/models"
	"myproject/ultils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MiddlewareAuthentication() gin.HandlerFunc {

	fmt.Printf("Middleware")
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		userId := c.GetHeader("X-Client-ID")

		objectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
			return
		}
		var key models.Key
		err = config.KeyCollection.FindOne(context.Background(), bson.M{"user_id": objectId}).Decode(&key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
			return
		}

		claims, err := ultils.ValidateToken(accessToken, key.PublicKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid", "key": key, "access": accessToken, "claims": claims})
			return
		}
		if claims != nil {
			c.Next()
		}
	}
}
