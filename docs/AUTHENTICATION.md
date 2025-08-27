# 身份验证功能使用指南

Claude Code Companion 现在支持管理界面的身份验证功能，确保只有授权用户才能访问管理后台。

## 功能特性

- **环境变量配置**：通过 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD` 环境变量配置管理员凭据
- **会话管理**：安全的会话管理，支持自动过期和清理
- **登录界面**：美观的响应式登录页面
- **智能重定向**：登录成功后自动重定向到原始请求页面
- **用户体验**：导航栏显示用户信息和登出功能
- **安全防护**：密码哈希存储、CSRF保护、常量时间比较防止时序攻击

## 快速开始

### 1. 启用身份验证

通过环境变量设置管理员凭据：

```bash
# 设置管理员用户名和密码
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your_secure_password

# 启动服务
./claude-code-companion -config config.yaml
```

或者在启动时直接指定：

```bash
ADMIN_USERNAME=admin ADMIN_PASSWORD=your_secure_password ./claude-code-companion -config config.yaml
```

### 2. 访问管理界面

1. 打开浏览器访问：`http://localhost:8080/admin/`
2. 系统会自动重定向到登录页面：`http://localhost:8080/admin/login`
3. 输入设置的用户名和密码
4. 登录成功后会重定向到管理界面

### 3. 登出

在管理界面的导航栏右上角，点击用户下拉菜单中的"登出"按钮。

## 配置选项

### 环境变量

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `ADMIN_USERNAME` | 管理员用户名 | `admin` | `admin` |
| `ADMIN_PASSWORD` | 管理员密码 | 无 | `your_secure_password` |

### 配置文件

身份验证配置也可以在 `config.yaml` 中设置：

```yaml
auth:
  enabled: true                    # 是否启用身份验证
  username: "admin"               # 管理员用户名
  password: ""                    # 管理员密码（建议使用环境变量）
  session_timeout: "24h"          # 会话超时时间
```

**注意**：
- 环境变量的优先级高于配置文件
- 如果设置了 `ADMIN_USERNAME` 或 `ADMIN_PASSWORD` 环境变量，身份验证会自动启用
- 密码建议通过环境变量设置，避免在配置文件中明文存储

## 安全特性

### 密码安全
- 密码使用 bcrypt 算法进行哈希存储
- 支持常量时间字符串比较，防止时序攻击
- 密码不会在日志中记录

### 会话安全
- 使用安全的随机会话ID
- 会话cookie设置了 HttpOnly、Secure、SameSite 属性
- 支持会话自动过期和清理
- 默认会话超时时间为24小时

### CSRF保护
- 所有管理API请求都受到CSRF保护
- 登录表单包含CSRF token验证

### 路径保护
- 所有 `/admin/*` 路径（除登录相关）都需要身份验证
- API端点 `/admin/api/*` 也受到保护
- 静态资源和公开页面不受影响

## 故障排除

### 常见问题

1. **忘记密码**
   - 停止服务器
   - 重新设置 `ADMIN_PASSWORD` 环境变量
   - 重启服务器

2. **无法访问管理界面**
   - 检查是否正确设置了环境变量
   - 确认用户名和密码正确
   - 查看服务器日志是否有错误信息

3. **会话过期**
   - 重新登录即可
   - 可以通过配置文件调整 `session_timeout` 设置

4. **CSRF错误**
   - 刷新页面重新获取CSRF token
   - 确保浏览器启用了cookie

### 禁用身份验证

如果需要临时禁用身份验证：

```bash
# 不设置环境变量启动
./claude-code-companion -config config.yaml
```

或者在配置文件中设置：

```yaml
auth:
  enabled: false
```

## 最佳实践

1. **使用强密码**：建议使用包含大小写字母、数字和特殊字符的复杂密码
2. **定期更换密码**：定期更新 `ADMIN_PASSWORD` 环境变量
3. **安全部署**：在生产环境中使用HTTPS确保传输安全
4. **监控访问**：定期检查访问日志，监控异常登录行为
5. **备份配置**：确保配置文件和环境变量设置有适当的备份

## 开发和测试

在开发环境中，可以使用简单的凭据进行测试：

```bash
# 开发环境示例
ADMIN_USERNAME=admin ADMIN_PASSWORD=test123 ./claude-code-companion -config config.yaml
```

在生产环境中，请务必使用强密码并通过安全的方式管理环境变量。
