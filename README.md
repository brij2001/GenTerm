# GenTerm

GenTerm is a terminal-based application that uses Generative AI to perform Retrieval-Augmented Generation (RAG) on your files. Upload documents, images, and other file types, then ask questions about them in a familiar terminal interface.

![GenTerm](https://github.com/genterm/genterm/raw/main/frontend/public/screenshot.png)

## Features

- 🖥️ Terminal-like interface for a familiar user experience
- 📄 Process and analyze text documents (PDF, TXT, DOCX, etc.)
- 🖼️ Image understanding capabilities
- 💬 Natural language queries against your uploaded content
- 🔍 RAG (Retrieval-Augmented Generation) for accurate answers
- 📜 Command history navigation (up/down arrows)
- 🔄 Session persistence

## Tech Stack

### Frontend
- React.js
- Axios for API requests
- PDF.js for PDF processing
- FontAwesome for icons

### Backend
- Go
- Standard library HTTP server
- Environment configuration with godotenv

## Installation

### Prerequisites
- Node.js (v14+)
- Go (v1.21+)
- npm or yarn

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/genterm.git
   cd genterm
   ```

2. Set up the frontend:
   ```bash
   cd frontend
   npm install
   # Configure your environment variables
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Set up the backend:
   ```bash
   cd backend
   # Configure your environment variables
   cp .env.example .env
   # Edit .env with your configuration
   ```

## Running the Application

1. Start the backend server:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. Start the frontend development server:
   ```bash
   cd frontend
   npm start
   ```

3. Open your browser and navigate to:
   ```
   http://localhost:3000
   ```

## Usage

1. Upload files using the panel on the right side of the interface
2. Type questions in the terminal about your uploaded content
3. View responses directly in the terminal interface

### Terminal Commands

- `help` - Display available commands
- `clear` - Clear the terminal screen
- `files` - List all uploaded files

## Development

### Frontend Structure
```
frontend/
├── public/
├── src/
│   ├── components/
│   │   └── FileUploader.js
│   ├── services/
│   ├── utils/
│   ├── App.js
│   ├── App.css
│   └── index.js
└── package.json
```

### Backend Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   ├── config/
│   └── session/
├── .env
└── go.mod
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all the open-source libraries that made this project possible
- Inspired by classic terminal interfaces and modern AI capabilities
