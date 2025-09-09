package ai

import (
	"encoding/json"
	"fmt"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func (a *AIClient) Prompt() string {
	prompt := `
You are an expert career counselor and job search specialist. Analyze the provided resume and help find the most relevant job opportunities.

Your task:
1. **Resume Analysis**: Carefully examine the resume to extract:
	- Technical skills and technologies
	- Years of experience and seniority level (Do not consider internship experience from resume as experience it comes under fresher only)
	- Previous job titles and roles
	- Industry experience
	- Education and certifications

2. **Job Search Strategy**: Based on the analysis:
	- Create targeted search queries for different job search platforms
	- Focus on the candidate's strongest skills and experience
	- Consider both current level and growth opportunities
	- Include relevant keywords that employers use
	- If a candidate have only internship experience then fetch some internship opportunities as well both general and skill specific
	- Fetch sone linkedin jobs as well

3. **Tool Usage**: Use the available job search tools effectively:
	- Use JSearch for broad market coverage and popular job boards
	- Use LinkUp for direct company postings and enterprise opportunities
	- Create multiple targeted searches rather than one generic search
	- Vary search terms to capture different types of opportunities

4. **Search Query Guidelines**:
	- Include experience level (Junior, Mid-level, Senior, Lead)
	- Combine role + key technology (e.g., "Senior Python Developer")
	- Use industry-standard terms and job titles
	- Consider related roles the candidate could transition to

Please analyze the resume thoroughly and then search for relevant job opportunities using both available tools.

Create multiple search queries:
- Broad query (e.g., "Software Engineer")  
- Specific to main skills (e.g., "Golang Developer")
- Determine the query that is suitable according to the resume
`
	return prompt
}

func (a *AIClient) PromptWithLocation(locationPreference dtos.LocationPreference) string {
	basePrompt := a.Prompt()

	if len(locationPreference.Types) == 0 {
		return basePrompt
	}

	locationGuidance := "\n\n5. **Location Preferences**: "

	if len(locationPreference.Types) == 1 {
		switch locationPreference.Types[0] {
		case "remote":
			locationGuidance += "The candidate is looking for REMOTE work opportunities only.\n\t- Focus on remote-friendly companies and positions\n\t- Include \"remote\" in search queries when relevant\n\t- Prioritize companies known for remote work culture"
		case "onsite":
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("The candidate is looking for ON-SITE work in %s.\n\t- Focus on companies and positions located in or near these areas\n\t- Include location-specific searches\n\t- Consider commuting distance and local job market", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "The candidate is looking for ON-SITE work opportunities.\n\t- Focus on local companies and positions\n\t- Include location-specific searches when possible"
			}
		case "hybrid":
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("The candidate is open to HYBRID work arrangements in %s.\n\t- Look for both remote and on-site opportunities in these areas\n\t- Focus on companies offering flexible work arrangements\n\t- Include both remote and location-specific searches", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "The candidate is open to HYBRID work arrangements.\n\t- Look for companies offering flexible work arrangements\n\t- Include both remote and location-specific searches"
			}
		}
	} else {
		locationGuidance += "The candidate is open to multiple work arrangements:\n"

		hasRemote := contains(locationPreference.Types, "remote")
		hasOnsite := contains(locationPreference.Types, "onsite")
		hasHybrid := contains(locationPreference.Types, "hybrid")

		if hasRemote {
			locationGuidance += "\t- REMOTE: Include remote-friendly positions and companies\n"
		}
		if hasOnsite {
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("\t- ON-SITE: Include positions in %s\n", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "\t- ON-SITE: Include local on-site positions\n"
			}
		}
		if hasHybrid {
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("\t- HYBRID: Include flexible arrangements in %s\n", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "\t- HYBRID: Include flexible work arrangements\n"
			}
		}

		locationGuidance += "\t- Cast a wide net to capture all preferred work arrangements\n\t- Use varied search terms to find opportunities matching any of these preferences"
	}

	return basePrompt + locationGuidance
}

func (r *RerankingClient) RerankingPrompt(jobs []dtos.Job) string {
	jobsJSON, _ := json.MarshalIndent(jobs, "", "  ")

	prompt := fmt.Sprintf(`
You are an expert career counselor and resume analyzer. Your task is to analyze the provided resume and rank the given job opportunities based on their relevance to the candidate's profile.

**Resume Analysis Instructions:**
1. Extract the candidate's:
	- Technical skills and technologies
	- Years of experience and seniority level
	- Previous job titles and roles
	- Industry experience
	- Education and certifications
	- Career progression and growth trajectory

**Job Ranking Instructions:**
2. For each job, evaluate:
	- **Skill Match**: How well do the required skills align with the candidate's skills?
	- **Experience Level Match**: Is the role appropriate for the candidate's experience level?
	- **Industry Relevance**: Does the role align with the candidate's industry background?
	- **Career Growth Potential**: Does this role offer appropriate career advancement?
	- **Role Transition Feasibility**: How realistic is it for the candidate to transition to this role?

**Scoring Guidelines:**
- Score each job from 0.0 to 10.0 (10.0 being perfect match)
- Consider both current qualifications and potential for growth
- Factor in transferable skills and learning potential
- Be realistic about experience requirements vs. candidate's background

**Output Format:**
Return a JSON array of ranked jobs in descending order of relevance score. Each job should include:
`+"```json"+`
[
	{
		"job_index": 0,
		"relevance_score": 8.5,
		"match_reason": "Strong alignment with candidate's Python and ML experience. Role matches seniority level.",
		"skills_matched": ["Python", "Machine Learning", "TensorFlow"],
		"experience_match": "Perfect match for mid-level candidate",
		"concerns": "May require some additional cloud experience"
	}
]
`+"```"+`

**Jobs to Rank:**
%s

Please analyze the resume thoroughly and provide a detailed ranking of these %d job opportunities.
`, string(jobsJSON), len(jobs))

	return prompt
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func joinLocations(locations []string) string {
	if len(locations) == 0 {
		return ""
	}
	if len(locations) == 1 {
		return locations[0]
	}
	if len(locations) == 2 {
		return locations[0] + " and " + locations[1]
	}

	result := ""
	for i, loc := range locations {
		if i == len(locations)-1 {
			result += "and " + loc
		} else if i > 0 {
			result += ", " + loc
		} else {
			result += loc
		}
	}
	return result
}
