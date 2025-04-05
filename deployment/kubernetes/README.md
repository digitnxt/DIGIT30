# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the DigitNxt platform on Google Cloud Platform (GCP).

## Resource Requirements

### Core Services

| Service | Replicas | CPU (limits/requests) | Memory (limits/requests) | Storage | Notes |
|---------|----------|----------------------|-------------------------|---------|-------|
| model-context-service | 1 | 100m/50m | 256Mi/128Mi | - | Service discovery |
| mcp-service | 1 | 100m/50m | 256Mi/128Mi | - | API gateway |
| llama-server | 1 | 1000m/500m | 4Gi/2Gi | - | Requires GPU |
| identity | 1 | 100m/50m | 256Mi/128Mi | - | Authentication |
| llm-client | 1 | 100m/50m | 256Mi/128Mi | - | LLM client |

### Infrastructure Services

| Service | Replicas | CPU (limits/requests) | Memory (limits/requests) | Storage | Notes |
|---------|----------|----------------------|-------------------------|---------|-------|
| PostgreSQL | 1 | 500m/250m | 1Gi/512Mi | 5Gi | Main database |
| Kafka | 1 | 500m/250m | 1Gi/512Mi | 5Gi | Message broker |
| Redis | 1 | 500m/250m | 1Gi/512Mi | 2Gi | Cache |
| Kong | 1 | 500m/250m | 1Gi/512Mi | - | API gateway |
| Keycloak | 1 | 500m/250m | 1Gi/512Mi | 5Gi | Identity provider |
| Jaeger | 1 | 500m/250m | 1Gi/512Mi | - | Tracing |
| Prometheus | 1 | 500m/250m | 1Gi/512Mi | 5Gi | Monitoring |
| Grafana | 1 | 500m/250m | 1Gi/512Mi | 5Gi | Visualization |

### Supporting Services

| Service | Replicas | CPU (limits/requests) | Memory (limits/requests) | Storage | Notes |
|---------|----------|----------------------|-------------------------|---------|-------|
| Kafka Exporter | 1 | 100m/50m | 128Mi/64Mi | - | Metrics |
| Redis Exporter | 1 | 100m/50m | 128Mi/64Mi | - | Metrics |

## Total Resource Requirements

- CPU: ~4.5 cores (4500m)
- Memory: ~12.5Gi
- Storage: ~32Gi
- GPU: 1 NVIDIA T4

## VM Recommendations

### Development/Staging Environment
- 2 x e2-standard-4 (4 vCPU, 16GB memory each)
- Total: 8 vCPU, 32GB memory
- Estimated Cost: ~$200/month

### Production Environment
- 3 x e2-standard-4 (4 vCPU, 16GB memory each)
- Total: 12 vCPU, 48GB memory
- Estimated Cost: ~$300/month

## Deployment Steps

1. Create a GCP project and enable required APIs:
   ```bash
   gcloud services enable container.googleapis.com
   gcloud services enable compute.googleapis.com
   ```

2. Create a GKE cluster with GPU support:
   ```bash
   gcloud container clusters create digitnxt-cluster \
     --num-nodes=3 \
     --machine-type=e2-standard-4 \
     --accelerator="type=nvidia-tesla-t4,count=1" \
     --region=us-central1
   ```

3. Install NVIDIA GPU drivers:
   ```bash
   kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/container-engine-accelerators/master/nvidia-driver-installer/cos/daemonset-preloaded.yaml
   ```

4. Create secrets:
   ```bash
   kubectl create secret generic postgres-secret --from-literal=password=your-password
   kubectl create secret generic kong-db-secret --from-literal=POSTGRES_USER=kong --from-literal=POSTGRES_PASSWORD=your-password
   kubectl create secret generic keycloak-db-secret --from-literal=POSTGRES_USER=keycloak --from-literal=POSTGRES_PASSWORD=your-password
   kubectl create secret generic keycloak-secret --from-literal=KEYCLOAK_ADMIN_PASSWORD=your-password
   kubectl create secret generic grafana-secret --from-literal=GF_SECURITY_ADMIN_USER=admin --from-literal=GF_SECURITY_ADMIN_PASSWORD=your-password
   ```

5. Deploy the manifests:
   ```bash
   kubectl apply -f manifests/
   ```

## Monitoring

The deployment includes:
- Prometheus for metrics collection
- Grafana for visualization
- Jaeger for distributed tracing
- Exporters for Kafka and Redis

Access the monitoring dashboards:
- Grafana: `http://<ingress-ip>/grafana`
- Prometheus: `http://<ingress-ip>/prometheus`
- Jaeger: `http://<ingress-ip>/jaeger`

## Scaling

To scale services:
```bash
# Scale a deployment
kubectl scale deployment <deployment-name> --replicas=<desired-replicas>

# Scale a statefulset
kubectl scale statefulset <statefulset-name> --replicas=<desired-replicas>
```

## Backup and Recovery

### Database Backups
```bash
# Create PostgreSQL backup
kubectl exec -it <postgres-pod> -- pg_dump -U admin registry > backup.sql

# Restore PostgreSQL backup
kubectl exec -i <postgres-pod> -- psql -U admin registry < backup.sql
```

### Persistent Volume Backups
```bash
# Create volume snapshot
kubectl create -f volume-snapshot.yaml

# Restore from snapshot
kubectl create -f volume-restore.yaml
```

## Troubleshooting

Common issues and solutions:

1. GPU not available:
   ```bash
   kubectl describe nodes | grep -A 10 "Capacity"
   kubectl describe nodes | grep -A 10 "Allocatable"
   ```

2. Pod scheduling issues:
   ```bash
   kubectl describe pod <pod-name>
   kubectl get events --sort-by='.lastTimestamp'
   ```

3. Service connectivity:
   ```bash
   kubectl get endpoints
   kubectl describe service <service-name>
   ```

## Maintenance

### Updates
```bash
# Update container images
kubectl set image deployment/<deployment-name> <container-name>=<new-image>

# Update configurations
kubectl apply -f updated-config.yaml
```

### Cleanup
```bash
# Delete all resources
kubectl delete -f manifests/

# Delete specific resources
kubectl delete deployment <deployment-name>
kubectl delete statefulset <statefulset-name>
kubectl delete service <service-name>
```

## Security Considerations

1. Use network policies to restrict pod-to-pod communication
2. Enable pod security policies
3. Use secrets for sensitive data
4. Enable audit logging
5. Regularly rotate credentials
6. Monitor for security vulnerabilities

## Cost Optimization and Resource Management

### Cost Optimization Script

The `cost-optimization.sh` script provides several functions for managing costs and resources:

```bash
# Monitor resource usage
./scripts/cost-optimization.sh monitor

# Optimize costs (scales down non-critical services during off-hours)
./scripts/cost-optimization.sh optimize

# Clean up unused resources
./scripts/cost-optimization.sh cleanup

# Generate cost report
./scripts/cost-optimization.sh report
```

### Cleanup Scripts

The `cleanup.sh` script provides thorough cleanup of resources:

```bash
# Clean up specific namespace
./scripts/cleanup.sh namespace <namespace-name>

# Clean up cluster resources
./scripts/cleanup.sh cluster

# Clean up GCP resources
./scripts/cleanup.sh gcp

# Clean up Helm releases
./scripts/cleanup.sh helm

# Clean up everything
./scripts/cleanup.sh all
```

### Resource Monitoring

#### Real-time Monitoring
```bash
# Monitor pod resource usage
kubectl top pods --all-namespaces

# Monitor node resource usage
kubectl top nodes

# Monitor storage usage
kubectl get pvc --all-namespaces
```

#### Cost Tracking
```bash
# Get cost report
./scripts/cost-optimization.sh report

# Monitor GCP billing
gcloud billing projects describe <project-id> --format="value(billingEnabled)"
```

#### Automated Cost Optimization

The system includes several automated cost optimization features:

1. **Auto-scaling**:
   - Cluster scales down to 1 node during low usage
   - Services scale down during off-hours
   - Preemptible nodes for non-critical workloads

2. **Resource Cleanup**:
   - Automatic cleanup of completed jobs
   - Removal of unused PVCs and ConfigMaps
   - Cleanup of old logs and images

3. **Storage Optimization**:
   - Minimal persistent volume sizes
   - Automatic cleanup of unused volumes
   - Regular cleanup of old snapshots

4. **Monitoring and Alerts**:
   - Resource usage monitoring
   - Cost threshold alerts
   - Automated scaling recommendations

### Cost Optimization Best Practices

1. **Resource Allocation**:
   - Use preemptible nodes for non-critical workloads
   - Implement horizontal pod autoscaling
   - Use cluster autoscaler
   - Set appropriate resource limits and requests

2. **Storage Management**:
   - Use appropriate storage classes
   - Implement retention policies
   - Regular cleanup of unused volumes
   - Use volume snapshots for backups

3. **Network Optimization**:
   - Use internal load balancers where possible
   - Implement network policies
   - Use appropriate ingress controllers
   - Monitor network usage

4. **Monitoring and Maintenance**:
   - Regular resource usage reviews
   - Automated cleanup of unused resources
   - Cost threshold monitoring
   - Regular optimization reviews

## Network Configuration

### Ingress Controller Setup
```bash
# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Create Ingress resource
kubectl apply -f manifests/ingress.yaml
```

### Network Policies
```bash
# Enable network policy enforcement
gcloud container clusters update digitnxt-cluster \
  --enable-network-policy \
  --region=us-central1

# Apply network policies
kubectl apply -f manifests/network-policies/
```

### DNS Configuration
```bash
# Create DNS records for services
gcloud dns record-sets create api.digitnxt.com \
  --type=A \
  --ttl=300 \
  --zone=digitnxt-zone \
  --rrdatas=<ingress-ip>

gcloud dns record-sets create grafana.digitnxt.com \
  --type=A \
  --ttl=300 \
  --zone=digitnxt-zone \
  --rrdatas=<ingress-ip>
```

## Load Balancing

### External Load Balancer
```bash
# Create load balancer service
kubectl apply -f manifests/loadbalancer.yaml

# Configure SSL/TLS
kubectl create secret tls digitnxt-tls \
  --cert=ssl/cert.pem \
  --key=ssl/key.pem
```

### Internal Load Balancer
```bash
# Create internal load balancer for database access
kubectl apply -f manifests/internal-lb.yaml
```

### Health Checks
```bash
# Configure health check endpoints
kubectl apply -f manifests/health-checks.yaml
```

## Security Configuration

### Pod Security Policies
```bash
# Enable PodSecurityPolicy admission controller
gcloud container clusters update digitnxt-cluster \
  --enable-pod-security-policy \
  --region=us-central1

# Create security policies
kubectl apply -f manifests/security/pod-security-policies.yaml
```

### RBAC Configuration
```bash
# Create service accounts
kubectl apply -f manifests/security/service-accounts.yaml

# Create role bindings
kubectl apply -f manifests/security/role-bindings.yaml
```

### Network Security
```bash
# Create network policies
kubectl apply -f manifests/security/network-policies.yaml

# Configure SSL/TLS
kubectl apply -f manifests/security/tls-config.yaml
```

### Container Security
```bash
# Enable container image scanning
gcloud container clusters update digitnxt-cluster \
  --enable-image-streaming \
  --region=us-central1

# Configure container runtime security
kubectl apply -f manifests/security/container-security.yaml
```

## Monitoring and Logging

### Stackdriver Integration
```bash
# Enable Stackdriver monitoring
gcloud container clusters update digitnxt-cluster \
  --enable-stackdriver-kubernetes \
  --region=us-central1

# Configure logging
kubectl apply -f manifests/monitoring/logging-config.yaml
```

### Audit Logging
```bash
# Enable audit logging
gcloud container clusters update digitnxt-cluster \
  --enable-audit-logging \
  --region=us-central1

# Configure audit policies
kubectl apply -f manifests/monitoring/audit-policies.yaml
```

## Backup and Disaster Recovery

### Cluster Backup
```bash
# Enable cluster backup
gcloud container clusters update digitnxt-cluster \
  --enable-backup-restore \
  --region=us-central1

# Configure backup schedule
kubectl apply -f manifests/backup/backup-schedule.yaml
```

### Disaster Recovery
```bash
# Create disaster recovery plan
kubectl apply -f manifests/backup/disaster-recovery.yaml

# Test recovery procedures
kubectl apply -f manifests/backup/recovery-test.yaml
```

## Cost Management

### Resource Quotas
```bash
# Create resource quotas
kubectl apply -f manifests/cost/resource-quotas.yaml

# Configure limit ranges
kubectl apply -f manifests/cost/limit-ranges.yaml
```

### Autoscaling
```bash
# Enable cluster autoscaling
gcloud container clusters update digitnxt-cluster \
  --enable-autoscaling \
  --min-nodes=2 \
  --max-nodes=5 \
  --region=us-central1

# Configure HPA
kubectl apply -f manifests/cost/horizontal-pod-autoscaler.yaml
```

## Maintenance Procedures

### Node Maintenance
```bash
# Drain nodes for maintenance
kubectl drain <node-name> --ignore-daemonsets

# Uncordon nodes after maintenance
kubectl uncordon <node-name>
```

### Cluster Updates
```bash
# Update cluster version
gcloud container clusters upgrade digitnxt-cluster \
  --cluster-version=<new-version> \
  --region=us-central1

# Update node pools
gcloud container node-pools upgrade <pool-name> \
  --cluster=digitnxt-cluster \
  --region=us-central1
```

## Troubleshooting

### Network Issues
```bash
# Check network connectivity
kubectl run -it --rm debug --image=nicolaka/netshoot -- /bin/bash

# Check DNS resolution
kubectl exec -it <pod-name> -- nslookup <service-name>
```

### Security Issues
```bash
# Check security policies
kubectl get podsecuritypolicies
kubectl get networkpolicies

# Check RBAC configuration
kubectl get roles
kubectl get rolebindings
```

### Performance Issues
```bash
# Check resource usage
kubectl top nodes
kubectl top pods

# Check cluster metrics
kubectl get --raw /apis/metrics.k8s.io/v1beta1/nodes
```

## Temporary Deployment Guide

For testing and development purposes, you can deploy a temporary cluster with minimal resource allocation. This setup is optimized for quick deployment and cleanup.

### Prerequisites
1. Install [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
2. Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
3. Install [Helm](https://helm.sh/docs/intro/install/)

### Quick Start

1. Configure your GCP project:
```bash
gcloud config set project YOUR_PROJECT_ID
```

2. Make the management script executable:
```bash
chmod +x scripts/manage-cluster.sh
```

3. Deploy the cluster:
```bash
./scripts/manage-cluster.sh create
```

4. Clean up when done:
```bash
./scripts/manage-cluster.sh delete
```

### Cost Optimization Features

The temporary deployment includes several cost-saving measures:

1. **Preemptible Nodes**: Using preemptible VMs for worker nodes (up to 80% cost savings)
2. **Ephemeral IPs**: Load balancers use ephemeral IPs instead of static ones
3. **Minimal Storage**: Reduced persistent volume sizes
4. **Auto-scaling**: Cluster scales down to 1 node when idle
5. **Resource Limits**: Conservative resource allocations for all services

### Important Notes

1. **Data Persistence**: This setup is not suitable for production data as it uses preemptible nodes and minimal storage
2. **Availability**: Preemptible nodes can be terminated at any time
3. **Cleanup**: Always run the cleanup script when done to avoid unnecessary charges
4. **Monitoring**: Stackdriver monitoring is enabled to track resource usage

### Cost Estimation

For a typical test deployment:
- 2 preemptible e2-standard-4 nodes: ~$0.20/hour
- Load balancer: ~$0.025/hour
- Storage: ~$0.10/GB/month
- Total estimated cost: ~$0.50/hour 