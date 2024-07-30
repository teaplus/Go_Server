package handlers

import (
	"bytes"
	"context"
	config "myproject/dbs"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
	"fmt"
)

const (
	timeCost   = 1
	memoryCost = 64 * 1024
	threads    = 4
	keyLength  = 32
	salt       = "3Le7Tow2"
)

func HashPassword(password string) []byte {
	hash := argon2.IDKey([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLength)
	return hash
}

func VerifyPassword(password string, hash []byte) bool {
	expectedHash := argon2.IDKey([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLength)
	return bytes.Equal(hash, expectedHash)
}

// func HashPassword(password string) []byte {
// 	hash := argon2.Key([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLength)
// 	return hash
// }

// func VerifyPassword(password string, hash []byte) bool {
// 	expectedHash := argon2.Key([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLength)
// 	return bytes.Equal(hash, expectedHash)
// }

// Register
func Register(c *gin.Context) {
	var newUser models.User
	fmt.Println("user", newUser)
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// hash pass
	hashedPassword := HashPassword(newUser.Password)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot has password"})
	// 	return
	// }
	newUser.Password = string(hashedPassword)

	// Check exist
	var existingUser models.User
	err := config.UserCollection.FindOne(context.TODO(), bson.M{"username": newUser.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User has already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking"})
		return
	}

	// Lưu người dùng
	_, err = config.UserCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login
func Login(c *gin.Context) {
	var loginData struct {
		UsernameOrEmail string `json:"username_or_email"`
		Password        string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// search
	var user models.User
	filter := bson.M{
		"$or": []bson.M{
			{"username": loginData.UsernameOrEmail},
			{"email": loginData.UsernameOrEmail},
		},
	}
	err := config.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User"})
		}
		return
	}

	// Check password
	if !VerifyPassword(loginData.Password, []byte(user.Password)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
