# PowerShell Script for Building GolangLab5 Docker Images
# Usage: .\build-images.ps1

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Building GolangLab5 Docker Images" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 服务列表
$services = @("api", "user", "video", "social", "interaction", "chat")

# 构建每个服务镜像
foreach ($service in $services) {
    Write-Host ""
    Write-Host "Building image for service: $service" -ForegroundColor Yellow
    Write-Host "-----------------------------------" -ForegroundColor Yellow
    
    docker build -t "golanglab5-${service}:v1.0" `
        --build-arg SERVICE_NAME=$service `
        -f Dockerfile.template .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Successfully built golanglab5-${service}:v1.0" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to build golanglab5-${service}:v1.0" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "All images built successfully!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Images created:" -ForegroundColor Yellow
foreach ($service in $services) {
    Write-Host "  - golanglab5-${service}:v1.0" -ForegroundColor White
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Load images to Kind cluster:" -ForegroundColor White
Write-Host "     .\load-images-to-kind.ps1" -ForegroundColor Gray
Write-Host "  2. Or run full deployment:" -ForegroundColor White
Write-Host "     .\deploy-all.ps1" -ForegroundColor Gray
