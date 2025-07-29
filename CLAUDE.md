# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Kubernetes Job Runner service that creates Kubernetes jobs from shell commands passed via HTTP requests.

## Development Commands

### Build
```bash
go build -o k8s-job-runner main.go
```

### Run
```bash
go run main.go
```

### Test
```bash
go test ./...
```

### Format
```bash
go fmt ./...
```

## Architecture

The project is currently a simple HTTP server with:
- **Main entry point**: `main.go` - HTTP server listening on port 8084
- **Endpoints**:
  - `/` - Returns "Hello, world!"
  - `/health` - Health check endpoint returning "ok"

**Note**: The Kubernetes job creation functionality described in the README is not yet implemented. The current codebase only contains a basic HTTP server skeleton.

## Next Steps for Implementation

When implementing the Kubernetes job runner functionality:
1. Add Kubernetes client-go dependency
2. Implement job creation endpoint that accepts shell commands
3. Add authentication/authorization for security
4. Implement job status tracking and retrieval
5. Add proper error handling and logging

## Commit Guidelines

- Use conventional commits and keep the messages short with very limited "bodies"