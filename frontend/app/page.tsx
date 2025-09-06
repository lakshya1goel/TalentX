'use client';

import { useState } from 'react';
import Header from '../components/Header';
import ResumeUploader from '../components/ResumeUploader';
import JobsList from '../components/JobsList';
import LoadingSpinner from '../components/LoadingSpinner';
import { Job } from '../types/job';

export default function Home() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleJobsReceived = (newJobs: Job[]) => {
    setJobs(newJobs);
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handleLoading = (isLoading: boolean) => {
    setLoading(isLoading);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      <main className="py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          {/* Hero Section */}
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-gray-900 sm:text-5xl md:text-6xl">
              Find Your Perfect Job
            </h1>
            <p className="mt-3 max-w-md mx-auto text-base text-gray-500 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
              Upload your resume and let our AI find the best job opportunities tailored to your skills and experience.
            </p>
          </div>

          {/* Resume Uploader */}
          <div className="mb-8">
            <ResumeUploader
              onJobsReceived={handleJobsReceived}
              onError={handleError}
              onLoading={handleLoading}
            />
          </div>

          {/* Error Message */}
          {error && (
            <div className="max-w-2xl mx-auto mb-8">
              <div className="bg-red-50 border border-red-200 rounded-md p-4">
                <div className="flex">
                  <div className="flex-shrink-0">
                    <svg
                      className="h-5 w-5 text-red-400"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm text-red-800">{error}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Loading Spinner */}
          {loading && (
            <div className="max-w-2xl mx-auto mb-8">
              <div className="bg-white rounded-lg shadow p-6">
                <div className="text-center">
                  <LoadingSpinner />
                  <p className="mt-4 text-sm text-gray-600">
                    Analyzing your resume and finding matching jobs...
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Jobs List */}
          {!loading && jobs.length > 0 && (
            <JobsList jobs={jobs} />
          )}

          {/* No Jobs Found */}
          {!loading && jobs.length === 0 && !error && (
            <div className="text-center mt-12">
              <svg
                className="mx-auto h-12 w-12 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900">No jobs yet</h3>
              <p className="mt-1 text-sm text-gray-500">
                Upload your resume to get started with personalized job recommendations.
              </p>
            </div>
          )}
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-white border-t mt-20">
        <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <p className="text-sm text-gray-500">
              &copy; 2024 Job Assistant. Powered by AI to help you find your dream job.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
