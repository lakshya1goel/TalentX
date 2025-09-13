'use client';

import { useState } from 'react';
// import Header from '../components/Header';
import ResumeUploader from '../components/ResumeUploader';
import JobsList from '../components/JobsList';
import LoadingSpinner from '../components/LoadingSpinner';

import { RankedJob } from '../types/job';

export default function Home() {
  const [jobs, setJobs] = useState<RankedJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showUploader, setShowUploader] = useState(true);

  const handleJobsReceived = (newJobs: RankedJob[]) => {
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
      {/* Background Vectors - Responsive */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-10 sm:top-20 left-4 sm:left-10 w-16 sm:w-24 lg:w-32 h-16 sm:h-24 lg:h-32 opacity-10">
          <div className="w-full h-full rounded-full" style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)'
          }}></div>
        </div>
        <div className="absolute top-20 sm:top-40 right-4 sm:right-20 w-12 sm:w-16 lg:w-24 h-12 sm:h-16 lg:h-24 opacity-10">
          <div className="w-full h-full" style={{
            background: 'linear-gradient(45deg, #16a085, #138f7a)',
            clipPath: 'polygon(50% 0%, 0% 100%, 100% 100%)'
          }}></div>
        </div>
        <div className="absolute bottom-10 sm:bottom-20 left-1/4 w-20 sm:w-32 lg:w-40 h-20 sm:h-32 lg:h-40 opacity-5">
          <div className="w-full h-full rounded-full" style={{
            background: 'radial-gradient(circle, #16a085, transparent)'
          }}></div>
        </div>
        <div className="absolute bottom-20 sm:bottom-40 right-4 sm:right-10 w-16 sm:w-20 lg:w-28 h-16 sm:h-20 lg:h-28 opacity-10">
          <div className="w-full h-full" style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)',
            borderRadius: '0 50% 50% 50%'
          }}></div>
        </div>
      </div>

      {/* Content */}
      <main className="py-6 sm:py-12 px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="max-w-7xl mx-auto">
          {showUploader && (
            <div className="text-center mb-8 sm:mb-12">
              {/* TalentX Logo/Brand */}
              <div className="mb-6 sm:mb-8">
                <div className="inline-flex items-center justify-center w-16 sm:w-20 h-16 sm:h-20 mb-4">
                  <img 
                    src="/talentx-icon.svg" 
                    alt="TalentX Logo" 
                    className="w-full h-full"
                    style={{
                      filter: 'drop-shadow(0 20px 40px rgba(29,205,159,.3))'
                    }}
                  />
                </div>
                <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-bold text-white mb-2">
                  <span style={{ background: 'linear-gradient(135deg, #16a085, #1dd1a1)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    Talent
                  </span>
                  <span className="text-white">X</span>
                </h1>
              </div>
              
              <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-white mb-3">
                Find Your Perfect Job
              </h2>
              <p className="mt-3 max-w-sm sm:max-w-md lg:max-w-3xl mx-auto text-sm sm:text-base md:text-lg lg:text-xl text-gray-300 md:mt-5">
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
            <div className="text-center mb-6 sm:mb-8">
              {/* Small TalentX branding on results page */}
              <div className="flex items-center justify-center mb-4">
                <div className="w-8 h-8 flex items-center justify-center mr-2">
                  <img 
                    src="/talentx-icon.svg" 
                    alt="TalentX Logo" 
                    className="w-full h-full"
                  />
                </div>
                <h1 className="text-xl font-bold text-white">
                  <span style={{ background: 'linear-gradient(135deg, #16a085, #1dd1a1)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    Talent
                  </span>
                  <span className="text-white">X</span>
                </h1>
              </div>
              
              <button
                onClick={handleSearchAgain}
                className="inline-flex items-center px-4 sm:px-6 py-2 sm:py-3 border-0 text-sm sm:text-base font-medium rounded-lg text-white transition-all duration-200 hover:scale-105"
                style={{
                  background: 'linear-gradient(135deg, #16a085, #138f7a)',
                  boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
                }}
              >
                <svg
                  className="mr-2 -ml-1 w-4 sm:w-5 h-4 sm:h-5"
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
            <div className="max-w-2xl mx-auto mb-6 sm:mb-8 px-4 sm:px-0">
              <div className="rounded-lg p-3 sm:p-4" style={{
                background: 'linear-gradient(135deg, #2d1b1b, #1a1a1a)',
                border: '1px solid rgba(220, 38, 38, 0.3)'
              }}>
                <div className="flex">
                  <div className="flex-shrink-0">
                    <svg
                      className="h-4 sm:h-5 w-4 sm:w-5 text-red-400"
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
                  <div className="ml-2 sm:ml-3">
                    <p className="text-xs sm:text-sm text-red-300">{error}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {loading && (
            <div className="max-w-2xl mx-auto mb-6 sm:mb-8 px-4 sm:px-0">
              <div className="rounded-lg p-4 sm:p-6" style={{
                background: 'linear-gradient(135deg, #0a0a0a, #1a1a1a)',
                border: '1px solid rgba(29,205,159,.2)'
              }}>
                <div className="text-center">
                  <LoadingSpinner />
                  <p className="mt-3 sm:mt-4 text-xs sm:text-sm text-gray-300">
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
