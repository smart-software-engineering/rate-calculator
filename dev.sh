#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Rate Calculator Development Environment${NC}"

# Check and install Air if needed
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}Air is not installed. Installing now...${NC}"
    go install github.com/cosmtrek/air@latest

    if ! command -v air &> /dev/null; then
        echo -e "${RED}Failed to install Air. Please install manually with: go install github.com/cosmtrek/air@latest${NC}"
        exit 1
    fi

    echo -e "${GREEN}Air installed successfully!${NC}"
fi

# Check if npm is available and frontend exists
if [ ! -d "frontend" ]; then
    echo -e "${RED}Frontend directory not found!${NC}"
    exit 1
fi

# Check if frontend dependencies are installed
if [ ! -d "frontend/node_modules" ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    cd frontend && npm install && cd ..
fi

mkdir -p tmp

# Set development environment variables
export AIR_SESSION_FLAG="-dev=true"
export GO_ENV="development"

# Function to cleanup background processes on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down development servers...${NC}"
    kill $(jobs -p) 2>/dev/null
    wait
    echo -e "${GREEN}Development servers stopped.${NC}"
}

# Set trap to cleanup on script exit
trap cleanup EXIT

echo -e "${GREEN}Starting Go backend with Air (port 8080)...${NC}"
air &

echo -e "${GREEN}Starting Svelte frontend with Vite (port 5173)...${NC}"
cd frontend && npm run dev &

echo -e "${BLUE}Development servers running:${NC}"
echo -e "  ${GREEN}• Go Backend:${NC}  http://localhost:8080"
echo -e "  ${GREEN}• Svelte SPA:${NC}  http://localhost:5173"
echo -e "${YELLOW}Press Ctrl+C to stop all servers${NC}"

# Wait for all background jobs
wait
