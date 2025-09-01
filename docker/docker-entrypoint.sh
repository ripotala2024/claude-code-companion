#!/bin/sh
set -e

# 打印启动信息
echo "=== Claude Code Companion Container Starting ==="
echo "Config file: ${CONFIG_FILE:-/data/config/config.yaml}"
echo "Log directory: ${LOG_DIR:-/data/logs}"

# 设置默认值
CONFIG_FILE=${CONFIG_FILE:-/data/config/config.yaml}
LOG_DIR=${LOG_DIR:-/data/logs}

# 确保配置目录存在并设置正确权限
CONFIG_DIR=$(dirname "$CONFIG_FILE")
mkdir -p "$CONFIG_DIR"
chmod 755 "$CONFIG_DIR"

# 如果配置文件不存在，从模板创建
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Creating default config file at $CONFIG_FILE"
    if [ -f "/app/config.yaml.example" ]; then
        cp /app/config.yaml.example "$CONFIG_FILE"
        echo "Default configuration created from template"
    else
        echo "Warning: No config template found, application will create default config"
    fi
fi

# 确保配置文件有正确的权限
if [ -f "$CONFIG_FILE" ]; then
    chmod 644 "$CONFIG_FILE"
fi

# 确保日志目录存在并设置正确权限
mkdir -p "$LOG_DIR"
chmod 755 "$LOG_DIR"

# 检查配置文件权限
if [ ! -r "$CONFIG_FILE" ]; then
    echo "Error: Cannot read config file $CONFIG_FILE"
    exit 1
fi

# 打印配置信息
echo "Configuration file: $CONFIG_FILE"
echo "Log directory: $LOG_DIR"
echo "Starting Claude Code Companion..."

# 构建启动参数
ARGS="-config $CONFIG_FILE"

# 如果有额外的命令行参数，添加到启动参数中
if [ $# -gt 0 ]; then
    ARGS="$ARGS $@"
fi

# 启动应用
echo "Executing: /app/claude-code-companion $ARGS"
exec /app/claude-code-companion $ARGS
