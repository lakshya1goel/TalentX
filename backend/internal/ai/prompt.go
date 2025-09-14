package ai

import (
	"encoding/json"
	"fmt"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func (a *AIClient) Prompt(locationContext string) string {
	prompt := fmt.Sprintf(`
You are a senior career strategist and resume analyzer with expertise in tech hiring and talent matching. Analyze the provided resume and execute strategic job searches with precise targeting.

**CRITICAL RESUME ANALYSIS FRAMEWORK:**

1. **Experience Level Calculation (VERY IMPORTANT):**
   - **ONLY count full-time professional work experience**
   - **INTERNSHIPS DO NOT COUNT as professional experience - treat as FRESHER**
   - **Current students/recent graduates with only internships = FRESHER category**
   - **Experience Categories:**
     * Fresher/Entry-level: 0-1 years professional experience (includes internship-only candidates)
     * Junior: 1-3 years professional experience
     * Mid-level: 3-6 years professional experience  
     * Senior: 6-10 years professional experience
     * Lead/Principal: 10+ years professional experience

2. **Education and Graduation Year Analysis:**
   - **Calculate expected graduation year from education timeline**
   - **Examples:**
     * "BTech 2022-2026" = Graduating in 2026 (current student)
     * "BTech 2020-2024" = Graduated in 2024
     * "MTech 2023-Present" = Currently pursuing, expected graduation 2025
   - **For current students (not yet graduated):**
     * Focus on internships, entry-level positions, and graduate programs
     * Include "fresher", "entry level", "new graduate" in searches
   - **For recent graduates (graduated within 2 years):**
     * Target entry-level and junior positions
     * Include "recent graduate", "fresher", "0-2 years experience"

3. **Technical Skills Proficiency Assessment:**
   - **Categorize skills by proficiency level:**
     * Expert (3+ years, project leadership, mentoring others)
     * Intermediate (1-3 years, independent work, complex projects)
     * Beginner (< 1 year, learning, basic projects, academic only)
   - **Extract technology stacks and frameworks**
   - **Identify primary programming languages and specializations**
   - **Note certifications and their recency/validity**

4. **Career Trajectory and Role Analysis:**
   - **Previous job titles and progression**
   - **Leadership and mentoring experience**
   - **Project complexity and team size managed**
   - **Industry domains and business contexts**
   - **Specialization areas (frontend, backend, full-stack, DevOps, data, etc.)**

5. **Location Preferences:**
	%s

**STRATEGIC JOB SEARCH EXECUTION USING TOOLS:**

**MANDATORY: You MUST use BOTH tools below for comprehensive job coverage**

1. **JSearch Tool (search_jsearch_jobs):**
	Execute 5-7 targeted job search queries using the search_jsearch_jobs function.
	Create queries based on:
	- Experience level + primary technology (e.g., "Senior Python Developer", "Junior React Engineer")
	- Technology combinations (e.g., "Python Django Developer", "React Node.js Engineer") 
	- Role-specific searches (e.g., "Full Stack Developer", "DevOps Engineer", "Data Scientist")
	- Industry-specific roles (e.g., "Software Engineer Fintech", "Developer Healthcare")
	- Career progression searches (e.g., "Lead Developer", "Principal Engineer")

2. **LinkUp Structured Tool (search_structured_jobs):**
	Execute ONE comprehensive job description using the search_structured_jobs function.
	Create a detailed job description that includes:
	- Complete professional profile with specific skills and experience level
	- Role preferences and technical expertise
	- Industry background and specializations
	- Work arrangement and location preferences
	
	Example: "Senior Software Engineer with 5+ years experience in Python, Django, React, and AWS. Looking for backend or full-stack roles in fintech or healthcare. Strong experience with microservices, API development, and cloud infrastructure."

**EXECUTION REQUIREMENTS:**
- MUST call search_jsearch_jobs multiple times (5-7 calls) with different strategic queries
- MUST call search_structured_jobs once with comprehensive job description
- Execute all function calls to maximize job discovery
- Vary search terms between tools to avoid duplication

Execute your comprehensive resume analysis and targeted job searches now using BOTH tools.
`, locationContext)

	return prompt
}

func (r *RerankingClient) RerankingPrompt(jobs []dtos.Job) string {
	jobsJSON, _ := json.MarshalIndent(jobs, "", "  ")

	prompt := fmt.Sprintf(`
	You are a professional career advisor and resume analyzer. Your task is to analyze the provided resume and score job opportunities based on their alignment with the candidate's profile. Fetch the job descriptions from the job portal url and analyze the job description based on the candidate's profile.
	
	**SCORING FRAMEWORK (0.0 - 1.0):**
	Use the FULL range from 0.0 to 1.0. Each score should reflect the true match quality:
	
	**Score Interpretation:**
	- 0.9-1.0: Perfect match - candidate exceeds requirements, ideal skills alignment
	- 0.8-0.89: Excellent match - meets all key requirements, strong skills overlap
	- 0.7-0.79: Very good match - meets most requirements, good experience fit
	- 0.6-0.69: Good match - solid alignment, minor gaps in skills/experience
	- 0.5-0.59: Moderate match - decent fit, some skill/experience gaps
	- 0.4-0.49: Fair match - basic alignment, notable gaps but possible
	- 0.3-0.39: Weak match - limited alignment, significant gaps
	- 0.2-0.29: Poor match - minimal relevance, major mismatches
	- 0.1-0.19: Very poor match - little to no alignment
	- 0.0-0.09: No match - completely irrelevant or the jobs which are not direct link to an opening rather its a job portal url, like 120000+ jobs in India etc.
	
	**EVALUATION CRITERIA:**
	
	1. **Technical Skills Assessment:**
		- Count exact skill matches vs. total required skills
		- Evaluate proficiency level alignment (junior/mid/senior requirements)
		- Consider related/transferable technologies
		- Account for learning curve for missing skills
	
	2. **Experience Level Matching:**
		- Compare candidate's years of experience with job requirements
		- Assess seniority level appropriateness (junior/mid/senior/lead)
		- Evaluate role complexity vs. candidate's background
		- Consider career progression readiness
	
	3. **Industry and Domain Fit:**
		- Match industry background and domain knowledge
		- Evaluate business context understanding
		- Consider transferability of experience across domains
		- Assess company size/culture fit
	
	4. **Role Responsibilities Alignment:**
		- Compare job duties with candidate's previous responsibilities
		- Evaluate leadership/management requirements vs. experience
		- Assess project scope and complexity match
		- Consider growth potential and career path
	
	**SCORING METHODOLOGY:**
	For each job, calculate the match score by honestly evaluating:
	- What percentage of required skills does the candidate have?
	- How well does their experience level match the role requirements?
	- How relevant is their industry/domain background?
	- How well do their previous responsibilities align with this role?
	
	**IMPORTANT GUIDELINES:**
	- Use the ENTIRE 0.0-1.0 range - don't cluster scores
	- Be honest about gaps and mismatches
	- Consider both current fit and growth potential
	- Vary your scores based on actual job-candidate alignment
	- Don't use predictable patterns or fixed decrements
	
	**OUTPUT FORMAT:**
	Return a JSON array of jobs ranked by match score (highest first):
	`+"```json"+`
	[
		{
			"job_index": 0,
			"match_score": 0.78,
			"match_reason": "Strong Python/Django skills match (7/8 required). 4 years experience fits mid-level role well. SaaS background aligns. Missing AWS experience reduces score.",
			"skills_matched": ["Python", "Django", "PostgreSQL", "REST APIs", "Git"],
			"experience_match": "Good fit - 4 years experience for mid-level role",
			"concerns": "Lacks AWS/cloud experience, no DevOps background"
		},
		{
			"job_index": 3,
			"match_score": 0.34,
			"match_reason": "Basic JavaScript knowledge but requires advanced React/Node.js. Experience level appropriate but major skill gaps in required technologies.",
			"skills_matched": ["JavaScript", "HTML", "CSS"],
			"experience_match": "Experience level matches but wrong technology focus",
			"concerns": "Significant frontend framework gaps, no Node.js experience"
		},
		{
			"job_index": 7,
			"match_score": 0.12,
			"match_reason": "Completely different technology stack (Java/Spring vs Python). Role requires 8+ years, candidate has 4. Different industry domain.",
			"skills_matched": [],
			"experience_match": "Under-qualified - requires 8+ years, candidate has 4",
			"concerns": "Wrong technology stack, insufficient experience, different domain"
		}
	]
	`+"```"+`
	
	**RESUME ANALYSIS STEPS:**
	1. Extract candidate's core skills, experience level, and background
	2. For each job, evaluate the 4 criteria above
	3. Calculate an honest match score using the full 0.0-1.0 range
	4. Provide detailed reasoning for each score
	
	**Jobs to Analyze:**
	%s
	
	Analyze the resume thoroughly and score these %d jobs using the full 0.0-1.0 range. Be honest and use varied scores that reflect true job-candidate alignment.
	`, string(jobsJSON), len(jobs))

	return prompt
}
