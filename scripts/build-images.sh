#!/bin/bash
set -e

echo "========================================="
echo "Building GolangLab5 Docker Images"
echo "========================================="

# 服务列表
services=("api" "user" "video" "social" "interaction" "chat")

# 构建每个服务镜像
for service in "${services[@]}"; do
    echo ""
    echo "Building image for service: $service"
    echo "-----------------------------------"
    
    docker build -t golanglab5-$service:v1.0 \
        --build-arg SERVICE_NAME=$service \
        -f Dockerfile.template .
    
    if [ $? -eq 0 ]; then
        echo "✓ Successfully built golanglab5-$service:v1.0"
    else
        echo "✗ Failed to build golanglab5-$service:v1.0"
        exit 1
    fi
done

echo ""
echo "========================================="
echo "All images built successfully!"
echo "========================================="
echo ""
echo "Images created:"
for service in "${services[@]}"; do
    echo "  - golanglab5-$service:v1.0"
done
