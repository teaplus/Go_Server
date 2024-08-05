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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with a username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Registration Data"
// @Success 200 {object} models.UserRegistrationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/register [post]
func Register(c *gin.Context) {
	var newUser models.User
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

// Login godoc
// @Summary Log in a user
// @Description Log in a user using username or email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginData body models.LoginRequest true "Login Data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/login [post]
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

	keyInsert, err := config.KeyCollection.InsertOne(context.TODO(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	// fmt.Println("key", keyInsert)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "client_id": keyInsert.InsertedID, "user_id": key.User, "token": token, "Public": publicKey, "private": privateKey})
}

// Logout godoc
// @Summary Log out a user
// @Description Log out a user by invalidating their session
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Client-ID header string true "Client ID"
// @Success 200 {object} map[string]string{"message": "Logout Success!"}
// @Failure 500 {object} map[string]string{"error": "Logout Error"}
// @Router /api/v1/logout [post]

func Logout(c *gin.Context) {
	clientId := c.GetHeader("X-Client-ID")

	objectId, err := primitive.ObjectIDFromHex(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}
	var key models.Key
	err = config.KeyCollection.FindOneAndDelete(context.Background(), bson.M{"_id": objectId}).Decode(&key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout Success!"})

}

// GetUser godoc
// @Summary Get user details
// @Description Retrieve details of a user based on user ID
// @Tags user
// @Accept json
// @Produce json
// @Param user header string true "User ID"
// @Success 200 {object} gin.H{"username": "username", "email": "email", "phoneNumber": "phoneNumber", "address": "address"}
// @Failure 400 {object} models.ErrorResponse{"error": "Invalid input"}
// @Failure 500 {object} models.ErrorResponse{"error": "Internal server error"}
// @Router /api/v1/user [get]

func GetUser(c *gin.Context) {
	userId := c.GetHeader("user")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}
	fmt.Println("object", objectId)
	var user models.User
	err = config.UserCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": user.Username, "email": user.Email, "phoneNumber": user.PhoneNumber, "address": user.Address})

}


// ChangePassword godoc
// @Summary Change user password
// @Description Change the password for the logged-in user
// @Tags user
// @Accept json
// @Produce json
// @Param current_password body string true "Current password"
// @Param new_password body string true "New password"
// @Success 200 {object} gin.H{"message": "Password has changed successfully"}
// @Failure 400 {object} models.ErrorResponse{"error": "Invalid input"}
// @Failure 500 {object} models.ErrorResponse{"error": "Internal server error"}
// @Router /api/v1/changepassword [post]
func ChangePassword(c *gin.Context) {
	var passworData struct {
		OldPassword string `json:"current_password"`
		NewPassword string `json:"new_password"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&passworData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetHeader("user")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}
	fmt.Println("object", objectId)
	var user models.User
	err = config.UserCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		return
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding hashed password"})
		return
	}
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding hashed salt"})
		return
	}
	fmt.Println("oldpass", passworData)
	if !VerifyPassword(passworData.OldPassword, []byte(hashedPassword), []byte(salt)) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password has wrong"})
		return
	}

	NewPassword := HashPassword(passworData.NewPassword, salt)
	base64HashedPassword := base64.StdEncoding.EncodeToString(NewPassword)
	user.Password = base64HashedPassword
	update := bson.M{
		"$set": bson.M{
			"password": base64HashedPassword,
		},
	}
	result, err := config.UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found with the provided ID"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has change successfully"})

}
