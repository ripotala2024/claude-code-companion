# 快速启动指南

## 使用 startup.sh 脚本

我们提供了一个便捷的启动脚本 `startup.sh`，让您可以快速测试身份验证功能。

### 基本用法

```bash
# 使用默认设置启动（用户名: admin, 密码: test123）
./startup.sh

# 查看帮助信息
./startup.sh -h
```

### 常用选项

```bash
# 自定义用户名和密码
./startup.sh -u myuser -p mypassword

# 禁用身份验证（开发调试用）
./startup.sh -n

# 先编译再运行
./startup.sh -b

# 使用不同端口
./startup.sh -P 9090

# 组合使用
./startup.sh -u admin -p secret123 -P 8081 -b
```

### 完整选项列表

| 选项 | 长选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `-u` | `--username` | 管理员用户名 | `admin` |
| `-p` | `--password` | 管理员密码 | `test123` |
| `-P` | `--port` | 服务端口 | `8080` |
| `-c` | `--config` | 配置文件路径 | `config.yaml` |
| `-n` | `--no-auth` | 禁用身份验证 | - |
| `-b` | `--build` | 先编译再运行 | - |
| `-h` | `--help` | 显示帮助信息 | - |

### 测试场景

#### 1. 测试身份验证功能
```bash
./startup.sh -u admin -p test123
```
然后访问 http://localhost:8080/admin/ 测试登录功能。

#### 2. 开发调试（无身份验证）
```bash
./startup.sh -n
```
直接访问管理界面，无需登录。

#### 3. 生产环境模拟
```bash
./startup.sh -u admin -p "StrongPassword123!" -P 8080 -b
```
使用强密码和编译后的二进制文件。

### 启动后的访问地址

- **管理界面**: http://localhost:8080/admin/
- **登录页面**: http://localhost:8080/admin/login （启用身份验证时）
- **API代理**: http://localhost:8080/v1/

### 停止服务

按 `Ctrl+C` 停止服务器。

### 故障排除

1. **权限错误**：
   ```bash
   chmod +x startup.sh
   ```

2. **端口被占用**：
   ```bash
   ./startup.sh -P 8081  # 使用其他端口
   ```

3. **Go环境问题**：
   确保已安装Go并设置了正确的环境变量。

4. **配置文件问题**：
   首次运行会自动生成默认配置文件。

### 快速测试流程

1. 启动服务：
   ```bash
   ./startup.sh
   ```

2. 打开浏览器访问：http://localhost:8080/admin/

3. 使用默认凭据登录：
   - 用户名：`admin`
   - 密码：`test123`

4. 测试各项功能：
   - 查看用户信息下拉菜单
   - 测试登出功能
   - 验证未登录时的重定向

这样您就可以快速验证身份验证功能是否正常工作了！
