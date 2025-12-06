#!/bin/bash
set -e

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker is not running. Please start Docker Desktop first."
  exit 1
fi

# Install kind if not present
if ! command -v kind &> /dev/null; then
    echo "Installing kind..."
    brew install kind
fi

echo "Creating clusters..."

# Create 3 clusters
clusters=("prod-cluster1" "prod-cluster2" "dev-cluster")

for cluster in "${clusters[@]}"; do
    if kind get clusters | grep -q "^$cluster$"; then
        echo "Cluster $cluster already exists, skipping creation."
    else
        echo "Creating cluster: $cluster..."
        if ! kind create cluster --name "$cluster"; then
            echo "Failed to create cluster $cluster. Check Docker resources."
            continue
        fi
        echo "Waiting for cluster to be ready..."
        sleep 5
    fi
    
    # Switch context to ensure we are targeting the right one
    kubectl config use-context "kind-$cluster" || true

    # Wait for default service account to be created
    echo "Waiting for default service account..."
    for i in {1..30}; do
        if kubectl get sa default > /dev/null 2>&1; then
            break
        fi
        sleep 2
    done
    
    # Create a dummy pod if it doesn't exist
    if ! kubectl get pod nginx > /dev/null 2>&1; then
        echo "Deploying nginx pod to $cluster..."
        kubectl run nginx --image=nginx --restart=Never
        
        # Wait for pod to be created
        sleep 2
        kubectl label pod nginx env=$cluster --overwrite
    else
        echo "Pod nginx already exists in $cluster."
    fi
    
    echo "----------------------------------------"
done

echo "Merging kubeconfigs..."
# Kind automatically merges into ~/.kube/config by default, so we are good.

echo "Done! You now have 3 clusters running."
echo "Try running:"
echo "  mk get pods"
echo "  mk get pods --context prod"
