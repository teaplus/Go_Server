package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	config "myproject/dbs"
	"myproject/models"
	"myproject/ultils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

type Claims struct {
	User string
	jwt.StandardClaims
}

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

const (
	timeCost   = 1
	memoryCost = 64 * 1024
	threads    = 4
	keyLength  = 32
)

func HashPassword(password string, salt []byte) []byte {
	hash := argon2.IDKey([]byte(password), []byte(salt), timeCost, memoryCost, threads, keyLength)
	return hash
}

func VerifyPassword(password string, hash, salt []byte) bool {
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
	salt, err := generateSalt(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating salt"})
		return
	}
	hashedPassword := HashPassword(newUser.Password, salt)
	base64HashedPassword := base64.StdEncoding.EncodeToString(hashedPassword)
	newUser.Salt = base64.StdEncoding.EncodeToString(salt)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot has password"})
	// 	return
	// }
	newUser.Password = string(base64HashedPassword)

	// Check exist
	var existingUser models.User
	err = config.UserCollection.FindOne(context.TODO(), bson.M{"username": string(newUser.Username)}).Decode(&existingUser)
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

func generateKey(size int) (string, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
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

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding hashed salt"})
		return
	}
	// fmt.Print("salt", salt)

	hashedPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding hashed password"})
		return
	}
	fmt.Print("salt", hashedPassword)

	// Check password
	if !VerifyPassword(loginData.Password, []byte(hashedPassword), []byte(salt)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	privateKey, err := generateKey(64)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}

	publicKey, err := generateKey(64)
	if err != nil {
		fmt.Println("Error generating public key:", err)
		return
	}

	fmt.Println("key", user)
	Claims := ultils.Claims{
		User: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * 24 * time.Hour).Unix(),
		},
	}
	token, err := ultils.CreateTokenpair(Claims, publicKey, privateKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error createTokenPair"})
		return
	}

	key := models.Key{
		User:         *user.ID,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		RefreshToken: token["refreshToken"],
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = config.KeyCollection.InsertOne(context.TODO(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user.ID, "token": token["accessToken"], "Public": publicKey, "private": privateKey})
}
