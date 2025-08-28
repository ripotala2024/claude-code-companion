package proxy

import (
	"fmt"
	"time"

	"claude-code-companion/internal/config"
	"claude-code-companion/internal/conversion"
	"claude-code-companion/internal/endpoint"
	"claude-code-companion/internal/health"
	"claude-code-companion/internal/i18n"
	"claude-code-companion/internal/logger"
	"claude-code-companion/internal/modelrewrite"
	"claude-code-companion/internal/security"
	"claude-code-companion/internal/tagging"
	"claude-code-companion/internal/validator"
	"claude-code-companion/internal/web"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config          *config.Config
	endpointManager *endpoint.Manager
	logger          *logger.Logger
	validator       *validator.ResponseValidator
	healthChecker   *health.Checker
	adminServer     *web.AdminServer
	taggingManager  *tagging.Manager         // 新增：tagging系统管理器
	modelRewriter   *modelrewrite.Rewriter   // 新增：模型重写器
	converter       conversion.Converter     // 新增：格式转换器
	i18nManager     *i18n.Manager            // 新增：国际化管理器
	sessionManager  *security.SessionManager // 新增：会话管理器
	authManager     *security.AuthManager    // 新增：身份验证管理器
	router          *gin.Engine
	configFilePath  string
}

func NewServer(cfg *config.Config, configFilePath string, version string) (*Server, error) {
	logConfig := logger.LogConfig{
		Level:           cfg.Logging.Level,
		LogRequestTypes: cfg.Logging.LogRequestTypes,
		LogRequestBody:  cfg.Logging.LogRequestBody,
		LogResponseBody: cfg.Logging.LogResponseBody,
		LogDirectory:    cfg.Logging.LogDirectory,
	}

	log, err := logger.NewLogger(logConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	endpointManager, err := endpoint.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize endpoint manager: %v", err)
	}
	responseValidator := validator.NewResponseValidator()

	// 初始化tagging系统
	taggingManager := tagging.NewManager()
	if err := taggingManager.Initialize(&cfg.Tagging); err != nil {
		return nil, fmt.Errorf("failed to initialize tagging system: %v", err)
	}

	// 初始化模型重写器
	modelRewriter := modelrewrite.NewRewriter(*log)

	// 初始化格式转换器
	converter := conversion.NewConverter(log)

	// 初始化健康检查器（需要在模型重写器和转换器之后）
	healthChecker := health.NewChecker(cfg.Timeouts.ToHealthCheckTimeoutConfig(), modelRewriter, converter)

	// 初始化国际化管理器
	i18nConfig := &i18n.Config{
		DefaultLanguage: i18n.Language(cfg.I18n.DefaultLanguage),
		LocalesPath:     cfg.I18n.LocalesPath,
		Enabled:         cfg.I18n.Enabled,
	}
	// 如果配置为空，使用默认配置
	if cfg.I18n.DefaultLanguage == "" {
		i18nConfig = i18n.DefaultConfig()
	}

	i18nManager, err := i18n.NewManager(i18nConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize i18n manager: %v", err)
	}

	// 创建会话管理器
	defaultTimeout, _ := time.ParseDuration(config.Default.Auth.SessionTimeout)
	sessionTimeout := config.GetTimeoutDuration(cfg.Auth.SessionTimeout, defaultTimeout)
	sessionManager := security.NewSessionManager(sessionTimeout)

	// 创建身份验证管理器
	authManager := security.NewAuthManager(sessionManager, cfg.Auth.Username, cfg.Auth.Password, cfg.Auth.Enabled)

	// 创建管理界面服务器（永远启用）
	adminServer := web.NewAdminServer(cfg, endpointManager, taggingManager, log, configFilePath, version, i18nManager, authManager)

	server := &Server{
		config:          cfg,
		endpointManager: endpointManager,
		logger:          log,
		validator:       responseValidator,
		healthChecker:   healthChecker,
		adminServer:     adminServer,
		taggingManager:  taggingManager, // 新增：设置tagging管理器
		modelRewriter:   modelRewriter,  // 新增：设置模型重写器
		converter:       converter,      // 新增：设置格式转换器
		i18nManager:     i18nManager,    // 新增：设置国际化管理器
		sessionManager:  sessionManager, // 新增：设置会话管理器
		authManager:     authManager,    // 新增：设置身份验证管理器
		configFilePath:  configFilePath,
	}

	// 设置热更新处理器
	adminServer.SetHotUpdateHandler(server)

	// 让端点管理器使用同一个健康检查器
	endpointManager.SetHealthChecker(healthChecker)

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	gin.SetMode(gin.ReleaseMode)

	s.router = gin.New()
	s.router.Use(gin.Recovery())

	// 注册管理界面路由（不需要认证）
	s.adminServer.RegisterRoutes(s.router)

	// 为 API 端点添加日志中间件
	apiGroup := s.router.Group("/v1")
	apiGroup.Use(s.loggingMiddleware())
	{
		apiGroup.Any("/*path", s.handleProxy)
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info(fmt.Sprintf("Starting proxy server on %s:%d", s.config.Server.Host, s.config.Server.Port))
	return s.router.Run(addr)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) GetEndpointManager() *endpoint.Manager {
	return s.endpointManager
}

func (s *Server) GetLogger() *logger.Logger {
	return s.logger
}

func (s *Server) GetHealthChecker() *health.Checker {
	return s.healthChecker
}

// HotUpdateConfig safely updates configuration without restarting the server
func (s *Server) HotUpdateConfig(newConfig *config.Config) error {
	// 验证新配置
	if err := s.validateConfigForHotUpdate(newConfig); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	s.logger.Info("Starting configuration hot update")

	// 更新端点配置
	if err := s.updateEndpoints(newConfig.Endpoints); err != nil {
		return fmt.Errorf("failed to update endpoints: %v", err)
	}

	// 更新日志配置（如果可能）
	if err := s.updateLoggingConfig(newConfig.Logging); err != nil {
		s.logger.Error("Failed to update logging config, continuing with endpoint updates", err)
	}

	// 更新验证器配置
	s.updateValidatorConfig(newConfig.Validation)

	// 更新身份验证配置
	if err := s.updateAuthConfig(newConfig.Auth); err != nil {
		s.logger.Error("Failed to update auth config, continuing with other updates", err)
	}

	// 更新内存中的配置
	s.config = newConfig

	s.logger.Info("Configuration hot update completed successfully")
	return nil
}

// validateConfigForHotUpdate validates the new configuration
func (s *Server) validateConfigForHotUpdate(newConfig *config.Config) error {
	// 检查是否尝试修改不可热更新的配置
	if newConfig.Server.Host != s.config.Server.Host {
		return fmt.Errorf("server host cannot be changed via hot update")
	}
	if newConfig.Server.Port != s.config.Server.Port {
		return fmt.Errorf("server port cannot be changed via hot update")
	}

	// 验证端点配置
	if len(newConfig.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be configured")
	}

	return nil
}

// updateEndpoints updates endpoint configuration
func (s *Server) updateEndpoints(newEndpoints []config.EndpointConfig) error {
	s.endpointManager.UpdateEndpoints(newEndpoints)
	return nil
}

// updateLoggingConfig updates logging configuration if possible
func (s *Server) updateLoggingConfig(newLogging config.LoggingConfig) error {
	// 目前只能更新日志级别和记录策略，不能更换日志目录
	if newLogging.LogDirectory != s.config.Logging.LogDirectory {
		return fmt.Errorf("log directory cannot be changed via hot update")
	}

	// 可以安全更新的日志配置
	s.config.Logging.Level = newLogging.Level
	s.config.Logging.LogRequestTypes = newLogging.LogRequestTypes
	s.config.Logging.LogRequestBody = newLogging.LogRequestBody
	s.config.Logging.LogResponseBody = newLogging.LogResponseBody

	return nil
}

// updateValidatorConfig updates response validator configuration
func (s *Server) updateValidatorConfig(newValidation config.ValidationConfig) {
	s.validator = validator.NewResponseValidator()
	s.config.Validation = newValidation
}

// saveConfigToFile 将当前配置保存到文件
func (s *Server) saveConfigToFile() error {
	return config.SaveConfig(s.config, s.configFilePath)
}

// createOAuthTokenRefreshCallback 创建 OAuth token 刷新后的回调函数
func (s *Server) createOAuthTokenRefreshCallback() func(*endpoint.Endpoint) error {
	return func(ep *endpoint.Endpoint) error {
		// 更新内存中的配置
		for i, cfgEndpoint := range s.config.Endpoints {
			if cfgEndpoint.Name == ep.Name {
				s.config.Endpoints[i].OAuthConfig = ep.OAuthConfig
				break
			}
		}

		// 保存到配置文件
		return s.saveConfigToFile()
	}
}

// updateAuthConfig 更新身份验证配置
func (s *Server) updateAuthConfig(newAuthConfig config.AuthConfig) error {
	// 更新身份验证管理器的凭据
	return s.authManager.UpdateCredentials(newAuthConfig.Username, newAuthConfig.Password, newAuthConfig.Enabled)
}
