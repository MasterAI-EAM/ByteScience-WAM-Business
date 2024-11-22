package main

import (
	"ByteScience-WAM-Business/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	engin := gin.Default()
	defer internal.ServerExit(engin)

	// 1. 尝试加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// 2. 获取环境变量，优先从 .env 文件中获取
	ginMode := os.Getenv("GIN_MODE_ADMIN")

	internal.ServerStart(engin, ginMode)
}
