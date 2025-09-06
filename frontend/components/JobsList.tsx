'use client';

import { Job } from '../types/job';

interface JobsListProps {
  jobs: Job[];
}

export default function JobsList({ jobs }: JobsListProps) {
  if (jobs.length === 0) {
    return null;
  }

  return (
    <div className="w-full max-w-6xl mx-auto mt-8">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">
          Found {jobs.length} matching job{jobs.length !== 1 ? 's' : ''}
        </h2>
        <p className="text-gray-600 mt-1">
          Jobs tailored to your resume and experience
        </p>
      </div>
      
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {jobs.map((job, index) => (
          <div
            key={index}
            className="bg-white rounded-lg shadow-md border hover:shadow-lg transition-shadow duration-200"
          >
            <div className="p-6">
              <div className="flex items-start justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-900 leading-tight">
                  {job.title}
                </h3>
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {job.source}
                </span>
              </div>
              
              {job.company && (
                <p className="text-sm font-medium text-gray-700 mb-2">
                  {job.company}
                </p>
              )}
              
              {job.location && (
                <p className="text-sm text-gray-500 mb-4 flex items-center">
                  <svg
                    className="w-4 h-4 mr-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                  </svg>
                  {job.location}
                </p>
              )}
              
              {job.description && (
                <p className="text-sm text-gray-600 mb-4 line-clamp-3">
                  {job.description.length > 150
                    ? `${job.description.substring(0, 150)}...`
                    : job.description}
                </p>
              )}
              
              <div className="mt-6">
                <a
                  href={job.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 w-full justify-center"
                >
                  Apply Now
                  <svg
                    className="ml-2 -mr-1 w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                    />
                  </svg>
                </a>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
} 