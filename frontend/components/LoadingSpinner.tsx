export default function LoadingSpinner() {
  return (
    <div className="flex items-center justify-center p-8">
      <div className="animate-spin rounded-full h-12 w-12 border-2 border-slate-600 border-t-blue-500 border-r-blue-500"></div>
    </div>
  );
} 