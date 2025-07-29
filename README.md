# Kubernetes Job Runner

A simple HTTP service that creates Kubernetes jobs from shell commands.

## Usage

Run the service:
```bash
go run .
```

Create a job:
```bash
curl -X POST http://localhost:8084/jobs \
  -H "Content-Type: application/json" \
  -d '{"name": "my-job", "command": ["echo", "hello world"]}'
```

## Endpoints

- `POST /jobs` - Create a new job
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## Configuration

The service uses your local kubeconfig or in-cluster config when deployed. Jobs are created in the `default` namespace.
