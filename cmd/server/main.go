package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/config"
	_ "server/docs"
	"server/internal/handler"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/router"
	"server/internal/service"
	"server/pkg/database"
	"server/pkg/logger"
	"server/pkg/template"

	"github.com/gin-gonic/gin"
)

// @title Server API
// @version 1.0
// @description Go 开发模板 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.File); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 连接数据库
	db, err := database.Connect(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Log.Fatal(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	// 初始化模板引擎
	tmplEngine, err := template.New(template.Config{
		TemplateDir: "web/templates",
		HotReload:   cfg.Server.Mode == "debug",
	})
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("初始化模板引擎失败: %v", err))
	}

	// 初始化依赖
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWT)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	webHandler := handler.NewWebHandler(tmplEngine)
	healthHandler := handler.NewHealthHandler(db)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	engine := gin.New()
	r := router.NewRouter(authHandler, userHandler, webHandler, healthHandler, authService)
	r.Setup(engine)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 启动服务器
	go func() {
		logger.Log.Info(fmt.Sprintf("服务启动在 %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal(fmt.Sprintf("启动服务失败: %v", err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("正在关闭服务器...")

	// 优雅关闭，最多等待 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal(fmt.Sprintf("服务器强制关闭: %v", err))
	}

	logger.Log.Info("服务器已退出")
}
