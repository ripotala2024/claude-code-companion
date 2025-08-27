#!/bin/bash

# Claude Code Companion 开发测试启动脚本
# 用于快速启动和测试身份验证功能

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 默认配置
DEFAULT_USERNAME="admin"
DEFAULT_PASSWORD="test123"
DEFAULT_PORT="8080"
CONFIG_FILE="config.yaml"

# 显示帮助信息
show_help() {
    echo -e "${CYAN}Claude Code Companion 开发测试启动脚本${NC}"
    echo ""
    echo -e "${YELLOW}用法:${NC}"
    echo "  $0 [选项]"
    echo ""
    echo -e "${YELLOW}选项:${NC}"
    echo "  -u, --username USERNAME    设置管理员用户名 (默认: admin)"
    echo "  -p, --password PASSWORD    设置管理员密码 (默认: test123)"
    echo "  -P, --port PORT           设置服务端口 (默认: 8080)"
    echo "  -c, --config CONFIG       指定配置文件 (默认: config.yaml)"
    echo "  -n, --no-auth            禁用身份验证"
    echo "  -b, --build              先编译再运行"
    echo "  -h, --help               显示此帮助信息"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0                                    # 使用默认设置启动"
    echo "  $0 -u myuser -p mypass               # 自定义用户名密码"
    echo "  $0 -n                                # 禁用身份验证"
    echo "  $0 -b                                # 编译后运行"
    echo "  $0 -P 9090                           # 使用端口9090"
}

# 解析命令行参数
USERNAME="$DEFAULT_USERNAME"
PASSWORD="$DEFAULT_PASSWORD"
PORT="$DEFAULT_PORT"
NO_AUTH=false
BUILD_FIRST=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--username)
            USERNAME="$2"
            shift 2
            ;;
        -p|--password)
            PASSWORD="$2"
            shift 2
            ;;
        -P|--port)
            PORT="$2"
            shift 2
            ;;
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -n|--no-auth)
            NO_AUTH=true
            shift
            ;;
        -b|--build)
            BUILD_FIRST=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}错误: 未知选项 $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 显示启动信息
echo -e "${CYAN}🚀 Claude Code Companion 开发测试启动${NC}"
echo -e "${PURPLE}================================================${NC}"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到Go环境，请先安装Go${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go环境检查通过${NC}"

# 检查项目文件
if [ ! -f "main.go" ]; then
    echo -e "${RED}❌ 错误: 未找到main.go文件，请确保在项目根目录运行此脚本${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 项目文件检查通过${NC}"

# 编译项目（如果需要）
if [ "$BUILD_FIRST" = true ]; then
    echo -e "${YELLOW}🔨 正在编译项目...${NC}"
    if go build -o claude-code-companion .; then
        echo -e "${GREEN}✅ 编译成功${NC}"
    else
        echo -e "${RED}❌ 编译失败${NC}"
        exit 1
    fi
fi

# 设置环境变量
if [ "$NO_AUTH" = false ]; then
    export ADMIN_USERNAME="$USERNAME"
    export ADMIN_PASSWORD="$PASSWORD"
    echo -e "${GREEN}🔐 身份验证已启用${NC}"
    echo -e "   用户名: ${YELLOW}$USERNAME${NC}"
    echo -e "   密码: ${YELLOW}$PASSWORD${NC}"
else
    unset ADMIN_USERNAME
    unset ADMIN_PASSWORD
    echo -e "${YELLOW}⚠️  身份验证已禁用${NC}"
fi

# 显示访问信息
echo -e "${PURPLE}================================================${NC}"
echo -e "${CYAN}📡 服务信息${NC}"
echo -e "   配置文件: ${YELLOW}$CONFIG_FILE${NC}"
echo -e "   服务端口: ${YELLOW}$PORT${NC}"
echo -e "   管理界面: ${BLUE}http://localhost:$PORT/admin/${NC}"
echo -e "   API代理: ${BLUE}http://localhost:$PORT/v1/${NC}"

if [ "$NO_AUTH" = false ]; then
    echo -e "   登录页面: ${BLUE}http://localhost:$PORT/admin/login${NC}"
fi

echo -e "${PURPLE}================================================${NC}"

# 直接启动，不等待用户确认

# 启动服务器
echo -e "${GREEN}🚀 正在启动服务器...${NC}"
echo ""

if [ "$BUILD_FIRST" = true ]; then
    # 使用编译后的二进制文件
    ./claude-code-companion --port "$PORT" -config "$CONFIG_FILE"
else
    # 直接运行Go代码
    go run . --port "$PORT" -config "$CONFIG_FILE"
fi
