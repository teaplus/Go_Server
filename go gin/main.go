package main

import (
	config "myproject/dbs"
	"myproject/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func protectedHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Hello, authenticated user!"})
}

func main() {
	// Kết nối với MongoDB
	config.ConnectMongoDB()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config))
	router.Use()

	// Định nghĩa các route cho đăng ký và đăng nhập
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.GET("/",handlers.MiddlewareAuthentication(), protectedHandler )
	router.Run(":8080") // Chạy server trên cổng 8080
}
