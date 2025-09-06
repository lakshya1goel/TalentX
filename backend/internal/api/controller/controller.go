package controller

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/job-assistance/internal/api/service"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobController struct {
	service service.JobService
}

func NewJobController() *JobController {
	return &JobController{
		service: service.NewJobService(),
	}
}

func (c *JobController) FetchJobs(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("resume")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error:     "Resume PDF file is required",
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error:     "Only PDF files are allowed",
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}

	if header.Size > 10*1024*1024 {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error:     "File size must be less than 10MB",
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}

	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Error:     "Failed to read PDF file",
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}

	jobs, err := c.service.FetchJobs(ctx, pdfBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Error:     err.Error(),
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}
