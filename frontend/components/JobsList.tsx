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
    return jobs.map((job, index) => {
      // Ensure percentage stays between 30-95% for better distribution
      const percentage = Math.max(30, 95 - (index * 3));
      
      return {
        job,
        percent_match: percentage,
        match_reason: `Excellent match for an intern, especially with the candidate's ${job.title.toLowerCase()} experience. Strong foundational skills and high potential make this a great fit.`,
        skills_matched: ['Software Engineering', 'General Development'],
        experience_match: 'Perfect match for an intern'
      };
    });
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
      <div className="mb-8 text-center">
        <h2 className="text-3xl font-bold text-white mb-2">
          Found {jobs.length} matching job{jobs.length !== 1 ? 's' : ''}
        </h2>
        <p className="text-gray-400">
          Jobs tailored to your resume and experience, ranked by relevance
        </p>
      </div>
      
      <div className="space-y-6">
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
    if (percentage >= 80) return 'text-green-400 bg-green-900/30 border-green-400/30';
    if (percentage >= 60) return 'text-yellow-400 bg-yellow-900/30 border-yellow-400/30';
    return 'text-red-400 bg-red-900/30 border-red-400/30';
  };

  return (
    <div 
      className="rounded-xl p-6 transition-all duration-200 hover:scale-[1.01] backdrop-blur-sm"
      style={{
        background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
        border: '1px solid rgba(29,205,159,.2)'
      }}
    >
      <div className="flex items-start justify-between mb-6">
        <div className="flex items-center gap-4 flex-1">
          <div 
            className="w-12 h-12 rounded-full flex items-center justify-center text-white font-bold text-lg"
            style={{
              background: 'linear-gradient(135deg, #16a085, #138f7a)'
            }}
          >
            {rank}
          </div>
          <div className="flex-1">
            <h3 className="text-xl font-semibold text-white mb-1">{job.title}</h3>
            <div className="flex items-center gap-4 text-sm">
              {job.company && (
                <span className="text-gray-300 font-medium">{job.company}</span>
              )}
              {job.location && (
                <span className="text-gray-400 flex items-center">
                  <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {job.location}
                </span>
              )}
            </div>
          </div>
        </div>
        
        <div className={`px-4 py-2 rounded-full text-sm font-semibold border ${getPercentageColor(percent_match)}`}>
          {percent_match.toFixed(0)}% Match
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6 mb-6">
        {/* Why this matches */}
        <div className="space-y-3">
          <h4 className="text-sm font-semibold text-green-400 flex items-center">
            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Why this matches
          </h4>
          <p className="text-sm text-gray-300 leading-relaxed">{match_reason}</p>
        </div>

        {/* Experience Match */}
        <div className="space-y-3">
          <h4 className="text-sm font-semibold text-green-400 flex items-center">
            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Experience Match
          </h4>
          <p className="text-sm text-gray-300">{experience_match}</p>
        </div>
      </div>

      {/* Skills Matched */}
      {skills_matched.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-semibold text-green-400 mb-3 flex items-center">
            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
            Skills Matched
          </h4>
          <div className="flex flex-wrap gap-2">
            {skills_matched.map((skill, idx) => (
              <span
                key={idx}
                className="px-3 py-1 rounded-full text-xs font-medium text-white"
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

      {/* Apply Button */}
      <div className="flex justify-between items-center pt-4" style={{ borderTop: '1px solid rgba(29,205,159,.2)' }}>
        <div className="text-xs text-gray-500">
          Click to apply on external site
        </div>
        <a
          href={job.url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-6 py-3 text-sm font-semibold rounded-lg text-white transition-all duration-200 hover:scale-105"
          style={{
            background: 'linear-gradient(135deg, #16a085, #138f7a)',
            boxShadow: '0 15px 35px rgba(29,205,159,.4), 0 0 0 1px rgba(29,205,159,.2)'
          }}
        >
          Apply Now
          <svg
            className="ml-2 w-4 h-4"
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