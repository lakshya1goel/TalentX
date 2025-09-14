# TalentX - AI-Powered Job Assistance Platform

TalentX is an intelligent job search platform that analyzes your resume using AI and finds the most relevant job opportunities tailored to your skills, experience, and career goals. Built with a modern tech stack, it leverages Google's Gemini AI to provide personalized job recommendations from multiple job boards.

## 🚀 Features

- **AI-Powered Resume Analysis**: Upload your resume and get intelligent analysis of your skills, experience level, and career trajectory
- **Smart Job Matching**: Uses advanced AI to find jobs that match your profile across multiple job boards
- **Dual Search Strategy**: Combines JSearch API and LinkUp structured search for comprehensive job coverage
- **Intelligent Ranking**: Re-ranks job results based on relevance to your specific resume
- **Location Preferences**: Supports flexible location preferences (remote, hybrid, on-site)
- **Modern UI**: Clean, responsive interface with real-time loading states
- **Experience Level Detection**: Automatically categorizes your experience level (Fresher, Junior, Mid-level, Senior, Lead/Principal)

## 📸 Screenshots

<!-- Add screenshots here -->
<img width="1857" height="960" alt="Pasted image (5)" src="https://github.com/user-attachments/assets/b7cde53e-64d3-49f7-9df5-6e71cb06e0cb" />


### Job Results
<!-- Screenshot of the main upload interface -->
<img width="1857" height="960" alt="Pasted image (4)" src="https://github.com/user-attachments/assets/877dcd3b-d59a-4745-8bf7-5b57a53f06e5" />

## 🛠️ Tech Stack

### Backend
- **Go 1.24.4** - High-performance backend server
- **Gin Framework** - Fast HTTP web framework
- **Google Gemini AI** - Advanced AI for resume analysis and job matching
- **JSearch API** - Comprehensive job listings from major job boards
- **LinkUp API** - Structured job search capabilities

### Frontend
- **Next.js 15.5.2** - React-based full-stack framework
- **React 19.1.0** - Modern React with latest features
- **TypeScript** - Type-safe development
- **Tailwind CSS 4** - Utility-first CSS framework
- **Turbopack** - Ultra-fast bundler for development

## 🏗️ Architecture

```
job-assistance/
├── backend/                 # Go backend server
│   ├── cmd/                # Application entry point
│   ├── config/             # Configuration management
│   ├── internal/           # Internal packages
│   │   ├── ai/            # AI client and job search logic
│   │   ├── api/           # HTTP handlers and routes
│   │   └── dtos/          # Data transfer objects
│   └── go.mod             # Go dependencies
└── frontend/               # Next.js frontend
    ├── app/               # Next.js app directory
    ├── components/        # Reusable React components
    ├── types/             # TypeScript type definitions
    └── utils/             # Utility functions
```

## 🚦 Getting Started

### Prerequisites

- **Go 1.24.4+**
- **Node.js 18+**
- **Google Gemini API Key** - Get from [Google AI Studio](https://aistudio.google.com/)
- **JSearch API Key** - Get from [RapidAPI JSearch](https://rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch)
- **LinkUp API Key** - Get from [LinkUp API](https://www.linkup.com/) (optional)

### Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/lakshya1goel/TalentX.git
cd job-assistance
```

2. Set up environment variables:

Create a `.env` file in the **backend** directory:
```env
GEMINI_API_KEY=your_gemini_api_key_here
ALLOWED_ORIGINS=add_rquired_origins
JSEARCH_API_KEY=your_jsearch_rapidapi_key_here
LINKUP_API_KEY=your_linkup_api_key_here
```

Create a `.env.local` file in the **frontend** directory:
```env
NEXT_PUBLIC_API_URL=backend_base_url
```

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Run the backend server:
```bash
go run ./cmd/
```

The backend server will start on `http://localhost:8084`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 👨‍💻 Author

**Lakshya Goel**
- GitHub: [@lakshya1goel](https://github.com/lakshya1goel)

## 🙏 Acknowledgments

- Google Gemini AI for advanced language processing
- JSearch API for comprehensive job data
- LinkUp for structured job search capabilities
- The Go and React communities for excellent frameworks

---

**Made with ❤️ by Lakshya Goel** 
