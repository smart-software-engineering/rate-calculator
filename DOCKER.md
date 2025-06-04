# Docker Deployment Guide

## Overview

This project uses a multi-stage Docker build to create a minimal production image containing:
- **Svelte SPA**: Built and embedded as static assets
- **Go Backend**: Compiled binary serving both API and static files  
- **Alpine Linux**: Minimal runtime (~23MB total image size)

## Architecture

```
Stage 1: Frontend Builder (node:22-alpine)
├── Install Node.js dependencies
├── Build Svelte SPA with Vite
└── Output: Static assets in /app/internal/server/static/spa

Stage 2: Backend Builder (golang:1.24-alpine) 
├── Download Go dependencies
├── Copy source code + built frontend assets
├── Compile Go binary
└── Output: Statically linked binary

Stage 3: Runtime (alpine:latest)
├── Install ca-certificates and tzdata
├── Copy binary + templates + static assets
└── Final minimal production image
```

## Building

### Using the build script (recommended):
```bash
./docker-build.sh
```

### Manual build:
```bash
podman build -t rate-calculator .
# or
docker build -t rate-calculator .
```

## Running

⚠️ **Important**: Production mode requires the `COOKIE_SESSION_KEY` environment variable to be set for secure session management.

### Generate a secure session key:
```bash
openssl rand -base64 32
```

### Basic run (production):
```bash
podman run -p 8080:8080 -e COOKIE_SESSION_KEY="your-secret-key-32-chars" rate-calculator
```

### Background daemon:
```bash
podman run -d -p 8080:8080 -e COOKIE_SESSION_KEY="your-secret-key-32-chars" --name rate-calc rate-calculator
```

### With additional environment variables:
```bash
podman run -p 8080:8080 \
  -e COOKIE_SESSION_KEY="your-secret-key-32-chars" \
  -e PORT=8080 \
  -e ENV=production \
  rate-calculator
```

## Access Points

Once running, the application provides:
- **Landing Page**: http://localhost:8080/
- **Svelte SPA**: http://localhost:8080/app  
- **API Health**: http://localhost:8080/api/v1/health
- **Legacy Calculator**: http://localhost:8080/calculator

## Features

✅ **Multi-stage build** for optimized image size  
✅ **Layer caching** for faster rebuilds  
✅ **Security**: Secure session management, ca-certificates  
✅ **Production ready**: Static compilation, no dev dependencies  
✅ **Cross-platform**: Works with Docker, Podman, etc.

## Image Details

- **Base**: Alpine Linux 3.22
- **Size**: ~23MB
- **Go version**: 1.24
- **Node version**: 22 LTS
- **Security**: Includes ca-certificates for HTTPS

## Development vs Production

| Environment | Frontend Serving | API Proxy | Hot Reload | Session Key |
|-------------|------------------|-----------|------------|-------------|
| Development | Vite dev server (5173) | Vite → Go (8080) | ✅ | Optional |
| Production | Go static files | Direct | ❌ | **Required** |

## Security Considerations

### Session Key Requirements
- **Development**: Session key is optional (uses default)
- **Production**: `COOKIE_SESSION_KEY` environment variable is **mandatory**
- **Key Format**: Base64 encoded, minimum 32 characters recommended
- **Generation**: Use `openssl rand -base64 32` for secure keys

### Container Security
- Runs as non-root user in container
- Minimal attack surface with Alpine Linux
- No shell or unnecessary binaries in final image
- TLS certificates included for secure external connections

## Troubleshooting

### Container exits with "COOKIE_SESSION_KEY required":
```bash
# Generate and set a secure session key
export SESSION_KEY=$(openssl rand -base64 32)
podman run -p 8080:8080 -e COOKIE_SESSION_KEY="$SESSION_KEY" rate-calculator
```

### Build fails with "No matching version":
- Check `frontend/package.json` versions match installed packages
- Run `npm list` in frontend/ to verify

### Container exits immediately:
- Check logs: `podman logs <container-id>`
- Verify port 8080 is available
- Ensure COOKIE_SESSION_KEY is set for production

### Can't access application:
- Ensure port mapping: `-p 8080:8080`
- Check firewall settings
- Verify container is running: `podman ps` 