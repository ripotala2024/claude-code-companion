# 多阶段构建 - 构建阶段  
FROM golang:1.23-alpine AS builder

# 安装构建依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git make ca-certificates

# 设置工作目录
WORKDIR /app

# 设置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用 - 静态链接，适合容器环境
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=container-$(date +%Y%m%d-%H%M%S) -w -s" \
    -a -installsuffix cgo \
    -o claude-code-companion .

# 运行时阶段 - 使用Alpine作为基础镜像
FROM alpine:3.18

# 安装运行时依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add \
    ca-certificates \
    curl \
    tzdata \
    && rm -rf /var/cache/apk/*

# 创建非root用户和组
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/claude-code-companion .

# 复制配置模板和启动脚本
COPY config.yaml.example ./config.yaml.example
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh

# 创建数据目录并设置权限
RUN mkdir -p /data/logs /data/config && \
    chmod +x /docker-entrypoint.sh && \
    chown -R appuser:appgroup /app /data && \
    chmod -R 755 /data

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

# 设置环境变量
ENV CONFIG_FILE=/data/config/config.yaml
ENV LOG_DIR=/data/logs
ENV GIN_MODE=release

# 数据卷
VOLUME ["/data/config", "/data/logs"]

# 启动命令 - 将通过entrypoint脚本处理
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD []
