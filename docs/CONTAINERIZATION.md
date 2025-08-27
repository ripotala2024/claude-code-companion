# Claude Code Companion 容器化部署指南

本文档详细介绍了如何使用Docker和Kubernetes部署Claude Code Companion。

## 概述

Claude Code Companion已完全支持容器化部署，提供以下部署方式：

- **Docker**: 单容器部署，适合开发和小规模生产环境
- **Docker Compose**: 完整的容器编排，包含数据持久化

## 前置要求

### Docker部署
- Docker 20.10+
- Docker Compose 2.0+

## Docker部署

### 1. 构建镜像

```bash
# 克隆项目
git clone <repository-url>
cd claude-code-companion

# 构建Docker镜像
docker build -f docker/Dockerfile -t claude-code-companion:latest .
```

### 2. 运行容器

```bash
# 创建数据目录
mkdir -p ./data/config ./data/logs

# 运行容器
docker run -d \
  --name claude-code-companion \
  -p 8080:8080 \
  -v $(pwd)/data/config:/data/config \
  -v $(pwd)/data/logs:/data/logs \
  -e CONFIG_FILE=/data/config/config.yaml \
  -e LOG_DIR=/data/logs \
  claude-code-companion:latest
```

### 3. 访问应用

- 代理服务: http://localhost:8080
- 管理界面: http://localhost:8080/admin/

## Docker Compose部署

### 1. 启动服务

```bash
# 使用Docker Compose启动
cd docker
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2. 配置管理

首次启动后，配置文件会自动创建在 `./data/config/config.yaml`。

编辑配置文件后重启服务：
```bash
cd docker
docker-compose restart
```



## 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| CONFIG_FILE | /data/config/config.yaml | 配置文件路径 |
| LOG_DIR | /data/logs | 日志目录路径 |
| GIN_MODE | release | Gin框架运行模式 |
| TZ | Asia/Shanghai | 时区设置 |

### 数据持久化

容器化部署中，以下数据需要持久化：

- **配置文件**: `/data/config/config.yaml`
- **日志数据**: `/data/logs/` (包含SQLite数据库)

### 端口说明

- **8080**: HTTP服务端口，提供API代理和Web管理界面

## 故障排除

### 常见问题

1. **容器启动失败**
   ```bash
   # 查看容器日志
   docker logs claude-code-companion
   
   # 或在Kubernetes中
   kubectl logs -n claude-code-companion deployment/claude-code-companion
   ```

2. **配置文件问题**
   ```bash
   # 检查配置文件权限
   ls -la ./data/config/
   
   # 重置配置文件
   rm ./data/config/config.yaml
   cd docker && docker-compose restart
   ```

3. **数据库锁定问题**
   ```bash
   # 检查SQLite数据库
   ls -la ./data/logs/logs.db*
   
   # 如果需要，删除数据库文件重新初始化
   rm ./data/logs/logs.db*
   ```

### 健康检查

容器提供健康检查端点：
```bash
# 检查应用健康状态
curl http://localhost:8080/

# 检查管理界面
curl http://localhost:8080/admin/
```

### 性能调优

1. **资源限制**
   - 内存: 建议64Mi-256Mi
   - CPU: 建议50m-200m

2. **SQLite优化**
   - 使用SSD存储提升I/O性能
   - 定期清理日志数据

## 安全建议

1. **运行用户**: 容器使用非root用户(uid:1001)运行
2. **网络隔离**: 使用Docker网络或Kubernetes网络策略
3. **敏感配置**: 使用Docker Secrets或Kubernetes Secrets管理API密钥
4. **镜像安全**: 定期更新基础镜像，扫描安全漏洞

## 监控和日志

### 日志收集
```bash
# Docker环境
docker logs claude-code-companion

# Kubernetes环境
kubectl logs -n claude-code-companion -l app=claude-code-companion
```

### 监控指标
- 容器资源使用率
- HTTP响应时间
- 错误率统计
- 数据库大小

## 升级指南

### Docker升级
```bash
# 拉取新镜像
docker pull claude-code-companion:latest

# 停止旧容器
docker stop claude-code-companion
docker rm claude-code-companion

# 启动新容器
docker run -d \
  --name claude-code-companion \
  -p 8080:8080 \
  -v $(pwd)/data/config:/data/config \
  -v $(pwd)/data/logs:/data/logs \
  claude-code-companion:latest
```


