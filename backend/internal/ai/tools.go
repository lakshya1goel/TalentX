package ai

import "google.golang.org/genai"

func (a *AIClient) Tools() []*genai.Tool {
	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name: "search_jsearch_jobs",
					Description: `Search for jobs using JSearch API (RapidAPI). This API provides comprehensive job listings from major job boards.
					Use this when you need to find jobs based on specific skills, roles, or keywords extracted from a resume.
					Best for: Software engineering, data science, product management, and tech roles.`,
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {
								Type: genai.TypeString,
								Description: `Job search query based on resume analysis. Should include:
								- Primary skills (e.g., "Python", "React", "Data Science")
								- Experience level (e.g., "Senior", "Junior", "Mid-level")
								- Role type (e.g., "Software Engineer", "Data Analyst")
								Example: "Senior Python Developer" or "Junior Data Scientist"`,
							},
						},
						Required: []string{"query"},
					},
				},
				{
					Name: "search_linkup_jobs",
					Description: `Search for jobs using LinkUp API. This provides high-quality job data directly from company career pages.
					Use this for comprehensive job market coverage and to find jobs that might not be on traditional job boards.
					Best for: Enterprise companies, specific company searches, and comprehensive market analysis.`,
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {
								Type: genai.TypeString,
								Description: `Job search query optimized for LinkUp API. Should focus on:
								- Core technical skills and technologies
								- Industry-specific terms
								- Job titles and roles
								Example: "Machine Learning Engineer" or "Full Stack Developer React Node.js"`,
							},
						},
						Required: []string{"query"},
					},
				},
			},
		},
	}
	return tools
}