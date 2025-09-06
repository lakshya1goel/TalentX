'use client';

import { useState, useCallback } from 'react';
import { uploadResumeAndGetJobs } from '../utils/api';
import { Job } from '../types/job';

interface ResumeUploaderProps {
  onJobsReceived: (jobs: Job[]) => void;
  onError: (error: string) => void;
  onLoading: (loading: boolean) => void;
}

export default function ResumeUploader({ onJobsReceived, onError, onLoading }: ResumeUploaderProps) {
  const [isDragOver, setIsDragOver] = useState(false);

  const handleFileUpload = useCallback(async (file: File) => {
    // Validate file type
    if (file.type !== 'application/pdf') {
      onError('Please upload a PDF file only.');
      return;
    }

    // Validate file size (10MB limit)
    if (file.size > 10 * 1024 * 1024) {
      onError('File size must be less than 10MB.');
      return;
    }

    try {
      onLoading(true);
      onError(''); // Clear any previous errors
      const jobs = await uploadResumeAndGetJobs(file);
      onJobsReceived(jobs);
    } catch (error) {
      onError(error instanceof Error ? error.message : 'An error occurred while processing your resume.');
    } finally {
      onLoading(false);
    }
  }, [onJobsReceived, onError, onLoading]);

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

  return (
    <div className="w-full max-w-2xl mx-auto">
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