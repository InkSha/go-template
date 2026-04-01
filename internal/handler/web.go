package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebHandler 前端页面处理器
type WebHandler struct {
	templateEngine interface {
		Render(w io.Writer, name string, data interface{}) error
	}
}

// NewWebHandler 创建前端处理器
func NewWebHandler(engine interface {
	Render(w io.Writer, name string, data interface{}) error
}) *WebHandler {
	return &WebHandler{templateEngine: engine}
}

// LoginPage 登录页面
func (h *WebHandler) LoginPage(c *gin.Context) {
	data := gin.H{"Title": "用户登录"}
	if err := h.templateEngine.Render(c.Writer, "login", data); err != nil {
		c.String(http.StatusInternalServerError, "渲染失败: %v", err)
	}
}

// RegisterPage 注册页面
func (h *WebHandler) RegisterPage(c *gin.Context) {
	data := gin.H{"Title": "用户注册"}
	if err := h.templateEngine.Render(c.Writer, "register", data); err != nil {
		c.String(http.StatusInternalServerError, "渲染失败: %v", err)
	}
}

// DashboardPage 仪表盘页面
func (h *WebHandler) DashboardPage(c *gin.Context) {
	data := gin.H{"Title": "仪表盘"}
	if err := h.templateEngine.Render(c.Writer, "dashboard", data); err != nil {
		c.String(http.StatusInternalServerError, "渲染失败: %v", err)
	}
}

// UsersPage 用户管理页面
func (h *WebHandler) UsersPage(c *gin.Context) {
	data := gin.H{"Title": "用户管理"}
	if err := h.templateEngine.Render(c.Writer, "users", data); err != nil {
		c.String(http.StatusInternalServerError, "渲染失败: %v", err)
	}
}
