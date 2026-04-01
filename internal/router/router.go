package router

import (
	"server/internal/handler"
	"server/internal/middleware"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	authHandler   *handler.AuthHandler
	userHandler   *handler.UserHandler
	webHandler    *handler.WebHandler
	healthHandler *handler.HealthHandler
	authService   *service.AuthService
}

func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, webHandler *handler.WebHandler, healthHandler *handler.HealthHandler, authService *service.AuthService) *Router {
	return &Router{
		authHandler:   authHandler,
		userHandler:   userHandler,
		webHandler:    webHandler,
		healthHandler: healthHandler,
		authService:   authService,
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.CORS())
	engine.Use(middleware.Logger())
	engine.Use(middleware.RateLimiter())
	engine.Use(gin.Recovery())

	// 静态文件
	engine.Static("/static", "./web/static")

	// 前端页面路由（公开）
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})
	engine.GET("/login", r.webHandler.LoginPage)
	engine.GET("/register", r.webHandler.RegisterPage)

	// 管理页面路由（需要认证）
	admin := engine.Group("/admin")
	admin.Use(middleware.JWTAuth(r.authService))
	{
		admin.GET("/dashboard", r.webHandler.DashboardPage)
		admin.GET("/users", r.webHandler.UsersPage)
	}

	// Swagger文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	engine.GET("/health", r.healthHandler.HealthCheck)

	// 版本信息
	engine.GET("/version", r.healthHandler.Version)

	// API v1
	v1 := engine.Group("/api/v1")
	{
		// 认证相关（无需JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}

		// 用户相关（需要JWT）
		users := v1.Group("/users")
		users.Use(middleware.JWTAuth(r.authService))
		{
			users.GET("/me", r.userHandler.GetCurrentUser)
			users.GET("", r.userHandler.ListUsers)
			users.POST("", r.userHandler.CreateUser)
			users.GET("/:id", r.userHandler.GetUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
			users.PATCH("/:id/status", r.userHandler.UpdateUserStatus)
		}
	}
}
