package main

import (
	config "myproject/dbs"
	"myproject/handlers"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/files"
	docs "myproject/docs"
)

func protectedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, authenticated user!"})
}

func main() {
	// Kết nối với MongoDB
	config.ConnectMongoDB()

	router := gin.Default()
	// Swagger API docs
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config))
	router.Use()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", handlers.Register)
		v1.POST("/login", handlers.Login)
		v1.POST("/logout", handlers.MiddlewareAuthentication(), handlers.Logout)
	}
	// Định nghĩa các route cho đăng ký và đăng nhập
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	

	router.Run(":8080") // Chạy server trên cổng 8080
}
