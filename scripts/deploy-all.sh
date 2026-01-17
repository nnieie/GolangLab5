#!/bin/bash
set -e

echo "========================================="
echo "GolangLab5 Kind Kubernetes Deployment"
echo "========================================="

# 1. 创建 Kind 集群
echo ""
echo "Step 1: Creating Kind cluster..."
if kind get clusters | grep -q "golanglab5-cluster"; then
    echo "Cluster already exists, skipping..."
else
    kind create cluster --config=kind-config.yaml
fi

# 2. 构建镜像
echo ""
echo "Step 2: Building Docker images..."
./build-images.sh

# 3. 加载镜像到 Kind
echo ""
echo "Step 3: Loading images to Kind..."
for service in api user video social interaction chat; do
    echo "Loading golanglab5-$service:v1.0..."
    kind load docker-image golanglab5-$service:v1.0 --name golanglab5-cluster
done

# 4. 创建 ConfigMap
echo ""
echo "Step 4: Creating ConfigMaps..."
if kubectl get configmap mysql-init-sql &> /dev/null; then
    echo "ConfigMap already exists, skipping..."
else
    kubectl create configmap mysql-init-sql --from-file=config/sql/init.sql
fi

# 5. 部署中间件
echo ""
echo "Step 5: Deploying middleware..."
kubectl apply -f k8s-manifests/middleware.yaml

# 6. 等待中间件就绪
echo ""
echo "Step 6: Waiting for middleware to be ready (this may take a few minutes)..."
echo "Waiting for MySQL..."
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
echo "Waiting for Etcd..."
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
echo "Waiting for Redis..."
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
echo "Waiting for Kafka..."
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s

# 7. 部署微服务
echo ""
echo "Step 7: Deploying microservices..."
kubectl apply -f k8s-manifests/services.yaml

# 8. 等待服务就绪
echo ""
echo "Step 8: Waiting for services to be ready..."
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
