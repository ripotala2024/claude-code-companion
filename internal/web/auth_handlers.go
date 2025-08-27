package web

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// handleLoginPage 显示登录页面
func (s *AdminServer) handleLoginPage(c *gin.Context) {
	// 如果身份验证未启用，重定向到首页
	if !s.authManager.IsEnabled() {
		c.Redirect(http.StatusFound, "/admin/")
		return
	}

	// 如果用户已经登录，重定向到目标页面或首页
	if _, loggedIn := s.authManager.GetCurrentUser(c); loggedIn {
		redirectTo := c.Query("redirect_to")
		if redirectTo == "" {
			redirectTo = "/admin/"
		}
		c.Redirect(http.StatusFound, redirectTo)
		return
	}

	// 获取重定向目标
	redirectTo := c.Query("redirect_to")
	
	// 获取错误信息
	errorMsg := c.Query("error")
	errorKey := c.Query("error_key")

	// 生成CSRF token
	csrfToken := s.csrfManager.GenerateToken()

	data := gin.H{
		"Title":       "管理员登录 - Claude Code 伴侣",
		"CSRFToken":   csrfToken,
		"RedirectTo":  redirectTo,
		"Error":       errorMsg,
		"ErrorKey":    errorKey,
		"CurrentPage": "login",
	}

	c.HTML(http.StatusOK, "login.html", data)
}

// handleLoginSubmit 处理登录表单提交
func (s *AdminServer) handleLoginSubmit(c *gin.Context) {
	// 如果身份验证未启用，重定向到首页
	if !s.authManager.IsEnabled() {
		c.Redirect(http.StatusFound, "/admin/")
		return
	}

	// 验证CSRF token
	csrfToken := c.PostForm("_csrf_token")
	if !s.csrfManager.ValidateToken(csrfToken) {
		s.redirectWithError(c, "CSRF token invalid", "csrf_invalid")
		return
	}

	// 获取表单数据
	username := c.PostForm("username")
	password := c.PostForm("password")
	redirectTo := c.PostForm("redirect_to")

	// 验证输入
	if username == "" || password == "" {
		s.redirectWithError(c, "用户名和密码不能为空", "empty_credentials")
		return
	}

	// 验证凭据
	if !s.authManager.ValidateCredentials(username, password) {
		s.redirectWithError(c, "用户名或密码错误", "invalid_credentials")
		return
	}

	// 创建会话
	session := s.authManager.CreateSession(username)
	s.authManager.SetSessionCookie(c, session.ID)

	// 重定向到目标页面
	if redirectTo == "" {
		redirectTo = "/admin/"
	}

	// 验证重定向URL的安全性
	if !s.isValidRedirectURL(redirectTo) {
		redirectTo = "/admin/"
	}

	c.Redirect(http.StatusFound, redirectTo)
}

// handleLogout 处理登出
func (s *AdminServer) handleLogout(c *gin.Context) {
	// 清除会话
	s.authManager.ClearSession(c)

	// 重定向到登录页面
	c.Redirect(http.StatusFound, "/admin/login")
}

// redirectWithError 重定向到登录页面并显示错误信息
func (s *AdminServer) redirectWithError(c *gin.Context, errorMsg, errorKey string) {
	redirectTo := c.PostForm("redirect_to")
	
	// 构建登录URL with error parameters
	loginURL := "/admin/login"
	params := url.Values{}
	
	if errorMsg != "" {
		params.Add("error", errorMsg)
	}
	if errorKey != "" {
		params.Add("error_key", errorKey)
	}
	if redirectTo != "" {
		params.Add("redirect_to", redirectTo)
	}
	
	if len(params) > 0 {
		loginURL += "?" + params.Encode()
	}

	c.Redirect(http.StatusFound, loginURL)
}

// isValidRedirectURL 验证重定向URL的安全性
func (s *AdminServer) isValidRedirectURL(redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	// 解析URL
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// 只允许相对路径或同域名的URL
	if parsedURL.IsAbs() {
		return false
	}

	// 确保路径以/admin开头
	if !isAdminPath(parsedURL.Path) {
		return false
	}

	return true
}

// isAdminPath 检查路径是否为管理路径
func isAdminPath(path string) bool {
	return path == "/admin" || path == "/admin/" || 
		   (len(path) > 6 && path[:7] == "/admin/")
}

// handleGetCurrentUser API端点：获取当前登录用户信息
func (s *AdminServer) handleGetCurrentUser(c *gin.Context) {
	if !s.authManager.IsEnabled() {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"auth_enabled":  false,
		})
		return
	}

	user, loggedIn := s.authManager.GetCurrentUser(c)
	if !loggedIn {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"auth_enabled":  true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"auth_enabled":  true,
		"user": gin.H{
			"username":    user.Username,
			"created_at":  user.CreatedAt,
			"last_access": user.LastAccess,
		},
	})
}

// handleAuthStatus API端点：获取身份验证状态
func (s *AdminServer) handleAuthStatus(c *gin.Context) {
	status := gin.H{
		"auth_enabled": s.authManager.IsEnabled(),
	}

	if s.authManager.IsEnabled() {
		_, authenticated := s.authManager.GetCurrentUser(c)
		status["authenticated"] = authenticated
	} else {
		status["authenticated"] = true // 如果身份验证未启用，视为已认证
	}

	c.JSON(http.StatusOK, status)
}
