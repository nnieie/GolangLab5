# PowerShell Script for Full Deployment to Kind Cluster
# Usage: .\deploy-all.ps1

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "GolangLab5 Kind Kubernetes Deployment" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

$clusterName = "lab5-cluster"

# 1. 创建 Kind 集群
Write-Host ""
Write-Host "Step 1: Creating Kind cluster..." -ForegroundColor Yellow
$clusters = kind get clusters 2>$null
if ($clusters -contains $clusterName) {
    Write-Host "Cluster already exists, skipping..." -ForegroundColor Gray
} else {
    kind create cluster --config=kind-config.yaml
    if ($LASTEXITCODE -ne 0) {
        Write-Host "✗ Failed to create cluster" -ForegroundColor Red
        exit 1
    }
}

# 2. 构建镜像
Write-Host ""
Write-Host "Step 2: Building Docker images..." -ForegroundColor Yellow
.\build-images.ps1
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to build images" -ForegroundColor Red
    exit 1
}

# 3. 加载镜像到 Kind
Write-Host ""
Write-Host "Step 3: Loading images to Kind..." -ForegroundColor Yellow
$services = @("api", "user", "video", "social", "interaction", "chat")
foreach ($service in $services) {
    Write-Host "Loading golanglab5-${service}:v1.0..." -ForegroundColor White
    kind load docker-image "golanglab5-${service}:v1.0" --name $clusterName
}

# kind load docker-image "golanglab5-api:v1.0" --name lab5-cluster
# kind load docker-image "golanglab5-user:v1.0" --name lab5-cluster
# kind load docker-image "golanglab5-video:v1.0" --name lab5-cluster
# kind load docker-image "golanglab5-social:v1.0" --name lab5-cluster
# kind load docker-image "golanglab5-interaction:v1.0" --name lab5-cluster
# kind load docker-image "golanglab5-chat:v1.0" --name lab5-cluster

# 4. 创建 ConfigMap
Write-Host ""
Write-Host "Step 4: Creating ConfigMaps..." -ForegroundColor Yellow
$configMapExists = kubectl get configmap mysql-init-sql 2>$null
if ($configMapExists) {
    Write-Host "ConfigMap already exists, skipping..." -ForegroundColor Gray
} else {
    kubectl create configmap mysql-init-sql --from-file=config/sql/init.sql
}

# 5. 部署中间件
Write-Host ""
Write-Host "Step 5: Deploying middleware..." -ForegroundColor Yellow
kubectl apply -f k8s-manifests/middleware.yaml

# 6. 等待中间件就绪
Write-Host ""
Write-Host "Step 6: Waiting for middleware to be ready (this may take a few minutes)..." -ForegroundColor Yellow
Write-Host "Waiting for MySQL..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
Write-Host "Waiting for Etcd..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
Write-Host "Waiting for Redis..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
Write-Host "Waiting for Kafka..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s

# 7. 部署微服务
Write-Host ""
Write-Host "Step 7: Deploying microservices..." -ForegroundColor Yellow
kubectl apply -f k8s-manifests/services.yaml

# 8. 等待服务就绪
Write-Host ""
Write-Host "Step 8: Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 30
kubectl get pods

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Deployment completed!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Access points:" -ForegroundColor Yellow
Write-Host "  - API Gateway: http://localhost:8080" -ForegroundColor White
Write-Host "  - Jaeger UI: http://localhost:16686" -ForegroundColor White
Write-Host ""
Write-Host "Useful commands:" -ForegroundColor Yellow
Write-Host "  - kubectl get pods              # Check pod status" -ForegroundColor Gray
Write-Host "  - kubectl get svc               # Check services" -ForegroundColor Gray
Write-Host "  - kubectl logs <pod-name>       # View logs" -ForegroundColor Gray
Write-Host "  - kubectl describe pod <name>   # Debug pod" -ForegroundColor Gray
