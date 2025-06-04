#!/bin/bash

# Build script for Rate Calculator container image (using Podman)

echo "🚀 Building Rate Calculator container image with Podman..."
echo "This will create a minimal production image with:"
echo "  - Svelte SPA built and embedded"
echo "  - Go backend compiled"
echo "  - Alpine Linux runtime (~15MB base)"
echo "  - PORT environment variable support (for Fly.io, etc.)"
echo ""

# Build the container image using Podman
podman build -t rate-calculator .

if [ $? -eq 0 ]; then
    echo "✅ Container image built successfully!"
    echo ""
    echo "To run the container (production mode requires session key):"
    echo "  podman run -p 8080:8080 -e COOKIE_SESSION_KEY=\"your-secret-key-32-chars\" rate-calculator"
    echo ""
    echo "To run with custom port (e.g., for Fly.io):"
    echo "  podman run -p 3000:3000 -e PORT=3000 -e COOKIE_SESSION_KEY=\"your-secret-key-32-chars\" rate-calculator"
    echo ""
    echo "To run in background:"
    echo "  podman run -d -p 8080:8080 -e COOKIE_SESSION_KEY=\"your-secret-key-32-chars\" --name rate-calc rate-calculator"
    echo ""
    echo "To generate a secure session key:"
    echo "  openssl rand -base64 32"
    echo ""
    echo "Image size:"
    podman images rate-calculator --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
else
    echo "❌ Container build failed"
    exit 1
fi 