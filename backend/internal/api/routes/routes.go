package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/job-assistance/internal/api/controller"
)

func JobRoutes(router *gin.RouterGroup, jobController *controller.JobController) {
	jobRouter := router.Group("/job")
	{
		jobRouter.POST("/", jobController.FetchJobs)
	}
}
