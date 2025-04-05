#!/bin/bash

# Exit on error
set -e

# Check if running on Google Cloud
if [ -z "$GOOGLE_CLOUD_PROJECT" ]; then
    echo "Error: GOOGLE_CLOUD_PROJECT environment variable not set"
    exit 1
fi

# Install required packages
echo "Installing required packages..."
sudo apt-get update
sudo apt-get install -y \
    curl \
    wget \
    git \
    python3 \
    python3-pip \
    docker.io \
    nvidia-container-toolkit

# Install NVIDIA Container Toolkit
echo "Setting up NVIDIA Container Toolkit..."
curl -s -L https://nvidia.github.io/nvidia-container-runtime/gpgkey | sudo apt-key add -
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-container-runtime/$distribution/nvidia-container-runtime.list | sudo tee /etc/apt/sources.list.d/nvidia-container-runtime.list
sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker

# Verify NVIDIA installation
echo "Verifying NVIDIA installation..."
nvidia-smi

# Create necessary directories
echo "Creating directories..."
sudo mkdir -p /models
sudo chmod 777 /models

# Download the model with retries
echo "Downloading model from Hugging Face..."
MAX_RETRIES=3
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if wget -O /models/qwen2-7b-instruct-q5_k_m.gguf https://huggingface.co/TheBloke/Qwen2-7B-Instruct-GGUF/resolve/main/qwen2-7b-instruct-q5_k_m.gguf; then
        echo "Model downloaded successfully"
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo "Download failed, retrying... (Attempt $RETRY_COUNT of $MAX_RETRIES)"
        sleep 10
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "Failed to download model after $MAX_RETRIES attempts"
    exit 1
fi

# Verify model file
echo "Verifying model file..."
if [ ! -f "/models/qwen2-7b-instruct-q5_k_m.gguf" ]; then
    echo "Error: Model file not found"
    exit 1
fi

# Install Docker Compose
echo "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Pull the llama.cpp image
echo "Pulling llama.cpp image..."
docker pull ghcr.io/ggerganov/llama.cpp:full

# Start the services
echo "Starting services..."
cd /deployment/docker-compose
docker-compose up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
sleep 30

# Check service health
echo "Checking service health..."
docker-compose ps

# Verify llama-server is running with GPU
echo "Verifying llama-server GPU access..."
docker exec llama-server nvidia-smi

echo "Setup completed successfully!" 