# Claude Code Companion Docker 部署

本目录包含 Claude Code Companion 的 Docker 容器化部署文件。

## 文件说明

- **`Dockerfile`**: 多阶段构建的 Docker 镜像定义
- **`docker-entrypoint.sh`**: 容器启动脚本，处理配置初始化
- **`docker-compose.yml`**: Docker Compose 服务编排配置
- **`docker-quick-start.sh`**: 一键启动脚本

## 快速开始

### 方法1: 使用一键启动脚本

```bash
# 在项目根目录执行
./docker/docker-quick-start.sh
```

### 方法2: 手动启动

```bash
# 进入 docker 目录
cd docker

# 构建并启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方法3: 仅使用 Docker

```bash
# 构建镜像
docker build -f docker/Dockerfile -t claude-code-companion:latest .

# 创建数据目录
mkdir -p ./data/config ./data/logs

# 运行容器
docker run -d \
  --name claude-code-companion \
  -p 8080:8080 \
  -v $(pwd)/data/config:/data/config \
  -v $(pwd)/data/logs:/data/logs \
  claude-code-companion:latest
```

## 访问应用

启动成功后，可以通过以下地址访问：

- **代理服务**: http://localhost:8080
- **管理界面**: http://localhost:8080/admin/

## 数据持久化

容器使用以下目录进行数据持久化：

- **配置文件**: `./data/config/config.yaml`
- **日志数据**: `./data/logs/` (包含 SQLite 数据库)

## 配置管理

首次启动时，容器会自动从模板创建默认配置文件。您可以：

1. 编辑 `./data/config/config.yaml` 文件
2. 重启服务使配置生效：`docker-compose restart`

## 故障排除

### 查看日志
```bash
cd docker
docker-compose logs -f
```

### 重置配置
```bash
rm ./data/config/config.yaml
cd docker && docker-compose restart
```

### 清理数据
```bash
rm -rf ./data/
cd docker && docker-compose down && docker-compose up -d
```

## 技术特性

- ✅ 多阶段构建，优化镜像大小
- ✅ 非 root 用户运行，提升安全性
- ✅ 健康检查机制
- ✅ 自动配置初始化
- ✅ 数据持久化
- ✅ 资源限制配置

更多详细信息请参考 [容器化部署文档](../docs/CONTAINERIZATION.md)。
