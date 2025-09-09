package controller

import (
	"io"
	"net/http"
	"strconv"
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

	locationTypes := ctx.PostFormArray("location_types")
	locations := ctx.PostFormArray("locations")

	if len(locationTypes) == 0 {
		if singleType := ctx.PostForm("location_type"); singleType != "" {
			locationTypes = []string{singleType}
		}
	}
	if len(locations) == 0 {
		if singleLocation := ctx.PostForm("location"); singleLocation != "" {
			locations = []string{singleLocation}
		}
	}

	locationPreference := dtos.LocationPreference{
		Types:     locationTypes,
		Locations: locations,
	}

	if len(locationPreference.Types) == 0 {
		locationPreference.Types = []string{"remote"}
	}

	validTypes := map[string]bool{"remote": true, "onsite": true, "hybrid": true}
	for _, locType := range locationPreference.Types {
		if !validTypes[locType] {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
				Error:     "Invalid location type. Must be 'remote', 'onsite', or 'hybrid'",
				Success:   false,
				Timestamp: time.Now(),
			})
			return
		}
	}

	needsLocation := false
	for _, locType := range locationPreference.Types {
		if locType == "onsite" || locType == "hybrid" {
			needsLocation = true
			break
		}
	}

	if needsLocation && len(locationPreference.Locations) == 0 {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error:     "At least one location is required for onsite and hybrid positions",
			Success:   false,
			Timestamp: time.Now(),
		})
		return
	}

	pagination := dtos.PaginationRequest{
		Page:     1,
		PageSize: 20,
	}

	if pageStr := ctx.PostForm("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := ctx.PostForm("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}

	jobs, err := c.service.FetchJobs(ctx, pdfBytes, locationPreference, pagination)
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
