# LLM Components

This directory contains the LLM-related components for the DIGIT 3.0 platform's backbone services. These components provide AI-powered capabilities for API documentation, service discovery, and intelligent service interactions.

## Components

### llama-server
A Dockerized instance of llama.cpp that provides a REST API for LLM inference. It uses the Qwen2-7B-Instruct model quantized to 5-bit precision. This service is part of the platform's backbone infrastructure, providing AI capabilities to other services.

### llm_client
A Python client that interacts with the llama-server to:
- Generate API calls based on Swagger documentation
- Explain API endpoints and their usage
- Provide context-aware responses about the platform's services
- Assist in service discovery and documentation

## Configuration

### llama-server
- Port: 8082
- Model: Qwen2-7B-Instruct (5-bit quantized)
- Context length: 4096 tokens
- GPU acceleration enabled
- Network: digit30_app-network

### llm_client
- Connects to llama-server on port 8082
- Integrates with the identity service for API documentation
- Provides a simple interface for generating API calls and explanations
- Part of the platform's backbone services

## Usage

1. Start the llama-server:
```bash
docker-compose up llama-server
```

2. Run the llm_client:
```bash
docker-compose run llm_client
```

## Integration with Backbone Services

The LLM components integrate with other backbone services:
- Consul: For service discovery
- Kong: For API gateway and routing
- Jaeger: For distributed tracing
- Prometheus: For metrics collection
- Grafana: For monitoring and visualization

## Dependencies

- llama.cpp
- Python 3.9+
- Docker
- NVIDIA GPU (for acceleration)
- Other backbone services (Consul, Kong, etc.) 