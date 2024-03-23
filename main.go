package main

import (
	"Learn/LittleRedBook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	web.RegisterUserRoutes(server)
	server.Run(":8080")
}
