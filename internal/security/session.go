package security

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Session 表示一个用户会话
type Session struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	LastAccess time.Time `json:"last_access"`
}

// SessionManager 管理用户会话
type SessionManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	timeout  time.Duration
}

// NewSessionManager 创建新的会话管理器
func NewSessionManager(timeout time.Duration) *SessionManager {
	manager := &SessionManager{
		sessions: make(map[string]*Session),
		timeout:  timeout,
	}

	// 启动清理过期会话的goroutine
	go manager.cleanupExpiredSessions()

	return manager
}

// CreateSession 创建新会话
func (sm *SessionManager) CreateSession(username string) *Session {
	sessionID := sm.generateSessionID()
	
	session := &Session{
		ID:         sessionID,
		Username:   username,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
	}

	sm.mutex.Lock()
	sm.sessions[sessionID] = session
	sm.mutex.Unlock()

	return session
}

// ValidateSession 验证会话是否有效
func (sm *SessionManager) ValidateSession(sessionID string) bool {
	if sessionID == "" {
		return false
	}

	sm.mutex.RLock()
	session, exists := sm.sessions[sessionID]
	sm.mutex.RUnlock()

	if !exists {
		return false
	}

	// 检查会话是否过期
	if time.Since(session.LastAccess) > sm.timeout {
		sm.DeleteSession(sessionID)
		return false
	}

	// 更新最后访问时间
	sm.mutex.Lock()
	session.LastAccess = time.Now()
	sm.mutex.Unlock()

	return true
}

// GetSession 获取会话信息
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	if sessionID == "" {
		return nil, false
	}

	sm.mutex.RLock()
	session, exists := sm.sessions[sessionID]
	sm.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	// 检查会话是否过期
	if time.Since(session.LastAccess) > sm.timeout {
		sm.DeleteSession(sessionID)
		return nil, false
	}

	// 更新最后访问时间
	sm.mutex.Lock()
	session.LastAccess = time.Now()
	sm.mutex.Unlock()

	return session, true
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mutex.Lock()
	delete(sm.sessions, sessionID)
	sm.mutex.Unlock()
}

// generateSessionID 生成安全的会话ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// cleanupExpiredSessions 定期清理过期会话
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时清理一次
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		sm.mutex.Lock()
		for sessionID, session := range sm.sessions {
			if now.Sub(session.LastAccess) > sm.timeout {
				delete(sm.sessions, sessionID)
			}
		}
		sm.mutex.Unlock()
	}
}

// SetSessionCookie 设置会话cookie
func (sm *SessionManager) SetSessionCookie(c *gin.Context, sessionID string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/admin",
		HttpOnly: true,
		Secure:   c.Request.TLS != nil, // 如果是HTTPS则设置Secure
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(sm.timeout.Seconds()),
	}
	http.SetCookie(c.Writer, cookie)
}

// GetSessionFromCookie 从cookie中获取会话ID
func (sm *SessionManager) GetSessionFromCookie(c *gin.Context) string {
	cookie, err := c.Request.Cookie("admin_session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ClearSessionCookie 清除会话cookie
func (sm *SessionManager) ClearSessionCookie(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/admin",
		HttpOnly: true,
		Secure:   c.Request.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // 立即过期
	}
	http.SetCookie(c.Writer, cookie)
}

// GetSessionCount 获取当前活跃会话数量（用于监控）
func (sm *SessionManager) GetSessionCount() int {
	sm.mutex.RLock()
	count := len(sm.sessions)
	sm.mutex.RUnlock()
	return count
}
