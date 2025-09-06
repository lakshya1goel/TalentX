import { Job, ErrorResponse } from '../types/job';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8084';

export async function uploadResumeAndGetJobs(file: File): Promise<Job[]> {
  const formData = new FormData();
  formData.append('resume', file);

  const response = await fetch(`${API_BASE_URL}/api/job/`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    throw new Error(errorData.error || 'Failed to fetch jobs');
  }

  return response.json();
} 