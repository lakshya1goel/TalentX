package dtos

import "time"

type LocationPreference struct {
	Types     []string `json:"types" binding:"required"`
	Locations []string `json:"locations,omitempty"`
}

type JobDescription struct {
	JobTitle        string   `json:"job_title"`
	ExperienceLevel string   `json:"experience_level"`
	RequiredSkills  []string `json:"required_skills"`
	Remote          bool     `json:"remote"`
	Location        *string  `json:"location,omitempty"`
	Salary          *int     `json:"salary,omitempty"`
	JobPostURL      string   `json:"job_post_url"`
	Company         string   `json:"company"`
}

type JobAnnouncements struct {
	Jobs []JobDescription `json:"jobs"`
}

type JobMatchEvaluation struct {
	MatchScore int    `json:"match_score"`
	Reasons    string `json:"reasons"`
}

type ResumeProfile struct {
	PotentialJobTitles []string `json:"potential_job_titles"`
	Seniority          string   `json:"seniority"`
	Skills             []string `json:"skills"`
	BasedIn            *string  `json:"based_in,omitempty"`
	WorkLocation       *string  `json:"work_location,omitempty"`
}

type DetailedJobMatch struct {
	Job        JobDescription `json:"job"`
	MatchScore int            `json:"match_score"`
	Reasons    string         `json:"reasons"`
	URL        string         `json:"url"`
	Company    string         `json:"company"`
}

type JobMatchResults struct {
	Matches map[string]DetailedJobMatch `json:"matches"`
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
	PercentMatch    float64  `json:"percent_match"`
	MatchReason     string   `json:"match_reason"`
	SkillsMatched   []string `json:"skills_matched"`
	ExperienceMatch string   `json:"experience_match"`
}

type JobSearchResponse struct {
	Jobs    []RankedJob `json:"jobs"`
	Total   int         `json:"total"`
	Success bool        `json:"success"`
}

type BatchResult struct {
	Jobs       []RankedJob
	BatchIndex int
}
