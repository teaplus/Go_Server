package main

import (
	"myproject/dbs"
	"myproject/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Kết nối với MongoDB
	config.ConnectMongoDB()

	router := gin.Default()

	// Định nghĩa các route cho đăng ký và đăng nhập
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.Run(":8080") // Chạy server trên cổng 8080
}
