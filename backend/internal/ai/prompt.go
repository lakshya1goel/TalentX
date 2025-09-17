package ai

import (
	"fmt"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func (a *ProfileClient) CandidateProfilePrompt(locationPreference dtos.LocationPreference) string {
	locationContext := a.ParseLocationPreference(locationPreference)

	prompt := fmt.Sprintf(`
You are an expert resume analyzer and career strategist. Your task is to extract and analyze a candidate's profile from their resume to create a comprehensive professional summary.

**CRITICAL ANALYSIS FRAMEWORK:**

1. **Experience Level Calculation (EXTREMELY IMPORTANT):**
   - **ONLY count full-time professional work experience**
   - **INTERNSHIPS DO NOT COUNT as professional experience - treat candidates with only internships as FRESHERS**
   - **Current students/recent graduates with only internships = FRESHER category**
   - **Experience Categories:**
     * Fresher/Entry-level: 0-1 years professional experience (includes internship-only candidates)
     * Junior: 1-3 years professional experience
     * Mid-level: 3-6 years professional experience  
     * Senior: 6-10 years professional experience
     * Lead/Principal: 10+ years professional experience

2. **Education Timeline Analysis:**
    - Calculate graduation year from education details
    - Identify current students vs graduates
    - Note degree types, institutions, and academic achievements
    - Extract relevant coursework and academic projects

3. **Technical Skills Assessment:**
    - **Categorize skills by proficiency level:**
      * Expert: 3+ years, project leadership, mentoring others
      * Intermediate: 1-3 years, independent work, complex projects  
      * Beginner: < 1 year, learning, basic projects, academic only
    - Extract programming languages, frameworks, tools, and technologies
    - Identify primary technology stacks and specializations
    - Note certifications and their validity/recency

4. **Career Trajectory Analysis:**
    - Previous job titles and career progression
    - Leadership and mentoring experience
    - Project complexity and impact
    - Industry domains and business contexts
    - Team sizes managed and responsibilities

5. **Location Preferences:**
	%s

**OUTPUT FORMAT:**
Provide a comprehensive, structured profile summary that includes:

- **Professional Level:** [Fresher/Junior/Mid-level/Senior/Lead] with years of experience
- **Primary Skills:** Top 5-7 technical skills with proficiency levels
- **Technology Stack:** Main programming languages, frameworks, and tools
- **Industry Experience:** Domains worked in or interested in
- **Education:** Degree, graduation year, institution
- **Suitable Job Titles:** 5-8 specific job titles they would be suitable for
- **Experience Highlights:** Key projects, achievements, and responsibilities
- **Work Preferences:** Location and work arrangement preferences
- **Career Stage:** Current career stage and trajectory

**IMPORTANT NOTES:**
- Be precise about experience calculation - do not inflate internship experience
- Focus on concrete skills with evidence from the resume
- Provide realistic job title suggestions based on actual experience level
- Consider both technical skills and soft skills/leadership abilities
- Format the output as clear, structured text (not JSON)
`, locationContext)

	return prompt
}

func (a *ProfileClient) ParseLocationPreference(locationPreference dtos.LocationPreference) string {
	locationGuidance := "\n\n5. **Location Preferences:** "

	if len(locationPreference.Types) == 1 {
		switch locationPreference.Types[0] {
		case "remote":
			locationGuidance += "The candidate is looking for REMOTE work opportunities only.\n" +
				"- Focus on remote-friendly companies and positions\n" +
				"- Include \"remote\" in search queries when relevant\n" +
				"- Prioritize companies known for remote work culture"
		case "onsite":
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("The candidate is looking for ON-SITE work in %s.\n", joinLocations(locationPreference.Locations)) +
					"- Focus on companies and positions located in or near these areas\n" +
					"- Include location-specific searches\n" +
					"- Consider commuting distance and local job market"
			} else {
				locationGuidance += "The candidate is looking for ON-SITE work opportunities.\n" +
					"- Focus on local companies and positions\n" +
					"- Include location-specific searches when possible"
			}
		case "hybrid":
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("The candidate is open to HYBRID work arrangements in %s.\n", joinLocations(locationPreference.Locations)) +
					"- Look for both remote and on-site opportunities in these areas\n" +
					"- Focus on companies offering flexible work arrangements\n" +
					"- Include both remote and location-specific searches"
			} else {
				locationGuidance += "The candidate is open to HYBRID work arrangements.\n" +
					"- Look for companies offering flexible work arrangements\n" +
					"- Include both remote and location-specific searches"
			}
		}
	} else {
		locationGuidance += "The candidate is open to multiple work arrangements:\n"

		hasRemote := contains(locationPreference.Types, "remote")
		hasOnsite := contains(locationPreference.Types, "onsite")
		hasHybrid := contains(locationPreference.Types, "hybrid")

		if hasRemote {
			locationGuidance += "- REMOTE: Include remote-friendly positions and companies\n"
		}
		if hasOnsite {
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("- ON-SITE: Include positions in %s\n", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "- ON-SITE: Include local on-site positions\n"
			}
		}
		if hasHybrid {
			if len(locationPreference.Locations) > 0 {
				locationGuidance += fmt.Sprintf("- HYBRID: Include flexible arrangements in %s\n", joinLocations(locationPreference.Locations))
			} else {
				locationGuidance += "- HYBRID: Include flexible work arrangements\n"
			}
		}
		locationGuidance += "- Cast a wide net to capture all preferred work arrangements\n" +
			"- Use varied search terms to find opportunities matching any of these preferences"
	}
	return locationGuidance
}

func (a *JobClient) JobSearchPrompt(profile string) string {
	prompt := fmt.Sprintf(`
You are a senior career strategist and job search specialist with expertise in tech hiring and talent matching. You have been provided with a comprehensive candidate profile. Use this profile to execute strategic job searches with precise targeting.

**CANDIDATE PROFILE:**
%s

**STRATEGIC JOB SEARCH EXECUTION USING TOOLS:**

Based on the candidate profile above, execute comprehensive job searches using the following strategy:

**MANDATORY: You MUST use BOTH tools below for comprehensive job coverage**

1. **LinkUp Structured Tool (search_structured_jobs):**
	Execute ONE comprehensive job description using the search_structured_jobs function.
	Create a detailed job description based on the candidate profile that includes:
	- Complete professional profile with specific skills and experience level from the profile
	- Role preferences and technical expertise identified in the profile
	- Industry background and specializations mentioned in the profile
	- Work arrangement and location preferences from the profile
NOTE: If location preferences are mentioned, make sure to include the location in the job description.

	Example: "Senior Software Engineer with 5+ years experience in Python, Django, React, and AWS. Looking for backend or full-stack roles in fintech or healthcare. Strong experience with microservices, API development, and cloud infrastructure."

2. **JSearch Tool (search_jsearch_jobs):**
	Execute 5 to 7 targeted job search queries using the search_jsearch_jobs function.
	Create queries based on the candidate profile:
	- Experience level + primary technology from profile (e.g., "Senior Python Developer", "Junior React Engineer")
	- Technology combinations mentioned in profile (e.g., "Python Django Developer", "React Node.js Engineer") 
	- Role-specific searches based on suitable job titles in profile (e.g., "Full Stack Developer", "DevOps Engineer", "Data Scientist")
	- Industry-specific roles based on industry experience in profile (e.g., "Software Engineer Fintech", "Developer Healthcare")
	- Career progression searches based on career stage in profile (e.g., "Lead Developer", "Principal Engineer")
  - IMPORTANT: Include different locations for each job search if location preferences are mentioned. If Remote and hybrid are also mntioned them make 5 queries for remote and 2 queries for onsite only.
	
**MANDATORY EXECUTION SEQUENCE - FOLLOW EXACTLY:**

STEP 1: Call search_jsearch_jobs function 5-7 times with targeted queries based on the profile:
1. [Experience Level from profile] + [Primary Technology from profile] 
2. [Technology Stack from profile] combination 
3. [Suitable Job Title from profile]
4. [Industry Experience from profile] specific role
5. [Career Level from profile] position
6. [Additional relevant search based on profile skills]
7. [Location-specific search if relevant]

STEP 2: Call search_structured_jobs function 1 time with comprehensive profile description

**IMPORTANT GUIDELINES:**
- Use the candidate's actual experience level and skills from the profile
- Target job searches based on the suitable job titles mentioned in the profile
- Consider the work preferences and location requirements from the profile
- Match the seniority level accurately (don't search for senior roles if candidate is fresher)
- Include relevant technologies and frameworks mentioned in the profile
- Consider industry experience and domain expertise from the profile

**YOU MUST EXECUTE ALL FUNCTION CALLS - DO NOT STOP AFTER THE FIRST ONE**

Execute now - make all function calls based on the candidate profile:
`, profile)

	return prompt
}

func (a *RankingClient) RankingPrompt() string {
	systemMessage := `You are a job matching assistant. Your task is to evaluate multiple jobs based on their match with the candidate's profile, taking into account the job title, the skills required, the seniority level, the physical location (where the company offering the work is based in) and the working location (remote/hybrid/on-site). You then have to produce a match score (between 0 and 100) for each job and justify that match score explaining your reasons.

Evaluation Criteria:
1. Job title alignment with candidate's potential roles
2. Required skills match with candidate's skills
3. Seniority level alignment (internship, entry level, junior, mid-level, senior)
4. Location preferences and work arrangement compatibility
5. Industry and domain experience relevance
6. Overall career trajectory fit

Provide your evaluation in the following JSON format for ALL jobs:
{
	"evaluations": [
		{
			"job_index": 0,
			"match_score": <integer between 0-100>,
			"reasons": "<detailed explanation of the match evaluation>",
			"skills_matched": ["<list of matched skills>"],
			"experience_match": "<assessment of experience level fit>"
		},
		{
			"job_index": 1,
			"match_score": <integer between 0-100>,
			"reasons": "<detailed explanation of the match evaluation>",
			"skills_matched": ["<list of matched skills>"],
			"experience_match": "<assessment of experience level fit>"
		}
	]
}

IMPORTANT: Provide evaluations for ALL jobs in the same order they are presented. Use job_index to match each evaluation to its corresponding job.`

	return systemMessage
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
