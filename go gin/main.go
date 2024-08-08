package main

import (
	config "myproject/dbs"
	"myproject/handlers"

	docs "myproject/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Kết nối với MongoDB
	config.ConnectMongoDB()

	router := gin.Default()
	// Swagger API docs
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization, X-Client-ID, user"}

	router.Use(cors.New(config))
	router.Use()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", handlers.Register)
		v1.POST("/login", handlers.Login)
		v1.POST("/logout", handlers.MiddlewareAuthentication(), handlers.Logout)
		v1.POST("/changepassword", handlers.MiddlewareAuthentication(), handlers.ChangePassword)
		v1.GET("/user", handlers.MiddlewareAuthentication(), handlers.GetUser)
	}
	// Định nghĩa các route cho đăng ký và đăng nhập
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(":8081") // Chạy server trên cổng 8080
}
