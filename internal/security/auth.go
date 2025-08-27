package security

import (
	"crypto/subtle"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthManager 身份验证管理器
type AuthManager struct {
	sessionManager *SessionManager
	username       string
	passwordHash   string
	enabled        bool
}

// NewAuthManager 创建新的身份验证管理器
func NewAuthManager(sessionManager *SessionManager, username, password string, enabled bool) *AuthManager {
	var passwordHash string
	if password != "" {
		// 对密码进行哈希处理
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			// 如果哈希失败，记录错误但不阻止启动
			passwordHash = ""
		} else {
			passwordHash = string(hash)
		}
	}

	return &AuthManager{
		sessionManager: sessionManager,
		username:       username,
		passwordHash:   passwordHash,
		enabled:        enabled,
	}
}

// AuthMiddleware 返回身份验证中间件
func (am *AuthManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果身份验证未启用，直接通过
		if !am.enabled {
			c.Next()
			return
		}

		// 检查是否为公开路径
		if am.isPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 验证会话
		sessionID := am.sessionManager.GetSessionFromCookie(c)
		if am.sessionManager.ValidateSession(sessionID) {
			// 会话有效，继续处理请求
			c.Next()
			return
		}

		// 会话无效，重定向到登录页面
		am.redirectToLogin(c)
	}
}

// ValidateCredentials 验证用户凭据
func (am *AuthManager) ValidateCredentials(username, password string) bool {
	if !am.enabled {
		return true // 如果身份验证未启用，总是返回true
	}

	// 检查用户名
	if !am.constantTimeStringCompare(username, am.username) {
		return false
	}

	// 检查密码
	if am.passwordHash == "" {
		return false // 没有设置密码哈希
	}

	err := bcrypt.CompareHashAndPassword([]byte(am.passwordHash), []byte(password))
	return err == nil
}

// CreateSession 创建新会话
func (am *AuthManager) CreateSession(username string) *Session {
	return am.sessionManager.CreateSession(username)
}

// SetSessionCookie 设置会话cookie
func (am *AuthManager) SetSessionCookie(c *gin.Context, sessionID string) {
	am.sessionManager.SetSessionCookie(c, sessionID)
}

// ClearSession 清除会话
func (am *AuthManager) ClearSession(c *gin.Context) {
	sessionID := am.sessionManager.GetSessionFromCookie(c)
	if sessionID != "" {
		am.sessionManager.DeleteSession(sessionID)
	}
	am.sessionManager.ClearSessionCookie(c)
}

// GetCurrentUser 获取当前登录用户信息
func (am *AuthManager) GetCurrentUser(c *gin.Context) (*Session, bool) {
	if !am.enabled {
		return nil, false
	}

	sessionID := am.sessionManager.GetSessionFromCookie(c)
	return am.sessionManager.GetSession(sessionID)
}

// IsEnabled 检查身份验证是否启用
func (am *AuthManager) IsEnabled() bool {
	return am.enabled
}

// isPublicPath 检查路径是否为公开路径（不需要身份验证）
func (am *AuthManager) isPublicPath(path string) bool {
	// 精确匹配根路径
	if path == "/" {
		return true
	}

	publicPaths := []string{
		"/admin/login",
		"/admin/api/csrf-token",
		"/static/",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

// redirectToLogin 重定向到登录页面
func (am *AuthManager) redirectToLogin(c *gin.Context) {
	// 保存原始请求URL用于登录后重定向
	originalURL := c.Request.URL.String()

	// 如果是API请求，返回JSON错误
	if strings.HasPrefix(c.Request.URL.Path, "/admin/api/") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		c.Abort()
		return
	}

	// 构建登录URL，包含重定向参数
	loginURL := "/admin/login"
	if originalURL != "" && originalURL != "/admin/login" {
		loginURL += "?redirect_to=" + url.QueryEscape(originalURL)
	}

	c.Redirect(http.StatusFound, loginURL)
	c.Abort()
}

// constantTimeStringCompare 常量时间字符串比较，防止时序攻击
func (am *AuthManager) constantTimeStringCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// UpdateCredentials 更新认证凭据（用于运行时配置更新）
func (am *AuthManager) UpdateCredentials(username, password string, enabled bool) error {
	var passwordHash string
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		passwordHash = string(hash)
	}

	am.username = username
	am.passwordHash = passwordHash
	am.enabled = enabled

	return nil
}
