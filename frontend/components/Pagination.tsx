import React from 'react';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
    pageSize: number;
    onPageSizeChange: (size: number) => void;
    totalItems: number;
}

const Pagination: React.FC<PaginationProps> = ({
    currentPage,
    totalPages,
    onPageChange,
    pageSize,
    onPageSizeChange,
    totalItems,
}) => {
  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

    const getPageNumbers = () => {
        const pages = [];
        const maxVisible = 5;
        
        let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
        const end = Math.min(totalPages, start + maxVisible - 1);
        
        if (end - start + 1 < maxVisible) {
            start = Math.max(1, end - maxVisible + 1);
        }
        
        for (let i = start; i <= end; i++) {
            pages.push(i);
        }
        
        return pages;
    };

    if (totalItems === 0) {
        return null;
    }

    return (
        <div 
        className="flex flex-col sm:flex-row items-center justify-between gap-4 p-6 rounded-xl mt-8 backdrop-blur-sm"
        style={{
            background: 'linear-gradient(135deg, rgba(10,10,10,0.8), rgba(26,26,26,0.8))',
            border: '1px solid rgba(29,205,159,.2)'
        }}
        >
        {/* Results info */}
        <div className="text-sm text-gray-300 flex items-center">
            <svg className="w-4 h-4 mr-2 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            Showing {startItem}-{endItem} of {totalItems} jobs
        </div>

        {/* Page size selector */}
        <div className="flex items-center gap-3">
            <span className="text-sm text-gray-300">Show:</span>
            <select
            value={pageSize}
            onChange={(e) => onPageSizeChange(Number(e.target.value))}
            className="rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-green-500 transition-all duration-200"
            style={{
                background: 'linear-gradient(135deg, rgba(22,160,133,0.2), rgba(19,143,122,0.2))',
                border: '1px solid rgba(29,205,159,.3)'
            }}
            >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={15}>15</option>
            <option value={20}>20</option>
            </select>
            <span className="text-sm text-gray-300">per page</span>
        </div>

        {/* Pagination controls */}
        <div className="flex items-center gap-2">
            {/* Previous button */}
            <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className="px-4 py-2 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 hover:scale-105 flex items-center"
            style={{
                background: currentPage === 1 
                ? 'rgba(29,205,159,.1)' 
                : 'linear-gradient(135deg, rgba(22,160,133,0.3), rgba(19,143,122,0.3))',
                border: '1px solid rgba(29,205,159,.3)'
            }}
            >
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Previous
            </button>

            {/* Page numbers */}
            <div className="flex items-center gap-1">
            {getPageNumbers().map((page) => (
                <button
                key={page}
                onClick={() => onPageChange(page)}
                className="w-10 h-10 rounded-lg transition-all duration-200 text-white font-medium hover:scale-105"
                style={{
                    background: currentPage === page
                    ? 'linear-gradient(135deg, #16a085, #138f7a)'
                    : 'rgba(29,205,159,.1)',
                    border: '1px solid rgba(29,205,159,.3)'
                }}
                >
                {page}
                </button>
            ))}
            </div>

            {/* Next button */}
            <button
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage === totalPages}
            className="px-4 py-2 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 hover:scale-105 flex items-center"
            style={{
                background: currentPage === totalPages 
                ? 'rgba(29,205,159,.1)' 
                : 'linear-gradient(135deg, rgba(22,160,133,0.3), rgba(19,143,122,0.3))',
                border: '1px solid rgba(29,205,159,.3)'
            }}
            >
            Next
            <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
            </button>
        </div>
        </div>
    );
};

export default Pagination; 