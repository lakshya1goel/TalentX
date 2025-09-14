'use client';

import { useState, useMemo } from 'react';
import { Job, RankedJob } from '../types/job';
import { paginateArray } from '../utils/api';
import Pagination from './Pagination';

interface JobsListProps {
  jobs: RankedJob[];
}

const CircularProgress: React.FC<{ percentage: number; size?: number }> = ({ 
  percentage, 
  size = 100 
}) => {
  const radius = (size - 16) / 2;
  const circumference = radius * 2 * Math.PI;
  const strokeDasharray = circumference;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  const getColor = (percent: number) => {
    if (percent >= 80) return '#10b981';
    if (percent >= 60) return '#f59e0b';
    return '#ef4444';
  };

  const getGradientId = (percent: number) => {
    if (percent >= 80) return 'greenGradient';
    if (percent >= 60) return 'amberGradient';
    return 'redGradient';
  };

  return (
    <div className="relative inline-flex items-center justify-center">
      <svg
        width={size}
        height={size}
        className="transform -rotate-90"
      >
        <defs>
          <linearGradient id="greenGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#34d399" stopOpacity="1" />
            <stop offset="50%" stopColor="#10b981" stopOpacity="1" />
            <stop offset="100%" stopColor="#059669" stopOpacity="1" />
          </linearGradient>
          <linearGradient id="amberGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#fbbf24" stopOpacity="1" />
            <stop offset="50%" stopColor="#f59e0b" stopOpacity="1" />
            <stop offset="100%" stopColor="#d97706" stopOpacity="1" />
          </linearGradient>
          <linearGradient id="redGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#f87171" stopOpacity="1" />
            <stop offset="50%" stopColor="#ef4444" stopOpacity="1" />
            <stop offset="100%" stopColor="#dc2626" stopOpacity="1" />
          </linearGradient>
        </defs>
        
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          stroke="rgba(29,205,159,0.1)"
          strokeWidth="8"
          fill="none"
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          stroke={`url(#${getGradientId(percentage)})`}
          strokeWidth="8"
          fill="none"
          strokeLinecap="round"
          strokeDasharray={strokeDasharray}
          strokeDashoffset={strokeDashoffset}
          className="transition-all duration-1000 ease-out"
          style={{
            strokeLinecap: 'round'
          }}
        />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center">
        <span className="text-sm font-bold text-white">
          {Math.round(percentage)}%
        </span>
      </div>
    </div>
  );
};

export default function JobsList({ jobs }: JobsListProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const rankedJobs: RankedJob[] = useMemo(() => {
    return jobs;
  }, [jobs]);

  const paginatedData = useMemo(() => {
    return paginateArray(rankedJobs, { page: currentPage, pageSize });
  }, [rankedJobs, currentPage, pageSize]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(1);
  };

  if (jobs?.length === 0) {
    return null;
  }

  return (
    <div className="w-full max-w-6xl mx-auto mt-6 sm:mt-8 px-4 sm:px-0">
      <div className="mb-6 sm:mb-8 text-center">
        <h2 className="text-2xl sm:text-3xl font-bold text-white mb-2">
          Found {jobs?.length} matching job{jobs?.length !== 1 ? 's' : ''}
        </h2>
        <p className="text-sm sm:text-base text-gray-400">
          Jobs tailored to your resume and experience, ranked by relevance
        </p>
      </div>
      
      <div className="space-y-4 sm:space-y-6">
        {paginatedData.data.map((rankedJob, index) => (
          <JobCard 
            key={`${rankedJob.job.url}-${index}`}
            rankedJob={rankedJob}
            rank={(currentPage - 1) * pageSize + index + 1}
          />
        ))}
      </div>

      <Pagination
        currentPage={currentPage}
        totalPages={paginatedData.pagination.totalPages}
        onPageChange={handlePageChange}
        pageSize={pageSize}
        onPageSizeChange={handlePageSizeChange}
        totalItems={jobs?.length || 0 }
      />
    </div>
  );
}

const JobCard: React.FC<{ rankedJob: RankedJob; rank: number }> = ({ rankedJob, rank }) => {
  const { job, percent_match, match_reason, skills_matched, experience_match } = rankedJob;
  const [showMore, setShowMore] = useState(false);

  const shouldTruncate = match_reason?.length > 100;
  const displayedReason = showMore || !shouldTruncate
    ? match_reason
    : match_reason.slice(0, 100) + (shouldTruncate ? '...' : '');

  return (
    <div 
      className="rounded-xl p-4 sm:p-6 transition-all duration-200 hover:scale-[1.01] backdrop-blur-sm"
      style={{
        background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
        border: '1px solid rgba(29,205,159,.2)'
      }}
    >
      <div className="flex flex-col sm:flex-row items-start justify-between mb-4 sm:mb-6 gap-4">
        <div className="flex items-start gap-3 sm:gap-4 flex-1 w-full">
          <div className="text-green-400 font-bold text-xl sm:text-2xl flex-shrink-0">
            #{rank}
          </div>
          <div className="flex-1 min-w-0">
            <h3 className="text-lg sm:text-xl font-semibold text-white mb-1">{job.title}</h3>
            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-4 text-xs sm:text-sm">
              {job.company && (
                <span className="text-gray-300 font-medium">{job.company}</span>
              )}
              {job.location && (
                <span className="text-gray-400 flex items-center">
                  <svg className="w-3 sm:w-4 h-3 sm:h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {job.location}
                </span>
              )}
            </div>
          </div>
        </div>
        
        <div className="flex flex-col items-center gap-2 flex-shrink-0">
          <CircularProgress percentage={percent_match} size={80} />
          <span className="text-xs text-gray-400 font-medium">Match</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 mb-4 sm:mb-6">
        <div className="space-y-2 sm:space-y-3">
          <h4 className="text-xs sm:text-sm font-semibold text-green-400 flex items-center">
            <svg className="w-3 sm:w-4 h-3 sm:h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Why this matches
          </h4>
          <p className="text-xs sm:text-sm text-gray-300 leading-relaxed">
            {displayedReason}
          </p>
          {shouldTruncate && (
            <button
              className="text-xs sm:text-sm text-gray-400"
              onClick={() => setShowMore((prev) => !prev)}
            >
              {showMore ? 'Show less' : 'Show more'}
            </button>
          )}
        </div>

        <div className="space-y-2 sm:space-y-3">
          <h4 className="text-xs sm:text-sm font-semibold text-green-400 flex items-center">
            <svg className="w-3 sm:w-4 h-3 sm:h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Experience Match
          </h4>
          <p className="text-xs sm:text-sm text-gray-300">{experience_match}</p>
        </div>
      </div>

      {skills_matched?.length > 0 && (
        <div className="mb-4 sm:mb-6">
          <h4 className="text-xs sm:text-sm font-semibold text-green-400 mb-2 sm:mb-3 flex items-center">
            <svg className="w-3 sm:w-4 h-3 sm:h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
            Skills Matched
          </h4>
          <div className="flex flex-wrap gap-1 sm:gap-2">
            {skills_matched.map((skill, idx) => (
              <span
                key={idx}
                className="px-2 sm:px-3 py-1 rounded-full text-xs font-medium text-white"
                style={{
                  background: 'linear-gradient(135deg, rgba(22,160,133,0.3), rgba(19,143,122,0.3))',
                  border: '1px solid rgba(29,205,159,.3)'
                }}
              >
                {skill}
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="flex flex-col sm:flex-row justify-between items-center pt-3 sm:pt-4 gap-2 sm:gap-0" style={{ borderTop: '1px solid rgba(29,205,159,.2)' }}>
        <div className="text-xs text-gray-500 text-center sm:text-left">
          Click to apply on external site
        </div>
        <a
          href={job.url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-4 sm:px-6 py-2 sm:py-3 text-xs sm:text-sm font-semibold rounded-lg text-white transition-all duration-200 hover:scale-105"
          style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)',
            boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
          }}
        >
          Apply Now
          <svg
            className="ml-2 w-3 sm:w-4 h-3 sm:h-4"
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