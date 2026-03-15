#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

CLUSTER_NAME="lab5-cluster"
REQUIRED_ENV_VARS=("R2_Endpoint" "R2_ACCESS_KEY_ID" "R2_SECRET_ACCESS_KEY")
SERVICES=(api user video social interaction chat)

assert_required_env_vars() {
  local env_var
  for env_var in "${REQUIRED_ENV_VARS[@]}"; do
    if [[ -z "${!env_var:-}" ]]; then
      echo "Missing required environment variable: $env_var"
      exit 1
    fi
  done
}

echo "================================================="
echo "GolangLab5 Kind Deployment Without Helm"
echo "================================================="

echo ""
echo "Step 1: Creating Kind cluster..."
if kind get clusters | grep -q "$CLUSTER_NAME"; then
  echo "Cluster already exists, skipping..."
else
  kind create cluster --config=kind-config.yaml
fi

echo ""
echo "Step 2: Building Docker images..."
"$SCRIPT_DIR/build-images.sh"

echo ""
echo "Step 3: Loading images to Kind..."
for service in "${SERVICES[@]}"; do
  echo "Loading golanglab5-$service:v1.0..."
  kind load docker-image "golanglab5-$service:v1.0" --name "$CLUSTER_NAME"
done

echo ""
echo "Step 4: Applying ConfigMaps and Secrets..."
assert_required_env_vars

kubectl create configmap app-config \
  --from-file=config.k8s.yaml=config/config.k8s.yaml \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap mysql-init-sql \
  --from-file=init.sql=config/sql/init.sql \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic r2-config \
  --from-literal="endpoint=${R2_Endpoint}" \
  --from-literal="accessKeyId=${R2_ACCESS_KEY_ID}" \
  --from-literal="secretAccessKey=${R2_SECRET_ACCESS_KEY}" \
  --dry-run=client -o yaml | kubectl apply -f -

echo ""
echo "Step 5: Applying Kubernetes manifests..."
kubectl apply -f k8s-manifests/middleware.yaml
kubectl apply -f k8s-manifests/services.yaml

echo ""
echo "Step 6: Waiting for middleware to be ready (this may take a few minutes)..."
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s

echo ""
echo "Step 7: Syncing MySQL schema..."
MYSQL_POD="$(kubectl get pod -l app=mysql -o jsonpath='{.items[0].metadata.name}')"
if [[ -z "$MYSQL_POD" ]]; then
  echo "Failed to find MySQL pod"
  exit 1
fi

kubectl exec "$MYSQL_POD" -- sh -lc 'export MYSQL_PWD="$MYSQL_ROOT_PASSWORD"; mysql -uroot "$MYSQL_DATABASE" < /docker-entrypoint-initdb.d/init.sql'

echo ""
echo "Step 8: Waiting for services to be ready..."
kubectl rollout status deployment/user-service --timeout=300s
kubectl rollout status deployment/video-service --timeout=300s
kubectl rollout status deployment/social-service --timeout=300s
kubectl rollout status deployment/interaction-service --timeout=300s
kubectl rollout status deployment/chat-service --timeout=300s
kubectl rollout status deployment/api-gateway --timeout=300s

echo ""
kubectl get pods

echo ""
echo "================================================="
echo "Deployment without Helm completed!"
echo "================================================="
echo ""
echo "Access points:"
echo "  - API Gateway: http://localhost:8080"
echo "  - Jaeger UI: http://localhost:16686"
