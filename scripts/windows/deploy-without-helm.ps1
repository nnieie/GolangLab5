# PowerShell Script for Full Deployment to Kind Cluster Without Helm
# Usage: .\deploy-without-helm.ps1

$ErrorActionPreference = "Stop"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Resolve-Path (Join-Path $scriptDir "../..")
Set-Location $repoRoot

Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "GolangLab5 Kind Deployment Without Helm" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Cyan

$clusterName = "lab5-cluster"
$requiredEnvVars = @("R2_Endpoint", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY")
$services = @("api", "user", "video", "social", "interaction", "chat")

function Assert-RequiredEnvVars {
    foreach ($envVar in $requiredEnvVars) {
        $value = (Get-Item "Env:$envVar" -ErrorAction SilentlyContinue).Value
        if ([string]::IsNullOrWhiteSpace($value)) {
            Write-Host "Missing required environment variable: $envVar" -ForegroundColor Red
            exit 1
        }
    }
}

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
foreach ($service in $services) {
    Write-Host "Loading golanglab5-${service}:v1.0..." -ForegroundColor White
    kind load docker-image "golanglab5-${service}:v1.0" --name $clusterName
    if ($LASTEXITCODE -ne 0) {
        Write-Host "✗ Failed to load image golanglab5-${service}:v1.0" -ForegroundColor Red
        exit 1
    }
}

# 4. 准备 ConfigMap 和 Secret
Write-Host ""
Write-Host "Step 4: Applying ConfigMaps and Secrets..." -ForegroundColor Yellow
Assert-RequiredEnvVars

kubectl create configmap app-config --from-file=config.k8s.yaml=config/config.k8s.yaml --dry-run=client -o yaml | kubectl apply -f -
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to apply ConfigMap app-config" -ForegroundColor Red
    exit 1
}

kubectl create configmap mysql-init-sql --from-file=init.sql=config/sql/init.sql --dry-run=client -o yaml | kubectl apply -f -
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to apply ConfigMap mysql-init-sql" -ForegroundColor Red
    exit 1
}

kubectl create secret generic r2-config `
    --from-literal="endpoint=$($env:R2_Endpoint)" `
    --from-literal="accessKeyId=$($env:R2_ACCESS_KEY_ID)" `
    --from-literal="secretAccessKey=$($env:R2_SECRET_ACCESS_KEY)" `
    --dry-run=client -o yaml | kubectl apply -f -
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to apply Secret r2-config" -ForegroundColor Red
    exit 1
}

# 5. 部署中间件和服务
Write-Host ""
Write-Host "Step 5: Applying Kubernetes manifests..." -ForegroundColor Yellow
kubectl apply -f k8s-manifests/middleware.yaml
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to apply middleware manifests" -ForegroundColor Red
    exit 1
}

kubectl apply -f k8s-manifests/services.yaml
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to apply service manifests" -ForegroundColor Red
    exit 1
}

# 6. 等待中间件就绪
Write-Host ""
Write-Host "Step 6: Waiting for middleware to be ready (this may take a few minutes)..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s
kubectl wait --for=condition=ready pod -l app=etcd --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
kubectl wait --for=condition=ready pod -l app=kafka --timeout=300s
kubectl wait --for=condition=ready pod -l app=otel-collector --timeout=300s
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Middleware pods did not become ready in time" -ForegroundColor Red
    exit 1
}

# 7. 同步数据库结构
Write-Host ""
Write-Host "Step 7: Syncing MySQL schema..." -ForegroundColor Yellow
$mysqlPod = kubectl get pod -l app=mysql -o jsonpath='{.items[0].metadata.name}'
if ([string]::IsNullOrWhiteSpace($mysqlPod)) {
    Write-Host "✗ Failed to find MySQL pod" -ForegroundColor Red
    exit 1
}

kubectl exec $mysqlPod -- sh -lc 'export MYSQL_PWD="$MYSQL_ROOT_PASSWORD"; mysql -uroot "$MYSQL_DATABASE" < /docker-entrypoint-initdb.d/init.sql'
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to sync MySQL schema" -ForegroundColor Red
    exit 1
}

# 8. 等待服务就绪
Write-Host ""
Write-Host "Step 8: Waiting for services to be ready..." -ForegroundColor Yellow
kubectl rollout status deployment/user-service --timeout=300s
kubectl rollout status deployment/video-service --timeout=300s
kubectl rollout status deployment/social-service --timeout=300s
kubectl rollout status deployment/interaction-service --timeout=300s
kubectl rollout status deployment/chat-service --timeout=300s
kubectl rollout status deployment/api-gateway --timeout=300s
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ One or more service deployments failed to become ready" -ForegroundColor Red
    exit 1
}

Write-Host ""
kubectl get pods

Write-Host ""
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "Deployment without Helm completed!" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Access points:" -ForegroundColor Yellow
Write-Host "  - API Gateway: http://localhost:8080" -ForegroundColor White
Write-Host "  - Jaeger UI: http://localhost:16686" -ForegroundColor White
