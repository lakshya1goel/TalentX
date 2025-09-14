import { Job, ErrorResponse } from '../types/job';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8084';

export interface LocationPreference {
  types: ('remote' | 'onsite' | 'hybrid')[];
  locations?: string[];
}

export interface RankedJob {
  job: Job;
  percent_match: number;
  match_reason: string;
  skills_matched: string[];
  experience_match: string;
}

export interface JobSearchResponse {
  jobs: RankedJob[];
  total: number;
  success: boolean;
}

export async function uploadResumeAndGetJobs(
  file: File, 
  locationPreference: LocationPreference,
  apiKey: string
): Promise<RankedJob[]> {
  const formData = new FormData();
  formData.append('resume', file);
  formData.append('api_key', apiKey);
  
  locationPreference.types.forEach(type => {
    formData.append('location_types', type);
  });
  
  if (locationPreference.locations && locationPreference.locations?.length > 0) {
    locationPreference.locations.forEach(location => {
      formData.append('locations', location);
    });
  }

  const response = await fetch(`${API_BASE_URL}/api/job/`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    throw new Error(errorData.error || 'Failed to fetch jobs');
  }

  const result: JobSearchResponse = await response.json();
  return result.jobs;
}

export interface PaginationParams {
  page: number;
  pageSize: number;
}

export interface PaginatedResult<T> {
  data: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
}

export function paginateArray<T>(
  array: T[], 
  { page, pageSize }: PaginationParams
): PaginatedResult<T> {
  const total = array?.length || 0;
  const totalPages = Math.ceil(total / pageSize);
  const startIndex = (page - 1) * pageSize;
  const endIndex = Math.min(startIndex + pageSize, total);
  
  return {
    data: array.slice(startIndex, endIndex),
    pagination: {
      page,
      pageSize,
      total,
      totalPages,
      hasNext: page < totalPages,
      hasPrev: page > 1,
    },
  };
}