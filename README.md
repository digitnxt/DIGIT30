# DIGIT 3.0

DIGIT 3.0 is a microservices-based platform built with Go, Gin, and several modern tools to demonstrate a scalable, observable, and well-orchestrated distributed system. It leverages Docker Compose for deployment, enabling rapid development, testing, and production readiness.

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Services Overview](#services-overview)
- [APIs](#apis)
- [Tools & Technologies](#tools--technologies)
- [Deployment](#deployment)
- [Configuration](#configuration)
- [Development Guidelines](#development-guidelines)
- [Security](#security)
- [Versioning & Dependencies](#versioning--dependencies)

## Architecture Overview

DIGIT 3.0 follows a microservices architecture where individual components handle specific responsibilities. The system integrates service discovery, API gateway management, tracing, logging, and monitoring to ensure robust communication and observability across services.

## Services Overview

The platform consists of multiple containerized services managed through Docker Compose. Key services include:

### Postgres (Registry)
- PostgreSQL database (version 14) used as a registry
- Includes health checks to ensure database readiness

### Consul
- Service discovery tool (HashiCorp Consul 1.15)
- Facilitates service registration and health monitoring

### Kafka & Kafka Exporter
- Kafka (using Bitnami's image) serves as the messaging backbone
- Kafka Exporter provides metrics for monitoring

### Kong API Gateway & Kong DB
- **Kong**: Acts as the API gateway with custom plugins (such as api-costing) for managing, monitoring, and securing API traffic
- **Kong DB**: A dedicated PostgreSQL instance (version 13) for Kong's internal database

### Kong Synchronizer
- Service built from a Node.js application
- Ensures data consistency across Kong and other system components

### Keycloak
- Provides authentication and authorization
- Enables secure access to the platform

### Redis & Redis Exporter
- Redis is used for caching
- Redis Exporter collects metrics from Redis for Prometheus monitoring

### Jaeger
- Provides distributed tracing
- Monitors and troubleshoots microservice interactions

### Prometheus
- Collects metrics from all services
- Configured to scrape endpoints such as Kong's /metrics and others

### Grafana
- Visualizes the metrics collected by Prometheus
- Displays dashboards for system health, API usage, and business metrics

### Identity Service
- Sample service built with Go and Gin
- Demonstrates tracing, metrics, and service registration with Consul and Jaeger

## APIs

DIGIT 3.0 exposes several APIs for internal and external use:

### Ping API (Health Check)
- Located in `services/common/identity/main.go`
- Provides a simple endpoint to verify that the service is running

### Account API
- Found in `services/account/cmd/main.go`
- Offers endpoints to create, read, update, and delete account-related data

## Tools & Technologies

DIGIT 3.0 leverages a variety of modern technologies:

| Technology | Version | Purpose |
|------------|---------|----------|
| Go & Gin | 1.24 | High-performance web services |
| Docker Compose | Latest | Container orchestration |
| Consul | 1.15 | Service discovery and registration |
| Kong | 3.9.0 | API gateway with custom plugins |
| PostgreSQL | 14.x/13.x | Registry and Kong databases |
| Kafka | Latest | Messaging system |
| Redis | 7.x | Caching service |
| Jaeger | Latest | Distributed tracing |
| Prometheus & Grafana | 11.6.0 | Monitoring and visualization |
| Keycloak | Latest | Authentication and authorization |

## Deployment

DIGIT 3.0 is deployed using Docker Compose, which defines the services, networks, and volumes in a single YAML file.

### Key Aspects
- **Containerization**: Each service runs in its own container
- **Networking**: All services connected through a common Docker network (app-network)
- **Volumes**: Data persistence via Docker volumes
- **Service Dependencies**: Managed through depends_on with health checks

## Configuration

Configuration is primarily managed via:
- Environment variables in Docker Compose file
- Additional configuration files mounted as volumes
- Dynamic updates without container rebuilds

## Development Guidelines

### 1. Getting Started
```bash
# Clone the repository
git clone <repository-url>
cd digitnxt

# Navigate to docker-compose directory
cd deployment/docker-compose

# Build and run services
docker-compose up --build -d
```

### 2. Service Verification
Check the Identity Service: http://localhost:8000/identity/ping

### 3. Monitoring and Debugging
- **Grafana**: http://localhost:3000
  - View system dashboards
- **Prometheus**: http://localhost:9090
  - Query metrics (e.g., "digit_http_requests_total")
- **Jaeger**: http://localhost:16686
  - Distributed tracing
  - Select "identity" and click "Find Traces"

### 4. Service Discovery
- **Consul**: http://localhost:8500
  - Monitor service registration
  - Check service health

### 5. API Gateway
- **Kong Admin**: http://localhost:8001/services
  - Manage API configurations
  - Monitor plugin status

## Security

DIGIT 3.0 implements comprehensive security measures:

### Automated Security
- Dependabot for automated security updates
- Regular govulncheck scans for Go dependencies

### Authentication & Authorization
- Keycloak integration
- Kong API Gateway security plugins
- Rate limiting and authentication

### Infrastructure Security
- HTTPS/TLS support for production
- Secure service communication
- Detailed security policy in `.github/SECURITY.md`

## Versioning & Dependencies

### Core Components
| Component | Version |
|-----------|---------|
| Go | 1.24 |
| Kong | 3.9.0 |
| PostgreSQL (Registry) | 14.x |
| PostgreSQL (Kong) | 13.x |
| Redis | 7.x |
| Grafana | 11.6.0 |

### Dependency Management
- Go modules via `go.mod`
- Automated updates via Dependabot
- Regular security patches
- Semantic versioning

For detailed version information and compatibility, check the respective `go.mod` files in each service directory.
	
