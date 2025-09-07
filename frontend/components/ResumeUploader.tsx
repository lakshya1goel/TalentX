'use client';

import { useState, useCallback } from 'react';
import { uploadResumeAndGetJobs, LocationPreference } from '../utils/api';
import { Job } from '../types/job';

interface ResumeUploaderProps {
  onJobsReceived: (jobs: Job[]) => void;
  onError: (error: string) => void;
  onLoading: (loading: boolean) => void;
}

export default function ResumeUploader({ onJobsReceived, onError, onLoading }: ResumeUploaderProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const [locationPreference, setLocationPreference] = useState<LocationPreference>({
    types: ['remote'],
    locations: []
  });
  const [locationInput, setLocationInput] = useState('');

  const handleFileUpload = useCallback(async (file: File) => {
    if (file.type !== 'application/pdf') {
      onError('Please upload a PDF file only.');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      onError('File size must be less than 10MB.');
      return;
    }

    const needsLocation = locationPreference.types.some(type => type === 'onsite' || type === 'hybrid');
    if (needsLocation && (!locationPreference.locations || locationPreference.locations.length === 0)) {
      onError('Please specify at least one location for onsite or hybrid positions.');
      return;
    }

    try {
      onLoading(true);
      onError('');
      const jobs = await uploadResumeAndGetJobs(file, locationPreference);
      onJobsReceived(jobs);
    } catch (error) {
      onError(error instanceof Error ? error.message : 'An error occurred while processing your resume.');
    } finally {
      onLoading(false);
    }
  }, [onJobsReceived, onError, onLoading, locationPreference]);

  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileUpload(files[0]);
    }
  }, [handleFileUpload]);

  const handleDragOver = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleFileInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileUpload(files[0]);
    }
  }, [handleFileUpload]);

  const handleLocationTypeChange = (type: LocationPreference['types'][0], checked: boolean) => {
    setLocationPreference(prev => ({
      ...prev,
      types: checked 
        ? [...prev.types, type]
        : prev.types.filter(t => t !== type)
    }));
  };

  const addLocation = () => {
    if (locationInput.trim() && !locationPreference.locations?.includes(locationInput.trim())) {
      setLocationPreference(prev => ({
        ...prev,
        locations: [...(prev.locations || []), locationInput.trim()]
      }));
      setLocationInput('');
    }
  };

  const removeLocation = (locationToRemove: string) => {
    setLocationPreference(prev => ({
      ...prev,
      locations: prev.locations?.filter(loc => loc !== locationToRemove) || []
    }));
  };

  const handleLocationInputKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addLocation();
    }
  };

  const needsLocation = locationPreference.types.some(type => type === 'onsite' || type === 'hybrid');

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <div className="bg-white p-6 rounded-lg border border-gray-200">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Location Preferences</h3>
        
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-700 block mb-2">
              Work Arrangement (Select all that apply)
            </label>
            <div className="space-y-2">
              {(['remote', 'onsite', 'hybrid'] as const).map((type) => (
                <label key={type} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={locationPreference.types.includes(type)}
                    onChange={(e) => handleLocationTypeChange(type, e.target.checked)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <span className="ml-2 text-sm text-gray-700 capitalize">
                    {type}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {needsLocation && (
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-2">
                Locations (Required for onsite and hybrid positions)
              </label>
              <div className="flex space-x-2 mb-2">
                <input
                  type="text"
                  value={locationInput}
                  onChange={(e) => setLocationInput(e.target.value)}
                  onKeyPress={handleLocationInputKeyPress}
                  placeholder="e.g., San Francisco, CA"
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-black caret-black"
                />
                <button
                  type="button"
                  onClick={addLocation}
                  disabled={!locationInput.trim()}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Add
                </button>
              </div>
              
              {locationPreference.locations && locationPreference.locations.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {locationPreference.locations.map((location, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
                    >
                      {location}
                      <button
                        type="button"
                        onClick={() => removeLocation(location)}
                        className="ml-2 text-blue-600 hover:text-blue-800"
                      >
                        ×
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <div
        className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
          isDragOver
            ? 'border-blue-400 bg-blue-50'
            : 'border-gray-300 hover:border-gray-400'
        }`}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
      >
        <div className="space-y-4">
          <div className="mx-auto w-12 h-12 text-gray-400">
            <svg
              className="w-full h-full"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          </div>
          <div>
            <p className="text-lg font-medium text-gray-900">
              Upload your resume
            </p>
            <p className="text-sm text-gray-500">
              Drag and drop your PDF resume here, or click to browse
            </p>
            <p className="text-xs text-gray-400 mt-2">
              Maximum file size: 10MB | Supported format: PDF
            </p>
          </div>
          <div>
            <input
              type="file"
              accept=".pdf"
              onChange={handleFileInputChange}
              className="hidden"
              id="resume-upload"
            />
            <label
              htmlFor="resume-upload"
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 cursor-pointer"
            >
              Choose File
            </label>
          </div>
        </div>
      </div>
    </div>
  );
} 