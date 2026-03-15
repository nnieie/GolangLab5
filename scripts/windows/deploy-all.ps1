# PowerShell Script for Full Deployment to Kind Cluster
# Usage: .\deploy-all.ps1

$ErrorActionPreference = "Stop"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Resolve-Path (Join-Path $scriptDir "../..")
Set-Location $repoRoot

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "GolangLab5 Kind Kubernetes Deployment" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

$clusterName = "lab5-cluster"
$releaseName = "golanglab5"
$requiredEnvVars = @("R2_Endpoint", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY")

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
& (Join-Path $scriptDir "build-images.ps1")
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

# 4. 使用 Helm 部署
Write-Host ""
Write-Host "Step 4: Deploying with Helm..." -ForegroundColor Yellow
foreach ($envVar in $requiredEnvVars) {
    if ([string]::IsNullOrWhiteSpace((Get-Item "Env:$envVar" -ErrorAction SilentlyContinue).Value)) {
        Write-Host "Missing required environment variable: $envVar" -ForegroundColor Red
        exit 1
    }
}

helm upgrade --install $releaseName ./chart `
    --set-file config.configK8sContent=config/config.k8s.yaml `
    --set-file config.mysqlInitSqlContent=config/sql/init.sql `
    --set-string r2.endpoint=$env:R2_Endpoint `
    --set-string r2.accessKeyId=$env:R2_ACCESS_KEY_ID `
    --set-string r2.secretAccessKey=$env:R2_SECRET_ACCESS_KEY

# 5. 等待中间件就绪
Write-Host ""
Write-Host "Step 5: Waiting for middleware to be ready (this may take a few minutes)..." -ForegroundColor Yellow
Write-Host "Waiting for MySQL..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
Write-Host "Waiting for Etcd..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
Write-Host "Waiting for Redis..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
Write-Host "Waiting for Kafka..." -ForegroundColor White
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s

# 6. 等待服务就绪
Write-Host ""
Write-Host "Step 6: Waiting for services to be ready..." -ForegroundColor Yellow
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
