# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Claude Code Companion 是一个多协议 API 代理服务，为 Claude Code 等客户端提供统一的 API 访问入口。主要实现多端点负载均衡、格式转换、智能路由等功能。

## 开发命令

### 构建与运行
```bash
# 构建项目
make build

# 运行测试
make test

# 格式化代码
make fmt

# 代码检查（需要 golangci-lint）
make lint

# 开发模式（需要 air 工具）
make dev

# 构建所有平台版本
make all

# 清理构建产物
make clean
```

### 运行服务
```bash
# 使用默认配置运行
make run

# 直接运行二进制文件
./claude-code-companion -config config.yaml

# 指定端口运行
./claude-code-companion --port 8080
```

### Docker 部署
```bash
# 快速启动
./docker/docker-quick-start.sh

# 手动启动
cd docker && docker-compose up -d
```

## 项目架构

### 核心组件层次
- **main.go** - 程序入口和命令行参数处理
- **internal/proxy/** - HTTP 代理服务核心逻辑
- **internal/endpoint/** - 端点管理和选择器
- **internal/conversion/** - OpenAI/Anthropic 格式转换
- **internal/config/** - 配置管理和验证
- **internal/web/** - Web 管理界面
- **internal/oauth/** - OAuth 认证管理
- **internal/tagging/** - 智能标签路由
- **internal/statistics/** - 使用统计和监控

### 关键设计模式
- **端点管理器 (endpoint.Manager)**: 管理多个上游 API 端点，支持优先级路由和故障转移
- **格式转换器 (conversion.Converter)**: 实现 OpenAI 与 Anthropic API 格式的双向转换
- **标签路由系统**: 基于请求特征的动态路由，支持 Starlark 脚本扩展
- **健康检查机制**: 定期检测端点状态，自动剔除和恢复故障端点

### 配置系统
- 配置文件位置：`config.yaml`
- 端点配置示例：`config.yaml.example`
- 端点预设配置：`internal/config/endpoint_profiles.yaml`
- 支持运行时配置热更新（通过 Web 界面）

### Web 管理界面
- 访问地址：`http://localhost:8080/admin/`
- 模板位置：`web/templates/`
- 静态资源：`web/static/`
- 多语言支持：`web/locales/`

## 关键实现细节

### 格式转换机制
项目实现了完整的 OpenAI 到 Anthropic 格式转换：
- 请求转换：`internal/conversion/request_converter.go`
- 响应转换：`internal/conversion/response_converter.go`
- 流式响应处理：`internal/conversion/sse_parser.go`

### 智能路由系统
支持基于请求内容的动态路由：
- 内置标签器：`internal/taggers/builtin/`
- Starlark 脚本支持：`internal/taggers/starlark/`
- 路由管道：`internal/tagging/pipeline.go`

### 统计和监控
- 内存统计：`internal/statistics/memory_manager.go`
- GORM 持久化：`internal/logger/gorm_*.go`
- 使用统计分析工具：`tools/token_usage_analyzer.py`

## 测试策略

### 测试文件位置
- 单元测试：各模块 `*_test.go` 文件
- 集成测试：`internal/conversion/integration_test.go`
- 性能测试：`internal/logger/storage_benchmark_test.go`

### 测试覆盖重点
- 格式转换的正确性和完整性
- 端点选择和故障转移逻辑
- 配置加载和验证机制
- Web API 的功能完整性
