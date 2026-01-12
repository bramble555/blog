# Frontend Developer Guide

## Setup & Run
1. Navigate to frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies (if not already done):
   ```bash
   npm install
   ```
3. Start development server:
   ```bash
   npm run dev
   ```
   Access at `http://localhost:5173`.

## Backend Integration
- The frontend expects the backend to be running at `http://localhost:8080`.
- API calls are proxied to `/api` or made directly depending on config. Current setup points to `http://localhost:8080/api`.

## Implemented Features
- **Article Management**: List, Create, Edit, Delete.
- **Visuals**: VS Code Dark Theme using Tailwind CSS.
- **Routing**: Vue Router setup.

## Missing / Mocked Features
- **Authentication**: Currently bypassed/not implemented in frontend UI.
- **Image Uploads**: Banner upload not yet integrated in Article Editor.
- **Complex Rich Text**: Using basic Markdown textarea.
