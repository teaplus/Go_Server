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
		clientId := c.GetHeader("X-Client-ID")
		fmt.Println("XXXXXXXXXX", clientId)
		fmt.Println("YYYY", accessToken)

		objectId, err := primitive.ObjectIDFromHex(clientId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
			return
		}
		var key models.Key
		err = config.KeyCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&key)
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
			fmt.Println("claim")
			c.Next()
		}
	}
}
