package models

type JSearchJob struct {
	Title       string `json:"job_title"`
	Company     string `json:"employer_name"`
	Location    string `json:"job_city"`
	Description string `json:"job_description"`
	URL         string `json:"job_apply_link"`
}

type LinkupJob struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
	Source  string `json:"source"`
}

type Job struct {
	Title       string `json:"title"`
	Company     string `json:"company,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	Source      string `json:"source"`
}
