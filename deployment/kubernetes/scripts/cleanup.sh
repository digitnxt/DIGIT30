#!/bin/bash

# Thorough cleanup script

function cleanup_namespace() {
    local namespace=$1
    echo "Cleaning up namespace: $namespace"
    
    # Delete all resources in namespace
    kubectl delete all --all -n $namespace
    kubectl delete pvc --all -n $namespace
    kubectl delete configmap --all -n $namespace
    kubectl delete secret --all -n $namespace
    kubectl delete serviceaccount --all -n $namespace
    kubectl delete rolebinding --all -n $namespace
    kubectl delete role --all -n $namespace
    kubectl delete networkpolicy --all -n $namespace
    kubectl delete poddisruptionbudget --all -n $namespace
    kubectl delete horizontalpodautoscaler --all -n $namespace
}

function cleanup_cluster() {
    echo "Cleaning up cluster resources..."
    
    # Delete all namespaces except system ones
    for ns in $(kubectl get ns -o jsonpath='{.items[*].metadata.name}'); do
        if [[ $ns != "kube-system" && $ns != "kube-public" && $ns != "kube-node-lease" ]]; then
            cleanup_namespace $ns
            kubectl delete namespace $ns
        fi
    done
    
    # Delete cluster-wide resources
    kubectl delete clusterrolebinding --all
    kubectl delete clusterrole --all
    kubectl delete customresourcedefinition --all
    
    # Clean up storage classes
    kubectl delete storageclass --all
    
    # Clean up persistent volumes
    kubectl delete pv --all
    
    # Clean up priority classes
    kubectl delete priorityclass --all
}

function cleanup_gcp_resources() {
    echo "Cleaning up GCP resources..."
    
    # Delete load balancers
    gcloud compute forwarding-rules delete --quiet --global $(gcloud compute forwarding-rules list --format="value(name)")
    gcloud compute target-http-proxies delete --quiet $(gcloud compute target-http-proxies list --format="value(name)")
    gcloud compute url-maps delete --quiet $(gcloud compute url-maps list --format="value(name)")
    
    # Delete firewall rules
    gcloud compute firewall-rules delete --quiet $(gcloud compute firewall-rules list --format="value(name)")
    
    # Delete disks
    gcloud compute disks delete --quiet $(gcloud compute disks list --format="value(name)")
    
    # Delete snapshots
    gcloud compute snapshots delete --quiet $(gcloud compute snapshots list --format="value(name)")
}

function cleanup_helm() {
    echo "Cleaning up Helm releases..."
    
    # Delete all Helm releases
    helm uninstall $(helm list --all --short)
    
    # Clean up Helm tiller
    kubectl delete deployment tiller-deploy -n kube-system
    kubectl delete service tiller-deploy -n kube-system
}

case "$1" in
    "namespace")
        if [ -z "$2" ]; then
            echo "Usage: $0 namespace <namespace-name>"
            exit 1
        fi
        cleanup_namespace $2
        ;;
    "cluster")
        cleanup_cluster
        ;;
    "gcp")
        cleanup_gcp_resources
        ;;
    "helm")
        cleanup_helm
        ;;
    "all")
        cleanup_cluster
        cleanup_gcp_resources
        cleanup_helm
        ;;
    *)
        echo "Usage: $0 {namespace <name>|cluster|gcp|helm|all}"
        exit 1
        ;;
esac 