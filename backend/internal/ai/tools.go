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
					Name: "search_structured_jobs",
					Description: `Search for jobs using structured LinkUp API with detailed job information extraction.
					This provides structured job data including experience level, required skills, salary, and remote work options.
					Use this when you need detailed job analysis and matching capabilities.
					Best for: Detailed job matching, skill analysis, and comprehensive job evaluation.`,
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {
								Type: genai.TypeString,
								Description: `Detailed job search query for structured extraction. Should include:
								- Specific role titles and responsibilities
								- Required technical skills and experience level
								- Industry context and domain expertise
								- Location preferences and work arrangements
								Example: "Senior Software Engineer with Python Django experience for fintech startup"`,
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
