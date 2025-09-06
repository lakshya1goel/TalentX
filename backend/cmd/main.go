package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/job-assistance/config"
	"github.com/lakshya1goel/job-assistance/internal/api/controller"
	"github.com/lakshya1goel/job-assistance/internal/api/routes"
)

func main() {
	config.LoadEnv()

	router := gin.Default()
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string

	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	}
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}

	corsConfig.AllowOrigins = allowedOrigins
	router.Use(cors.New(corsConfig))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Health Check!"})
	})

	apiRouter := router.Group("/api")
	{
		routes.JobRoutes(apiRouter, controller.NewJobController())
	}

	router.Run(":8084")
}
