export interface Job {
  title: string;
  company?: string;
  location?: string;
  description?: string;
  url: string;
  source: string;
}

export interface RankedJob {
  job: Job;
  percent_match: number;
  match_reason: string;
  skills_matched: string[];
  experience_match: string;
}

export interface ErrorResponse {
  error: string;
  success: boolean;
  timestamp: string;
} 