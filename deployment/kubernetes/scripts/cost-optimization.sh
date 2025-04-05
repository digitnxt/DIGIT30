#!/bin/bash

# Cost optimization and monitoring script

function monitor_resources() {
    echo "Monitoring resource usage..."
    
    # Get cluster metrics
    echo "Cluster Metrics:"
    kubectl get --raw /apis/metrics.k8s.io/v1beta1/nodes | jq '.items[] | {name: .metadata.name, cpu: .usage.cpu, memory: .usage.memory}'
    
    # Get pod metrics
    echo "Pod Metrics:"
    kubectl top pods --all-namespaces
    
    # Get node metrics
    echo "Node Metrics:"
    kubectl top nodes
    
    # Get storage usage
    echo "Storage Usage:"
    kubectl get pvc --all-namespaces -o json | jq '.items[] | {name: .metadata.name, namespace: .metadata.namespace, capacity: .status.capacity.storage}'
}

function optimize_costs() {
    echo "Optimizing costs..."
    
    # Scale down non-critical services during off-hours
    if [[ $(date +%H) -ge 22 || $(date +%H) -le 6 ]]; then
        echo "Scaling down non-critical services..."
        kubectl scale deployment model-context-service --replicas=0
        kubectl scale deployment mcp-service --replicas=0
        kubectl scale deployment llm-client --replicas=0
    else
        echo "Scaling up services..."
        kubectl scale deployment model-context-service --replicas=1
        kubectl scale deployment mcp-service --replicas=1
        kubectl scale deployment llm-client --replicas=1
    fi
    
    # Clean up completed jobs
    echo "Cleaning up completed jobs..."
    kubectl delete jobs --field-selector status.successful=1 --all-namespaces
    
    # Clean up failed pods
    echo "Cleaning up failed pods..."
    kubectl delete pods --field-selector status.phase=Failed --all-namespaces
    
    # Clean up old logs
    echo "Cleaning up old logs..."
    kubectl logs --since=24h --all-containers=true --all-namespaces > logs-$(date +%Y%m%d).log
    kubectl delete pods --field-selector status.phase=Succeeded --all-namespaces
}

function cleanup_resources() {
    echo "Cleaning up unused resources..."
    
    # Delete unused PVCs
    echo "Deleting unused PVCs..."
    kubectl get pvc --all-namespaces -o json | jq '.items[] | select(.status.phase=="Released") | .metadata.name' | xargs -I {} kubectl delete pvc {} --namespace={}
    
    # Delete unused ConfigMaps
    echo "Deleting unused ConfigMaps..."
    kubectl get configmaps --all-namespaces -o json | jq '.items[] | select(.metadata.ownerReferences==null) | .metadata.name' | xargs -I {} kubectl delete configmap {} --namespace={}
    
    # Delete unused Secrets
    echo "Deleting unused Secrets..."
    kubectl get secrets --all-namespaces -o json | jq '.items[] | select(.metadata.ownerReferences==null) | .metadata.name' | xargs -I {} kubectl delete secret {} --namespace={}
    
    # Clean up old images
    echo "Cleaning up old images..."
    kubectl get pods --all-namespaces -o json | jq '.items[].spec.containers[].image' | sort | uniq -c | sort -nr
}

function generate_cost_report() {
    echo "Generating cost report..."
    
    # Get resource usage
    CPU_USAGE=$(kubectl get nodes -o json | jq '.items[].status.allocatable.cpu' | tr -d '"' | awk '{sum += $1} END {print sum}')
    MEMORY_USAGE=$(kubectl get nodes -o json | jq '.items[].status.allocatable.memory' | tr -d '"' | awk '{sum += $1} END {print sum}')
    STORAGE_USAGE=$(kubectl get pvc --all-namespaces -o json | jq '.items[].status.capacity.storage' | tr -d '"' | awk '{sum += $1} END {print sum}')
    
    # Calculate estimated costs
    NODE_COST=$(echo "$CPU_USAGE * 0.0335 + $MEMORY_USAGE * 0.004445" | bc)
    STORAGE_COST=$(echo "$STORAGE_USAGE * 0.04" | bc)
    TOTAL_COST=$(echo "$NODE_COST + $STORAGE_COST" | bc)
    
    echo "Cost Report:"
    echo "------------"
    echo "CPU Usage: $CPU_USAGE cores"
    echo "Memory Usage: $MEMORY_USAGE"
    echo "Storage Usage: $STORAGE_USAGE"
    echo "Node Cost: \$$NODE_COST/hour"
    echo "Storage Cost: \$$STORAGE_COST/hour"
    echo "Total Cost: \$$TOTAL_COST/hour"
}

case "$1" in
    "monitor")
        monitor_resources
        ;;
    "optimize")
        optimize_costs
        ;;
    "cleanup")
        cleanup_resources
        ;;
    "report")
        generate_cost_report
        ;;
    *)
        echo "Usage: $0 {monitor|optimize|cleanup|report}"
        exit 1
        ;;
esac 