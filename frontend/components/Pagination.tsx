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
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 p-6 bg-slate-800 rounded-lg border border-slate-700 mt-6">
        {/* Results info */}
        <div className="text-sm text-slate-300">
            Showing {startItem}-{endItem} of {totalItems} jobs
        </div>

        {/* Page size selector */}
        <div className="flex items-center gap-2">
            <span className="text-sm text-slate-300">Show:</span>
            <select
                value={pageSize}
                onChange={(e) => onPageSizeChange(Number(e.target.value))}
                className="bg-slate-700 border border-slate-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                <option value={5}>5</option>
                <option value={10}>10</option>
                <option value={15}>15</option>
                <option value={20}>20</option>
            </select>
            <span className="text-sm text-slate-300">per page</span>
        </div>

        {/* Pagination controls */}
        <div className="flex items-center gap-1">
            {/* Previous button */}
            <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className="px-3 py-1 border border-slate-600 bg-slate-700 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-600 transition-colors"
            >
            Previous
            </button>

            {/* Page numbers */}
            {getPageNumbers().map((page) => (
            <button
                key={page}
                onClick={() => onPageChange(page)}
                className={`px-3 py-1 border rounded transition-colors ${
                currentPage === page
                    ? 'bg-blue-600 text-white border-blue-600'
                    : 'border-slate-600 bg-slate-700 text-white hover:bg-slate-600'
                }`}
            >
                {page}
            </button>
            ))}

            {/* Next button */}
            <button
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage === totalPages}
            className="px-3 py-1 border border-slate-600 bg-slate-700 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-600 transition-colors"
            >
            Next
            </button>
        </div>
        </div>
    );
};

export default Pagination; 