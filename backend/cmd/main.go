package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/job-assistance/config"
)

func main() {
	config.LoadEnv()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Health Check!"})
	})
	router.Run(":8084")
}
