import { Job, ErrorResponse } from '../types/job';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8084';

export interface LocationPreference {
  types: ('remote' | 'onsite' | 'hybrid')[];
  locations?: string[];
}

export async function uploadResumeAndGetJobs(
  file: File, 
  locationPreference: LocationPreference
): Promise<Job[]> {
  const formData = new FormData();
  formData.append('resume', file);
  
  // Append multiple location types
  locationPreference.types.forEach(type => {
    formData.append('location_types', type);
  });
  
  // Append multiple locations if provided
  if (locationPreference.locations && locationPreference.locations.length > 0) {
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

  return response.json();
} 