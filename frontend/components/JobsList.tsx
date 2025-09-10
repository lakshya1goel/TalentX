'use client';

import { useState, useMemo } from 'react';
import { Job, RankedJob } from '../types/job';
import { paginateArray } from '../utils/api';
import Pagination from './Pagination';

interface JobsListProps {
  jobs: Job[];
}

export default function JobsList({ jobs }: JobsListProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const rankedJobs: RankedJob[] = useMemo(() => {
    return jobs.map((job, index) => ({
      job,
      percent_match: 95 - (index * 5), // Mock percentages
      match_reason: `Excellent match for an intern, especially with the candidate's ${job.title.toLowerCase()} experience. Strong foundational skills and high potential make this a great fit.`,
      skills_matched: ['Software Engineering', 'General Development'],
      experience_match: 'Perfect match for an intern'
    }));
  }, [jobs]);

  // Client-side pagination
  const paginatedData = useMemo(() => {
    return paginateArray(rankedJobs, { page: currentPage, pageSize });
  }, [rankedJobs, currentPage, pageSize]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    // Scroll to top when page changes
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(1); // Reset to first page when changing page size
  };

  if (jobs.length === 0) {
    return null;
  }

  return (
    <div className="w-full max-w-6xl mx-auto mt-8">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-white">
          Found {jobs.length} matching job{jobs.length !== 1 ? 's' : ''}
        </h2>
        <p className="text-slate-400 mt-1">
          Jobs tailored to your resume and experience
        </p>
      </div>
      
      <div className="space-y-4">
        {paginatedData.data.map((rankedJob, index) => (
          <JobCard 
            key={`${rankedJob.job.url}-${index}`}
            rankedJob={rankedJob}
            rank={(currentPage - 1) * pageSize + index + 1}
          />
        ))}
      </div>

      {/* Pagination */}
      <Pagination
        currentPage={currentPage}
        totalPages={paginatedData.pagination.totalPages}
        onPageChange={handlePageChange}
        pageSize={pageSize}
        onPageSizeChange={handlePageSizeChange}
        totalItems={jobs.length}
      />
    </div>
  );
}

const JobCard: React.FC<{ rankedJob: RankedJob; rank: number }> = ({ rankedJob, rank }) => {
  const { job, percent_match, match_reason, skills_matched, experience_match } = rankedJob;

  const getPercentageColor = (percentage: number) => {
    if (percentage >= 80) return 'text-green-400 bg-green-900/30 border-green-700';
    if (percentage >= 60) return 'text-yellow-400 bg-yellow-900/30 border-yellow-700';
    return 'text-red-400 bg-red-900/30 border-red-700';
  };

  return (
    <div className="bg-slate-800 rounded-lg shadow-lg border border-slate-700 hover:shadow-xl hover:border-slate-600 transition-all duration-200 p-6">
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center gap-3">
          <span className="bg-blue-900/50 text-blue-300 border border-blue-700 px-3 py-1 rounded-full text-sm font-medium">
            #{rank}
          </span>
          <h3 className="text-xl font-semibold text-white">
            {job.title}
          </h3>
        </div>
        
        <div className={`px-3 py-1 rounded-full text-sm font-medium border ${getPercentageColor(percent_match)}`}>
          {percent_match.toFixed(1)}%
        </div>
      </div>

      <div className="mb-4">
        {job.company && (
          <p className="text-slate-300 font-medium mb-1">{job.company}</p>
        )}
        {job.location && (
          <p className="text-slate-400 text-sm">{job.location}</p>
        )}
      </div>

      {/* Why this matches */}
      <div className="mb-4">
        <h4 className="text-sm font-medium text-slate-300 mb-2">Why this matches:</h4>
        <p className="text-sm text-slate-400">{match_reason}</p>
      </div>

      {/* Experience Match */}
      <div className="mb-4">
        <h4 className="text-sm font-medium text-slate-300 mb-1">Experience Match:</h4>
        <p className="text-sm text-slate-400">{experience_match}</p>
      </div>

      {/* Skills Matched */}
      {skills_matched.length > 0 && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-slate-300 mb-2">Skills Matched:</h4>
          <div className="flex flex-wrap gap-2">
            {skills_matched.map((skill, idx) => (
              <span
                key={idx}
                className="bg-blue-900/50 text-blue-300 px-3 py-1 rounded-full text-xs border border-blue-700"
              >
                {skill}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Apply Button */}
      <div className="pt-4 border-t border-slate-700">
        <a
          href={job.url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-4 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 focus:ring-offset-slate-800 transition-colors duration-200 shadow-lg hover:shadow-xl"
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
  );
}; 