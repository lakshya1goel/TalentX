export default function Header() {
  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-6">
          <div className="flex items-center">
            <h1 className="text-2xl font-bold text-gray-900">
              Job Assistant
            </h1>
          </div>
          <div className="text-sm text-gray-500">
            Upload your resume to find matching jobs
          </div>
        </div>
      </div>
    </header>
  );
} 