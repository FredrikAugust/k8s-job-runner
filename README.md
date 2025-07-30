# Kubernetes Job Runner

[![CI](https://github.com/fredrikaugust/k8s-job-runner/actions/workflows/ci.yml/badge.svg)](https://github.com/fredrikaugust/k8s-job-runner/actions/workflows/ci.yml)

> [!WARNING]
> This is a learning project and should not be used in production. It has no authentication, authorization, or input validation, allowing anyone to execute arbitrary commands in your Kubernetes cluster.

A simple HTTP service that creates Kubernetes jobs from shell commands. Includes Prometheus metrics.

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

With custom image and namespace:
```bash
curl -X POST http://localhost:8084/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-job",
    "command": ["python", "-c", "print(\"Hello\")"],
    "image": "python:3.9-slim",
    "k8s_namespace": "custom-namespace"
  }'
```

## Endpoints

- `POST /jobs` - Create a new job
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## Metrics

- `http_requests_total` - Total HTTP requests by path and method
- `running_jobs` - Current number of running jobs (gauge)

## Configuration

The service uses your local kubeconfig or in-cluster config when deployed. Defaults: `alpine` image, `default` namespace.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

> [!NOTE]
> AI was used in this project for writing commit messages, README documentation, and tests, along with continuous code review. This was done to maximize learning and focus on the interesting architectural and implementation challenges.
