package dtos

import "time"

type LocationPreference struct {
	Types     []string `json:"types" binding:"required"`
	Locations []string `json:"locations,omitempty"`
}

type JSearchJob struct {
	Title       string `json:"job_title"`
	Company     string `json:"employer_name"`
	IsRemote    bool   `json:"job_is_remote"`
	Location    string `json:"job_city"`
	Description string `json:"job_description"`
	URL         string `json:"job_apply_link"`
}

type LinkupJob struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	URL     string `json:"url"`
}

type Job struct {
	Title       string `json:"title"`
	Company     string `json:"company,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	Source      string `json:"source"`
}

type ErrorResponse struct {
	Error     string    `json:"error"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

type JobSearchResult struct {
	Jobs   []Job
	Error  error
	Source string
}

type RankedJob struct {
	Job             Job      `json:"job"`
	RelevanceScore  float64  `json:"relevance_score"`
	MatchReason     string   `json:"match_reason"`
	SkillsMatched   []string `json:"skills_matched"`
	ExperienceMatch string   `json:"experience_match"`
}

type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

type PaginatedJobResponse struct {
	Jobs       []RankedJob `json:"jobs"`
	TotalJobs  int         `json:"total_jobs"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	Success    bool        `json:"success"`
}
