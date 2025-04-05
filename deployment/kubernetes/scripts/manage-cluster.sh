#!/bin/bash

# Configuration
PROJECT_ID="your-project-id"
CLUSTER_NAME="digitnxt-test"
REGION="us-central1"
ZONE="us-central1-a"
MACHINE_TYPE="e2-standard-4"
NUM_NODES=2
DISK_SIZE=50

function create_cluster() {
    echo "Creating GKE cluster..."
    gcloud container clusters create $CLUSTER_NAME \
        --project=$PROJECT_ID \
        --zone=$ZONE \
        --machine-type=$MACHINE_TYPE \
        --num-nodes=$NUM_NODES \
        --disk-size=$DISK_SIZE \
        --enable-autoscaling \
        --min-nodes=1 \
        --max-nodes=3 \
        --enable-ip-alias \
        --create-subnetwork="" \
        --enable-network-policy \
        --no-enable-autoupgrade \
        --enable-shielded-nodes \
        --enable-stackdriver-kubernetes \
        --preemptible
}

function deploy_helm() {
    echo "Deploying Helm chart..."
    helm install digitnxt ./helm-chart \
        --namespace digitnxt \
        --create-namespace \
        --set global.environment=test \
        --set postgresql.persistence.size=5Gi \
        --set redis.persistence.size=2Gi \
        --set kafka.persistence.size=5Gi \
        --set kong.persistence.size=5Gi \
        --set keycloak.persistence.size=5Gi \
        --set prometheus.persistence.size=5Gi \
        --set grafana.persistence.size=5Gi
}

function cleanup() {
    echo "Cleaning up resources..."
    helm uninstall digitnxt --namespace digitnxt
    kubectl delete namespace digitnxt
    gcloud container clusters delete $CLUSTER_NAME \
        --project=$PROJECT_ID \
        --zone=$ZONE \
        --quiet
}

function get_credentials() {
    echo "Getting cluster credentials..."
    gcloud container clusters get-credentials $CLUSTER_NAME \
        --project=$PROJECT_ID \
        --zone=$ZONE
}

case "$1" in
    "create")
        create_cluster
        get_credentials
        deploy_helm
        ;;
    "delete")
        cleanup
        ;;
    *)
        echo "Usage: $0 {create|delete}"
        exit 1
        ;;
esac 