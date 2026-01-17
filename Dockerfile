FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 构建参数：服务名称
ARG SERVICE_NAME

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
    -ldflags '-w -s' \
    -o /app/service ./cmd/${SERVICE_NAME}

# 最终镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/service .

# 复制配置文件
COPY config/config.yaml /root/config/config.yaml

EXPOSE 8080

# 启动服务
CMD ["./service"]
