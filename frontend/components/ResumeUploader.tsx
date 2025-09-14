'use client';

import { useState, useCallback } from 'react';
import { uploadResumeAndGetJobs, LocationPreference, RankedJob } from '../utils/api';
import { Job } from '../types/job';

interface ResumeUploaderProps {
  onJobsReceived: (jobs: RankedJob[]) => void;
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
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [apiKey, setApiKey] = useState('');

  const handleFileUpload = useCallback((file: File) => {
    if (file.type !== 'application/pdf') {
      onError('Please upload a PDF file only.');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      onError('File size must be less than 10MB.');
      return;
    }

    setSelectedFile(file);
    onError(''); // Clear any previous errors
  }, [onError]);

  const handleSubmit = useCallback(async () => {
    if (!selectedFile) {
      onError('Please upload a resume file first.');
      return;
    }

    if (!apiKey.trim()) {
      onError('Please enter your Gemini API key.');
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
      const rankedJobs = await uploadResumeAndGetJobs(selectedFile, locationPreference, apiKey.trim());
      onJobsReceived(rankedJobs);
    } catch (error) {
      onError(error instanceof Error ? error.message : 'An error occurred while processing your resume.');
    } finally {
      onLoading(false);
    }
  }, [selectedFile, locationPreference, apiKey, onJobsReceived, onError, onLoading]);

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
    <div className="w-full max-w-2xl mx-auto space-y-4 sm:space-y-6 px-4 sm:px-0">
      <div className="p-4 sm:p-6 rounded-xl backdrop-blur-sm" style={{
        background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
        border: '1px solid rgba(29,205,159,.2)'
      }}>
        <h3 className="text-base sm:text-lg font-semibold text-white mb-3 sm:mb-4 flex items-center">
          <svg className="w-4 sm:w-5 h-4 sm:h-5 mr-2 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m0 0a2 2 0 012 2m-2-2v6a2 2 0 01-2 2H9a2 2 0 01-2-2V9a2 2 0 012-2h6z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7V5a2 2 0 012-2h2a2 2 0 012 2v2" />
          </svg>
          Gemini API Key
        </h3>
        <div className="space-y-2">
          <label className="text-xs sm:text-sm font-medium text-gray-300 block">
            Enter your Google Gemini API Key
          </label>
          <input
            type="password"
            value={apiKey}
            onChange={(e) => setApiKey(e.target.value)}
            placeholder="AIza..."
            className="w-full px-3 sm:px-4 py-2 sm:py-3 rounded-lg text-white placeholder-gray-400 border-0 focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
            style={{
              background: 'rgba(10,10,10,0.6)',
              border: '1px solid rgba(29,205,159,.2)'
            }}
          />
          <p className="text-xs text-gray-400">
            Get your API key from{' '}
            <a 
              href="https://makersuite.google.com/app/apikey" 
              target="_blank" 
              rel="noopener noreferrer" 
              className="text-green-400 hover:text-green-300 underline"
            >
              Google AI Studio
            </a>
          </p>
        </div>
      </div>

      {/* Location Preferences */}
      <div className="p-4 sm:p-6 rounded-xl backdrop-blur-sm" style={{
        background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
        border: '1px solid rgba(29,205,159,.2)'
      }}>
        <h3 className="text-base sm:text-lg font-semibold text-white mb-3 sm:mb-4 flex items-center">
          <svg className="w-4 sm:w-5 h-4 sm:h-5 mr-2 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Location Preferences
        </h3>
        
        <div className="space-y-4 sm:space-y-6">
          <div>
            <label className="text-xs sm:text-sm font-medium text-gray-300 block mb-2 sm:mb-3">
              Work Arrangement
            </label>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-2 sm:gap-3">
              {(['remote', 'onsite', 'hybrid'] as const).map((type) => (
                <label key={type} className="flex items-center justify-center p-2 sm:p-3 rounded-lg cursor-pointer transition-all duration-200 hover:scale-105" style={{
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
                  <span className="text-xs sm:text-sm font-medium text-white capitalize">
                    {type}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {needsLocation && (
            <div>
              <label className="text-xs sm:text-sm font-medium text-gray-300 block mb-2 sm:mb-3">
                Specific Locations
              </label>
              <div className="flex flex-col sm:flex-row gap-2 mb-3">
                <input
                  type="text"
                  value={locationInput}
                  onChange={(e) => setLocationInput(e.target.value)}
                  onKeyPress={handleLocationInputKeyPress}
                  placeholder="e.g., San Francisco, CA"
                  className="flex-1 px-3 sm:px-4 py-2 rounded-lg text-white placeholder-gray-400 border-0 focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                  style={{
                    background: 'rgba(10,10,10,0.6)',
                    border: '1px solid rgba(29,205,159,.2)'
                  }}
                />
                <button
                  type="button"
                  onClick={addLocation}
                  disabled={!locationInput.trim()}
                  className="px-3 sm:px-4 py-2 text-white rounded-lg font-medium transition-all duration-200 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
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
        className={`border-2 border-dashed rounded-xl p-4 sm:p-8 text-center transition-all duration-300 cursor-pointer ${
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
        <div className="space-y-3 sm:space-y-4">
          <div className="mx-auto w-12 sm:w-16 h-12 sm:h-16 flex items-center justify-center rounded-full" style={{
            background: 'linear-gradient(135deg, rgba(22,160,133,0.2), rgba(19,143,122,0.2))'
          }}>
            <svg
              className="w-6 sm:w-8 h-6 sm:h-8 text-green-400"
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
            <p className="text-lg sm:text-xl font-semibold text-white mb-1 sm:mb-2">
              Upload your resume
            </p>
            <p className="text-xs sm:text-sm text-gray-400 mb-3 sm:mb-4">
              Drag and drop your PDF resume here, or click to browse
            </p>
            <p className="text-xs text-gray-500 mb-4">
              Maximum file size: 10MB | Supported format: PDF
            </p>
            
            {/* Selected File Display */}
            {selectedFile && (
              <div className="mb-4 p-3 rounded-lg" style={{
                background: 'rgba(22,160,133,0.1)',
                border: '1px solid rgba(29,205,159,.3)'
              }}>
                <div className="flex items-center">
                  <svg className="w-4 h-4 text-green-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span className="text-sm text-white font-medium">{selectedFile.name}</span>
                  <span className="text-xs text-gray-400 ml-2">({(selectedFile.size / 1024 / 1024).toFixed(2)} MB)</span>
                </div>
              </div>
            )}
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
              className="inline-flex items-center justify-center px-4 sm:px-6 py-2 sm:py-3 text-xs sm:text-sm font-semibold rounded-lg text-white cursor-pointer transition-all duration-200 hover:scale-105"
              style={{
                background: selectedFile 
                  ? 'rgba(29,205,159,.2)' 
                  : 'linear-gradient(135deg, #16a085, #138f7a)',
                boxShadow: selectedFile 
                  ? '0 0 0 1px rgba(29,205,159,.3)' 
                  : '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)',
                border: selectedFile ? '1px solid rgba(29,205,159,.3)' : 'none'
              }}
            >
              <svg className="w-3 sm:w-4 h-3 sm:h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              {selectedFile ? 'Change File' : 'Choose File'}
            </label>
          </div>
        </div>
      </div>

      {/* Submit Button - Outside the upload box */}
      <div className="text-center mt-6">
        <button
          onClick={handleSubmit}
          disabled={!selectedFile}
          className="inline-flex items-center justify-center px-8 sm:px-12 py-3 sm:py-4 text-sm sm:text-base font-semibold rounded-lg text-white transition-all duration-200 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
          style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)',
            boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
          }}
        >
          <svg className="w-4 sm:w-5 h-4 sm:h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          Find Jobs
        </button>
      </div>
    </div>
  );
} 