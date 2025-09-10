'use client';

import { useState } from 'react';
// import Header from '../components/Header';
import ResumeUploader from '../components/ResumeUploader';
import JobsList from '../components/JobsList';
import LoadingSpinner from '../components/LoadingSpinner';
import { Job } from '../types/job';

export default function Home() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showUploader, setShowUploader] = useState(true);

  const handleJobsReceived = (newJobs: Job[]) => {
    setJobs(newJobs);
    setShowUploader(false);
  };

  const handleError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handleLoading = (isLoading: boolean) => {
    setLoading(isLoading);
  };

  const handleSearchAgain = () => {
    setShowUploader(true);
    setJobs([]);
    setError('');
  };

  return (
    <div className="min-h-screen relative overflow-hidden" style={{
      background: 'linear-gradient(135deg, #0a0a0a, #1a1a1a)'
    }}>
      {/* Background Vectors */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-32 h-32 opacity-10">
          <div className="w-full h-full rounded-full" style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)'
          }}></div>
        </div>
        <div className="absolute top-40 right-20 w-24 h-24 opacity-10">
          <div className="w-full h-full" style={{
            background: 'linear-gradient(45deg, #16a085, #138f7a)',
            clipPath: 'polygon(50% 0%, 0% 100%, 100% 100%)'
          }}></div>
        </div>
        <div className="absolute bottom-20 left-1/4 w-40 h-40 opacity-5">
          <div className="w-full h-full rounded-full" style={{
            background: 'radial-gradient(circle, #16a085, transparent)'
          }}></div>
        </div>
        <div className="absolute bottom-40 right-10 w-28 h-28 opacity-10">
          <div className="w-full h-full" style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)',
            borderRadius: '0 50% 50% 50%'
          }}></div>
        </div>
      </div>

      {/* Content */}
      <main className="py-12 px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="max-w-7xl mx-auto">
          {showUploader && (
            <div className="text-center mb-12">
              <h1 className="text-4xl font-bold text-white sm:text-5xl md:text-6xl">
                Find Your Perfect Job
              </h1>
              <p className="mt-3 max-w-md mx-auto text-base text-gray-300 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
                Upload your resume and let our AI find the best job opportunities tailored to your skills and experience.
              </p>
            </div>
          )}

          {showUploader && !loading && (
            <div className="mb-8">
              <ResumeUploader
                onJobsReceived={handleJobsReceived}
                onError={handleError}
                onLoading={handleLoading}
              />
            </div>
          )}

          {!showUploader && !loading && (
            <div className="text-center mb-8">
              <button
                onClick={handleSearchAgain}
                className="inline-flex items-center px-6 py-3 border-0 text-base font-medium rounded-lg text-white transition-all duration-200 hover:scale-105"
                style={{
                  background: 'linear-gradient(135deg, #16a085, #138f7a)',
                  boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
                }}
              >
                <svg
                  className="mr-2 -ml-1 w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 4v16m8-8H4"
                  />
                </svg>
                Search Again
              </button>
            </div>
          )}

          {error && (
            <div className="max-w-2xl mx-auto mb-8">
              <div className="rounded-lg p-4" style={{
                background: 'linear-gradient(135deg, #2d1b1b, #1a1a1a)',
                border: '1px solid rgba(220, 38, 38, 0.3)'
              }}>
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
                    <p className="text-sm text-red-300">{error}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {loading && (
            <div className="max-w-2xl mx-auto mb-8">
              <div className="rounded-lg p-6" style={{
                background: 'linear-gradient(135deg, #0a0a0a, #1a1a1a)',
                border: '1px solid rgba(29,205,159,.2)'
              }}>
                <div className="text-center">
                  <LoadingSpinner />
                  <p className="mt-4 text-sm text-gray-300">
                    Analyzing your resume and finding matching jobs...
                  </p>
                </div>
              </div>
            </div>
          )}

          {!loading && jobs.length > 0 && (
            <JobsList jobs={jobs} />
          )}

          {!loading && jobs.length === 0 && !error && !showUploader && (
            <div className="text-center mt-12">
              <svg
                className="mx-auto h-12 w-12 text-gray-500"
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
              <h3 className="mt-2 text-sm font-medium text-gray-200">No jobs found</h3>
              <p className="mt-1 text-sm text-gray-400">
                Try adjusting your search criteria or upload a different resume.
              </p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
