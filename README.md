DIGIT30

DIGIT30 is a microservices-based platform built with Go, Gin, and several modern tools to demonstrate a scalable, observable, and well-orchestrated distributed system. It leverages Docker Compose for deployment, enabling rapid development, testing, and production readiness.

Table of Contents
	•	Architecture Overview
	•	Services Overview
	•	APIs
	•	Tools & Technologies
	•	Deployment
	•	Configuration
	•	Development Guidelines

Architecture Overview

DIGIT30 follows a microservices architecture where individual components handle specific responsibilities. The system integrates service discovery, API gateway management, tracing, logging, and monitoring to ensure robust communication and observability across services.

Services Overview

The platform consists of multiple containerized services managed through Docker Compose. Key services include:
	•	Postgres (Registry)
A PostgreSQL database (version 14) used as a registry, with health checks to ensure database readiness.
	•	Consul
A service discovery tool (HashiCorp Consul 1.15) that facilitates service registration and health monitoring.
	•	Kafka & Kafka Exporter
Kafka (using Bitnami’s image) serves as the messaging backbone, while Kafka Exporter provides metrics for monitoring.
	•	Kong API Gateway & Kong DB
	•	Kong: Acts as the API gateway with custom plugins (such as api-costing) for managing, monitoring, and securing API traffic.
	•	Kong DB: A dedicated PostgreSQL instance (version 13) for Kong’s internal database.
	•	Kong Synchronizer
A service (built from a Node.js application) ensuring data consistency across Kong and other system components.
	•	Keycloak
Provides authentication and authorization via Keycloak, enabling secure access to the platform.
	•	Redis & Redis Exporter
Redis is used for caching, and Redis Exporter collects metrics from Redis for Prometheus monitoring.
	•	Jaeger
Provides distributed tracing to monitor and troubleshoot microservice interactions.
	•	Prometheus
Collects metrics from all services. Prometheus is configured to scrape endpoints such as Kong’s /metrics and others.
	•	Grafana
Visualizes the metrics collected by Prometheus, displaying dashboards for system health, API usage, and business metrics.
	•	Identity Service
A sample service (built with Go and Gin) demonstrating tracing, metrics, and service registration with Consul and Jaeger.

APIs

DIGIT30 exposes several APIs for internal and external use:
	•	Ping API (Health Check)
Located in DIGIT30/services/common/identity/main.go, it provides a simple endpoint to verify that the service is running.
	•	Account API
Found in DIGIT30/services/account/cmd/main.go, it offers endpoints to create, read, update, and delete account-related data.

Tools & Technologies

DIGIT30 leverages a variety of modern technologies:
	•	Go & Gin: For building high-performance web services.
	•	Docker Compose: To orchestrate multiple containerized services.
	•	Consul: For service discovery and registration.
	•	Kong: As an API gateway, with custom plugins like api-costing for API-based costing.
	•	PostgreSQL: Used for both registry and Kong’s internal database.
	•	Kafka: As a messaging system, along with its exporter for monitoring.
	•	Redis: For caching, supported by a Redis Exporter.
	•	Jaeger: For distributed tracing across services.
	•	Prometheus & Grafana: For monitoring, alerting, and visualizing system metrics.
	•	Keycloak: To handle authentication and authorization.

Deployment

DIGIT30 is deployed using Docker Compose, which defines the services, networks, and volumes in a single YAML file. Key aspects include:
	•	Containerization: Each service runs in its own container, making it easier to scale and manage dependencies.
	•	Networking: All services are connected through a common Docker network (app-network), enabling seamless inter-service communication.
	•	Volumes: Data persistence is managed via Docker volumes (e.g., postgres-data, kong-db-data, kafka-data, etc.).
	•	Service Dependencies: The Compose file uses depends_on with health and start conditions to ensure that services start in the correct order.

Configuration

Configuration is primarily managed via environment variables specified in the Docker Compose file. This includes database credentials, API gateway settings, and plugin configurations. Additional configuration files (e.g., for Prometheus and Grafana) are mounted as volumes to enable dynamic updates without rebuilding containers.

Development Guidelines

To get started with DIGIT30:
	1.	Clone the Repository:
Clone the repository and navigate to the desired service directory.
	2.	Build and Run Services:
Use go build and go run for individual services during development. For a full-stack environment, run:

docker-compose up -d


	3.	Monitoring and Debugging:
	•	Access Grafana at http://localhost:3000 to view system dashboards.
	•	Use Prometheus at http://localhost:9090 for querying metrics.
	•	Jaeger is available at http://localhost:16686 for distributed tracing.
	4.	Service Discovery:
Ensure that Consul (running on port 8500) is correctly registering services, allowing for dynamic discovery and routing.
	5.	API Gateway Management:
The Kong API gateway (available on ports 8000 and 8001) manages all API requests and integrates custom plugins like api-costing for detailed metrics and costing.
	6.	Contribution:
Follow the contribution guidelines provided in the repository. Make sure to adhere to coding standards, write tests, and document your changes.
