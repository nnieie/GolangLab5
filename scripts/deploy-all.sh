#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

CLUSTER_NAME="lab5-cluster"
RELEASE_NAME="golanglab5"
REQUIRED_ENV_VARS=("R2_Endpoint" "R2_ACCESS_KEY_ID" "R2_SECRET_ACCESS_KEY")

echo "========================================="
echo "GolangLab5 Kind Kubernetes Deployment"
echo "========================================="

# 1. 创建 Kind 集群
echo ""
echo "Step 1: Creating Kind cluster..."
if kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "Cluster already exists, skipping..."
else
    kind create cluster --config=kind-config.yaml
fi

# 2. 构建镜像
echo ""
echo "Step 2: Building Docker images..."
"$SCRIPT_DIR/build-images.sh"

# 3. 加载镜像到 Kind
echo ""
echo "Step 3: Loading images to Kind..."
for service in api user video social interaction chat; do
    echo "Loading golanglab5-$service:v1.0..."
    kind load docker-image golanglab5-$service:v1.0 --name "$CLUSTER_NAME"
done

# 4. 使用 Helm 部署
echo ""
echo "Step 4: Deploying with Helm..."
for env_var in "${REQUIRED_ENV_VARS[@]}"; do
    if [ -z "${!env_var}" ]; then
        echo "Missing required environment variable: $env_var"
        exit 1
    fi
done

helm upgrade --install "$RELEASE_NAME" ./chart \
    --set-file config.configK8sContent=config/config.k8s.yaml \
    --set-file config.mysqlInitSqlContent=config/sql/init.sql \
    --set-string r2.endpoint="$R2_Endpoint" \
    --set-string r2.accessKeyId="$R2_ACCESS_KEY_ID" \
    --set-string r2.secretAccessKey="$R2_SECRET_ACCESS_KEY"

# 5. 等待中间件就绪
echo ""
echo "Step 5: Waiting for middleware to be ready (this may take a few minutes)..."
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s
kubectl wait --for=condition=ready pod -l app=otel-collector --timeout=300s

# 6. 等待服务就绪
echo ""
echo "Step 6: Waiting for services to be ready..."
sleep 30
kubectl get pods

echo ""
echo "========================================="
echo "Deployment completed!"
echo "========================================="
echo ""
echo "Access points:"
echo "  - API Gateway: http://localhost:8080"
echo "  - Jaeger UI: http://localhost:16686"
echo ""
echo "Useful commands:"
echo "  - kubectl get pods              # Check pod status"
echo "  - kubectl get svc               # Check services"
echo "  - kubectl logs <pod-name>       # View logs"
echo "  - kubectl describe pod <name>   # Debug pod"
