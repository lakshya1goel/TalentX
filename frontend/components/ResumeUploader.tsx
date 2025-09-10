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
      const rankedJobs = await uploadResumeAndGetJobs(file, locationPreference);
      const jobs: Job[] = rankedJobs.map(rankedJob => rankedJob.job);
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
      {/* Location Preferences */}
      <div className="p-6 rounded-xl backdrop-blur-sm" style={{
        background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
        border: '1px solid rgba(29,205,159,.2)'
      }}>
        <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
          <svg className="w-5 h-5 mr-2 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Location Preferences
        </h3>
        
        <div className="space-y-6">
          <div>
            <label className="text-sm font-medium text-gray-300 block mb-3">
              Work Arrangement
            </label>
            <div className="grid grid-cols-3 gap-3">
              {(['remote', 'onsite', 'hybrid'] as const).map((type) => (
                <label key={type} className="flex items-center justify-center p-3 rounded-lg cursor-pointer transition-all duration-200 hover:scale-105" style={{
                  background: locationPreference.types.includes(type) 
                    ? 'linear-gradient(135deg, #16a085, #138f7a)'
                    : 'rgba(29,205,159,.1)',
                  border: `1px solid ${locationPreference.types.includes(type) ? 'rgba(29,205,159,.4)' : 'rgba(29,205,159,.2)'}`
                }}>
                  <input
                    type="checkbox"
                    checked={locationPreference.types.includes(type)}
                    onChange={(e) => handleLocationTypeChange(type, e.target.checked)}
                    className="sr-only"
                  />
                  <span className="text-sm font-medium text-white capitalize">
                    {type}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {needsLocation && (
            <div>
              <label className="text-sm font-medium text-gray-300 block mb-3">
                Specific Locations
              </label>
              <div className="flex gap-2 mb-3">
                <input
                  type="text"
                  value={locationInput}
                  onChange={(e) => setLocationInput(e.target.value)}
                  onKeyPress={handleLocationInputKeyPress}
                  placeholder="e.g., San Francisco, CA"
                  className="flex-1 px-4 py-2 rounded-lg text-white placeholder-gray-400 border-0 focus:outline-none focus:ring-2 focus:ring-green-500"
                  style={{
                    background: 'rgba(10,10,10,0.6)',
                    border: '1px solid rgba(29,205,159,.2)'
                  }}
                />
                <button
                  type="button"
                  onClick={addLocation}
                  disabled={!locationInput.trim()}
                  className="px-4 py-2 text-white rounded-lg font-medium transition-all duration-200 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
                  style={{
                    background: 'linear-gradient(135deg, #16a085, #138f7a)'
                  }}
                >
                  Add
                </button>
              </div>
              
              {locationPreference.locations && locationPreference.locations.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {locationPreference.locations.map((location, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center px-3 py-1 rounded-full text-sm text-white"
                      style={{
                        background: 'linear-gradient(135deg, rgba(22,160,133,0.3), rgba(19,143,122,0.3))',
                        border: '1px solid rgba(29,205,159,.3)'
                      }}
                    >
                      {location}
                      <button
                        type="button"
                        onClick={() => removeLocation(location)}
                        className="ml-2 text-green-300 hover:text-white transition-colors duration-200"
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

      {/* File Upload Area */}
      <div
        className={`border-2 border-dashed rounded-xl p-8 text-center transition-all duration-300 cursor-pointer ${
          isDragOver ? 'scale-105' : 'hover:scale-[1.02]'
        }`}
        style={{
          background: isDragOver 
            ? 'linear-gradient(135deg, rgba(22,160,133,0.1), rgba(19,143,122,0.1))'
            : 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
          borderColor: isDragOver ? '#16a085' : 'rgba(29,205,159,.3)'
        }}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
      >
        <div className="space-y-4">
          <div className="mx-auto w-16 h-16 flex items-center justify-center rounded-full" style={{
            background: 'linear-gradient(135deg, rgba(22,160,133,0.2), rgba(19,143,122,0.2))'
          }}>
            <svg
              className="w-8 h-8 text-green-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
          </div>
          <div>
            <p className="text-xl font-semibold text-white mb-2">
              Upload your resume
            </p>
            <p className="text-sm text-gray-400 mb-4">
              Drag and drop your PDF resume here, or click to browse
            </p>
            <p className="text-xs text-gray-500">
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
              className="inline-flex items-center px-6 py-3 text-sm font-semibold rounded-lg text-white cursor-pointer transition-all duration-200 hover:scale-105"
              style={{
                background: 'linear-gradient(135deg, #16a085, #138f7a)',
                boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
              }}
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Choose File
            </label>
          </div>
        </div>
      </div>
    </div>
  );
} 