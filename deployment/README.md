# Deployment Guide for Google Cloud

This guide will help you deploy the application on Google Cloud with GPU support.

## Prerequisites

1. Google Cloud account with billing enabled
2. Google Cloud SDK installed
3. Docker installed on your local machine

## Step 1: Create a GPU-enabled VM

```bash
gcloud compute instances create llama-server \
    --machine-type=n1-standard-4 \
    --zone=us-central1-a \
    --boot-disk-size=200GB \
    --accelerator="type=nvidia-tesla-t4,count=1" \
    --maintenance-policy=TERMINATE \
    --image-family=ubuntu-2004-lts \
    --image-project=ubuntu-os-cloud \
    --metadata="install-nvidia-driver=True"
```

## Step 2: Set up the environment

1. SSH into the VM:
```bash
gcloud compute ssh llama-server
```

2. Clone the repository:
```bash
git clone <your-repository-url>
cd digitnxt
```

3. Make the setup script executable and run it:
```bash
chmod +x deployment/scripts/setup-cloud.sh
./deployment/scripts/setup-cloud.sh
```

## Step 3: Verify the deployment

1. Check if all services are running:
```bash
docker-compose ps
```

2. Test the llama-server:
```bash
curl http://localhost:8082/health
```

3. Test the MCP service:
```bash
curl http://localhost:8086/health
```

## Step 4: Set up monitoring (optional)

1. Access Grafana at `http://<vm-ip>:3000`
2. Default credentials:
   - Username: admin
   - Password: admin

## Troubleshooting

1. If GPU is not detected:
```bash
nvidia-smi
```

2. Check Docker logs:
```bash
docker logs llama-server
```

3. Check NVIDIA Container Toolkit:
```bash
nvidia-container-cli --version
```

## Security Considerations

1. Update default passwords for:
   - Grafana
   - Keycloak
   - PostgreSQL

2. Configure firewall rules to restrict access to necessary ports only

3. Enable HTTPS for all services

## Cost Optimization

1. Use preemptible instances for non-critical workloads
2. Set up auto-scaling based on load
3. Use spot instances for development environments

## Maintenance

1. Regular backups of:
   - PostgreSQL databases
   - Model files
   - Configuration files

2. Monitor GPU utilization and adjust resources as needed

3. Keep all components updated to the latest stable versions 